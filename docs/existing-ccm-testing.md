# Testing Existing CCM Deployments

## Overview

This document describes how to test Cloud Controller Manager (CCM) implementations that are already running in a Kubernetes cluster, using indirect testing through Kubernetes resources.

## Architecture

### Two Testing Approaches

#### 1. Direct Cloud Provider Testing (Original)
- Tests the cloud provider interface directly
- Requires `cloudProvider.Interface` access
- Used with `--provider mock`, `--provider aws`, etc.
- Calls methods like `LoadBalancer()`, `Instances()`, `Routes()` directly

#### 2. Existing CCM Testing (New - Option 3)
- Tests a CCM that's already running in the cluster
- No direct cloud provider interface access
- Used with `--provider existing`
- Tests by creating Kubernetes resources and observing CCM behavior

## Test Suites for Existing CCM

### Load Balancer Test Suite
Tests CCM load balancer functionality by creating and managing Services:

1. **CreateLoadBalancer**
   - Creates a `LoadBalancer` type Service
   - Waits for CCM to provision the load balancer
   - Verifies ingress points are populated
   - Cleans up the service

2. **UpdateLoadBalancer**
   - Creates a Service with one port
   - Updates it to add additional ports
   - Verifies CCM handles the update correctly
   - Ensures load balancer remains available

3. **DeleteLoadBalancer**
   - Creates and provisions a load balancer
   - Deletes the service
   - Verifies CCM cleans up the load balancer

### Node Management Test Suite
Tests CCM node management by examining existing nodes:

1. **GetExistingNodes**
   - Lists all nodes in the cluster
   - Displays node information (name, status, provider ID)
   - Verifies nodes are being managed

2. **VerifyNodeMetadata**
   - Checks each node for CCM-populated metadata:
     - Provider ID (required)
     - Node addresses (required)
     - Instance type label (optional)
   - Fails if any node is missing required metadata

## Usage

### Running All Tests
```bash
./bin/e2e-test-runner --provider existing --kubeconfig ~/.kube/config --suite all --verbose
```

### Running Specific Test Suites
```bash
# Load balancer tests only
./bin/e2e-test-runner --provider existing --kubeconfig ~/.kube/config --suite loadbalancer --verbose

# Node management tests only
./bin/e2e-test-runner --provider existing --kubeconfig ~/.kube/config --suite nodes --verbose
```

### Test Output Example
```
=== CCM E2E Test Results ===
Total Duration: 39.46938925s
Test Summary: 5 total, 5 passed, 0 failed, 0 skipped

Detailed Results:
  PASSED: CreateLoadBalancer (5.441154375s)
  PASSED: UpdateLoadBalancer (15.876363166s)
  PASSED: DeleteLoadBalancer (17.5934565s)
  PASSED: GetExistingNodes (410.713458ms)
  PASSED: VerifyNodeMetadata (147.577917ms)

✅ All tests passed: 5 passed out of 5 total
```

## Available Test Suites by Provider

### Direct Cloud Provider Testing
Available suites: `all`, `loadbalancer`, `nodes`, `routes`, `instances`, `zones`, `clusters`

### Existing CCM Testing
Available suites: `all`, `loadbalancer`, `nodes`

**Note:** Routes, Instances, Zones, and Clusters test suites require direct cloud provider interface access and are not available for existing CCM testing.

## Implementation Details

### File Structure
```
pkg/testing/
├── existing_ccm_test_interface.go     # Interface for existing CCM testing
├── existing_ccm_test_suites.go        # Test suites for existing CCM (NEW)
├── test_suites.go                      # Test suites for direct cloud provider testing
├── ccm_test_interface.go              # Interface for direct testing
└── real_cloud_provider.go             # Real cloud provider adapter

cmd/e2e-test-runner/
└── main.go                             # Updated to route to correct test suites
```

### Key Components

#### ExistingCCMTestInterface
```go
// GetCloudProvider returns nil for existing CCM testing
func (e *ExistingCCMTestInterface) GetCloudProvider() cloudprovider.Interface {
    return nil  // No direct access - CCM is separate!
}
```

#### Test Suite Selection
```go
func addTestSuites(runner *ccmtesting.TestRunner, suite, provider string) {
    if provider == "existing" {
        // Use test suites that test indirectly via Kubernetes resources
        addExistingCCMTestSuites(runner, suite)
        return
    }
    
    // Use test suites that access cloud provider interface directly
    // ... standard test suites ...
}
```

## Requirements

### Prerequisites
- A Kubernetes cluster with CCM already running
- `kubeconfig` file with access to the cluster
- CCM must support the features being tested (e.g., LoadBalancer services)

### Timeouts
- CreateLoadBalancer: 5 minutes
- UpdateLoadBalancer: 5 minutes
- DeleteLoadBalancer: 3 minutes
- GetExistingNodes: 2 minutes
- VerifyNodeMetadata: 2 minutes

## Troubleshooting

### Common Issues

#### "nodes already exists" error
**Cause:** Leftover test resources from previous failed runs

**Solution:**
```bash
# Clean up test nodes
kubectl get nodes -l test-prefix=e2e-test
kubectl delete nodes -l test-prefix=e2e-test

# Clean up test namespaces
kubectl get namespaces | grep ccm-test
kubectl delete namespace <namespace-name>
```

#### "CCM failed to provision load balancer"
**Cause:** CCM may not support LoadBalancer services or is misconfigured

**Solution:**
1. Check CCM logs: `kubectl logs -n kube-system <ccm-pod>`
2. Verify CCM has correct cloud credentials
3. Check cloud provider quota/limits

#### "nodes are missing CCM metadata"
**Cause:** CCM may not have initialized nodes properly

**Solution:**
1. Check if CCM is running: `kubectl get pods -n kube-system | grep cloud-controller`
2. Review CCM logs for errors
3. Verify nodes have correct provider IDs

## Best Practices

1. **Run tests in non-production clusters** - Tests create and delete resources
2. **Use `--cleanup` flag** - Ensures resources are cleaned up even on failure
3. **Monitor CCM logs** - Run `kubectl logs -f -n kube-system <ccm-pod>` in another terminal
4. **Start with specific suites** - Test `loadbalancer` or `nodes` individually before running `all`
5. **Review test logs** - Use `--verbose` flag to see detailed test execution

## Comparison: Direct vs Existing CCM Testing

| Aspect | Direct Testing | Existing CCM Testing |
|--------|---------------|---------------------|
| Cloud Provider Access | Direct interface access | No direct access |
| Test Method | Call provider methods | Create K8s resources |
| Use Case | Unit/integration tests | E2E tests in real clusters |
| Available Suites | All 6 suites | 2 suites (LB, Nodes) |
| Setup Required | Mock/test provider | Running CCM in cluster |
| Test Fidelity | Tests provider code | Tests full CCM behavior |

## Future Enhancements

Potential additional test suites for existing CCM:
- **Routes Testing** - Test route management via pod networking
- **Zones Testing** - Verify zone/region labels on nodes
- **Service Annotations** - Test cloud-specific service annotations
- **Multiple Load Balancers** - Test concurrent LB provisioning
- **Load Balancer Health Checks** - Verify health check configuration

## Related Documentation

- [E2E Testing Guide](e2e-testing-guide.md) - General E2E testing documentation
- [Cloud Provider Testing Interface](https://github.com/miyadav/cloud-provider-testing-interface) - Testing framework used

