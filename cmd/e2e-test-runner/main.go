/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"

	"github.com/kubernetes/ccm-cloudagnostic-tests/pkg/testing"
	ccmtesting "github.com/miyadav/cloud-provider-testing-interface"
)

var (
	// Test configuration
	kubeconfig     = flag.String("kubeconfig", "", "Path to kubeconfig file")
	provider       = flag.String("provider", "", "Cloud provider (aws, gcp, azure, mock, existing)")
	region         = flag.String("region", "", "Cloud provider region")
	zone           = flag.String("zone", "", "Cloud provider zone")
	clusterName    = flag.String("cluster", "", "Cluster name")
	resourcePrefix = flag.String("prefix", "e2e-test", "Prefix for test resources")

	// Test execution
	suite   = flag.String("suite", "all", "Test suite to run")
	timeout = flag.Duration("timeout", 30*time.Minute, "Test timeout")
	verbose = flag.Bool("verbose", false, "Enable verbose output")
	cleanup = flag.Bool("cleanup", true, "Clean up resources after tests")

	// Output
	outputFormat = flag.String("output", "text", "Output format (text, json)")

	// Credentials (for real cloud providers)
	credentialsFile = flag.String("credentials", "", "Path to credentials file")

	// Logging
	logLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	flag.Parse()

	// Set log level
	setLogLevel(*logLevel)

	// Validate required flags
	if *provider == "" {
		klog.Fatal("--provider flag is required")
	}

	if *provider != "mock" && *provider != "existing" && *kubeconfig == "" {
		klog.Fatal("--kubeconfig flag is required for real cloud providers (aws, gcp, azure)")
	}

	// Create Kubernetes client
	var kubeClient kubernetes.Interface
	var err error

	if *provider == "mock" {
		klog.Info("Using mock cloud provider")
	} else {
		klog.Infof("Connecting to cluster using kubeconfig: %s", *kubeconfig)
		kubeClient, err = createKubeClient(*kubeconfig)
		if err != nil {
			klog.Fatalf("Failed to create Kubernetes client: %v", err)
		}

		// Verify cluster connectivity
		if err := verifyClusterConnection(kubeClient); err != nil {
			klog.Fatalf("Failed to connect to cluster: %v", err)
		}
	}

	// Create test configuration
	config := &ccmtesting.TestConfig{
		ProviderName:         *provider,
		ClusterName:          *clusterName,
		Region:               *region,
		Zone:                 *zone,
		TestTimeout:          *timeout,
		CleanupResources:     *cleanup,
		MockExternalServices: *provider == "mock",
		TestData: map[string]interface{}{
			"resource-prefix": *resourcePrefix,
			"test-mode":       "e2e",
		},
	}

	// Create cloud provider adapter
	cloudProvider, err := createCloudProvider(*provider, kubeClient)
	if err != nil {
		klog.Fatalf("Failed to create cloud provider: %v", err)
	}

	// Create test interface based on provider type
	var testImpl ccmtesting.TestInterface
	if *provider == "existing" {
		testImpl = testing.NewExistingCCMTestInterface(kubeClient, config)
	} else {
		testImpl = testing.NewCCMTestInterface(cloudProvider)
	}

	// Setup test environment
	klog.Info("Setting up test environment...")
	err = testImpl.SetupTestEnvironment(config)
	if err != nil {
		klog.Fatalf("Failed to setup test environment: %v", err)
	}
	defer func() {
		klog.Info("Tearing down test environment...")
		testImpl.TeardownTestEnvironment()
	}()

	// Create test runner
	runner := ccmtesting.NewTestRunner(testImpl)

	// Add test suites based on provider capabilities
	addTestSuites(runner, *suite, *provider)

	// Run tests
	klog.Info("Starting e2e tests...")
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	err = runner.RunTests(ctx)
	if err != nil {
		klog.Fatalf("Test execution failed: %v", err)
	}

	endTime := time.Now()

	// Print results
	results := runner.GetResults()
	summary := runner.GetSummary()

	printResults(results, summary, startTime, endTime, *outputFormat, *verbose)

	// Exit with appropriate code
	if summary.FailedTests > 0 {
		os.Exit(1)
	}
}

func createKubeClient(kubeconfigPath string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return clientset, nil
}

func verifyClusterConnection(client kubernetes.Interface) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	klog.Info("Successfully connected to cluster")
	return nil
}

func createCloudProvider(providerName string, kubeClient kubernetes.Interface) (cloudprovider.Interface, error) {
	switch strings.ToLower(providerName) {
	case "mock":
		return testing.NewMockCloudProvider(), nil
	case "existing":
		// For existing CCM testing, we don't need a cloud provider interface
		// The test interface will handle everything through the Kubernetes API
		return nil, nil
	case "aws":
		return createAWSCloudProvider(kubeClient)
	case "gcp":
		return createGCPCloudProvider(kubeClient)
	case "azure":
		return createAzureCloudProvider(kubeClient)
	default:
		return nil, fmt.Errorf("unsupported cloud provider: %s", providerName)
	}
}

func createAWSCloudProvider(kubeClient kubernetes.Interface) (cloudprovider.Interface, error) {
	// Load AWS credentials
	credentials, err := loadCredentials(*credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS credentials: %w", err)
	}

	config := &testing.RealCloudProviderConfig{
		ProviderName:     "aws",
		Region:           *region,
		Zone:             *zone,
		ClusterName:      *clusterName,
		Credentials:      credentials,
		TestTimeout:      int(timeout.Minutes()),
		CleanupResources: *cleanup,
		ResourcePrefix:   *resourcePrefix,
	}

	adapter, err := testing.NewAWSCloudProviderAdapter(kubeClient, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS cloud provider adapter: %w", err)
	}

	// Initialize the adapter
	if err := adapter.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize AWS cloud provider: %w", err)
	}

	return adapter.GetCloudProvider(), nil
}

func createGCPCloudProvider(kubeClient kubernetes.Interface) (cloudprovider.Interface, error) {
	// Load GCP credentials
	credentials, err := loadCredentials(*credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load GCP credentials: %w", err)
	}

	config := &testing.RealCloudProviderConfig{
		ProviderName:     "gcp",
		Region:           *region,
		Zone:             *zone,
		ClusterName:      *clusterName,
		Credentials:      credentials,
		TestTimeout:      int(timeout.Minutes()),
		CleanupResources: *cleanup,
		ResourcePrefix:   *resourcePrefix,
	}

	adapter, err := testing.NewGCPCloudProviderAdapter(kubeClient, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCP cloud provider adapter: %w", err)
	}

	// Initialize the adapter
	if err := adapter.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize GCP cloud provider: %w", err)
	}

	return adapter.GetCloudProvider(), nil
}

func createAzureCloudProvider(kubeClient kubernetes.Interface) (cloudprovider.Interface, error) {
	// Load Azure credentials
	credentials, err := loadCredentials(*credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load Azure credentials: %w", err)
	}

	config := &testing.RealCloudProviderConfig{
		ProviderName:     "azure",
		Region:           *region,
		Zone:             *zone,
		ClusterName:      *clusterName,
		Credentials:      credentials,
		TestTimeout:      int(timeout.Minutes()),
		CleanupResources: *cleanup,
		ResourcePrefix:   *resourcePrefix,
	}

	adapter, err := testing.NewAzureCloudProviderAdapter(kubeClient, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure cloud provider adapter: %w", err)
	}

	// Initialize the adapter
	if err := adapter.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize Azure cloud provider: %w", err)
	}

	return adapter.GetCloudProvider(), nil
}

func loadCredentials(credentialsFile string) (map[string]string, error) {
	if credentialsFile == "" {
		return make(map[string]string), nil
	}

	// Implementation would depend on your credential format
	// This is a placeholder - you'd implement based on your needs
	return make(map[string]string), nil
}

func addTestSuites(runner *ccmtesting.TestRunner, suite, provider string) {
	switch strings.ToLower(suite) {
	case "all":
		runner.AddTestSuite(testing.CreateLoadBalancerTestSuite())
		runner.AddTestSuite(testing.CreateNodeTestSuite())
		runner.AddTestSuite(testing.CreateRouteTestSuite())
		runner.AddTestSuite(testing.CreateInstancesTestSuite())
		runner.AddTestSuite(testing.CreateZonesTestSuite())
		runner.AddTestSuite(testing.CreateClustersTestSuite())
	case "loadbalancer":
		runner.AddTestSuite(testing.CreateLoadBalancerTestSuite())
	case "nodes":
		runner.AddTestSuite(testing.CreateNodeTestSuite())
	case "routes":
		runner.AddTestSuite(testing.CreateRouteTestSuite())
	case "instances":
		runner.AddTestSuite(testing.CreateInstancesTestSuite())
	case "zones":
		runner.AddTestSuite(testing.CreateZonesTestSuite())
	case "clusters":
		runner.AddTestSuite(testing.CreateClustersTestSuite())
	default:
		klog.Fatalf("Unknown test suite: %s", suite)
	}
}

func printResults(results []ccmtesting.TestResult, summary ccmtesting.TestSummary, startTime, endTime time.Time, format string, verbose bool) {
	totalDuration := endTime.Sub(startTime)

	switch format {
	case "json":
		printJSONResults(results, summary, totalDuration)
	default:
		printTextResults(results, summary, totalDuration, verbose)
	}
}

func printTextResults(results []ccmtesting.TestResult, summary ccmtesting.TestSummary, totalDuration time.Duration, verbose bool) {
	fmt.Printf("\n=== CCM E2E Test Results ===\n")
	fmt.Printf("Total Duration: %v\n", totalDuration)
	fmt.Printf("Test Summary: %d total, %d passed, %d failed, %d skipped\n",
		summary.TotalTests, summary.PassedTests, summary.FailedTests, summary.SkippedTests)

	if verbose {
		fmt.Printf("\nDetailed Results:\n")
		for _, result := range results {
			status := "PASSED"
			if !result.Success {
				status = "FAILED"
			}
			if result.Test.Skip {
				status = "SKIPPED"
			}
			fmt.Printf("  %s: %s (%v)\n", status, result.Test.Name, result.Duration)

			// Note: TestResult doesn't have a Logs field in the current interface
			// Logs are handled through the test interface's GetTestResults() method
		}
	}

	if summary.FailedTests > 0 {
		fmt.Printf("\n❌ Some tests failed: %d failed out of %d total\n", summary.FailedTests, summary.TotalTests)
	} else {
		fmt.Printf("\n✅ All tests passed: %d passed out of %d total\n", summary.PassedTests, summary.TotalTests)
	}
}

func printJSONResults(results []ccmtesting.TestResult, summary ccmtesting.TestSummary, totalDuration time.Duration) {
	// Implementation for JSON output
	fmt.Printf("JSON output not yet implemented\n")
}

func setLogLevel(level string) {
	// Note: klog.SetLevel is not available in klog/v2
	// Log level is controlled by environment variables or flags
	switch strings.ToLower(level) {
	case "debug":
		// Set debug level via environment variable
		os.Setenv("KLOG_V", "4")
	case "info":
		// Default info level
		os.Setenv("KLOG_V", "2")
	case "warn":
		os.Setenv("KLOG_V", "1")
	case "error":
		os.Setenv("KLOG_V", "0")
	default:
		os.Setenv("KLOG_V", "2")
	}
}
