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
limifications under the License.
*/

package testing

import (
	"context"
	"fmt"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"

	ccmtesting "github.com/miyadav/cloud-provider-testing-interface"
)

// CCMTestInterface implements the cloud-provider-testing-interface TestInterface
// to provide cloud-agnostic testing for Cloud Controller Manager functionality.
type CCMTestInterface struct {
	// Cloud provider instance to be tested
	cloudProvider cloudprovider.Interface

	// Kubernetes client for test operations
	kubeClient kubernetes.Interface

	// Informer factory for watching resources
	informerFactory informers.SharedInformerFactory

	// Test configuration
	config *ccmtesting.TestConfig

	// Test results
	results *ccmtesting.TestResults

	// Track created resources for cleanup
	createdResources map[string][]string
	mu               sync.RWMutex

	// Mock services for testing
	mockServices map[string]interface{}
}

// NewCCMTestInterface creates a new CCM test interface instance.
func NewCCMTestInterface(cloudProvider cloudprovider.Interface) *CCMTestInterface {
	return &CCMTestInterface{
		cloudProvider:    cloudProvider,
		kubeClient:       fake.NewSimpleClientset(),
		createdResources: make(map[string][]string),
		mockServices:     make(map[string]interface{}),
		results: &ccmtesting.TestResults{
			ResourceCounts: make(map[string]int),
			Metrics:        make(map[string]interface{}),
		},
	}
}

// SetupTestEnvironment initializes the test environment with the given configuration.
func (c *CCMTestInterface) SetupTestEnvironment(config *ccmtesting.TestConfig) error {
	c.config = config
	c.results = &ccmtesting.TestResults{
		ResourceCounts: make(map[string]int),
		Metrics:        make(map[string]interface{}),
	}

	// Initialize informer factory
	if config.InformerFactory != nil {
		c.informerFactory = config.InformerFactory
	} else {
		c.informerFactory = informers.NewSharedInformerFactory(c.kubeClient, 0)
	}

	// Start informers
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	c.informerFactory.Start(ctx.Done())

	// Wait for informers to sync
	syncResult := c.informerFactory.WaitForCacheSync(ctx.Done())
	for _, synced := range syncResult {
		if !synced {
			return fmt.Errorf("failed to sync informers")
		}
	}

	c.results.AddLog(fmt.Sprintf("Test environment setup completed for provider: %s", config.ProviderName))
	return nil
}

// TeardownTestEnvironment cleans up the test environment and removes any test resources.
func (c *CCMTestInterface) TeardownTestEnvironment() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clean up created resources if cleanup is enabled
	if c.config != nil && c.config.CleanupResources {
		for resourceType, resources := range c.createdResources {
			for _, resourceName := range resources {
				c.results.AddLog(fmt.Sprintf("Cleaning up %s: %s", resourceType, resourceName))
				// Cleanup logic would be implemented here based on resource type
			}
		}
	}

	c.results.AddLog("Test environment teardown completed")
	return nil
}

// GetCloudProvider returns the cloud provider instance to be tested.
func (c *CCMTestInterface) GetCloudProvider() cloudprovider.Interface {
	return c.cloudProvider
}

// CreateTestNode creates a test node with the specified configuration.
func (c *CCMTestInterface) CreateTestNode(ctx context.Context, nodeConfig *ccmtesting.TestNodeConfig) (*v1.Node, error) {
	node := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:        nodeConfig.Name,
			Labels:      nodeConfig.Labels,
			Annotations: nodeConfig.Annotations,
		},
		Spec: v1.NodeSpec{
			ProviderID: nodeConfig.ProviderID,
		},
		Status: v1.NodeStatus{
			Addresses:  nodeConfig.Addresses,
			Conditions: nodeConfig.Conditions,
		},
	}

	createdNode, err := c.kubeClient.CoreV1().Nodes().Create(ctx, node, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create test node: %w", err)
	}

	// Track created resource
	c.mu.Lock()
	c.createdResources["nodes"] = append(c.createdResources["nodes"], nodeConfig.Name)
	c.results.IncrementResourceCount("nodes")
	c.mu.Unlock()

	c.results.AddLog(fmt.Sprintf("Created test node: %s", nodeConfig.Name))
	return createdNode, nil
}

// DeleteTestNode deletes a test node.
func (c *CCMTestInterface) DeleteTestNode(ctx context.Context, nodeName string) error {
	err := c.kubeClient.CoreV1().Nodes().Delete(ctx, nodeName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete test node: %w", err)
	}

	c.results.AddLog(fmt.Sprintf("Deleted test node: %s", nodeName))
	return nil
}

// CreateTestService creates a test service with the specified configuration.
func (c *CCMTestInterface) CreateTestService(ctx context.Context, serviceConfig *ccmtesting.TestServiceConfig) (*v1.Service, error) {
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        serviceConfig.Name,
			Namespace:   serviceConfig.Namespace,
			Labels:      serviceConfig.Labels,
			Annotations: serviceConfig.Annotations,
		},
		Spec: v1.ServiceSpec{
			Type:                  serviceConfig.Type,
			Ports:                 serviceConfig.Ports,
			LoadBalancerIP:        serviceConfig.LoadBalancerIP,
			ExternalTrafficPolicy: serviceConfig.ExternalTrafficPolicy,
			InternalTrafficPolicy: serviceConfig.InternalTrafficPolicy,
		},
	}

	createdService, err := c.kubeClient.CoreV1().Services(serviceConfig.Namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create test service: %w", err)
	}

	// Track created resource
	c.mu.Lock()
	key := fmt.Sprintf("services/%s", serviceConfig.Namespace)
	c.createdResources[key] = append(c.createdResources[key], serviceConfig.Name)
	c.results.IncrementResourceCount("services")
	c.mu.Unlock()

	c.results.AddLog(fmt.Sprintf("Created test service: %s/%s", serviceConfig.Namespace, serviceConfig.Name))
	return createdService, nil
}

// DeleteTestService deletes a test service.
func (c *CCMTestInterface) DeleteTestService(ctx context.Context, serviceName string) error {
	// For simplicity, we'll delete from default namespace
	// In a real implementation, you'd need to track the namespace
	err := c.kubeClient.CoreV1().Services("default").Delete(ctx, serviceName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete test service: %w", err)
	}

	c.results.AddLog(fmt.Sprintf("Deleted test service: %s", serviceName))
	return nil
}

// CreateTestRoute creates a test route with the specified configuration.
func (c *CCMTestInterface) CreateTestRoute(ctx context.Context, routeConfig *ccmtesting.TestRouteConfig) (*cloudprovider.Route, error) {
	route := &cloudprovider.Route{
		Name:            routeConfig.Name,
		TargetNode:      routeConfig.TargetNode,
		DestinationCIDR: routeConfig.DestinationCIDR,
		Blackhole:       routeConfig.Blackhole,
	}

	// In a real implementation, you would create the route through the cloud provider
	// For now, we'll just track it
	c.mu.Lock()
	c.createdResources["routes"] = append(c.createdResources["routes"], routeConfig.Name)
	c.results.IncrementResourceCount("routes")
	c.mu.Unlock()

	c.results.AddLog(fmt.Sprintf("Created test route: %s", routeConfig.Name))
	return route, nil
}

// DeleteTestRoute deletes a test route.
func (c *CCMTestInterface) DeleteTestRoute(ctx context.Context, routeName string) error {
	// In a real implementation, you would delete the route through the cloud provider
	c.results.AddLog(fmt.Sprintf("Deleted test route: %s", routeName))
	return nil
}

// WaitForCondition waits for a condition to be met.
func (c *CCMTestInterface) WaitForCondition(ctx context.Context, condition ccmtesting.TestCondition) error {
	timeout := condition.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("condition wait timed out: %s", condition.Type)
		case <-ticker.C:
			if condition.CheckFunction != nil {
				met, err := condition.CheckFunction()
				if err != nil {
					klog.Warningf("Error checking condition: %v", err)
					continue
				}
				if met {
					c.results.AddLog(fmt.Sprintf("Condition met: %s", condition.Type))
					return nil
				}
			}
		}
	}
}

// GetTestResults returns the results of the test execution.
func (c *CCMTestInterface) GetTestResults() *ccmtesting.TestResults {
	return c.results
}

// ResetTestState resets the test state to a clean state.
func (c *CCMTestInterface) ResetTestState() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clear created resources tracking
	c.createdResources = make(map[string][]string)

	// Reset test results
	c.results = &ccmtesting.TestResults{
		ResourceCounts: make(map[string]int),
		Metrics:        make(map[string]interface{}),
	}

	c.results.AddLog("Test state reset completed")
	return nil
}

// GetKubeClient returns the Kubernetes client used for testing.
func (c *CCMTestInterface) GetKubeClient() kubernetes.Interface {
	return c.kubeClient
}

// GetInformerFactory returns the informer factory used for testing.
func (c *CCMTestInterface) GetInformerFactory() informers.SharedInformerFactory {
	return c.informerFactory
}

// GetConfig returns the test configuration.
func (c *CCMTestInterface) GetConfig() *ccmtesting.TestConfig {
	return c.config
}
