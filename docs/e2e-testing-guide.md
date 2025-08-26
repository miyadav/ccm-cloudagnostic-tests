# CCM Cloud-Agnostic E2E Testing Guide

This guide explains how to use the CCM Cloud-Agnostic Testing Framework for end-to-end (e2e) testing on real Kubernetes clusters with actual cloud providers.

## Overview

The framework provides a cloud-agnostic way to test Cloud Controller Manager (CCM) functionality across different cloud providers (AWS, GCP, Azure) while maintaining the same test interface and expectations.

## Prerequisites

### 1. Cluster Access
- A running Kubernetes cluster with CCM enabled
- Valid kubeconfig file for cluster access
- Appropriate RBAC permissions for testing

### 2. Cloud Provider Credentials
- Valid credentials for your cloud provider
- Proper IAM/service account permissions
- Network access to cloud provider APIs

### 3. Build the E2E Test Runner
```bash
# Build the e2e test runner
go build -o bin/e2e-test-runner cmd/e2e-test-runner/main.go
```

## Usage

### Basic Command Structure
```bash
./bin/e2e-test-runner \
  --provider <cloud-provider> \
  --kubeconfig <path-to-kubeconfig> \
  --region <cloud-region> \
  --zone <cloud-zone> \
  --cluster <cluster-name> \
  [additional-options]
```

### Example Commands

#### 1. Test with Mock Provider (Local Testing)
```bash
./bin/e2e-test-runner \
  --provider mock \
  --suite loadbalancer \
  --verbose
```

#### 2. Test AWS EKS Cluster
```bash
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

#### 3. Test GCP GKE Cluster
```bash
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

#### 4. Test Azure AKS Cluster
```bash
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

## Command Line Options

### Required Flags
- `--provider`: Cloud provider name (`aws`, `gcp`, `azure`, `mock`)
- `--kubeconfig`: Path to kubeconfig file (not required for mock)

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

### Credentials (Optional)
- `--credentials`: Path to credentials file

## Cloud Provider Setup

### AWS EKS Setup

1. **Install AWS CLI and configure credentials**:
```bash
aws configure
```

2. **Update kubeconfig for EKS cluster**:
```bash
aws eks update-kubeconfig --region us-west-2 --name my-eks-cluster
```

3. **Verify cluster access**:
```bash
kubectl get nodes
```

4. **Run e2e tests**:
```bash
./bin/e2e-test-runner \
  --provider aws \
  --kubeconfig ~/.kube/config \
  --region us-west-2 \
  --zone us-west-2a \
  --cluster my-eks-cluster \
  --suite all \
  --verbose
```

### GCP GKE Setup

1. **Install gcloud CLI and authenticate**:
```bash
gcloud auth login
gcloud config set project my-project-id
```

2. **Get cluster credentials**:
```bash
gcloud container clusters get-credentials my-gke-cluster --zone us-central1-a
```

3. **Verify cluster access**:
```bash
kubectl get nodes
```

4. **Run e2e tests**:
```bash
./bin/e2e-test-runner \
  --provider gcp \
  --kubeconfig ~/.kube/config \
  --region us-central1 \
  --zone us-central1-a \
  --cluster my-gke-cluster \
  --suite all \
  --verbose
```

### Azure AKS Setup

1. **Install Azure CLI and login**:
```bash
az login
az account set --subscription my-subscription-id
```

2. **Get cluster credentials**:
```bash
az aks get-credentials --resource-group my-resource-group --name my-aks-cluster
```

3. **Verify cluster access**:
```bash
kubectl get nodes
```

4. **Run e2e tests**:
```bash
./bin/e2e-test-runner \
  --provider azure \
  --kubeconfig ~/.kube/config \
  --region eastus \
  --zone eastus-1 \
  --cluster my-aks-cluster \
  --suite all \
  --verbose
```

## Test Suites

### LoadBalancer Tests
Tests cloud provider load balancer functionality:
- Create load balancer
- Update load balancer
- Delete load balancer
- Load balancer status
- Health checks

### Node Management Tests
Tests cloud provider node management:
- Node initialization
- Node addresses
- Provider ID
- Instance type
- Zone information

### Route Management Tests
Tests cloud provider route management:
- Create routes
- Delete routes
- List routes

### Instances Tests
Tests cloud provider instances functionality:
- Instance existence
- Instance shutdown detection
- Instance metadata

### Zones Tests
Tests cloud provider zones functionality:
- Get zone information
- Get zone by provider ID

### Clusters Tests
Tests cloud provider clusters functionality:
- List clusters
- Master node detection

## Resource Management

### Resource Prefix
All test resources are created with a configurable prefix to avoid conflicts:
- Services: `{prefix}-lb-{name}`
- Routes: `{prefix}-route-{name}`
- Nodes: `{prefix}-node-{name}`

### Cleanup
The framework automatically cleans up test resources after tests complete:
- Load balancers are deleted
- Routes are removed
- Test services are cleaned up

To disable cleanup (for debugging):
```bash
./bin/e2e-test-runner --cleanup=false
```

## Troubleshooting

### Common Issues

1. **Cluster Connection Failed**
   ```
   Error: Failed to connect to cluster
   ```
   - Verify kubeconfig path is correct
   - Check cluster is running and accessible
   - Ensure proper RBAC permissions

2. **Cloud Provider Authentication Failed**
   ```
   Error: Failed to initialize cloud provider
   ```
   - Verify cloud provider credentials
   - Check IAM/service account permissions
   - Ensure network access to cloud APIs

3. **Resource Creation Failed**
   ```
   Error: Failed to create load balancer
   ```
   - Check cloud provider quotas
   - Verify network configuration
   - Ensure proper resource permissions

4. **Test Timeout**
   ```
   Error: Test execution failed: context deadline exceeded
   ```
   - Increase timeout with `--timeout` flag
   - Check cluster performance
   - Verify cloud provider API responsiveness

### Debug Mode
Enable debug logging for detailed troubleshooting:
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

## CI/CD Integration

### GitHub Actions Example
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

### Jenkins Pipeline Example
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

## Best Practices

1. **Use Resource Prefixes**: Always use unique prefixes to avoid conflicts
2. **Set Appropriate Timeouts**: Cloud operations can take time, set realistic timeouts
3. **Enable Cleanup**: Always enable cleanup unless debugging
4. **Monitor Resources**: Keep an eye on cloud provider quotas and limits
5. **Use Verbose Logging**: Enable verbose output for better debugging
6. **Test in Isolation**: Run tests in dedicated test clusters when possible
7. **Monitor Costs**: Be aware of cloud provider costs for resource creation

## Extending the Framework

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

### Custom Test Configuration
1. Extend the `RealCloudProviderConfig` struct
2. Add new command line flags
3. Update provider initialization logic
4. Modify test functions to use new configuration
