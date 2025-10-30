# Automatic Test Resource Cleanup

## Overview

The CCM test framework includes multiple layers of automatic cleanup to ensure test resources are properly removed, even when tests fail partway through execution.

## Cleanup Layers

### Layer 1: Test-Level Cleanup (defer blocks)

Each individual test function uses Go's `defer` statement to ensure resources are cleaned up when the test completes, regardless of success or failure.

#### Example: Load Balancer Test
```go
func testExistingCCMCreateLoadBalancer(ti ccmtesting.TestInterface) error {
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

    // ... rest of test ...
}
```

**Benefits:**
- ✅ Cleanup happens immediately after test completes
- ✅ Works even if test panics or returns error early
- ✅ Specific to the resources that test created

### Layer 2: Test Suite Teardown

Each test suite has a `Teardown` function that runs after all tests in the suite complete. This catches any resources that individual tests missed.

#### Load Balancer Suite Teardown
```go
func teardownExistingCCMLoadBalancerTestSuite(ti ccmtesting.TestInterface) error {
    existingCCM := ti.(*ExistingCCMTestInterface)
    ctx := context.Background()
    
    // List all LoadBalancer services in test namespace
    services, err := existingCCM.kubeClient.CoreV1().Services(
        existingCCM.GetNamespace()
    ).List(ctx, metav1.ListOptions{})
    
    // Delete any LoadBalancer services that remain
    for _, svc := range services.Items {
        if svc.Spec.Type == v1.ServiceTypeLoadBalancer {
            _ = existingCCM.kubeClient.CoreV1().Services(
                existingCCM.GetNamespace()
            ).Delete(ctx, svc.Name, metav1.DeleteOptions{})
        }
    }
    
    return nil
}
```

#### Node Suite Teardown
```go
func teardownExistingCCMNodeTestSuite(ti ccmtesting.TestInterface) error {
    existingCCM := ti.(*ExistingCCMTestInterface)
    ctx := context.Background()
    
    // Delete any nodes with test label
    nodes, err := existingCCM.kubeClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{
        LabelSelector: "test-prefix=e2e-test",
    })
    
    for _, node := range nodes.Items {
        _ = existingCCM.kubeClient.CoreV1().Nodes().Delete(
            ctx, node.Name, metav1.DeleteOptions{}
        )
    }
    
    return nil
}
```

**Benefits:**
- ✅ Catches resources missed by individual tests
- ✅ Bulk cleanup of all resources in test namespace
- ✅ Label-based cleanup for cluster-scoped resources (nodes)

### Layer 3: Test Environment Teardown

The test environment's `TeardownTestEnvironment()` function runs at the very end, cleaning up the entire test namespace.

```go
func (e *ExistingCCMTestInterface) TeardownTestEnvironment() error {
    // Delete test namespace (cascade deletes all namespaced resources)
    err := e.kubeClient.CoreV1().Namespaces().Delete(
        context.Background(), 
        e.namespace, 
        metav1.DeleteOptions{GracePeriodSeconds: func() *int64 { v := int64(0); return &v }()}
    )
    
    // Wait for namespace to be fully deleted
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return nil
        case <-ticker.C:
            _, err := e.kubeClient.CoreV1().Namespaces().Get(ctx, e.namespace, metav1.GetOptions{})
            if err != nil {
                // Namespace deleted successfully
                return nil
            }
        }
    }
}
```

**Benefits:**
- ✅ Deletes entire test namespace
- ✅ Kubernetes cascade deletion removes all namespaced resources
- ✅ Waits for deletion to complete
- ✅ Always runs (via `defer` in main.go)

## Cleanup Order

```
Test Execution
    ↓
[1] Test-level defer cleanup (individual resources)
    ↓
[2] Test suite teardown (remaining suite resources)
    ↓
[3] Environment teardown (namespace + all resources)
```

## Resource Tracking

### Namespace-Scoped Resources
All test resources are created in a timestamped test namespace:
```
ccm-test-20251030-103035
```

**Automatically cleaned up by:**
- Test-level defers (immediate)
- Suite teardown (bulk)
- Namespace deletion (cascade)

### Cluster-Scoped Resources (Nodes)
Nodes cannot be in a namespace, so they use labels:
```yaml
labels:
  test-prefix: e2e-test
```

**Automatically cleaned up by:**
- Suite teardown (label selector)

## Cleanup in Different Scenarios

### Scenario 1: Test Passes ✅
```
Test runs → Resources created → Test completes → defer cleanup → Suite teardown → Namespace deleted
```
**Result:** All resources cleaned up automatically

### Scenario 2: Test Fails ❌
```
Test runs → Resources created → Test FAILS → defer cleanup (still runs!) → Suite teardown → Namespace deleted
```
**Result:** All resources cleaned up automatically

### Scenario 3: Test Panics 💥
```
Test runs → Resources created → PANIC → defer cleanup (still runs!) → Suite teardown → Namespace deleted
```
**Result:** All resources cleaned up automatically

### Scenario 4: Test Times Out ⏰
```
Test runs → Resources created → TIMEOUT → defer cleanup → Suite teardown → Namespace deleted
```
**Result:** All resources cleaned up automatically

### Scenario 5: Program Crashes/Killed 🔥
```
Test runs → Resources created → CRASH/KILL
```
**Result:** Resources remain (but labeled for easy manual cleanup)

**Manual cleanup for crash scenario:**
```bash
# Clean up test namespaces
kubectl get namespaces | grep ccm-test
kubectl delete namespace <namespace-name>

# Clean up test nodes
kubectl get nodes -l test-prefix=e2e-test
kubectl delete nodes -l test-prefix=e2e-test
```

## Verification

### Check for Leftover Resources
```bash
# Check for test namespaces
kubectl get namespaces | grep ccm-test

# Check for test nodes
kubectl get nodes -l test-prefix=e2e-test

# Check for test services in namespace (if namespace still exists)
kubectl get services -n <test-namespace> --field-selector spec.type=LoadBalancer
```

### Expected Output After Successful Test Run
```bash
$ kubectl get namespaces | grep ccm-test
(no output - all cleaned up)

$ kubectl get nodes -l test-prefix=e2e-test
No resources found
```

## Best Practices

### 1. Always Use defer for Resource Cleanup
```go
resource, err := createResource(ctx, config)
if err != nil {
    return err
}

// GOOD: Defer cleanup immediately after creation
defer func() {
    if err := deleteResource(ctx, resource); err != nil {
        log.Printf("Warning: cleanup failed: %v", err)
    }
}()
```

### 2. Label All Test Resources
```go
// Namespace-scoped resources
metadata:
  labels:
    test-prefix: e2e-test

// Cluster-scoped resources (nodes)
metadata:
  labels:
    test-prefix: e2e-test
```

### 3. Use Unique Names with Timestamps
```go
// Test namespace
namespace := "ccm-test-" + time.Now().Format("20060102-150405")

// Test services
serviceName := "test-lb-create"  // Unique per test
```

### 4. Log Cleanup Actions
```go
defer func() {
    if err := cleanup(); err != nil {
        ti.GetTestResults().AddLog(
            fmt.Sprintf("Warning: failed to cleanup: %v", err)
        )
    }
}()
```

### 5. Implement Idempotent Cleanup
```go
// Don't fail if resource already deleted
err := deleteResource()
if err != nil && !errors.IsNotFound(err) {
    return err  // Only fail on real errors
}
```

## Testing the Cleanup

### Test Normal Flow
```bash
./bin/e2e-test-runner --provider existing --kubeconfig ~/.kube/config --suite all --verbose

# Verify cleanup
kubectl get namespaces | grep ccm-test
kubectl get nodes -l test-prefix=e2e-test
```

### Test Failure Scenario
To verify cleanup works even on failure, you can:

1. **Temporarily break a test:**
   ```go
   // Add this to cause a test to fail
   return fmt.Errorf("intentional failure")
   ```

2. **Run tests:**
   ```bash
   ./bin/e2e-test-runner --provider existing --kubeconfig ~/.kube/config --suite loadbalancer
   ```

3. **Verify cleanup still happened:**
   ```bash
   kubectl get namespaces | grep ccm-test
   kubectl get services --all-namespaces | grep test-lb
   ```

## Troubleshooting

### Resources Not Being Cleaned Up

**Problem:** Resources remain after tests complete

**Possible Causes:**
1. Test panicked before defer could run
2. Cleanup function has a bug
3. Kubernetes API errors during deletion
4. Finalizers preventing deletion

**Solutions:**
```bash
# Force delete namespace
kubectl delete namespace <namespace> --force --grace-period=0

# Remove finalizers from stuck resources
kubectl patch service <service-name> -n <namespace> -p '{"metadata":{"finalizers":[]}}' --type=merge

# Delete stuck nodes
kubectl delete node <node-name> --force --grace-period=0
```

### Cleanup Taking Too Long

**Problem:** Namespace deletion hangs

**Cause:** Load balancer cleanup takes time (cloud provider must delete actual LB)

**Solution:**
- Wait up to 5 minutes for load balancer cleanup
- Check CCM logs to see if it's processing the deletion
- Verify cloud provider credentials are valid

## Implementation Files

- **Test Functions with defer:** `pkg/testing/existing_ccm_test_suites.go`
- **Suite Teardown:** `pkg/testing/existing_ccm_test_suites.go`
- **Environment Teardown:** `pkg/testing/existing_ccm_test_interface.go`
- **Main Cleanup Coordination:** `cmd/e2e-test-runner/main.go`

## Summary

The multi-layer cleanup approach ensures:
- ✅ Resources are cleaned up immediately after use (defer)
- ✅ Leftover resources are caught by suite teardown
- ✅ Namespace deletion provides final cleanup
- ✅ Test failures don't leave resources behind
- ✅ Manual cleanup is easy (labels + namespaces)
- ✅ No manual intervention required for normal operation

