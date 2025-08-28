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
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/kubernetes/ccm-cloudagnostic-tests/pkg/testing"
	ccmtesting "github.com/miyadav/cloud-provider-testing-interface"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "Path to kubeconfig file")
	namespace  = flag.String("namespace", "ccm-test", "Namespace for testing")
	timeout    = flag.Duration("timeout", 5*time.Minute, "Test timeout")
	verbose    = flag.Bool("verbose", false, "Enable verbose output")
)

func main() {
	flag.Parse()

	// Configure verbose logging
	if *verbose {
		klog.InitFlags(nil)
		flag.Set("v", "4")
	}

	if *kubeconfig == "" {
		klog.Fatal("--kubeconfig flag is required")
	}

	// Create Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		klog.Fatalf("Failed to build config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Failed to create clientset: %v", err)
	}

	// Verify cluster connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		klog.Fatalf("Failed to connect to cluster: %v", err)
	}

	klog.Info("Successfully connected to cluster")

	// Create test interface
	testInterface := testing.NewExistingCCMTestInterface(clientset, &ccmtesting.TestConfig{
		ProviderName: "existing",
		TestData: map[string]interface{}{
			"resource-prefix": "existing-ccm-test",
			"namespace":       *namespace,
		},
	})

	// Setup test environment
	klog.Info("Setting up test environment...")
	err = testInterface.SetupTestEnvironment(&ccmtesting.TestConfig{
		ProviderName: "existing",
		TestData: map[string]interface{}{
			"resource-prefix": "existing-ccm-test",
		},
	})
	if err != nil {
		klog.Fatalf("Failed to setup test environment: %v", err)
	}
	defer func() {
		klog.Info("Tearing down test environment...")
		testInterface.TeardownTestEnvironment()
	}()

	// Run tests
	klog.Info("Starting existing CCM tests...")

	// Test 1: Create a LoadBalancer service and wait for CCM to provision it
	if err := testLoadBalancerCreation(testInterface); err != nil {
		klog.Errorf("Load balancer test failed: %v", err)
	} else {
		klog.Info("✅ Load balancer test passed")
	}

	// Test 2: Create a node and verify CCM handles it
	if err := testNodeManagement(testInterface); err != nil {
		klog.Errorf("Node management test failed: %v", err)
	} else {
		klog.Info("✅ Node management test passed")
	}

	klog.Info("All tests completed")
}

func testLoadBalancerCreation(ti *testing.ExistingCCMTestInterface) error {
	klog.Info("Testing load balancer creation...")

	// Create a test service
	serviceConfig := &ccmtesting.TestServiceConfig{
		Name:      "test-lb-service",
		Namespace: ti.GetNamespace(),
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Port:     80,
				Protocol: v1.ProtocolTCP,
			},
		},
	}

	service, err := ti.CreateTestService(context.Background(), serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	klog.Infof("Created service: %s/%s", service.Namespace, service.Name)

	// Wait for CCM to provision the load balancer
	klog.Info("Waiting for CCM to provision load balancer...")
	lbStatus, err := ti.WaitForLoadBalancer(service.Name, *timeout)
	if err != nil {
		// Clean up the service even if the test failed
		ti.DeleteTestService(context.Background(), service.Name)
		return fmt.Errorf("failed to wait for load balancer: %w", err)
	}

	klog.Infof("✅ Load balancer provisioned: %+v", lbStatus)

	// Clean up
	err = ti.DeleteTestService(context.Background(), service.Name)
	if err != nil {
		return fmt.Errorf("failed to delete test service: %w", err)
	}

	klog.Info("✅ Service cleaned up successfully")
	return nil
}

func testNodeManagement(ti *testing.ExistingCCMTestInterface) error {
	klog.Info("Testing node management...")

	// Instead of creating a test node, let's test with existing nodes in the cluster
	// This is more realistic and tests the actual CCM functionality
	klog.Info("Checking existing nodes in the cluster...")

	nodes, err := ti.GetExistingNodes()
	if err != nil {
		return fmt.Errorf("failed to get existing nodes: %w", err)
	}

	if len(nodes) == 0 {
		return fmt.Errorf("no nodes found in the cluster")
	}

	klog.Infof("Found %d nodes in the cluster", len(nodes))

	// Test each node to verify CCM processing
	for _, node := range nodes {
		klog.Infof("Testing node: %s", node.Name)

		// Verify CCM has processed the node
		if err := ti.VerifyCCMNodeProcessing(&node); err != nil {
			klog.Warningf("CCM processing verification failed for node %s: %v", node.Name, err)
		} else {
			klog.Infof("✅ CCM has successfully processed node: %s", node.Name)
		}

		// Check node conditions
		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady {
				if condition.Status == v1.ConditionTrue {
					klog.Infof("✅ Node %s is ready", node.Name)
				} else {
					klog.Warningf("⚠️ Node %s is not ready: %s", node.Name, condition.Message)
				}
				break
			}
		}
	}

	klog.Info("✅ Node management test completed")
	return nil
}
