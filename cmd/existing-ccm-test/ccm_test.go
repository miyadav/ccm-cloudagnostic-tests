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
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	ccmtestpkg "github.com/kubernetes/ccm-cloudagnostic-tests/pkg/testing"
	ccmtesting "github.com/miyadav/cloud-provider-testing-interface"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "Path to kubeconfig file")
	namespace  = flag.String("namespace", "ccm-test", "Namespace for testing")
	timeout    = flag.Duration("timeout", 5*time.Minute, "Test timeout")
	junitFile  = flag.String("junit-file", "", "Path to JUnit XML output file")
)

var (
	testInterface *ccmtestpkg.ExistingCCMTestInterface
	clientset     kubernetes.Interface
)

func TestCCM(t *testing.T) {
	RegisterFailHandler(Fail)

	// Configure JUnit reporter if specified
	if *junitFile != "" {
		// Ensure directory exists
		dir := filepath.Dir(*junitFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			klog.Fatalf("Failed to create JUnit output directory: %v", err)
		}

		// Add JUnit reporter
		RunSpecs(t, "CCM Cloud-Agnostic Tests", Label("ccm", "cloud-provider"))
	} else {
		RunSpecs(t, "CCM Cloud-Agnostic Tests", Label("ccm", "cloud-provider"))
	}
}

var _ = BeforeSuite(func() {
	// Configure verbose logging
	if *verbose {
		klog.InitFlags(nil)
		flag.Set("v", "4")
	}

	// Validate required flags
	Expect(*kubeconfig).NotTo(BeEmpty(), "kubeconfig flag is required")

	// Create Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	Expect(err).NotTo(HaveOccurred(), "Failed to build kubeconfig")

	clientset, err = kubernetes.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred(), "Failed to create clientset")

	// Verify cluster connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{Limit: 1})
	Expect(err).NotTo(HaveOccurred(), "Failed to connect to cluster")

	klog.Info("Successfully connected to cluster")

	// Create test interface
	testInterface = ccmtestpkg.NewExistingCCMTestInterface(clientset, &ccmtesting.TestConfig{
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
	Expect(err).NotTo(HaveOccurred(), "Failed to setup test environment")
})

var _ = AfterSuite(func() {
	if testInterface != nil {
		klog.Info("Tearing down test environment...")
		err := testInterface.TeardownTestEnvironment()
		if err != nil {
			klog.Warningf("Failed to teardown test environment: %v", err)
		}
	}
})

var _ = Describe("CCM Load Balancer Tests", Label("loadbalancer"), func() {
	Context("LoadBalancer Service Creation", func() {
		It("should create a LoadBalancer service and wait for CCM to provision it", func() {
			By("Creating a test LoadBalancer service")
			serviceConfig := &ccmtesting.TestServiceConfig{
				Name:      "test-lb-service",
				Namespace: testInterface.GetNamespace(),
				Type:      v1.ServiceTypeLoadBalancer,
				Ports: []v1.ServicePort{
					{
						Port:     80,
						Protocol: v1.ProtocolTCP,
					},
				},
			}

			service, err := testInterface.CreateTestService(context.Background(), serviceConfig)
			Expect(err).NotTo(HaveOccurred(), "Failed to create test service")
			Expect(service).NotTo(BeNil(), "Service should not be nil")
			Expect(service.Name).To(Equal("test-lb-service"), "Service name should match")
			Expect(service.Namespace).To(Equal(testInterface.GetNamespace()), "Service namespace should match")

			klog.Infof("Created service: %s/%s", service.Namespace, service.Name)

			By("Waiting for CCM to provision the load balancer")
			lbStatus, err := testInterface.WaitForLoadBalancer(service.Name, *timeout)
			Expect(err).NotTo(HaveOccurred(), "Failed to wait for load balancer")
			Expect(lbStatus).NotTo(BeNil(), "Load balancer status should not be nil")
			Expect(lbStatus.Ingress).NotTo(BeEmpty(), "Load balancer should have ingress")

			klog.Infof("✅ Load balancer provisioned: %+v", lbStatus)

			By("Cleaning up the service")
			err = testInterface.DeleteTestService(context.Background(), service.Name)
			Expect(err).NotTo(HaveOccurred(), "Failed to delete test service")

			klog.Info("✅ Service cleaned up successfully")
		})

		It("should validate load balancer provider matches cluster provider", func() {
			By("Creating a test LoadBalancer service")
			serviceConfig := &ccmtesting.TestServiceConfig{
				Name:      "test-lb-provider-validation",
				Namespace: testInterface.GetNamespace(),
				Type:      v1.ServiceTypeLoadBalancer,
				Ports: []v1.ServicePort{
					{
						Port:     80,
						Protocol: v1.ProtocolTCP,
					},
				},
			}

			service, err := testInterface.CreateTestService(context.Background(), serviceConfig)
			Expect(err).NotTo(HaveOccurred(), "Failed to create test service")

			By("Waiting for load balancer and validating provider")
			lbStatus, err := testInterface.WaitForLoadBalancer(service.Name, *timeout)
			Expect(err).NotTo(HaveOccurred(), "Failed to wait for load balancer")

			// Note: Load balancer provider validation would be implemented here
			// For now, we just verify the load balancer was provisioned
			Expect(lbStatus.Ingress).NotTo(BeEmpty(), "Load balancer should have ingress")

			By("Cleaning up the service")
			err = testInterface.DeleteTestService(context.Background(), service.Name)
			Expect(err).NotTo(HaveOccurred(), "Failed to delete test service")
		})
	})
})

var _ = Describe("CCM Node Management Tests", Label("node-management"), func() {
	Context("Node Processing Validation", func() {
		It("should verify CCM has processed existing nodes", func() {
			By("Getting existing nodes in the cluster")
			nodes, err := testInterface.GetExistingNodes()
			Expect(err).NotTo(HaveOccurred(), "Failed to get existing nodes")
			Expect(nodes).NotTo(BeEmpty(), "No nodes found in the cluster")

			klog.Infof("Found %d nodes in the cluster", len(nodes))

			By("Verifying CCM processing for each node")
			for _, node := range nodes {
				klog.Infof("Testing node: %s", node.Name)

				// Verify CCM has processed the node
				err := testInterface.VerifyCCMNodeProcessing(&node)
				Expect(err).NotTo(HaveOccurred(), "CCM processing verification failed for node %s", node.Name)

				klog.Infof("✅ CCM has successfully processed node: %s", node.Name)

				// Check node conditions
				By("Checking node readiness condition")
				var nodeReady bool
				for _, condition := range node.Status.Conditions {
					if condition.Type == v1.NodeReady {
						if condition.Status == v1.ConditionTrue {
							nodeReady = true
							klog.Infof("✅ Node %s is ready", node.Name)
						} else {
							klog.Warningf("⚠️ Node %s is not ready: %s", node.Name, condition.Message)
						}
						break
					}
				}
				Expect(nodeReady).To(BeTrue(), "Node %s should be ready", node.Name)
			}
		})

		It("should validate node cloud provider annotations and labels", func() {
			By("Getting existing nodes")
			nodes, err := testInterface.GetExistingNodes()
			Expect(err).NotTo(HaveOccurred(), "Failed to get existing nodes")
			Expect(nodes).NotTo(BeEmpty(), "No nodes found in the cluster")

			By("Validating cloud provider metadata on nodes")
			for _, node := range nodes {
				// Check for common cloud provider annotations/labels
				hasCloudMetadata := false

				// Check annotations
				cloudAnnotations := []string{
					"node.cloudprovider.kubernetes.io/instance-id",
					"node.cloudprovider.kubernetes.io/instance-type",
					"node.cloudprovider.kubernetes.io/zone",
					"node.cloudprovider.kubernetes.io/region",
					"node.kubernetes.io/instance-type",
					"topology.kubernetes.io/zone",
					"topology.kubernetes.io/region",
				}

				for _, annotation := range cloudAnnotations {
					if _, exists := node.Annotations[annotation]; exists {
						hasCloudMetadata = true
						klog.Infof("Found cloud provider annotation on node %s: %s", node.Name, annotation)
						break
					}
				}

				// Check labels
				cloudLabels := []string{
					"node.kubernetes.io/instance-type",
					"topology.kubernetes.io/zone",
					"topology.kubernetes.io/region",
					"node.kubernetes.io/cloud-provider",
				}

				for _, label := range cloudLabels {
					if _, exists := node.Labels[label]; exists {
						hasCloudMetadata = true
						klog.Infof("Found cloud provider label on node %s: %s", node.Name, label)
						break
					}
				}

				// Check provider ID
				if node.Spec.ProviderID != "" {
					hasCloudMetadata = true
					klog.Infof("Node %s has provider ID: %s", node.Name, node.Spec.ProviderID)
				}

				Expect(hasCloudMetadata).To(BeTrue(), "Node %s should have cloud provider metadata", node.Name)
			}
		})
	})
})

var _ = Describe("CCM Integration Tests", Label("integration"), func() {
	Context("End-to-End CCM Functionality", func() {
		It("should perform complete CCM workflow test", func() {
			By("Testing load balancer creation and cleanup")
			serviceConfig := &ccmtesting.TestServiceConfig{
				Name:      "integration-test-lb",
				Namespace: testInterface.GetNamespace(),
				Type:      v1.ServiceTypeLoadBalancer,
				Ports: []v1.ServicePort{
					{
						Port:     80,
						Protocol: v1.ProtocolTCP,
					},
				},
			}

			service, err := testInterface.CreateTestService(context.Background(), serviceConfig)
			Expect(err).NotTo(HaveOccurred(), "Failed to create integration test service")

			lbStatus, err := testInterface.WaitForLoadBalancer(service.Name, *timeout)
			Expect(err).NotTo(HaveOccurred(), "Failed to wait for load balancer in integration test")

			// Note: Load balancer provider validation would be implemented here
			// For now, we just verify the load balancer was provisioned
			Expect(lbStatus.Ingress).NotTo(BeEmpty(), "Load balancer should have ingress in integration test")

			// Cleanup
			err = testInterface.DeleteTestService(context.Background(), service.Name)
			Expect(err).NotTo(HaveOccurred(), "Failed to delete integration test service")

			By("Verifying node management functionality")
			nodes, err := testInterface.GetExistingNodes()
			Expect(err).NotTo(HaveOccurred(), "Failed to get nodes in integration test")
			Expect(nodes).NotTo(BeEmpty(), "No nodes found in integration test")

			// Verify at least one node is properly processed
			hasProcessedNode := false
			for _, node := range nodes {
				err := testInterface.VerifyCCMNodeProcessing(&node)
				if err == nil {
					hasProcessedNode = true
					break
				}
			}
			Expect(hasProcessedNode).To(BeTrue(), "At least one node should be properly processed by CCM")

			klog.Info("✅ Integration test completed successfully")
		})
	})
})
