# CCM Cloud-Agnostic Testing Framework

A comprehensive, cloud-agnostic testing framework for Kubernetes Cloud Controller Manager (CCM) functionality across different cloud providers.

## 🚀 Quick Start

### 1. Build the Framework
```bash
# Build the main e2e test runner
make build

# Build all binaries (e2e test runner + existing CCM test)
make build-all
```

### 2. Test with Mock Provider (Local Testing)
```bash
# Test with mock provider (no cluster needed)
./bin/e2e-test-runner --provider mock --suite all --verbose
```

### 3. Test with Real Cluster (No Cloud Credentials Needed!)
```bash
# Test existing CCM in your cluster (recommended)
./bin/existing-ccm-test --kubeconfig ~/.kube/config --verbose

# Or use the full e2e test runner
./bin/e2e-test-runner --provider existing --kubeconfig ~/.kube/config --suite loadbalancer --verbose
```

### 4. Test with Real Cloud Provider (Requires Credentials)
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

## 🎯 Key Features

### ✅ **Cloud-Agnostic Testing**
- Same test interface across all cloud providers
- Consistent test expectations and validation
- Easy to add new cloud providers

### ✅ **Multiple Testing Modes**
- **Mock Provider**: Local testing without cluster (fast, free)
- **Existing CCM**: Test your running CCM (no credentials needed)
- **Real Providers**: Full e2e testing with cloud credentials

### ✅ **Comprehensive Test Suites**
- **LoadBalancer**: Creation, updates, deletion, status
- **Node Management**: Initialization, addresses, provider IDs
- **Route Management**: Creation, deletion, listing
- **Instances**: Existence, shutdown detection, metadata
- **Zones**: Information retrieval
- **Clusters**: Listing and master node detection

### ✅ **Production Ready**
- Resource cleanup and management
- Error handling and reporting
- CI/CD integration examples
- Comprehensive documentation

## 🏗️ Architecture

### **Components**
1. **E2E Test Runner** (`cmd/e2e-test-runner/`): Main testing tool supporting all providers
2. **Existing CCM Test** (`cmd/existing-ccm-test/`): Simple tool for testing running CCM
3. **Test Interface** (`pkg/testing/`): Cloud-agnostic testing interface
4. **Cloud Provider Adapters** (`pkg/testing/`): Provider-specific implementations
5. **Mock Provider**: Simulated cloud provider for local testing

### **Supported Cloud Providers**
- **Mock**: Simulated provider for local testing
- **Existing**: Test your running CCM (no credentials needed)
- **AWS**: Amazon EKS support
- **GCP**: Google GKE support
- **Azure**: Azure AKS support
- **IBM Cloud**: IBM IKS support (structure ready)

## 📋 Prerequisites

### **For Local Testing (Mock Provider)**
- Go 1.24+
- No external dependencies

### **For Real Cluster Testing (Existing CCM)**
- Running Kubernetes cluster with CCM enabled
- Valid kubeconfig file
- No cloud credentials needed!

### **For Real Cloud Provider Testing**
- Running Kubernetes cluster with CCM enabled
- Valid kubeconfig file
- Cloud provider credentials with appropriate permissions

## 🔧 Usage Examples

### **Mock Provider Testing**
```bash
# Test all suites
./bin/e2e-test-runner --provider mock --suite all --verbose

# Test specific suite
./bin/e2e-test-runner --provider mock --suite loadbalancer --verbose
```

### **Existing CCM Testing (Recommended)**
```bash
# Simple test of your running CCM
./bin/existing-ccm-test --kubeconfig ~/.kube/config --verbose

# Full e2e test runner with existing CCM
./bin/e2e-test-runner --provider existing --kubeconfig ~/.kube/config --suite all --verbose
```

### **Real Cloud Provider Testing**
```bash
# AWS EKS
./bin/e2e-test-runner \
  --provider aws \
  --kubeconfig ~/.kube/config \
  --region us-west-2 \
  --zone us-west-2a \
  --cluster my-eks-cluster \
  --suite loadbalancer \
  --verbose

# GCP GKE
./bin/e2e-test-runner \
  --provider gcp \
  --kubeconfig ~/.kube/config \
  --region us-central1 \
  --zone us-central1-a \
  --cluster my-gke-cluster \
  --suite loadbalancer \
  --verbose

# Azure AKS
./bin/e2e-test-runner \
  --provider azure \
  --kubeconfig ~/.kube/config \
  --region eastus \
  --zone eastus-1 \
  --cluster my-aks-cluster \
  --suite loadbalancer \
  --verbose
```

## ⚙️ Configuration Options

### **Required Flags**
- `--provider`: Cloud provider (`mock`, `existing`, `aws`, `gcp`, `azure`, `ibmcloud`)
- `--kubeconfig`: Path to kubeconfig (not required for mock)

### **Cloud Provider Configuration**
- `--region`: Cloud provider region
- `--zone`: Cloud provider zone/availability zone
- `--cluster`: Cluster name
- `--prefix`: Resource prefix for test resources (default: `e2e-test`)

### **Test Execution**
- `--suite`: Test suite to run (`all`, `loadbalancer`, `nodes`, `routes`, `instances`, `zones`, `clusters`)
- `--timeout`: Test timeout (default: 30m)
- `--verbose`: Enable verbose output
- `--cleanup`: Clean up resources after tests (default: true)

## 🔄 CI/CD Integration

### **GitHub Actions Example**
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
      run: make build
    
    - name: Test with mock provider
      run: ./bin/e2e-test-runner --provider mock --suite all --verbose
    
    - name: Test existing CCM (if cluster available)
      run: |
        if [ -n "$KUBECONFIG" ]; then
          ./bin/existing-ccm-test --kubeconfig $KUBECONFIG --verbose
        fi
```

## 🛠️ Development

### **Local Development**
```bash
# Setup development environment
make dev-setup

# Run tests
make test

# Run e2e tests with mock
make test-e2e

# Build all binaries
make build-all

# Run linting
make lint

# Run all checks
make check
```

### **Adding New Cloud Providers**
1. Implement the `cloudprovider.Interface` for your provider
2. Create a provider adapter in `pkg/testing/real_cloud_provider.go`
3. Add provider initialization in `cmd/e2e-test-runner/main.go`
4. Update the `createCloudProvider` function

### **Adding New Test Suites**
1. Create test functions in `pkg/testing/test_suites.go`
2. Add test suite creation function
3. Update the `addTestSuites` function in the test runner

## 📚 Documentation

- **[E2E Testing Guide](docs/e2e-testing-guide.md)**: Detailed guide for e2e testing
- **[IBM Cloud Implementation Guide](docs/ibmcloud-implementation-guide.md)**: IBM Cloud specific implementation

## 🎯 Best Practices

1. **Start with Mock Provider**: Use mock provider for local development and testing
2. **Test Existing CCM**: Use existing CCM testing for validation without credentials
3. **Use Resource Prefixes**: Always use unique prefixes to avoid conflicts
4. **Enable Cleanup**: Always enable cleanup unless debugging
5. **Monitor Resources**: Keep an eye on cloud provider quotas and limits
6. **Use Verbose Logging**: Enable verbose output for better debugging

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
