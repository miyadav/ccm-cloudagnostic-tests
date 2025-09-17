# CCM Cloud-Agnostic Testing Framework

A comprehensive, cloud-agnostic testing framework for Kubernetes Cloud Controller Manager (CCM) functionality across different cloud providers.
It uses - https://github.com/miyadav/cloud-provider-testing-interface/ as interface to create tests

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

#### **New Ginkgo-Based Tests (Recommended)**
```bash
# Using Ginkgo CLI (best experience)
cd cmd/existing-ccm-test
ginkgo run -v -- --kubeconfig ~/.kube/config

# Using Go test (standard)
cd cmd/existing-ccm-test
go test -v --kubeconfig ~/.kube/config

# Using Makefile targets
make test-ginkgo-verbose --kubeconfig ~/.kube/config
```

#### **Legacy E2E Test Runner**
```bash
# Use the full e2e test runner
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

### ✅ **Modern Testing Framework**
- **Ginkgo/Gomega**: Industry-standard BDD testing framework
- **Structured Tests**: Clear test hierarchy with Describe/Context/It/By
- **Rich Assertions**: Powerful Gomega matchers for validation
- **JUnit Reports**: Prow-compatible XML output for CI/CD

### ✅ **Cloud-Agnostic Testing**
- Same test interface across all cloud providers
- Consistent test expectations and validation
- Easy to add new cloud providers

### ✅ **Multiple Testing Modes**
- **Mock Provider**: Local testing without cluster (fast, free)
- **Existing CCM**: Test your running CCM (no credentials needed)
- **Real Providers**: Full e2e testing with cloud credentials

### ✅ **Comprehensive Test Suites**
- **LoadBalancer**: Creation, updates, deletion, status, provider validation
- **Node Management**: Initialization, addresses, provider IDs, CCM processing
- **Route Management**: Creation, deletion, listing
- **Instances**: Existence, shutdown detection, metadata
- **Zones**: Information retrieval
- **Clusters**: Listing and master node detection

### ✅ **Production Ready**
- Resource cleanup and management
- Error handling and reporting
- CI/CD integration with Prow support
- Comprehensive documentation

## 🏗️ Architecture

### **Components**
1. **E2E Test Runner** (`cmd/e2e-test-runner/`): Legacy testing tool supporting all providers
2. **Ginkgo Tests** (`cmd/existing-ccm-test/`): Modern BDD tests using Ginkgo/Gomega framework
3. **Test Interface** (`pkg/testing/`): Cloud-agnostic testing interface
4. **Cloud Provider Adapters** (`pkg/testing/`): Provider-specific implementations
5. **Mock Provider**: Simulated cloud provider for local testing

### **Supported Cloud Providers**
- **Mock**: Simulated provider for local testing
- **Existing**: Test your running CCM (no credentials needed)
- **AWS**: Amazon EKS support
- **GCP**: Google GKE support
- **Azure**: Azure AKS support

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

### **Ginkgo-Based Testing (Recommended)**
```bash
# Using Ginkgo CLI (best experience with rich output)
cd cmd/existing-ccm-test
ginkgo run -v -- --kubeconfig ~/.kube/config

# Using Go test (standard output)
cd cmd/existing-ccm-test
go test -v --kubeconfig ~/.kube/config

# With JUnit output for CI/CD
cd cmd/existing-ccm-test
ginkgo run --junit-report=../../test-results/junit.xml -- --kubeconfig ~/.kube/config

# Using Makefile targets
make test-ginkgo-verbose --kubeconfig ~/.kube/config
make test-ginkgo-junit --kubeconfig ~/.kube/config
```

### **Legacy E2E Test Runner**
```bash
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

## 🧪 Test Structure & Features

### **Ginkgo Test Organization**
```go
Describe("CCM Load Balancer Tests", Label("loadbalancer"), func() {
    Context("LoadBalancer Service Creation", func() {
        It("should create a LoadBalancer service and wait for CCM to provision it", func() {
            By("Creating a test LoadBalancer service")
            By("Waiting for CCM to provision the load balancer")
            By("Cleaning up the service")
        })
    })
})
```

### **Test Categories**
- **LoadBalancer Tests**: Service creation, provisioning, provider validation
- **Node Management Tests**: CCM processing, provider metadata validation
- **Integration Tests**: End-to-end CCM workflow validation

### **Key Features**
- **Structured Output**: Clear test hierarchy with colored output
- **Step-by-Step Reporting**: `By()` statements show test progression
- **Label Filtering**: Run specific test categories with `--label-filter`
- **JUnit Reports**: Generate XML reports for CI/CD integration
- **Provider Validation**: Automatic cloud provider consistency checking

## ⚙️ Configuration Options

### **Ginkgo Test Flags**
- `--kubeconfig`: Path to kubeconfig file (required)
- `--namespace`: Test namespace (default: `ccm-test`)
- `--timeout`: Test timeout (default: 5m)
- `--verbose`: Enable verbose output
- `--junit-file`: Path to JUnit XML output file

### **Legacy E2E Test Runner Flags**
- `--provider`: Cloud provider (`mock`, `existing`, `aws`, `gcp`, `azure`)
- `--kubeconfig`: Path to kubeconfig (not required for mock)
- `--region`: Cloud provider region
- `--zone`: Cloud provider zone/availability zone
- `--cluster`: Cluster name
- `--prefix`: Resource prefix for test resources (default: `e2e-test`)
- `--suite`: Test suite to run (`all`, `loadbalancer`, `nodes`, `routes`, `instances`, `zones`, `clusters`)
- `--timeout`: Test timeout (default: 30m)
- `--verbose`: Enable verbose output
- `--cleanup`: Clean up resources after tests (default: true)

## 🔄 CI/CD Integration

### **Prow Integration (Kubernetes CI)**
```yaml
# prow-config.yaml
presubmits:
  kubernetes/ccm-cloudagnostic-tests:
  - name: ccm-ginkgo-tests
    always_run: true
    spec:
      containers:
      - image: golang:1.24
        command:
        - make
        - test-ginkgo-prow
        - --kubeconfig=/etc/kubeconfig/config
        volumeMounts:
        - name: kubeconfig
          mountPath: /etc/kubeconfig
          readOnly: true
      volumes:
      - name: kubeconfig
        secret:
          secretName: kubeconfig
```

### **GitHub Actions Example**
```yaml
name: CCM E2E Tests
on: [push, pull_request]

jobs:
  ginkgo-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.24'
    
    - name: Install Ginkgo
      run: go install github.com/onsi/ginkgo/v2/ginkgo@latest
    
    - name: Run Ginkgo tests with JUnit output
      run: |
        cd cmd/existing-ccm-test
        ginkgo run --junit-report=../../test-results/junit.xml -- --kubeconfig ${{ secrets.KUBECONFIG }}
    
    - name: Upload test results
      uses: actions/upload-artifact@v3
      if: always()
      with:
        name: test-results
        path: test-results/
  
  legacy-tests:
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
```

## 🛠️ Development

### **Local Development**
```bash
# Setup development environment
make dev-setup

# Run Ginkgo tests (recommended)
make test-ginkgo-verbose --kubeconfig ~/.kube/config
make test-ginkgo-junit --kubeconfig ~/.kube/config

# Run legacy tests
make test
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

#### **For Ginkgo Tests (Recommended)**
1. Add new test cases in `cmd/existing-ccm-test/ccm_test.go`
2. Use `Describe()`, `Context()`, `It()`, and `By()` for structure
3. Add appropriate labels for filtering
4. Use Gomega matchers for assertions

#### **For Legacy E2E Tests**
1. Create test functions in `pkg/testing/test_suites.go`
2. Add test suite creation function
3. Update the `addTestSuites` function in the test runner

## 📚 Documentation

- **[Ginkgo Refactoring Guide](docs/ginkgo-refactoring-guide.md)**: Complete guide to the new Ginkgo-based testing framework
- **[E2E Testing Guide](docs/e2e-testing-guide.md)**: Detailed guide for legacy e2e testing

## 🎯 Best Practices

1. **Use Ginkgo Tests**: Prefer the new Ginkgo-based tests for better structure and reporting
2. **Start with Mock Provider**: Use mock provider for local development and testing
3. **Test Existing CCM**: Use existing CCM testing for validation without credentials
4. **Use JUnit Reports**: Generate XML reports for CI/CD integration
5. **Use Resource Prefixes**: Always use unique prefixes to avoid conflicts
6. **Enable Cleanup**: Always enable cleanup unless debugging
7. **Monitor Resources**: Keep an eye on cloud provider quotas and limits
8. **Use Verbose Logging**: Enable verbose output for better debugging
9. **Label Your Tests**: Use appropriate labels for test filtering and organization

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
