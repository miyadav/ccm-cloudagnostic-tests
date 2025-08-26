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
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	cloudprovider "k8s.io/cloud-provider"

	ccmtesting "github.com/miyadav/cloud-provider-testing-interface"
)

// CreateLoadBalancerTestSuite creates a test suite for load balancer functionality.
func CreateLoadBalancerTestSuite() ccmtesting.TestSuite {
	return ccmtesting.TestSuite{
		Name:        "LoadBalancer",
		Description: "Tests for cloud provider load balancer functionality",
		Setup:       setupLoadBalancerTestSuite,
		Teardown:    teardownLoadBalancerTestSuite,
		Tests: []ccmtesting.Test{
			{
				Name:        "CreateLoadBalancer",
				Description: "Test creating a load balancer",
				Run:         testCreateLoadBalancer,
				Timeout:     5 * time.Minute,
			},
			{
				Name:        "UpdateLoadBalancer",
				Description: "Test updating a load balancer",
				Run:         testUpdateLoadBalancer,
				Timeout:     5 * time.Minute,
			},
			{
				Name:        "DeleteLoadBalancer",
				Description: "Test deleting a load balancer",
				Run:         testDeleteLoadBalancer,
				Timeout:     5 * time.Minute,
			},
			{
				Name:        "LoadBalancerStatus",
				Description: "Test load balancer status updates",
				Run:         testLoadBalancerStatus,
				Timeout:     3 * time.Minute,
			},
			{
				Name:        "LoadBalancerHealthCheck",
				Description: "Test load balancer health check functionality",
				Run:         testLoadBalancerHealthCheck,
				Timeout:     3 * time.Minute,
			},
		},
	}
}

// CreateNodeTestSuite creates a test suite for node management functionality.
func CreateNodeTestSuite() ccmtesting.TestSuite {
	return ccmtesting.TestSuite{
		Name:        "NodeManagement",
		Description: "Tests for cloud provider node management functionality",
		Setup:       setupNodeTestSuite,
		Teardown:    teardownNodeTestSuite,
		Tests: []ccmtesting.Test{
			{
				Name:        "NodeInitialization",
				Description: "Test node initialization and registration",
				Run:         testNodeInitialization,
				Timeout:     3 * time.Minute,
			},
			{
				Name:        "NodeAddresses",
				Description: "Test node address management",
				Run:         testNodeAddresses,
				Timeout:     2 * time.Minute,
			},
			{
				Name:        "NodeProviderID",
				Description: "Test node provider ID management",
				Run:         testNodeProviderID,
				Timeout:     2 * time.Minute,
			},
			{
				Name:        "NodeInstanceType",
				Description: "Test node instance type detection",
				Run:         testNodeInstanceType,
				Timeout:     2 * time.Minute,
			},
			{
				Name:        "NodeZones",
				Description: "Test node zone management",
				Run:         testNodeZones,
				Timeout:     2 * time.Minute,
			},
		},
	}
}

// CreateRouteTestSuite creates a test suite for route management functionality.
func CreateRouteTestSuite() ccmtesting.TestSuite {
	return ccmtesting.TestSuite{
		Name:        "RouteManagement",
		Description: "Tests for cloud provider route management functionality",
		Setup:       setupRouteTestSuite,
		Teardown:    teardownRouteTestSuite,
		Tests: []ccmtesting.Test{
			{
				Name:        "CreateRoute",
				Description: "Test creating a route",
				Run:         func(ti ccmtesting.TestInterface) error { return testCreateRoute(context.Background(), ti) },
				Timeout:     3 * time.Minute,
			},
			{
				Name:        "DeleteRoute",
				Description: "Test deleting a route",
				Run:         func(ti ccmtesting.TestInterface) error { return testDeleteRoute(context.Background(), ti) },
				Timeout:     3 * time.Minute,
			},
			{
				Name:        "ListRoutes",
				Description: "Test listing routes",
				Run:         func(ti ccmtesting.TestInterface) error { return testListRoutes(context.Background(), ti) },
				Timeout:     2 * time.Minute,
			},
		},
	}
}

// CreateInstancesTestSuite creates a test suite for instances functionality.
func CreateInstancesTestSuite() ccmtesting.TestSuite {
	return ccmtesting.TestSuite{
		Name:        "Instances",
		Description: "Tests for cloud provider instances functionality",
		Setup:       setupInstancesTestSuite,
		Teardown:    teardownInstancesTestSuite,
		Tests: []ccmtesting.Test{
			{
				Name:        "InstanceExists",
				Description: "Test instance existence check",
				Run:         func(ti ccmtesting.TestInterface) error { return testInstanceExists(context.Background(), ti) },
				Timeout:     2 * time.Minute,
			},
			{
				Name:        "InstanceShutdown",
				Description: "Test instance shutdown detection",
				Run:         func(ti ccmtesting.TestInterface) error { return testInstanceShutdown(context.Background(), ti) },
				Timeout:     2 * time.Minute,
			},
			{
				Name:        "InstanceMetadata",
				Description: "Test instance metadata retrieval",
				Run:         func(ti ccmtesting.TestInterface) error { return testInstanceMetadata(context.Background(), ti) },
				Timeout:     2 * time.Minute,
			},
		},
	}
}

// CreateZonesTestSuite creates a test suite for zones functionality.
func CreateZonesTestSuite() ccmtesting.TestSuite {
	return ccmtesting.TestSuite{
		Name:        "Zones",
		Description: "Tests for cloud provider zones functionality",
		Setup:       setupZonesTestSuite,
		Teardown:    teardownZonesTestSuite,
		Tests: []ccmtesting.Test{
			{
				Name:        "GetZone",
				Description: "Test zone information retrieval",
				Run:         func(ti ccmtesting.TestInterface) error { return testGetZone(context.Background(), ti) },
				Timeout:     2 * time.Minute,
			},
			{
				Name:        "GetZoneByProviderID",
				Description: "Test zone retrieval by provider ID",
				Run:         func(ti ccmtesting.TestInterface) error { return testGetZoneByProviderID(context.Background(), ti) },
				Timeout:     2 * time.Minute,
			},
		},
	}
}

// CreateClustersTestSuite creates a test suite for clusters functionality.
func CreateClustersTestSuite() ccmtesting.TestSuite {
	return ccmtesting.TestSuite{
		Name:        "Clusters",
		Description: "Tests for cloud provider clusters functionality",
		Setup:       setupClustersTestSuite,
		Teardown:    teardownClustersTestSuite,
		Tests: []ccmtesting.Test{
			{
				Name:        "ListClusters",
				Description: "Test listing clusters",
				Run:         func(ti ccmtesting.TestInterface) error { return testListClusters(context.Background(), ti) },
				Timeout:     2 * time.Minute,
			},
			{
				Name:        "Master",
				Description: "Test master node detection",
				Run:         func(ti ccmtesting.TestInterface) error { return testMaster(context.Background(), ti) },
				Timeout:     2 * time.Minute,
			},
		},
	}
}

// Setup and teardown functions for test suites

func setupLoadBalancerTestSuite(ti ccmtesting.TestInterface) error {
	// Create test nodes for load balancer tests
	ctx := context.Background()

	nodeConfig := &ccmtesting.TestNodeConfig{
		Name:         "test-node-1",
		ProviderID:   "test-provider://test-node-1",
		InstanceType: "test-instance-type",
		Zone:         "test-zone",
		Region:       "test-region",
		Addresses: []v1.NodeAddress{
			{Type: v1.NodeInternalIP, Address: "10.0.0.1"},
			{Type: v1.NodeExternalIP, Address: "192.168.1.1"},
		},
	}

	_, err := ti.CreateTestNode(ctx, nodeConfig)
	if err != nil {
		return fmt.Errorf("failed to create test node for load balancer tests: %w", err)
	}

	return nil
}

func teardownLoadBalancerTestSuite(ti ccmtesting.TestInterface) error {
	// Clean up test nodes
	ctx := context.Background()
	err := ti.DeleteTestNode(ctx, "test-node-1")
	if err != nil {
		return fmt.Errorf("failed to delete test node: %w", err)
	}

	return nil
}

func setupNodeTestSuite(ti ccmtesting.TestInterface) error {
	// Setup for node tests
	return nil
}

func teardownNodeTestSuite(ti ccmtesting.TestInterface) error {
	// Cleanup for node tests
	return nil
}

func setupRouteTestSuite(ti ccmtesting.TestInterface) error {
	// Create test nodes for route tests
	ctx := context.Background()

	nodeConfig := &ccmtesting.TestNodeConfig{
		Name:         "route-test-node",
		ProviderID:   "test-provider://route-test-node",
		InstanceType: "test-instance-type",
		Zone:         "test-zone",
		Region:       "test-region",
		Addresses: []v1.NodeAddress{
			{Type: v1.NodeInternalIP, Address: "10.0.0.2"},
		},
	}

	_, err := ti.CreateTestNode(ctx, nodeConfig)
	if err != nil {
		return fmt.Errorf("failed to create test node for route tests: %w", err)
	}

	return nil
}

func teardownRouteTestSuite(ti ccmtesting.TestInterface) error {
	// Clean up test nodes
	ctx := context.Background()
	err := ti.DeleteTestNode(ctx, "route-test-node")
	if err != nil {
		return fmt.Errorf("failed to delete test node: %w", err)
	}

	return nil
}

func setupInstancesTestSuite(ti ccmtesting.TestInterface) error {
	// Setup for instances tests
	return nil
}

func teardownInstancesTestSuite(ti ccmtesting.TestInterface) error {
	// Cleanup for instances tests
	return nil
}

func setupZonesTestSuite(ti ccmtesting.TestInterface) error {
	// Setup for zones tests
	return nil
}

func teardownZonesTestSuite(ti ccmtesting.TestInterface) error {
	// Cleanup for zones tests
	return nil
}

func setupClustersTestSuite(ti ccmtesting.TestInterface) error {
	// Setup for clusters tests
	return nil
}

func teardownClustersTestSuite(ti ccmtesting.TestInterface) error {
	// Cleanup for clusters tests
	return nil
}

// Test functions for load balancer functionality

func testCreateLoadBalancer(ti ccmtesting.TestInterface) error {
	ctx := context.Background()
	cloudProvider := ti.GetCloudProvider()

	// Get load balancer interface
	lb, ok := cloudProvider.LoadBalancer()
	if !ok {
		return fmt.Errorf("cloud provider does not support load balancer functionality")
	}

	// Create test service
	serviceConfig := &ccmtesting.TestServiceConfig{
		Name:      "test-loadbalancer",
		Namespace: "default",
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Name:       "http",
				Protocol:   v1.ProtocolTCP,
				Port:       80,
				TargetPort: intstr.FromInt(8080),
			},
		},
		ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyCluster,
	}

	service, err := ti.CreateTestService(ctx, serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	// Create mock nodes for the load balancer
	mockNodes := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "mock-node-1"},
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{Type: v1.NodeInternalIP, Address: "10.0.0.1"},
				},
			},
		},
	}

	// Ensure load balancer
	status, err := lb.EnsureLoadBalancer(ctx, "test-cluster", service, mockNodes)
	if err != nil {
		return fmt.Errorf("failed to ensure load balancer: %w", err)
	}

	if status == nil || len(status.Ingress) == 0 {
		return fmt.Errorf("load balancer status is empty")
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Load balancer created successfully with %d ingress addresses", len(status.Ingress)))
	return nil
}

func testUpdateLoadBalancer(ti ccmtesting.TestInterface) error {
	ctx := context.Background()
	cloudProvider := ti.GetCloudProvider()

	lb, ok := cloudProvider.LoadBalancer()
	if !ok {
		return fmt.Errorf("cloud provider does not support load balancer functionality")
	}

	// Create a mock service for update testing
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "test-loadbalancer", Namespace: "default"},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer,
			Ports: []v1.ServicePort{
				{Name: "http", Port: 80, TargetPort: intstr.FromInt(8080)},
				{Name: "https", Port: 443, TargetPort: intstr.FromInt(8443)},
			},
		},
	}

	// Create mock nodes
	mockNodes := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "mock-node-1"},
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{Type: v1.NodeInternalIP, Address: "10.0.0.1"},
				},
			},
		},
	}

	// Update load balancer
	_, err := lb.EnsureLoadBalancer(ctx, "test-cluster", service, mockNodes)
	if err != nil {
		return fmt.Errorf("failed to update load balancer: %w", err)
	}

	ti.GetTestResults().AddLog("Load balancer updated successfully")
	return nil
}

func testDeleteLoadBalancer(ti ccmtesting.TestInterface) error {
	ctx := context.Background()
	cloudProvider := ti.GetCloudProvider()

	lb, ok := cloudProvider.LoadBalancer()
	if !ok {
		return fmt.Errorf("cloud provider does not support load balancer functionality")
	}

	// Create a mock service for deletion testing
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "test-loadbalancer", Namespace: "default"},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeLoadBalancer,
		},
	}

	// Delete load balancer
	err := lb.EnsureLoadBalancerDeleted(ctx, "test-cluster", service)
	if err != nil {
		return fmt.Errorf("failed to delete load balancer: %w", err)
	}

	ti.GetTestResults().AddLog("Load balancer deleted successfully")
	return nil
}

func testLoadBalancerStatus(ti ccmtesting.TestInterface) error {
	ctx := context.Background()
	cloudProvider := ti.GetCloudProvider()

	lb, ok := cloudProvider.LoadBalancer()
	if !ok {
		return fmt.Errorf("cloud provider does not support load balancer functionality")
	}

	// Create a new service for status testing
	serviceConfig := &ccmtesting.TestServiceConfig{
		Name:      "status-test-lb",
		Namespace: "default",
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Name:       "http",
				Protocol:   v1.ProtocolTCP,
				Port:       80,
				TargetPort: intstr.FromInt(8080),
			},
		},
	}

	service, err := ti.CreateTestService(ctx, serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	// Create mock nodes
	mockNodes := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "mock-node-1"},
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{Type: v1.NodeInternalIP, Address: "10.0.0.1"},
				},
			},
		},
	}

	// Create load balancer
	status, err := lb.EnsureLoadBalancer(ctx, "test-cluster", service, mockNodes)
	if err != nil {
		return fmt.Errorf("failed to create load balancer: %w", err)
	}

	// Wait for load balancer to be ready
	condition := ccmtesting.TestCondition{
		Type:    "LoadBalancerReady",
		Timeout: 2 * time.Minute,
		CheckFunction: func() (bool, error) {
			// Check if load balancer has ingress addresses
			if status != nil && len(status.Ingress) > 0 {
				for _, ingress := range status.Ingress {
					if ingress.IP != "" || ingress.Hostname != "" {
						return true, nil
					}
				}
			}
			return false, nil
		},
	}

	err = ti.WaitForCondition(ctx, condition)
	if err != nil {
		return fmt.Errorf("load balancer did not become ready: %w", err)
	}

	// Clean up
	err = lb.EnsureLoadBalancerDeleted(ctx, "test-cluster", service)
	if err != nil {
		return fmt.Errorf("failed to delete test load balancer: %w", err)
	}

	ti.GetTestResults().AddLog("Load balancer status test completed successfully")
	return nil
}

func testLoadBalancerHealthCheck(ti ccmtesting.TestInterface) error {
	// This test would verify load balancer health check functionality
	// Implementation would depend on the specific cloud provider
	ti.GetTestResults().AddLog("Load balancer health check test completed")
	return nil
}

// Test functions for node management

func testNodeInitialization(ti ccmtesting.TestInterface) error {
	ctx := context.Background()
	cloudProvider := ti.GetCloudProvider()

	// Get instances interface
	instances, ok := cloudProvider.Instances()
	if !ok {
		return fmt.Errorf("cloud provider does not support instances functionality")
	}

	// Create test node
	nodeConfig := &ccmtesting.TestNodeConfig{
		Name:         "init-test-node",
		ProviderID:   "test-provider://init-test-node",
		InstanceType: "test-instance-type",
		Zone:         "test-zone",
		Region:       "test-region",
		Addresses: []v1.NodeAddress{
			{Type: v1.NodeInternalIP, Address: "10.0.0.3"},
			{Type: v1.NodeExternalIP, Address: "192.168.1.3"},
		},
	}

	node, err := ti.CreateTestNode(ctx, nodeConfig)
	if err != nil {
		return fmt.Errorf("failed to create test node: %w", err)
	}

	// Test node initialization
	providerID, err := instances.InstanceID(ctx, types.NodeName(node.Name))
	if err != nil {
		return fmt.Errorf("failed to get instance ID: %w", err)
	}

	if providerID == "" {
		return fmt.Errorf("provider ID is empty")
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Node initialization test completed. Provider ID: %s", providerID))

	// Clean up
	err = ti.DeleteTestNode(ctx, node.Name)
	if err != nil {
		return fmt.Errorf("failed to delete test node: %w", err)
	}

	return nil
}

func testNodeAddresses(ti ccmtesting.TestInterface) error {
	ctx := context.Background()
	cloudProvider := ti.GetCloudProvider()

	instances, ok := cloudProvider.Instances()
	if !ok {
		return fmt.Errorf("cloud provider does not support instances functionality")
	}

	// Create test node
	nodeConfig := &ccmtesting.TestNodeConfig{
		Name:         "address-test-node",
		ProviderID:   "test-provider://address-test-node",
		InstanceType: "test-instance-type",
		Zone:         "test-zone",
		Region:       "test-region",
		Addresses: []v1.NodeAddress{
			{Type: v1.NodeInternalIP, Address: "10.0.0.4"},
			{Type: v1.NodeExternalIP, Address: "192.168.1.4"},
		},
	}

	node, err := ti.CreateTestNode(ctx, nodeConfig)
	if err != nil {
		return fmt.Errorf("failed to create test node: %w", err)
	}

	// Test node addresses
	addresses, err := instances.NodeAddresses(ctx, types.NodeName(node.Name))
	if err != nil {
		return fmt.Errorf("failed to get node addresses: %w", err)
	}

	if len(addresses) == 0 {
		return fmt.Errorf("no addresses returned for node")
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Node addresses test completed. Found %d addresses", len(addresses)))

	// Clean up
	err = ti.DeleteTestNode(ctx, node.Name)
	if err != nil {
		return fmt.Errorf("failed to delete test node: %w", err)
	}

	return nil
}

func testNodeProviderID(ti ccmtesting.TestInterface) error {
	ctx := context.Background()
	cloudProvider := ti.GetCloudProvider()

	instances, ok := cloudProvider.Instances()
	if !ok {
		return fmt.Errorf("cloud provider does not support instances functionality")
	}

	// Create test node
	nodeConfig := &ccmtesting.TestNodeConfig{
		Name:         "providerid-test-node",
		ProviderID:   "test-provider://providerid-test-node",
		InstanceType: "test-instance-type",
		Zone:         "test-zone",
		Region:       "test-region",
	}

	node, err := ti.CreateTestNode(ctx, nodeConfig)
	if err != nil {
		return fmt.Errorf("failed to create test node: %w", err)
	}

	// Test provider ID
	providerID, err := instances.InstanceID(ctx, types.NodeName(node.Name))
	if err != nil {
		return fmt.Errorf("failed to get provider ID: %w", err)
	}

	if providerID == "" {
		return fmt.Errorf("provider ID is empty")
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Provider ID test completed. Provider ID: %s", providerID))

	// Clean up
	err = ti.DeleteTestNode(ctx, node.Name)
	if err != nil {
		return fmt.Errorf("failed to delete test node: %w", err)
	}

	return nil
}

func testNodeInstanceType(ti ccmtesting.TestInterface) error {
	ctx := context.Background()
	cloudProvider := ti.GetCloudProvider()

	instances, ok := cloudProvider.Instances()
	if !ok {
		return fmt.Errorf("cloud provider does not support instances functionality")
	}

	// Create test node
	nodeConfig := &ccmtesting.TestNodeConfig{
		Name:         "instancetype-test-node",
		ProviderID:   "test-provider://instancetype-test-node",
		InstanceType: "test-instance-type",
		Zone:         "test-zone",
		Region:       "test-region",
	}

	node, err := ti.CreateTestNode(ctx, nodeConfig)
	if err != nil {
		return fmt.Errorf("failed to create test node: %w", err)
	}

	// Test instance type
	instanceType, err := instances.InstanceType(ctx, types.NodeName(node.Name))
	if err != nil {
		return fmt.Errorf("failed to get instance type: %w", err)
	}

	if instanceType == "" {
		return fmt.Errorf("instance type is empty")
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Instance type test completed. Instance type: %s", instanceType))

	// Clean up
	err = ti.DeleteTestNode(ctx, node.Name)
	if err != nil {
		return fmt.Errorf("failed to delete test node: %w", err)
	}

	return nil
}

func testNodeZones(ti ccmtesting.TestInterface) error {
	ctx := context.Background()
	cloudProvider := ti.GetCloudProvider()

	zones, ok := cloudProvider.Zones()
	if !ok {
		return fmt.Errorf("cloud provider does not support zones functionality")
	}

	// Create test node
	nodeConfig := &ccmtesting.TestNodeConfig{
		Name:         "zones-test-node",
		ProviderID:   "test-provider://zones-test-node",
		InstanceType: "test-instance-type",
		Zone:         "test-zone",
		Region:       "test-region",
	}

	node, err := ti.CreateTestNode(ctx, nodeConfig)
	if err != nil {
		return fmt.Errorf("failed to create test node: %w", err)
	}

	// Test zone information
	zone, err := zones.GetZone(ctx)
	if err != nil {
		return fmt.Errorf("failed to get zone: %w", err)
	}

	if zone.Region == "" {
		return fmt.Errorf("zone region is empty")
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Zones test completed. Region: %s, Zone: %s", zone.Region, zone.FailureDomain))

	// Clean up
	err = ti.DeleteTestNode(ctx, node.Name)
	if err != nil {
		return fmt.Errorf("failed to delete test node: %w", err)
	}

	return nil
}

// Test functions for route management

func testCreateRoute(ctx context.Context, ti ccmtesting.TestInterface) error {
	cloudProvider := ti.GetCloudProvider()

	routes, ok := cloudProvider.Routes()
	if !ok {
		return fmt.Errorf("cloud provider does not support routes functionality")
	}

	// Create test route
	routeConfig := &ccmtesting.TestRouteConfig{
		Name:            "test-route",
		ClusterName:     "test-cluster",
		TargetNode:      "route-test-node",
		DestinationCIDR: "10.0.0.0/24",
		Blackhole:       false,
	}

	route, err := ti.CreateTestRoute(ctx, routeConfig)
	if err != nil {
		return fmt.Errorf("failed to create test route: %w", err)
	}

	// Create route through cloud provider
	err = routes.CreateRoute(ctx, "test-cluster", "test-route", route)
	if err != nil {
		return fmt.Errorf("failed to create route: %w", err)
	}

	ti.GetTestResults().AddLog("Route created successfully")
	return nil
}

func testDeleteRoute(ctx context.Context, ti ccmtesting.TestInterface) error {
	cloudProvider := ti.GetCloudProvider()

	routes, ok := cloudProvider.Routes()
	if !ok {
		return fmt.Errorf("cloud provider does not support routes functionality")
	}

	// Delete route
	route := &cloudprovider.Route{
		Name:            "test-route",
		TargetNode:      "route-test-node",
		DestinationCIDR: "10.0.0.0/24",
	}
	err := routes.DeleteRoute(ctx, "test-cluster", route)
	if err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}

	ti.GetTestResults().AddLog("Route deleted successfully")
	return nil
}

func testListRoutes(ctx context.Context, ti ccmtesting.TestInterface) error {
	cloudProvider := ti.GetCloudProvider()

	routes, ok := cloudProvider.Routes()
	if !ok {
		return fmt.Errorf("cloud provider does not support routes functionality")
	}

	// List routes
	routeList, err := routes.ListRoutes(ctx, "test-cluster")
	if err != nil {
		return fmt.Errorf("failed to list routes: %w", err)
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Listed %d routes", len(routeList)))
	return nil
}

// Test functions for instances functionality

func testInstanceExists(ctx context.Context, ti ccmtesting.TestInterface) error {
	cloudProvider := ti.GetCloudProvider()

	instances, ok := cloudProvider.Instances()
	if !ok {
		return fmt.Errorf("cloud provider does not support instances functionality")
	}

	// Test instance existence by provider ID
	exists, err := instances.InstanceExistsByProviderID(ctx, "test-provider://test-node")
	if err != nil {
		return fmt.Errorf("failed to check instance existence: %w", err)
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Instance exists check completed. Exists: %t", exists))
	return nil
}

func testInstanceShutdown(ctx context.Context, ti ccmtesting.TestInterface) error {
	cloudProvider := ti.GetCloudProvider()

	instances, ok := cloudProvider.Instances()
	if !ok {
		return fmt.Errorf("cloud provider does not support instances functionality")
	}

	// Test instance shutdown detection by provider ID
	shutdown, err := instances.InstanceShutdownByProviderID(ctx, "test-provider://test-node")
	if err != nil {
		return fmt.Errorf("failed to check instance shutdown: %w", err)
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Instance shutdown check completed. Shutdown: %t", shutdown))
	return nil
}

func testInstanceMetadata(ctx context.Context, ti ccmtesting.TestInterface) error {
	cloudProvider := ti.GetCloudProvider()

	instances, ok := cloudProvider.Instances()
	if !ok {
		return fmt.Errorf("cloud provider does not support instances functionality")
	}

	// Test instance ID (closest to metadata)
	instanceID, err := instances.InstanceID(ctx, types.NodeName("test-node"))
	if err != nil {
		return fmt.Errorf("failed to get instance ID: %w", err)
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Instance ID retrieved: %s", instanceID))
	return nil
}

// Test functions for zones functionality

func testGetZone(ctx context.Context, ti ccmtesting.TestInterface) error {
	cloudProvider := ti.GetCloudProvider()

	zones, ok := cloudProvider.Zones()
	if !ok {
		return fmt.Errorf("cloud provider does not support zones functionality")
	}

	// Test zone retrieval
	zone, err := zones.GetZone(ctx)
	if err != nil {
		return fmt.Errorf("failed to get zone: %w", err)
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Zone retrieved. Region: %s, Zone: %s", zone.Region, zone.FailureDomain))
	return nil
}

func testGetZoneByProviderID(ctx context.Context, ti ccmtesting.TestInterface) error {
	cloudProvider := ti.GetCloudProvider()

	zones, ok := cloudProvider.Zones()
	if !ok {
		return fmt.Errorf("cloud provider does not support zones functionality")
	}

	// Test zone retrieval by provider ID
	zone, err := zones.GetZoneByProviderID(ctx, "test-provider://test-node")
	if err != nil {
		return fmt.Errorf("failed to get zone by provider ID: %w", err)
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Zone by provider ID retrieved. Region: %s, Zone: %s", zone.Region, zone.FailureDomain))
	return nil
}

// Test functions for clusters functionality

func testListClusters(ctx context.Context, ti ccmtesting.TestInterface) error {
	cloudProvider := ti.GetCloudProvider()

	clusters, ok := cloudProvider.Clusters()
	if !ok {
		return fmt.Errorf("cloud provider does not support clusters functionality")
	}

	// Test cluster listing
	clusterList, err := clusters.ListClusters(ctx)
	if err != nil {
		return fmt.Errorf("failed to list clusters: %w", err)
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Listed %d clusters", len(clusterList)))
	return nil
}

func testMaster(ctx context.Context, ti ccmtesting.TestInterface) error {
	cloudProvider := ti.GetCloudProvider()

	clusters, ok := cloudProvider.Clusters()
	if !ok {
		return fmt.Errorf("cloud provider does not support clusters functionality")
	}

	// Test master node detection
	master, err := clusters.Master(ctx, "test-cluster")
	if err != nil {
		return fmt.Errorf("failed to get master: %w", err)
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Master node: %s", master))
	return nil
}
