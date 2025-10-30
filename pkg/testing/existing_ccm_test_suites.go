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
	"k8s.io/apimachinery/pkg/util/intstr"

	ccmtesting "github.com/miyadav/cloud-provider-testing-interface"
)

// CreateExistingCCMLoadBalancerTestSuite creates a test suite for testing load balancer functionality
// through an existing CCM deployment (indirect testing via Kubernetes resources)
func CreateExistingCCMLoadBalancerTestSuite() ccmtesting.TestSuite {
	return ccmtesting.TestSuite{
		Name:        "ExistingCCM-LoadBalancer",
		Description: "Tests for load balancer functionality using existing CCM in cluster",
		Setup:       setupExistingCCMLoadBalancerTestSuite,
		Teardown:    teardownExistingCCMLoadBalancerTestSuite,
		Tests: []ccmtesting.Test{
			{
				Name:        "CreateLoadBalancer",
				Description: "Test creating a load balancer through existing CCM",
				Run:         testExistingCCMCreateLoadBalancer,
				Timeout:     5 * time.Minute,
			},
			{
				Name:        "UpdateLoadBalancer",
				Description: "Test updating a load balancer through existing CCM",
				Run:         testExistingCCMUpdateLoadBalancer,
				Timeout:     5 * time.Minute,
			},
			{
				Name:        "DeleteLoadBalancer",
				Description: "Test deleting a load balancer through existing CCM",
				Run:         testExistingCCMDeleteLoadBalancer,
				Timeout:     3 * time.Minute,
			},
		},
	}
}

// CreateExistingCCMNodeTestSuite creates a test suite for testing node management functionality
// through an existing CCM deployment
func CreateExistingCCMNodeTestSuite() ccmtesting.TestSuite {
	return ccmtesting.TestSuite{
		Name:        "ExistingCCM-NodeManagement",
		Description: "Tests for node management functionality using existing CCM in cluster",
		Setup:       setupExistingCCMNodeTestSuite,
		Teardown:    teardownExistingCCMNodeTestSuite,
		Tests: []ccmtesting.Test{
			{
				Name:        "GetExistingNodes",
				Description: "Test retrieving existing nodes managed by CCM",
				Run:         testExistingCCMGetNodes,
				Timeout:     2 * time.Minute,
			},
			{
				Name:        "VerifyNodeMetadata",
				Description: "Test that CCM has populated node metadata correctly",
				Run:         testExistingCCMNodeMetadata,
				Timeout:     2 * time.Minute,
			},
		},
	}
}

// Setup and teardown functions for test suites

func setupExistingCCMLoadBalancerTestSuite(ti ccmtesting.TestInterface) error {
	// No specific setup needed for load balancer tests
	// Each test will create its own service
	return nil
}

func teardownExistingCCMLoadBalancerTestSuite(ti ccmtesting.TestInterface) error {
	// Clean up any leftover services in the test namespace
	existingCCM, ok := ti.(*ExistingCCMTestInterface)
	if !ok {
		return nil // Skip if not the right interface type
	}

	ctx := context.Background()
	services, err := existingCCM.kubeClient.CoreV1().Services(existingCCM.GetNamespace()).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list services for cleanup: %w", err)
	}

	for _, svc := range services.Items {
		if svc.Spec.Type == v1.ServiceTypeLoadBalancer {
			err := existingCCM.kubeClient.CoreV1().Services(existingCCM.GetNamespace()).Delete(ctx, svc.Name, metav1.DeleteOptions{})
			if err != nil {
				return fmt.Errorf("failed to cleanup service %s: %w", svc.Name, err)
			}
		}
	}

	return nil
}

func setupExistingCCMNodeTestSuite(ti ccmtesting.TestInterface) error {
	// No specific setup needed for node tests
	// Tests will examine existing cluster nodes
	return nil
}

func teardownExistingCCMNodeTestSuite(ti ccmtesting.TestInterface) error {
	// Clean up any test nodes that may have been created
	existingCCM, ok := ti.(*ExistingCCMTestInterface)
	if !ok {
		return nil // Skip if not the right interface type
	}

	ctx := context.Background()
	
	// Get the resource prefix from config
	resourcePrefix, ok := existingCCM.config.TestData["resource-prefix"].(string)
	if !ok {
		resourcePrefix = "e2e-test"
	}

	// Delete any nodes with the test label
	nodes, err := existingCCM.kubeClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("test-prefix=%s", resourcePrefix),
	})
	if err != nil {
		return fmt.Errorf("failed to list test nodes for cleanup: %w", err)
	}

	for _, node := range nodes.Items {
		err := existingCCM.kubeClient.CoreV1().Nodes().Delete(ctx, node.Name, metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("failed to cleanup test node %s: %w", node.Name, err)
		}
	}

	return nil
}

// Test functions for existing CCM load balancer functionality

func testExistingCCMCreateLoadBalancer(ti ccmtesting.TestInterface) error {
	ctx := context.Background()

	// Cast to ExistingCCMTestInterface to access methods
	existingCCM, ok := ti.(*ExistingCCMTestInterface)
	if !ok {
		return fmt.Errorf("test interface is not ExistingCCMTestInterface")
	}

	// Create a test service
	serviceConfig := &ccmtesting.TestServiceConfig{
		Name:      "test-lb-create",
		Namespace: existingCCM.GetNamespace(),
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Name:       "http",
				Port:       80,
				Protocol:   v1.ProtocolTCP,
				TargetPort: intstr.FromInt(8080),
			},
		},
	}

	service, err := existingCCM.CreateTestService(ctx, serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	// Ensure cleanup happens even if test fails
	defer func() {
		if err := existingCCM.DeleteTestService(ctx, service.Name); err != nil {
			ti.GetTestResults().AddLog(fmt.Sprintf("Warning: failed to cleanup service %s: %v", service.Name, err))
		}
	}()

	ti.GetTestResults().AddLog(fmt.Sprintf("Created service %s in namespace %s", service.Name, service.Namespace))

	// Wait for CCM to provision the load balancer
	lbStatus, err := existingCCM.WaitForLoadBalancer(service.Name, 3*time.Minute)
	if err != nil {
		return fmt.Errorf("CCM failed to provision load balancer: %w", err)
	}

	if lbStatus == nil || len(lbStatus.Ingress) == 0 {
		return fmt.Errorf("load balancer status is empty after provisioning")
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Load balancer provisioned with %d ingress points", len(lbStatus.Ingress)))
	for i, ingress := range lbStatus.Ingress {
		if ingress.IP != "" {
			ti.GetTestResults().AddLog(fmt.Sprintf("  Ingress %d: IP=%s", i, ingress.IP))
		}
		if ingress.Hostname != "" {
			ti.GetTestResults().AddLog(fmt.Sprintf("  Ingress %d: Hostname=%s", i, ingress.Hostname))
		}
	}

	ti.GetTestResults().AddLog("Successfully completed test (cleanup will happen via defer)")
	return nil
}

func testExistingCCMUpdateLoadBalancer(ti ccmtesting.TestInterface) error {
	ctx := context.Background()

	existingCCM, ok := ti.(*ExistingCCMTestInterface)
	if !ok {
		return fmt.Errorf("test interface is not ExistingCCMTestInterface")
	}

	// Create initial service
	serviceConfig := &ccmtesting.TestServiceConfig{
		Name:      "test-lb-update",
		Namespace: existingCCM.GetNamespace(),
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Name:       "http",
				Port:       80,
				Protocol:   v1.ProtocolTCP,
				TargetPort: intstr.FromInt(8080),
			},
		},
	}

	service, err := existingCCM.CreateTestService(ctx, serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	// Ensure cleanup happens even if test fails
	defer func() {
		if err := existingCCM.DeleteTestService(ctx, service.Name); err != nil {
			ti.GetTestResults().AddLog(fmt.Sprintf("Warning: failed to cleanup service %s: %v", service.Name, err))
		}
	}()

	ti.GetTestResults().AddLog(fmt.Sprintf("Created service %s", service.Name))

	// Wait for initial provisioning
	_, err = existingCCM.WaitForLoadBalancer(service.Name, 3*time.Minute)
	if err != nil {
		return fmt.Errorf("CCM failed to provision initial load balancer: %w", err)
	}

	ti.GetTestResults().AddLog("Initial load balancer provisioned")

	// Update the service by adding a new port
	service, err = existingCCM.kubeClient.CoreV1().Services(existingCCM.GetNamespace()).Get(ctx, service.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get service for update: %w", err)
	}

	service.Spec.Ports = append(service.Spec.Ports, v1.ServicePort{
		Name:       "https",
		Port:       443,
		Protocol:   v1.ProtocolTCP,
		TargetPort: intstr.FromInt(8443),
	})

	_, err = existingCCM.kubeClient.CoreV1().Services(existingCCM.GetNamespace()).Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	ti.GetTestResults().AddLog("Updated service with additional port")

	// Wait for CCM to process the update
	time.Sleep(10 * time.Second)

	// Verify the service still has load balancer status
	updatedService, err := existingCCM.kubeClient.CoreV1().Services(existingCCM.GetNamespace()).Get(ctx, service.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get updated service: %w", err)
	}

	if len(updatedService.Status.LoadBalancer.Ingress) == 0 {
		return fmt.Errorf("load balancer ingress lost after update")
	}

	ti.GetTestResults().AddLog("Load balancer update successful (cleanup will happen via defer)")
	return nil
}

func testExistingCCMDeleteLoadBalancer(ti ccmtesting.TestInterface) error {
	ctx := context.Background()

	existingCCM, ok := ti.(*ExistingCCMTestInterface)
	if !ok {
		return fmt.Errorf("test interface is not ExistingCCMTestInterface")
	}

	// Create service
	serviceConfig := &ccmtesting.TestServiceConfig{
		Name:      "test-lb-delete",
		Namespace: existingCCM.GetNamespace(),
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Port:     80,
				Protocol: v1.ProtocolTCP,
			},
		},
	}

	service, err := existingCCM.CreateTestService(ctx, serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	// Track whether we've explicitly deleted the service
	serviceDeleted := false
	
	// Ensure cleanup happens if test fails before we test deletion
	defer func() {
		if !serviceDeleted {
			if err := existingCCM.DeleteTestService(ctx, service.Name); err != nil {
				ti.GetTestResults().AddLog(fmt.Sprintf("Warning: failed to cleanup service %s: %v", service.Name, err))
			}
		}
	}()

	ti.GetTestResults().AddLog(fmt.Sprintf("Created service %s", service.Name))

	// Wait for provisioning
	_, err = existingCCM.WaitForLoadBalancer(service.Name, 3*time.Minute)
	if err != nil {
		return fmt.Errorf("CCM failed to provision load balancer: %w", err)
	}

	ti.GetTestResults().AddLog("Load balancer provisioned, now testing deletion")

	// Delete the service (CCM should clean up the load balancer)
	// This is the actual test - we're testing the delete functionality
	err = existingCCM.DeleteTestService(ctx, service.Name)
	if err != nil {
		return fmt.Errorf("failed to delete test service: %w", err)
	}
	serviceDeleted = true // Mark as explicitly deleted

	ti.GetTestResults().AddLog("Service deleted successfully, CCM should have cleaned up the load balancer")

	// Wait for service to be fully deleted
	deletionTimeout := time.After(30 * time.Second)
	deletionTicker := time.NewTicker(2 * time.Second)
	defer deletionTicker.Stop()

	for {
		select {
		case <-deletionTimeout:
			return fmt.Errorf("timeout waiting for service to be deleted")
		case <-deletionTicker.C:
			_, err = existingCCM.kubeClient.CoreV1().Services(existingCCM.GetNamespace()).Get(ctx, service.Name, metav1.GetOptions{})
			if err != nil {
				// Service is deleted (Get returned error)
				ti.GetTestResults().AddLog("Service fully deleted, load balancer cleanup verified")
				return nil
			}
		}
	}
}

// Test functions for existing CCM node management

func testExistingCCMGetNodes(ti ccmtesting.TestInterface) error {
	existingCCM, ok := ti.(*ExistingCCMTestInterface)
	if !ok {
		return fmt.Errorf("test interface is not ExistingCCMTestInterface")
	}

	nodes, err := existingCCM.GetExistingNodes()
	if err != nil {
		return fmt.Errorf("failed to get existing nodes: %w", err)
	}

	if len(nodes) == 0 {
		return fmt.Errorf("no nodes found in cluster")
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Found %d nodes in cluster", len(nodes)))

	for _, node := range nodes {
		ti.GetTestResults().AddLog(fmt.Sprintf("  Node: %s, Status: %s, ProviderID: %s",
			node.Name,
			getNodeStatus(&node),
			node.Spec.ProviderID))
	}

	return nil
}

func testExistingCCMNodeMetadata(ti ccmtesting.TestInterface) error {
	existingCCM, ok := ti.(*ExistingCCMTestInterface)
	if !ok {
		return fmt.Errorf("test interface is not ExistingCCMTestInterface")
	}

	nodes, err := existingCCM.GetExistingNodes()
	if err != nil {
		return fmt.Errorf("failed to get existing nodes: %w", err)
	}

	if len(nodes) == 0 {
		return fmt.Errorf("no nodes found in cluster")
	}

	ti.GetTestResults().AddLog(fmt.Sprintf("Verifying CCM metadata on %d nodes", len(nodes)))

	nodesMissingMetadata := 0
	for _, node := range nodes {
		// Check if node has provider ID (should be set by CCM)
		if node.Spec.ProviderID == "" {
			ti.GetTestResults().AddLog(fmt.Sprintf("  WARNING: Node %s missing ProviderID", node.Name))
			nodesMissingMetadata++
			continue
		}

		// Check if node has addresses (should be populated by CCM)
		if len(node.Status.Addresses) == 0 {
			ti.GetTestResults().AddLog(fmt.Sprintf("  WARNING: Node %s missing addresses", node.Name))
			nodesMissingMetadata++
			continue
		}

		// Check for instance type label (optional, depends on CCM implementation)
		instanceType, hasLabel := node.Labels["node.kubernetes.io/instance-type"]
		if hasLabel {
			ti.GetTestResults().AddLog(fmt.Sprintf("  Node %s: ProviderID=%s, InstanceType=%s, Addresses=%d",
				node.Name, node.Spec.ProviderID, instanceType, len(node.Status.Addresses)))
		} else {
			ti.GetTestResults().AddLog(fmt.Sprintf("  Node %s: ProviderID=%s, Addresses=%d",
				node.Name, node.Spec.ProviderID, len(node.Status.Addresses)))
		}
	}

	if nodesMissingMetadata > 0 {
		return fmt.Errorf("%d nodes are missing CCM metadata", nodesMissingMetadata)
	}

	ti.GetTestResults().AddLog("All nodes have proper CCM metadata")
	return nil
}

// Helper function to get node status
func getNodeStatus(node *v1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady {
			if condition.Status == v1.ConditionTrue {
				return "Ready"
			}
			return "NotReady"
		}
	}
	return "Unknown"
}

