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

package testing

import (
	"context"
	"testing"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"

	ccmtesting "github.com/miyadav/cloud-provider-testing-interface"
)

// TestLoadBalancerIntegration tests load balancer integration with mock cloud provider
func TestLoadBalancerIntegration(t *testing.T) {
	// Create mock cloud provider
	cloudProvider := NewMockCloudProvider()

	// Create test interface
	testImpl := NewCCMTestInterface(cloudProvider)

	// Create test configuration
	config := &ccmtesting.TestConfig{
		ProviderName:         "mock-cloud-provider",
		ClusterName:          "test-cluster",
		Region:               "test-region",
		Zone:                 "test-zone",
		TestTimeout:          5 * time.Minute,
		CleanupResources:     true,
		MockExternalServices: true,
		TestData: map[string]interface{}{
			"test-mode": "integration",
		},
	}

	// Setup test environment
	err := testImpl.SetupTestEnvironment(config)
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}
	defer testImpl.TeardownTestEnvironment()

	// Create test runner
	runner := ccmtesting.NewTestRunner(testImpl)

	// Add load balancer test suite
	runner.AddTestSuite(CreateLoadBalancerTestSuite())

	// Run tests
	ctx := context.Background()
	err = runner.RunTests(ctx)
	if err != nil {
		t.Fatalf("Test execution failed: %v", err)
	}

	// Verify results
	summary := runner.GetSummary()
	results := runner.GetResults()

	t.Logf("Test Summary: %d total, %d passed, %d failed, %d skipped",
		summary.TotalTests, summary.PassedTests, summary.FailedTests, summary.SkippedTests)

	// Check that all tests passed
	if summary.FailedTests > 0 {
		t.Errorf("Some tests failed: %d failed out of %d total", summary.FailedTests, summary.TotalTests)
	}

	// Print detailed results
	for _, result := range results {
		status := "PASSED"
		if !result.Success {
			status = "FAILED"
		}
		if result.Test.Skip {
			status = "SKIPPED"
		}
		t.Logf("  %s: %s (%v)", status, result.Test.Name, result.Duration)
	}
}

// TestNodeManagementIntegration tests node management integration
func TestNodeManagementIntegration(t *testing.T) {
	// Create mock cloud provider
	cloudProvider := NewMockCloudProvider()

	// Create test interface
	testImpl := NewCCMTestInterface(cloudProvider)

	// Create test configuration
	config := &testing.TestConfig{
		ProviderName:         "mock-cloud-provider",
		ClusterName:          "test-cluster",
		Region:               "test-region",
		Zone:                 "test-zone",
		TestTimeout:          3 * time.Minute,
		CleanupResources:     true,
		MockExternalServices: true,
	}

	// Setup test environment
	err := testImpl.SetupTestEnvironment(config)
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}
	defer testImpl.TeardownTestEnvironment()

	// Create test runner
	runner := testing.NewTestRunner(testImpl)

	// Add node management test suite
	runner.AddTestSuite(CreateNodeTestSuite())

	// Run tests
	ctx := context.Background()
	err = runner.RunTests(ctx)
	if err != nil {
		t.Fatalf("Test execution failed: %v", err)
	}

	// Verify results
	summary := runner.GetSummary()

	t.Logf("Node Management Test Summary: %d total, %d passed, %d failed, %d skipped",
		summary.TotalTests, summary.PassedTests, summary.FailedTests, summary.SkippedTests)

	if summary.FailedTests > 0 {
		t.Errorf("Some node management tests failed: %d failed out of %d total", summary.FailedTests, summary.TotalTests)
	}
}

// TestRouteManagementIntegration tests route management integration
func TestRouteManagementIntegration(t *testing.T) {
	// Create mock cloud provider
	cloudProvider := NewMockCloudProvider()

	// Create test interface
	testImpl := NewCCMTestInterface(cloudProvider)

	// Create test configuration
	config := &testing.TestConfig{
		ProviderName:         "mock-cloud-provider",
		ClusterName:          "test-cluster",
		Region:               "test-region",
		Zone:                 "test-zone",
		TestTimeout:          3 * time.Minute,
		CleanupResources:     true,
		MockExternalServices: true,
	}

	// Setup test environment
	err := testImpl.SetupTestEnvironment(config)
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}
	defer testImpl.TeardownTestEnvironment()

	// Create test runner
	runner := testing.NewTestRunner(testImpl)

	// Add route management test suite
	runner.AddTestSuite(CreateRouteTestSuite())

	// Run tests
	ctx := context.Background()
	err = runner.RunTests(ctx)
	if err != nil {
		t.Fatalf("Test execution failed: %v", err)
	}

	// Verify results
	summary := runner.GetSummary()

	t.Logf("Route Management Test Summary: %d total, %d passed, %d failed, %d skipped",
		summary.TotalTests, summary.PassedTests, summary.FailedTests, summary.SkippedTests)

	if summary.FailedTests > 0 {
		t.Errorf("Some route management tests failed: %d failed out of %d total", summary.FailedTests, summary.TotalTests)
	}
}

// TestAllSuitesIntegration tests all test suites together
func TestAllSuitesIntegration(t *testing.T) {
	// Create mock cloud provider
	cloudProvider := NewMockCloudProvider()

	// Create test interface
	testImpl := NewCCMTestInterface(cloudProvider)

	// Create test configuration
	config := &testing.TestConfig{
		ProviderName:         "mock-cloud-provider",
		ClusterName:          "test-cluster",
		Region:               "test-region",
		Zone:                 "test-zone",
		TestTimeout:          10 * time.Minute,
		CleanupResources:     true,
		MockExternalServices: true,
	}

	// Setup test environment
	err := testImpl.SetupTestEnvironment(config)
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}
	defer testImpl.TeardownTestEnvironment()

	// Create test runner
	runner := testing.NewTestRunner(testImpl)

	// Add all test suites
	runner.AddTestSuite(CreateLoadBalancerTestSuite())
	runner.AddTestSuite(CreateNodeTestSuite())
	runner.AddTestSuite(CreateRouteTestSuite())
	runner.AddTestSuite(CreateInstancesTestSuite())
	runner.AddTestSuite(CreateZonesTestSuite())
	runner.AddTestSuite(CreateClustersTestSuite())

	// Run tests
	ctx := context.Background()
	err = runner.RunTests(ctx)
	if err != nil {
		t.Fatalf("Test execution failed: %v", err)
	}

	// Verify results
	summary := runner.GetSummary()

	t.Logf("All Suites Test Summary: %d total, %d passed, %d failed, %d skipped",
		summary.TotalTests, summary.PassedTests, summary.FailedTests, summary.SkippedTests)

	if summary.FailedTests > 0 {
		t.Errorf("Some tests failed: %d failed out of %d total", summary.FailedTests, summary.TotalTests)
	}

	// Print test results
	results := runner.GetResults()
	for _, result := range results {
		status := "PASSED"
		if !result.Success {
			status = "FAILED"
		}
		if result.Test.Skip {
			status = "SKIPPED"
		}
		t.Logf("  %s: %s (%v)", status, result.Test.Name, result.Duration)
	}
}

// TestCloudProviderInterface tests the cloud provider interface directly
func TestCloudProviderInterface(t *testing.T) {
	// Create mock cloud provider
	cloudProvider := NewMockCloudProvider()

	// Test LoadBalancer interface
	lb, ok := cloudProvider.LoadBalancer()
	if !ok {
		t.Fatal("LoadBalancer interface not available")
	}

	// Test Instances interface
	instances, ok := cloudProvider.Instances()
	if !ok {
		t.Fatal("Instances interface not available")
	}

	// Test Zones interface
	zones, ok := cloudProvider.Zones()
	if !ok {
		t.Fatal("Zones interface not available")
	}

	// Test Routes interface
	routes, ok := cloudProvider.Routes()
	if !ok {
		t.Fatal("Routes interface not available")
	}

	// Test Clusters interface
	clusters, ok := cloudProvider.Clusters()
	if !ok {
		t.Fatal("Clusters interface not available")
	}

	// Test basic functionality
	ctx := context.Background()

	// Test instances functionality
	providerID, err := instances.InstanceID(ctx, "test-node")
	if err != nil {
		t.Errorf("Failed to get instance ID: %v", err)
	}
	if providerID == "" {
		t.Error("Provider ID is empty")
	}

	// Test zones functionality
	zone, err := zones.GetZone(ctx)
	if err != nil {
		t.Errorf("Failed to get zone: %v", err)
	}
	if zone.Region == "" {
		t.Error("Zone region is empty")
	}

	// Test routes functionality
	routeList, err := routes.ListRoutes(ctx, "test-cluster")
	if err != nil {
		t.Errorf("Failed to list routes: %v", err)
	}
	if len(routeList) == 0 {
		t.Error("No routes returned")
	}

	// Test clusters functionality
	clusterList, err := clusters.ListClusters(ctx)
	if err != nil {
		t.Errorf("Failed to list clusters: %v", err)
	}
	if len(clusterList) == 0 {
		t.Error("No clusters returned")
	}

	t.Logf("Cloud provider interface test completed successfully")
}

// TestTestInterfaceMethods tests the test interface methods directly
func TestTestInterfaceMethods(t *testing.T) {
	// Create mock cloud provider
	cloudProvider := NewMockCloudProvider()

	// Create test interface
	testImpl := NewCCMTestInterface(cloudProvider)

	// Test configuration
	config := &testing.TestConfig{
		ProviderName:         "mock-cloud-provider",
		ClusterName:          "test-cluster",
		Region:               "test-region",
		Zone:                 "test-zone",
		TestTimeout:          1 * time.Minute,
		CleanupResources:     true,
		MockExternalServices: true,
	}

	// Setup test environment
	err := testImpl.SetupTestEnvironment(config)
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}
	defer testImpl.TeardownTestEnvironment()

	ctx := context.Background()

	// Test node creation
	nodeConfig := &testing.TestNodeConfig{
		Name:         "test-node",
		ProviderID:   "test-provider://test-node",
		InstanceType: "test-instance-type",
		Zone:         "test-zone",
		Region:       "test-region",
	}

	node, err := testImpl.CreateTestNode(ctx, nodeConfig)
	if err != nil {
		t.Errorf("Failed to create test node: %v", err)
	}
	if node.Name != "test-node" {
		t.Errorf("Expected node name 'test-node', got '%s'", node.Name)
	}

	// Test service creation
	serviceConfig := &testing.TestServiceConfig{
		Name:      "test-service",
		Namespace: "default",
		Type:      "ClusterIP",
	}

	service, err := testImpl.CreateTestService(ctx, serviceConfig)
	if err != nil {
		t.Errorf("Failed to create test service: %v", err)
	}
	if service.Name != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", service.Name)
	}

	// Test route creation
	routeConfig := &testing.TestRouteConfig{
		Name:            "test-route",
		ClusterName:     "test-cluster",
		TargetNode:      "test-node",
		DestinationCIDR: "10.0.0.0/24",
	}

	route, err := testImpl.CreateTestRoute(ctx, routeConfig)
	if err != nil {
		t.Errorf("Failed to create test route: %v", err)
	}
	if route.Name != "test-route" {
		t.Errorf("Expected route name 'test-route', got '%s'", route.Name)
	}

	// Test condition waiting
	condition := testing.TestCondition{
		Type:    "TestCondition",
		Timeout: 5 * time.Second,
		CheckFunction: func() (bool, error) {
			return true, nil
		},
	}

	err = testImpl.WaitForCondition(ctx, condition)
	if err != nil {
		t.Errorf("Failed to wait for condition: %v", err)
	}

	// Test results
	results := testImpl.GetTestResults()
	if results == nil {
		t.Error("Test results are nil")
	}

	// Test state reset
	err = testImpl.ResetTestState()
	if err != nil {
		t.Errorf("Failed to reset test state: %v", err)
	}

	t.Logf("Test interface methods test completed successfully")
}

// TestWithRealKubeClient tests with a real Kubernetes client
func TestWithRealKubeClient(t *testing.T) {
	// Create mock cloud provider
	cloudProvider := NewMockCloudProvider()

	// Create fake Kubernetes client
	fakeClient := fake.NewSimpleClientset()

	// Create informer factory
	informerFactory := informers.NewSharedInformerFactory(fakeClient, 0)

	// Create test interface
	testImpl := NewCCMTestInterface(cloudProvider)

	// Override the Kubernetes client
	testImpl.kubeClient = fakeClient
	testImpl.informerFactory = informerFactory

	// Create test configuration
	config := &testing.TestConfig{
		ProviderName:         "mock-cloud-provider",
		ClusterName:          "test-cluster",
		Region:               "test-region",
		Zone:                 "test-zone",
		TestTimeout:          1 * time.Minute,
		CleanupResources:     true,
		MockExternalServices: true,
		InformerFactory:      informerFactory,
	}

	// Setup test environment
	err := testImpl.SetupTestEnvironment(config)
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}
	defer testImpl.TeardownTestEnvironment()

	// Test that the informer factory is properly set
	if testImpl.GetInformerFactory() != informerFactory {
		t.Error("Informer factory not properly set")
	}

	// Test that the Kubernetes client is properly set
	if testImpl.GetKubeClient() != fakeClient {
		t.Error("Kubernetes client not properly set")
	}

	t.Logf("Real Kubernetes client test completed successfully")
}

// BenchmarkLoadBalancerTests benchmarks the load balancer tests
func BenchmarkLoadBalancerTests(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Create mock cloud provider
		cloudProvider := NewMockCloudProvider()

		// Create test interface
		testImpl := NewCCMTestInterface(cloudProvider)

		// Create test configuration
		config := &testing.TestConfig{
			ProviderName:         "mock-cloud-provider",
			ClusterName:          "test-cluster",
			Region:               "test-region",
			Zone:                 "test-zone",
			TestTimeout:          1 * time.Minute,
			CleanupResources:     true,
			MockExternalServices: true,
		}

		// Setup test environment
		err := testImpl.SetupTestEnvironment(config)
		if err != nil {
			b.Fatalf("Failed to setup test environment: %v", err)
		}
		defer testImpl.TeardownTestEnvironment()

		// Create test runner
		runner := testing.NewTestRunner(testImpl)

		// Add load balancer test suite
		runner.AddTestSuite(CreateLoadBalancerTestSuite())

		// Run tests
		ctx := context.Background()
		err = runner.RunTests(ctx)
		if err != nil {
			b.Fatalf("Test execution failed: %v", err)
		}

		// Verify results
		summary := runner.GetSummary()
		if summary.FailedTests > 0 {
			b.Errorf("Some tests failed: %d failed out of %d total", summary.FailedTests, summary.TotalTests)
		}
	}
}
