# CCM E2E Testing Guide - Legacy Framework

This guide provides detailed technical information for the **legacy e2e test runner** in the CCM Cloud-Agnostic Testing Framework.

> **⚠️ Note**: This guide covers the legacy e2e test runner. For the modern Ginkgo-based tests, see the [Ginkgo Refactoring Guide](ginkgo-refactoring-guide.md) and [main README](../README.md).

## When to Use This Guide

- **Legacy Support**: Maintaining existing CI/CD pipelines that use the e2e test runner
- **Provider Development**: Adding new cloud providers to the legacy framework
- **Advanced Configuration**: Complex test configurations not available in Ginkgo tests
- **Migration Reference**: Understanding the legacy framework before migrating to Ginkgo

## Quick Comparison

| Feature | Legacy E2E Runner | Modern Ginkgo Tests |
|---------|------------------|-------------------|
| **Test Organization** | Custom framework | Standard BDD structure |
| **Assertions** | Manual error handling | Rich Gomega matchers |
| **CI/CD Integration** | Basic | JUnit reports, Prow support |
| **Test Discovery** | Manual | Automatic with labels |
| **Maintenance** | Custom code | Industry standard |
| **Recommendation** | Legacy support only | ✅ **Use this** |

## 🏗️ Architecture Deep Dive

### **Test Interface Design**

The framework uses a cloud-agnostic test interface that abstracts away provider-specific details:

```go
type TestInterface interface {
    SetupTestEnvironment(config *TestConfig) error
    TeardownTestEnvironment() error
    GetCloudProvider() cloudprovider.Interface
    CreateTestNode(ctx context.Context, nodeConfig *TestNodeConfig) (*v1.Node, error)
    DeleteTestNode(ctx context.Context, nodeName string) error
    CreateTestService(ctx context.Context, serviceConfig *TestServiceConfig) (*v1.Service, error)
    DeleteTestService(ctx context.Context, serviceName string) error
    CreateTestRoute(ctx context.Context, routeConfig *TestRouteConfig) (*cloudprovider.Route, error)
    DeleteTestRoute(ctx context.Context, routeName string) error
    WaitForCondition(ctx context.Context, condition TestCondition) error
    GetTestResults() *TestResults
    ResetTestState() error
}
```

### **Provider Adapter Pattern**

Each cloud provider implements the same interface through adapters:

```go
// Mock Provider (for local testing)
type MockCloudProvider struct {
    // Simulated cloud resources
}

// Real Provider Adapter (for cloud testing)
type RealCloudProviderAdapter struct {
    cloudProvider cloudprovider.Interface
    kubeClient    kubernetes.Interface
    config        *RealCloudProviderConfig
}

// Existing CCM Adapter (for testing running CCM)
type ExistingCCMTestInterface struct {
    kubeClient kubernetes.Interface
    config     *TestConfig
    namespace  string
}
```

## 🔧 Advanced Configuration

### **Test Configuration Options**

```go
type TestConfig struct {
    ProviderName         string                 // Cloud provider name
    ClusterName          string                 // Cluster name
    Region               string                 // Cloud region
    Zone                 string                 // Cloud zone
    TestTimeout          time.Duration          // Test timeout
    CleanupResources     bool                   // Auto-cleanup
    MockExternalServices bool                   // Use mocks
    TestData             map[string]interface{} // Custom data
}
```

### **Custom Test Data**

You can pass custom configuration through the `TestData` field:

```bash
# Example: Custom resource prefix and test mode
./bin/e2e-test-runner \
  --provider aws \
  --kubeconfig ~/.kube/config \
  --prefix my-custom-test \
  --test-data '{"custom-setting": "value", "test-mode": "production"}'
```

## 🧪 Advanced Testing Patterns

### **1. Testing Existing CCM (Legacy)**

This approach tests your running CCM using the legacy e2e test runner:

```bash
# Full e2e test runner with existing CCM
./bin/e2e-test-runner --provider existing --kubeconfig ~/.kube/config --suite all --verbose
```

> **Note**: For modern testing, use the Ginkgo-based tests instead of the legacy e2e test runner.

### **2. Testing with Mock Provider**

For local development and CI/CD:

```bash
# Test framework functionality
./bin/e2e-test-runner --provider mock --suite all --verbose

# Test specific functionality
./bin/e2e-test-runner --provider mock --suite loadbalancer --verbose
```

### **3. Testing with Real Cloud Providers**

For comprehensive e2e testing:

```bash
# AWS EKS
./bin/e2e-test-runner \
  --provider aws \
  --kubeconfig ~/.kube/config \
  --region us-west-2 \
  --zone us-west-2a \
  --cluster my-eks-cluster \
  --suite loadbalancer \
  --timeout 45m \
  --verbose

# GCP GKE
./bin/e2e-test-runner \
  --provider gcp \
  --kubeconfig ~/.kube/config \
  --region us-central1 \
  --zone us-central1-a \
  --cluster my-gke-cluster \
  --suite loadbalancer \
  --timeout 45m \
  --verbose

# Azure AKS
./bin/e2e-test-runner \
  --provider azure \
  --kubeconfig ~/.kube/config \
  --region eastus \
  --zone eastus-1 \
  --cluster my-aks-cluster \
  --suite loadbalancer \
  --timeout 45m \
  --verbose
```

## 🔍 Test Suites Deep Dive

### **LoadBalancer Test Suite**

Tests cloud provider load balancer functionality:

```go
func CreateLoadBalancerTestSuite() *TestSuite {
    return &TestSuite{
        Name: "LoadBalancer",
        Tests: []Test{
            {
                Name: "CreateLoadBalancer",
                Run: func(ti TestInterface) error {
                    // Create service with LoadBalancer type
                    // Wait for CCM to provision load balancer
                    // Verify load balancer status
                    return nil
                },
            },
            {
                Name: "UpdateLoadBalancer",
                Run: func(ti TestInterface) error {
                    // Update service configuration
                    // Verify load balancer updates
                    return nil
                },
            },
            {
                Name: "DeleteLoadBalancer",
                Run: func(ti TestInterface) error {
                    // Delete service
                    // Verify load balancer cleanup
                    return nil
                },
            },
        },
    }
}
```

### **Node Management Test Suite**

Tests cloud provider node management:

```go
func CreateNodeTestSuite() *TestSuite {
    return &TestSuite{
        Name: "NodeManagement",
        Tests: []Test{
            {
                Name: "NodeInitialization",
                Run: func(ti TestInterface) error {
                    // Create test node
                    // Verify node initialization
                    return nil
                },
            },
            {
                Name: "NodeAddresses",
                Run: func(ti TestInterface) error {
                    // Verify node addresses
                    return nil
                },
            },
        },
    }
}
```

## 🔄 CI/CD Integration Patterns

### **GitHub Actions - Multi-Provider Testing**

```yaml
name: CCM E2E Tests
on: [push, pull_request]

jobs:
  test-mock:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Build
      run: make build
    - name: Test with mock
      run: ./bin/e2e-test-runner --provider mock --suite all --verbose

  test-existing-ccm:
    runs-on: ubuntu-latest
    if: github.event_name == 'pull_request'
    steps:
    - uses: actions/checkout@v3
    - name: Build
      run: make build
    - name: Test existing CCM
      run: |
        if [ -n "$KUBECONFIG" ]; then
          ./bin/existing-ccm-test --kubeconfig $KUBECONFIG --verbose
        fi

  test-aws:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
    - uses: actions/checkout@v3
    - name: Build
      run: make build
    - name: Configure AWS
      uses: aws-actions/configure-aws-credentials@v1
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: us-west-2
    - name: Test AWS
      run: |
        aws eks update-kubeconfig --region us-west-2 --name my-eks-cluster
        ./bin/e2e-test-runner \
          --provider aws \
          --kubeconfig ~/.kube/config \
          --region us-west-2 \
          --zone us-west-2a \
          --cluster my-eks-cluster \
          --suite loadbalancer \
          --timeout 45m
```

### **Jenkins Pipeline - Multi-Stage Testing**

```groovy
pipeline {
    agent any
    
    stages {
        stage('Build') {
            steps {
                sh 'make build'
            }
        }
        
        stage('Test Mock') {
            steps {
                sh './bin/e2e-test-runner --provider mock --suite all --verbose'
            }
        }
        
        stage('Test Existing CCM') {
            when {
                expression { env.KUBECONFIG != null }
            }
            steps {
                sh './bin/existing-ccm-test --kubeconfig $KUBECONFIG --verbose'
            }
        }
        
        stage('Test AWS') {
            when {
                branch 'main'
            }
            steps {
                withCredentials([
                    string(credentialsId: 'aws-access-key', variable: 'AWS_ACCESS_KEY_ID'),
                    string(credentialsId: 'aws-secret-key', variable: 'AWS_SECRET_ACCESS_KEY')
                ]) {
                    sh '''
                        export AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
                        export AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY
                        aws eks update-kubeconfig --region us-west-2 --name my-eks-cluster
                        ./bin/e2e-test-runner \
                          --provider aws \
                          --kubeconfig ~/.kube/config \
                          --region us-west-2 \
                          --zone us-west-2a \
                          --cluster my-eks-cluster \
                          --suite all \
                          --timeout 45m
                    '''
                }
            }
        }
    }
}
```

## 🛠️ Extending the Framework

### **Adding New Cloud Providers**

1. **Implement Cloud Provider Interface**:
```go
type MyCloudProvider struct {
    // Your cloud provider implementation
}

func (p *MyCloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
    return &MyLoadBalancer{}, true
}

func (p *MyCloudProvider) Instances() (cloudprovider.Instances, bool) {
    return &MyInstances{}, true
}

// ... implement other interfaces
```

2. **Create Provider Adapter**:
```go
type MyCloudProviderAdapter struct {
    *RealCloudProviderAdapter
    // Provider-specific fields
}

func NewMyCloudProviderAdapter(kubeClient kubernetes.Interface, config *RealCloudProviderConfig) (*MyCloudProviderAdapter, error) {
    // Initialize your cloud provider
    myProvider, err := initializeMyCloudProvider(config)
    if err != nil {
        return nil, err
    }
    
    adapter := &MyCloudProviderAdapter{
        RealCloudProviderAdapter: NewRealCloudProviderAdapter(myProvider, kubeClient, config),
        // Set provider-specific fields
    }
    
    return adapter, nil
}
```

3. **Add to Test Runner**:
```go
func createCloudProvider(providerName string, kubeClient kubernetes.Interface) (cloudprovider.Interface, error) {
    switch strings.ToLower(providerName) {
    case "mycloud":
        return createMyCloudProvider(kubeClient)
    // ... other cases
    }
}
```

### **Adding New Test Suites**

1. **Create Test Functions**:
```go
func testMyFeature(ti TestInterface) error {
    // Your test implementation
    return nil
}
```

2. **Create Test Suite**:
```go
func CreateMyFeatureTestSuite() *TestSuite {
    return &TestSuite{
        Name: "MyFeature",
        Tests: []Test{
            {
                Name: "TestMyFeature",
                Run: testMyFeature,
            },
        },
    }
}
```

3. **Add to Test Runner**:
```go
func addTestSuites(runner *TestRunner, suite, provider string) {
    switch strings.ToLower(suite) {
    case "myfeature":
        runner.AddTestSuite(CreateMyFeatureTestSuite())
    // ... other cases
    }
}
```

## 🔍 Troubleshooting

### **Common Issues and Solutions**

#### **1. Test Timeout**
```
Error: Test execution failed: context deadline exceeded
```
**Solution**: Increase timeout with `--timeout` flag
```bash
./bin/e2e-test-runner --provider aws --timeout 60m --verbose
```

#### **2. Resource Cleanup Failed**
```
Error: Failed to delete load balancer
```
**Solution**: Check cloud provider permissions and quotas
```bash
# Disable cleanup for debugging
./bin/e2e-test-runner --provider aws --cleanup=false --verbose
```

#### **3. Cluster Connection Failed**
```
Error: Failed to connect to cluster
```
**Solution**: Verify kubeconfig and cluster status
```bash
# Test cluster connectivity
kubectl get nodes
```

#### **4. Cloud Provider Authentication Failed**
```
Error: Failed to initialize cloud provider
```
**Solution**: Verify credentials and permissions
```bash
# Test cloud provider access
aws sts get-caller-identity  # for AWS
gcloud auth list             # for GCP
az account show              # for Azure
```

### **Debug Mode**

Enable debug logging for detailed troubleshooting:

```bash
./bin/e2e-test-runner \
  --provider aws \
  --kubeconfig ~/.kube/config \
  --log-level debug \
  --verbose
```

### **Manual Cleanup**

If tests fail and leave resources behind:

```bash
# List resources with test prefix
kubectl get services -l test-prefix=e2e-test
kubectl get routes -l test-prefix=e2e-test

# Delete manually if needed
kubectl delete service -l test-prefix=e2e-test
```

## 📊 Performance Testing

### **Benchmark Tests**

Run performance benchmarks:

```bash
# Run benchmark tests
make test-benchmark

# Or directly
go test -bench=. -benchmem ./pkg/testing/...
```

### **Load Testing**

Test CCM performance under load:

```bash
# Create multiple services simultaneously
for i in {1..10}; do
  kubectl create service loadbalancer test-lb-$i --tcp=80:80 &
done
wait

# Monitor CCM performance
kubectl logs -n kube-system deployment/cloud-controller-manager -f
```

## 🔒 Security Testing

### **Security Best Practices**

1. **Use Service Accounts**: Create dedicated service accounts for testing
2. **Limit Permissions**: Use minimal required permissions
3. **Network Policies**: Implement network policies for test namespaces
4. **Secret Management**: Use secure secret management for credentials

### **Security Testing Examples**

```bash
# Test with minimal permissions
kubectl create serviceaccount ccm-test
kubectl create clusterrolebinding ccm-test --clusterrole=view --serviceaccount=default:ccm-test

# Run tests with limited permissions
kubectl run ccm-test --serviceaccount=ccm-test --image=your-test-image
```

## 📚 Additional Resources

- [Cloud Controller Manager Documentation](https://kubernetes.io/docs/concepts/architecture/cloud-controller/)
- [AWS Cloud Provider](https://github.com/kubernetes/cloud-provider-aws)
- [GCP Cloud Provider](https://github.com/kubernetes/cloud-provider-gcp)
- [Azure Cloud Provider](https://github.com/kubernetes/cloud-provider-azure)
- [Kubernetes Testing Framework](https://github.com/kubernetes/kubernetes/tree/master/test/e2e)
