# CCM Cloud-Agnostic E2E Testing Framework

This framework provides cloud-agnostic end-to-end (e2e) testing for Kubernetes Cloud Controller Manager (CCM) functionality across different cloud providers.

## 🚀 Quick Start

### 1. Build the E2E Test Runner
```bash
# Build both test runners
make build build-e2e

# Or build individually
go build -o bin/test-runner cmd/test-runner/main.go
go build -o bin/e2e-test-runner cmd/e2e-test-runner/main.go
```

### 2. Test with Mock Provider (Local Testing)
```bash
./bin/e2e-test-runner --provider mock --suite all --verbose
```

### 3. Test with Real Cloud Provider
```bash
# AWS EKS Example
./bin/e2e-test-runner \
  --provider aws \
  --kubeconfig ~/.kube/config \
  --region us-west-2 \
  --zone us-west-2a \
  --cluster my-eks-cluster \
  --suite loadbalancer \
  --verbose
```

## 📋 Prerequisites

### For Real Cluster Testing
- **Kubernetes Cluster**: Running cluster with CCM enabled
- **Kubeconfig**: Valid kubeconfig file for cluster access
- **Cloud Provider Credentials**: Valid credentials with appropriate permissions
- **Network Access**: Access to cloud provider APIs

### For Local Testing
- **Go 1.24+**: For building and running tests
- **No External Dependencies**: Mock provider works offline

## 🏗️ Architecture

### Components
1. **Test Interface**: Cloud-agnostic interface for testing CCM functionality
2. **Cloud Provider Adapters**: Wrappers for real cloud providers (AWS, GCP, Azure)
3. **Mock Provider**: Simulated cloud provider for local testing
4. **Test Suites**: Organized test collections for different CCM features
5. **E2E Test Runner**: Command-line tool for running tests on real clusters

### Test Suites
- **LoadBalancer**: Load balancer creation, updates, deletion, status
- **Node Management**: Node initialization, addresses, provider IDs
- **Route Management**: Route creation, deletion, listing
- **Instances**: Instance existence, shutdown detection, metadata
- **Zones**: Zone information retrieval
- **Clusters**: Cluster listing, master node detection

## 🔧 Usage Examples

### Mock Provider Testing
```bash
# Test all suites
./bin/e2e-test-runner --provider mock --suite all --verbose

# Test specific suite
./bin/e2e-test-runner --provider mock --suite loadbalancer --verbose

# Test with custom timeout
./bin/e2e-test-runner --provider mock --suite all --timeout 10m --verbose
```

### AWS EKS Testing
```bash
# Setup AWS credentials
aws configure
aws eks update-kubeconfig --region us-west-2 --name my-eks-cluster

# Run e2e tests
./bin/e2e-test-runner \
  --provider aws \
  --kubeconfig ~/.kube/config \
  --region us-west-2 \
  --zone us-west-2a \
  --cluster my-eks-cluster \
  --prefix aws-e2e-test \
  --timeout 45m \
  --verbose
```

### GCP GKE Testing
```bash
# Setup GCP credentials
gcloud auth login
gcloud config set project my-project-id
gcloud container clusters get-credentials my-gke-cluster --zone us-central1-a

# Run e2e tests
./bin/e2e-test-runner \
  --provider gcp \
  --kubeconfig ~/.kube/config \
  --region us-central1 \
  --zone us-central1-a \
  --cluster my-gke-cluster \
  --prefix gcp-e2e-test \
  --timeout 45m \
  --verbose
```

### Azure AKS Testing
```bash
# Setup Azure credentials
az login
az account set --subscription my-subscription-id
az aks get-credentials --resource-group my-resource-group --name my-aks-cluster

# Run e2e tests
./bin/e2e-test-runner \
  --provider azure \
  --kubeconfig ~/.kube/config \
  --region eastus \
  --zone eastus-1 \
  --cluster my-aks-cluster \
  --prefix azure-e2e-test \
  --timeout 45m \
  --verbose
```

## ⚙️ Configuration Options

### Required Flags
- `--provider`: Cloud provider (`aws`, `gcp`, `azure`, `mock`)
- `--kubeconfig`: Path to kubeconfig (not required for mock)

### Cloud Provider Configuration
- `--region`: Cloud provider region
- `--zone`: Cloud provider zone/availability zone
- `--cluster`: Cluster name
- `--prefix`: Resource prefix for test resources (default: `e2e-test`)

### Test Execution
- `--suite`: Test suite to run (`all`, `loadbalancer`, `nodes`, `routes`, `instances`, `zones`, `clusters`)
- `--timeout`: Test timeout (default: 30m)
- `--verbose`: Enable verbose output
- `--cleanup`: Clean up resources after tests (default: true)

### Output and Logging
- `--output`: Output format (`text`, `json`)
- `--log-level`: Log level (`debug`, `info`, `warn`, `error`)

## 🔄 CI/CD Integration

### GitHub Actions
```yaml
name: CCM E2E Tests
on: [push, pull_request]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.24'
    
    - name: Build e2e test runner
      run: go build -o bin/e2e-test-runner cmd/e2e-test-runner/main.go
    
    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v1
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: us-west-2
    
    - name: Update kubeconfig
      run: aws eks update-kubeconfig --region us-west-2 --name my-eks-cluster
    
    - name: Run e2e tests
      run: |
        ./bin/e2e-test-runner \
          --provider aws \
          --kubeconfig ~/.kube/config \
          --region us-west-2 \
          --zone us-west-2a \
          --cluster my-eks-cluster \
          --suite all \
          --timeout 45m
```

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    
    stages {
        stage('Build') {
            steps {
                sh 'go build -o bin/e2e-test-runner cmd/e2e-test-runner/main.go'
            }
        }
        
        stage('E2E Tests') {
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

## 🛠️ Development

### Adding New Cloud Providers
1. Implement the `cloudprovider.Interface` for your provider
2. Create a provider adapter in `pkg/testing/real_cloud_provider.go`
3. Add provider initialization in `cmd/e2e-test-runner/main.go`
4. Update the `createCloudProvider` function

### Adding New Test Suites
1. Create test functions in `pkg/testing/test_suites.go`
2. Add test suite creation function
3. Update the `addTestSuites` function in the test runner
4. Add command line flag support if needed

### Local Development
```bash
# Setup development environment
make dev-setup

# Run tests
make test

# Run e2e tests with mock
make test-e2e

# Build both runners
make build build-e2e

# Run linting
make lint

# Run all checks
make check
```

## 🔍 Troubleshooting

### Common Issues

1. **Cluster Connection Failed**
   - Verify kubeconfig path is correct
   - Check cluster is running and accessible
   - Ensure proper RBAC permissions

2. **Cloud Provider Authentication Failed**
   - Verify cloud provider credentials
   - Check IAM/service account permissions
   - Ensure network access to cloud APIs

3. **Resource Creation Failed**
   - Check cloud provider quotas
   - Verify network configuration
   - Ensure proper resource permissions

4. **Test Timeout**
   - Increase timeout with `--timeout` flag
   - Check cluster performance
   - Verify cloud provider API responsiveness

### Debug Mode
```bash
./bin/e2e-test-runner \
  --provider aws \
  --kubeconfig ~/.kube/config \
  --log-level debug \
  --verbose
```

### Manual Cleanup
If tests fail and leave resources behind:
```bash
# List resources with test prefix
kubectl get services -l test-prefix=e2e-test
kubectl get routes -l test-prefix=e2e-test

# Delete manually if needed
kubectl delete service -l test-prefix=e2e-test
```

## 📊 Test Results

### Sample Output
```
=== CCM E2E Test Results ===
Total Duration: 2m15s
Test Summary: 15 total, 14 passed, 1 failed, 0 skipped

Detailed Results:
  PASSED: CreateLoadBalancer (45.2s)
  PASSED: UpdateLoadBalancer (12.1s)
  PASSED: DeleteLoadBalancer (8.3s)
  PASSED: LoadBalancerStatus (30.5s)
  PASSED: LoadBalancerHealthCheck (5.2s)
  PASSED: NodeInitialization (15.8s)
  PASSED: NodeAddresses (3.2s)
  PASSED: NodeProviderID (2.1s)
  PASSED: NodeInstanceType (1.9s)
  PASSED: NodeZones (2.3s)
  PASSED: CreateRoute (18.7s)
  PASSED: DeleteRoute (12.4s)
  PASSED: ListRoutes (1.2s)
  PASSED: InstanceExists (2.8s)
  FAILED: InstanceShutdown (timeout)

❌ Some tests failed: 1 failed out of 15 total
```

## 🎯 Best Practices

1. **Use Resource Prefixes**: Always use unique prefixes to avoid conflicts
2. **Set Appropriate Timeouts**: Cloud operations can take time, set realistic timeouts
3. **Enable Cleanup**: Always enable cleanup unless debugging
4. **Monitor Resources**: Keep an eye on cloud provider quotas and limits
5. **Use Verbose Logging**: Enable verbose output for better debugging
6. **Test in Isolation**: Run tests in dedicated test clusters when possible
7. **Monitor Costs**: Be aware of cloud provider costs for resource creation
8. **Version Control**: Use specific versions of cloud provider SDKs
9. **Error Handling**: Implement proper error handling and retry logic
10. **Documentation**: Document cloud provider specific configurations

## 📚 Additional Resources

- [Cloud Controller Manager Documentation](https://kubernetes.io/docs/concepts/architecture/cloud-controller/)
- [AWS Cloud Provider](https://github.com/kubernetes/cloud-provider-aws)
- [GCP Cloud Provider](https://github.com/kubernetes/cloud-provider-gcp)
- [Azure Cloud Provider](https://github.com/kubernetes/cloud-provider-azure)
- [Kubernetes Testing Framework](https://github.com/kubernetes/kubernetes/tree/master/test/e2e)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
