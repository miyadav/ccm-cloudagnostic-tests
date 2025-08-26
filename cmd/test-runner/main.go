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
	"log"
	"os"
	"strings"
	"time"

	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"

	"github.com/kubernetes/ccm-cloudagnostic-tests/pkg/testing"
	ccmtesting "github.com/miyadav/cloud-provider-testing-interface"
)

var (
	// Test suite flags
	testSuite = flag.String("suite", "all", "Test suite to run (all, loadbalancer, nodes, routes, instances, zones, clusters)")
	verbose   = flag.Bool("verbose", false, "Enable verbose output")
	timeout   = flag.Duration("timeout", 10*time.Minute, "Test timeout")

	// Cloud provider flags
	providerName = flag.String("provider", "mock", "Cloud provider name")
	clusterName  = flag.String("cluster", "test-cluster", "Cluster name")
	region       = flag.String("region", "test-region", "Region")
	zone         = flag.String("zone", "test-zone", "Zone")

	// Configuration flags
	cleanupResources     = flag.Bool("cleanup", true, "Clean up resources after tests")
	mockExternalServices = flag.Bool("mock-external", true, "Use mock external services")
	outputFormat         = flag.String("output", "text", "Output format (text, json)")

	// Test selection flags
	runSpecificTests = flag.String("tests", "", "Comma-separated list of specific tests to run")
	skipTests        = flag.String("skip", "", "Comma-separated list of tests to skip")

	// Logging flags
	logLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	flag.Parse()

	// Set log level
	setLogLevel(*logLevel)

	// Create cloud provider instance
	cloudProvider, err := createCloudProvider(*providerName)
	if err != nil {
		log.Fatalf("Failed to create cloud provider: %v", err)
	}

	// Create test interface
	testImpl := testing.NewCCMTestInterface(cloudProvider)

	// Create test configuration
	config := &ccmtesting.TestConfig{
		ProviderName:         *providerName,
		ClusterName:          *clusterName,
		Region:               *region,
		Zone:                 *zone,
		TestTimeout:          *timeout,
		CleanupResources:     *cleanupResources,
		MockExternalServices: *mockExternalServices,
		TestData: map[string]interface{}{
			"output-format": *outputFormat,
			"verbose":       *verbose,
		},
	}

	// Setup test environment
	err = testImpl.SetupTestEnvironment(config)
	if err != nil {
		log.Fatalf("Failed to setup test environment: %v", err)
	}
	defer testImpl.TeardownTestEnvironment()

	// Create test runner
	runner := ccmtesting.NewTestRunner(testImpl)

	// Add test suites based on flag
	err = addTestSuites(runner, *testSuite)
	if err != nil {
		log.Fatalf("Failed to add test suites: %v", err)
	}

	// Filter tests if specified
	if *runSpecificTests != "" || *skipTests != "" {
		filterTests(runner, *runSpecificTests, *skipTests)
	}

	// Run tests
	ctx := context.Background()
	startTime := time.Now()
	err = runner.RunTests(ctx)
	endTime := time.Now()

	if err != nil {
		log.Printf("Test execution failed: %v", err)
		os.Exit(1)
	}

	// Get results
	results := runner.GetResults()
	summary := runner.GetSummary()

	// Print results
	printResults(results, summary, startTime, endTime, *outputFormat, *verbose)

	// Exit with error code if any tests failed
	if summary.FailedTests > 0 {
		os.Exit(1)
	}
}

// createCloudProvider creates a cloud provider instance based on the provider name
func createCloudProvider(providerName string) (cloudprovider.Interface, error) {
	switch strings.ToLower(providerName) {
	case "mock":
		return testing.NewMockCloudProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported cloud provider: %s", providerName)
	}
}

// addTestSuites adds test suites to the runner based on the suite flag
func addTestSuites(runner *ccmtesting.TestRunner, suite string) error {
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
		return fmt.Errorf("unknown test suite: %s", suite)
	}
	return nil
}

// filterTests filters tests based on run and skip flags
func filterTests(runner *ccmtesting.TestRunner, runTests, skipTests string) {
	// This is a placeholder for test filtering functionality
	// In a real implementation, you would filter the tests in the runner
	if runTests != "" {
		klog.Infof("Running specific tests: %s", runTests)
	}
	if skipTests != "" {
		klog.Infof("Skipping tests: %s", skipTests)
	}
}

// printResults prints test results in the specified format
func printResults(results []ccmtesting.TestResult, summary ccmtesting.TestSummary, startTime, endTime time.Time, format string, verbose bool) {
	totalDuration := endTime.Sub(startTime)

	switch format {
	case "json":
		printJSONResults(results, summary, totalDuration)
	case "text":
		fallthrough
	default:
		printTextResults(results, summary, totalDuration, verbose)
	}
}

// printTextResults prints results in text format
func printTextResults(results []ccmtesting.TestResult, summary ccmtesting.TestSummary, totalDuration time.Duration, verbose bool) {
	fmt.Printf("\n=== CCM Cloud-Agnostic Test Results ===\n")
	fmt.Printf("Total Duration: %v\n", totalDuration)
	fmt.Printf("Test Summary: %d total, %d passed, %d failed, %d skipped\n",
		summary.TotalTests, summary.PassedTests, summary.FailedTests, summary.SkippedTests)
	fmt.Printf("Test Suite Duration: %v\n", summary.TotalDuration)

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
			if result.Error != nil {
				fmt.Printf("    Error: %v\n", result.Error)
			}
		}
	}

	// Print summary
	if summary.FailedTests > 0 {
		fmt.Printf("\n❌ Some tests failed: %d failed out of %d total\n", summary.FailedTests, summary.TotalTests)
	} else {
		fmt.Printf("\n✅ All tests passed: %d passed out of %d total\n", summary.PassedTests, summary.TotalTests)
	}
}

// printJSONResults prints results in JSON format
func printJSONResults(results []ccmtesting.TestResult, summary ccmtesting.TestSummary, totalDuration time.Duration) {
	// This is a placeholder for JSON output
	// In a real implementation, you would marshal the results to JSON
	fmt.Printf("{\n")
	fmt.Printf("  \"totalDuration\": \"%v\",\n", totalDuration)
	fmt.Printf("  \"summary\": {\n")
	fmt.Printf("    \"totalTests\": %d,\n", summary.TotalTests)
	fmt.Printf("    \"passedTests\": %d,\n", summary.PassedTests)
	fmt.Printf("    \"failedTests\": %d,\n", summary.FailedTests)
	fmt.Printf("    \"skippedTests\": %d,\n", summary.SkippedTests)
	fmt.Printf("    \"testSuiteDuration\": \"%v\"\n", summary.TotalDuration)
	fmt.Printf("  }\n")
	fmt.Printf("}\n")
}

// setLogLevel sets the log level
func setLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		klog.InitFlags(nil)
		flag.Set("v", "4")
	case "info":
		klog.InitFlags(nil)
		flag.Set("v", "2")
	case "warn":
		klog.InitFlags(nil)
		flag.Set("v", "1")
	case "error":
		klog.InitFlags(nil)
		flag.Set("v", "0")
	default:
		klog.InitFlags(nil)
		flag.Set("v", "2")
	}
}
