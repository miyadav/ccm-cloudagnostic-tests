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
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

// MockCloudProvider implements the cloudprovider.Interface for testing purposes.
type MockCloudProvider struct {
	// Mock interfaces
	instances    *MockInstances
	zones        *MockZones
	loadBalancer *MockLoadBalancer
	routes       *MockRoutes
	clusters     *MockClusters

	// Test data
	nodes    map[string]*v1.Node
	services map[string]*v1.Service
	routeMap map[string]*cloudprovider.Route
}

// NewMockCloudProvider creates a new mock cloud provider.
func NewMockCloudProvider() *MockCloudProvider {
	return &MockCloudProvider{
		instances:    NewMockInstances(),
		zones:        NewMockZones(),
		loadBalancer: NewMockLoadBalancer(),
		routes:       NewMockRoutes(),
		clusters:     NewMockClusters(),
		nodes:        make(map[string]*v1.Node),
		services:     make(map[string]*v1.Service),
		routeMap:     make(map[string]*cloudprovider.Route),
	}
}

// Initialize initializes the cloud provider.
func (m *MockCloudProvider) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
	klog.Info("Mock cloud provider initialized")
}

// LoadBalancer returns the load balancer interface.
func (m *MockCloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return m.loadBalancer, true
}

// Instances returns the instances interface.
func (m *MockCloudProvider) Instances() (cloudprovider.Instances, bool) {
	return m.instances, true
}

// InstancesV2 returns the instances v2 interface.
func (m *MockCloudProvider) InstancesV2() (cloudprovider.InstancesV2, bool) {
	return nil, false
}

// Zones returns the zones interface.
func (m *MockCloudProvider) Zones() (cloudprovider.Zones, bool) {
	return m.zones, true
}

// Clusters returns the clusters interface.
func (m *MockCloudProvider) Clusters() (cloudprovider.Clusters, bool) {
	return m.clusters, true
}

// Routes returns the routes interface.
func (m *MockCloudProvider) Routes() (cloudprovider.Routes, bool) {
	return m.routes, true
}

// ProviderName returns the provider name.
func (m *MockCloudProvider) ProviderName() string {
	return "mock-cloud-provider"
}

// HasClusterID returns whether the cluster has a cluster ID.
func (m *MockCloudProvider) HasClusterID() bool {
	return true
}

// MockInstances implements the cloudprovider.Instances interface.
type MockInstances struct {
	mu sync.RWMutex
}

// NewMockInstances creates a new mock instances interface.
func NewMockInstances() *MockInstances {
	return &MockInstances{}
}

// NodeAddresses returns the addresses of the specified instance.
func (m *MockInstances) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return mock addresses
	return []v1.NodeAddress{
		{Type: v1.NodeInternalIP, Address: "10.0.0.1"},
		{Type: v1.NodeExternalIP, Address: "192.168.1.1"},
		{Type: v1.NodeHostName, Address: string(name)},
	}, nil
}

// NodeAddressesByProviderID returns the addresses of the specified instance.
func (m *MockInstances) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	return m.NodeAddresses(ctx, types.NodeName("mock-node"))
}

// InstanceID returns the cloud provider ID of the specified instance.
func (m *MockInstances) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	return fmt.Sprintf("mock-provider://%s", nodeName), nil
}

// InstanceType returns the type of the specified instance.
func (m *MockInstances) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	return "mock-instance-type", nil
}

// InstanceTypeByProviderID returns the type of the specified instance.
func (m *MockInstances) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	return "mock-instance-type", nil
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances.
func (m *MockInstances) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return nil
}

// CurrentNodeName returns the name of the node we are currently running on.
func (m *MockInstances) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

// InstanceExistsByProviderID returns true if the instance for the given provider ID still exists.
func (m *MockInstances) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	return true, nil
}

// InstanceShutdownByProviderID returns true if the instance is shutdown in cloudprovider.
func (m *MockInstances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, nil
}

// MockZones implements the cloudprovider.Zones interface.
type MockZones struct {
}

// NewMockZones creates a new mock zones interface.
func NewMockZones() *MockZones {
	return &MockZones{}
}

// GetZone returns the Zone containing the current failure zone and locality region.
func (m *MockZones) GetZone(ctx context.Context) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{
		FailureDomain: "mock-zone",
		Region:        "mock-region",
	}, nil
}

// GetZoneByProviderID returns the Zone containing the current failure zone and locality region.
func (m *MockZones) GetZoneByProviderID(ctx context.Context, providerID string) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{
		FailureDomain: "mock-zone",
		Region:        "mock-region",
	}, nil
}

// GetZoneByNodeName returns the Zone containing the current failure zone and locality region.
func (m *MockZones) GetZoneByNodeName(ctx context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	return cloudprovider.Zone{
		FailureDomain: "mock-zone",
		Region:        "mock-region",
	}, nil
}

// MockLoadBalancer implements the cloudprovider.LoadBalancer interface.
type MockLoadBalancer struct {
	mu sync.RWMutex
}

// NewMockLoadBalancer creates a new mock load balancer interface.
func NewMockLoadBalancer() *MockLoadBalancer {
	return &MockLoadBalancer{}
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one.
func (m *MockLoadBalancer) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return mock load balancer status
	status := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{IP: "192.168.1.100"},
			{Hostname: "mock-lb.example.com"},
		},
	}

	return status, nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
func (m *MockLoadBalancer) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	return nil
}

// EnsureLoadBalancerDeleted deletes the specified load balancer if it exists.
func (m *MockLoadBalancer) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	return nil
}

// GetLoadBalancerName returns the name of the load balancer.
func (m *MockLoadBalancer) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	return fmt.Sprintf("mock-lb-%s", service.Name)
}

// GetLoadBalancer returns whether the specified load balancer exists, and if so, what its status is.
func (m *MockLoadBalancer) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (*v1.LoadBalancerStatus, bool, error) {
	status := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{IP: "192.168.1.100"},
		},
	}
	return status, true, nil
}

// MockRoutes implements the cloudprovider.Routes interface.
type MockRoutes struct {
	mu sync.RWMutex
}

// NewMockRoutes creates a new mock routes interface.
func NewMockRoutes() *MockRoutes {
	return &MockRoutes{}
}

// ListRoutes lists all managed routes that belong to the specified clusterName.
func (m *MockRoutes) ListRoutes(ctx context.Context, clusterName string) ([]*cloudprovider.Route, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return []*cloudprovider.Route{
		{
			Name:            "mock-route-1",
			TargetNode:      "mock-node-1",
			DestinationCIDR: "10.0.0.0/24",
		},
	}, nil
}

// CreateRoute creates the described managed route.
func (m *MockRoutes) CreateRoute(ctx context.Context, clusterName string, nameHint string, route *cloudprovider.Route) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}

// DeleteRoute deletes the specified managed route.
func (m *MockRoutes) DeleteRoute(ctx context.Context, clusterName string, route *cloudprovider.Route) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}

// MockClusters implements the cloudprovider.Clusters interface.
type MockClusters struct {
}

// NewMockClusters creates a new mock clusters interface.
func NewMockClusters() *MockClusters {
	return &MockClusters{}
}

// ListClusters lists the names of the available clusters.
func (m *MockClusters) ListClusters(ctx context.Context) ([]string, error) {
	return []string{"mock-cluster-1", "mock-cluster-2"}, nil
}

// Master gets back the address (either DNS name or IP address) of the master node for the cluster.
func (m *MockClusters) Master(ctx context.Context, clusterName string) (string, error) {
	return "mock-master.example.com", nil
}

// Clusters is a convenience function that returns a clusters interface.
func (m *MockClusters) Clusters() (cloudprovider.Clusters, bool) {
	return m, true
}
