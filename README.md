# Cloud-Agnostic CCM Testing Framework

A comprehensive, cloud-agnostic testing framework for Kubernetes Cloud Controller Manager (CCM) implementations. This framework provides a standardized way to test cloud provider functionality across different cloud platforms while maintaining independence from specific cloud provider implementations.

## Overview

The Cloud-Agnostic CCM Testing Framework is built on top of the [cloud-provider-testing-interface](https://github.com/miyadav/cloud-provider-testing-interface) and provides:

- **Cloud-Agnostic Testing**: Test cloud provider functionality without being tied to specific cloud implementations
- **Comprehensive Test Suites**: Pre-built test suites for all major CCM functionality
- **Mock Cloud Provider**: A complete mock implementation for testing and development
- **Flexible Configuration**: Configurable test environments and parameters
- **Extensible Architecture**: Easy to extend with custom test cases and cloud providers

## Features

### Test Suites

The framework includes comprehensive test suites for all major CCM functionality:

- **Load Balancer Tests**: Create, update, delete, and manage load balancers
- **Node Management Tests**: Node initialization, addressing, provider IDs, and zone management
- **Route Management Tests**: Route creation, deletion, and listing
- **Instances Tests**: Instance existence, shutdown detection, and metadata
- **Zones Tests**: Zone information retrieval and management
- **Clusters Tests**: Cluster listing and master node detection

### Cloud Provider Independence

The framework is designed to work with any cloud provider that implements the standard Kubernetes cloud provider interface:

```go
type Interface interface {
    Initialize(clientBuilder ControllerClientBuilder, stop <-chan struct{})
    LoadBalancer() (LoadBalancer, bool)
    Instances() (Instances, bool)
    InstancesV2() (InstancesV2, bool)
    Zones() (Zones, bool)
    Clusters() (Clusters, bool)
    Routes() (Routes, bool)
    ProviderName() string
    HasClusterID() bool
}
```

### Mock Cloud Provider

A complete mock cloud provider implementation is included for:

- Development and testing without real cloud resources
- CI/CD pipeline testing
- Learning and understanding cloud provider interfaces
- Validating test framework functionality

## Installation

### Prerequisites

- Go 1.24 or later
- Kubernetes client-go libraries
- Access to a Kubernetes cluster (optional, for integration testing)

### Building

```bash
# Clone the repository
git clone https://github.com/kubernetes/ccm-cloudagnostic-tests.git
cd ccm-cloudagnostic-tests

# Build the test runner
go build -o bin/test-runner cmd/test-runner/main.go

# Run tests
go test ./pkg/testing/...
```

## Usage

### Command Line Test Runner

The framework includes a command-line test runner with extensive configuration options:

```bash
# Run all test suites
./bin/test-runner

# Run specific test suite
./bin/test-runner -suite=loadbalancer

# Run with verbose output
./bin/test-runner -suite=all -verbose

# Run with custom configuration
./bin/test-runner \
  -suite=loadbalancer \
  -provider=mock \
  -cluster=my-cluster \
  -region=us-west-1 \
  -zone=us-west-1a \
  -timeout=5m \
  -verbose

# Run specific tests
./bin/test-runner -tests=CreateLoadBalancer,UpdateLoadBalancer

# Skip specific tests
./bin/test-runner -skip=LoadBalancerHealthCheck

# Output in JSON format
./bin/test-runner -output=json
```

### Available Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-suite` | Test suite to run (all, loadbalancer, nodes, routes, instances, zones, clusters) | `all` |
| `-verbose` | Enable verbose output | `false` |
| `-timeout` | Test timeout | `10m` |
| `-provider` | Cloud provider name | `mock` |
| `-cluster` | Cluster name | `test-cluster` |
| `-region` | Region | `test-region` |
| `-zone` | Zone | `test-zone` |
| `-cleanup` | Clean up resources after tests | `true` |
| `-mock-external` | Use mock external services | `true` |
| `-output` | Output format (text, json) | `text` |
| `-tests` | Comma-separated list of specific tests to run | `` |
| `-skip` | Comma-separated list of tests to skip | `` |
| `-log-level` | Log level (debug, info, warn, error) | `info` |

### Programmatic Usage

You can also use the framework programmatically in your own tests:

```go
package main

import (
    "context"
    "time"
    
    "github.com/kubernetes/ccm-cloudagnostic-tests/pkg/testing"
    testing "github.com/miyadav/cloud-provider-testing-interface"
)

func main() {
    // Create your cloud provider instance
    cloudProvider := &YourCloudProvider{}
    
    // Create test interface
    testImpl := testing.NewCCMTestInterface(cloudProvider)
    
    // Create test configuration
    config := &testing.TestConfig{
        ProviderName:         "your-cloud-provider",
        ClusterName:          "test-cluster",
        Region:               "us-west-1",
        Zone:                 "us-west-1a",
        TestTimeout:          5 * time.Minute,
        CleanupResources:     true,
        MockExternalServices: false,
    }
    
    // Setup test environment
    err := testImpl.SetupTestEnvironment(config)
    if err != nil {
        log.Fatalf("Failed to setup test environment: %v", err)
    }
    defer testImpl.TeardownTestEnvironment()
    
    // Create test runner
    runner := testing.NewTestRunner(testImpl)
    
    // Add test suites
    runner.AddTestSuite(testing.CreateLoadBalancerTestSuite())
    runner.AddTestSuite(testing.CreateNodeTestSuite())
    
    // Run tests
    ctx := context.Background()
    err = runner.RunTests(ctx)
    if err != nil {
        log.Fatalf("Test execution failed: %v", err)
    }
    
    // Get results
    summary := runner.GetSummary()
    fmt.Printf("Tests completed: %d passed, %d failed\n", 
        summary.PassedTests, summary.FailedTests)
}
```

## Test Suites

### Load Balancer Test Suite

Tests load balancer functionality including creation, updates, deletion, and status management.

**Tests included:**
- `CreateLoadBalancer`: Test creating a new load balancer
- `UpdateLoadBalancer`: Test updating an existing load balancer
- `DeleteLoadBalancer`: Test deleting a load balancer
- `LoadBalancerStatus`: Test load balancer status updates
- `LoadBalancerHealthCheck`: Test health check functionality

### Node Management Test Suite

Tests node management functionality including initialization, addressing, and metadata.

**Tests included:**
- `NodeInitialization`: Test node initialization and registration
- `NodeAddresses`: Test node address management
- `NodeProviderID`: Test provider ID management
- `NodeInstanceType`: Test instance type detection
- `NodeZones`: Test zone management

### Route Management Test Suite

Tests route management functionality including creation, deletion, and listing.

**Tests included:**
- `CreateRoute`: Test creating a new route
- `DeleteRoute`: Test deleting a route
- `ListRoutes`: Test listing existing routes

### Instances Test Suite

Tests instance-related functionality including existence checks and metadata.

**Tests included:**
- `InstanceExists`: Test instance existence check
- `InstanceShutdown`: Test instance shutdown detection
- `InstanceMetadata`: Test instance metadata retrieval

### Zones Test Suite

Tests zone-related functionality including zone information retrieval.

**Tests included:**
- `GetZone`: Test zone information retrieval
- `GetZoneByProviderID`: Test zone retrieval by provider ID

### Clusters Test Suite

Tests cluster-related functionality including cluster listing and master node detection.

**Tests included:**
- `ListClusters`: Test listing clusters
- `Master`: Test master node detection

## Extending the Framework

### Adding Custom Test Suites

You can create custom test suites for your specific needs:

```go
func CreateCustomTestSuite() testing.TestSuite {
    return testing.TestSuite{
        Name:        "CustomTests",
        Description: "Custom test suite for specific functionality",
        Setup:       setupCustomTestSuite,
        Teardown:    teardownCustomTestSuite,
        Tests: []testing.Test{
            {
                Name:        "CustomTest1",
                Description: "Custom test 1",
                Run:         testCustomFunctionality1,
                Timeout:     2 * time.Minute,
            },
            {
                Name:        "CustomTest2",
                Description: "Custom test 2",
                Run:         testCustomFunctionality2,
                Timeout:     1 * time.Minute,
            },
        },
    }
}

func setupCustomTestSuite(ti testing.TestInterface) error {
    // Setup code for custom test suite
    return nil
}

func teardownCustomTestSuite(ti testing.TestInterface) error {
    // Cleanup code for custom test suite
    return nil
}

func testCustomFunctionality1(ti testing.TestInterface) error {
    // Test implementation
    return nil
}

func testCustomFunctionality2(ti testing.TestInterface) error {
    // Test implementation
    return nil
}
```

### Adding Custom Cloud Providers

To add support for a new cloud provider:

1. Implement the `cloudprovider.Interface`
2. Add the provider to the `createCloudProvider` function in `cmd/test-runner/main.go`
3. Create provider-specific configuration if needed

```go
func createCloudProvider(providerName string) (cloudprovider.Interface, error) {
    switch strings.ToLower(providerName) {
    case "mock":
        return testing.NewMockCloudProvider(), nil
    case "your-provider":
        return &YourCloudProvider{}, nil
    default:
        return nil, fmt.Errorf("unsupported cloud provider: %s", providerName)
    }
}
```

## Configuration

### Environment Variables

The framework supports configuration through environment variables:

- `CLOUD_PROVIDER_CONFIG`: Path to cloud provider configuration file
- `TEST_TIMEOUT`: Default test timeout
- `CLEANUP_RESOURCES`: Whether to clean up resources after tests
- `MOCK_EXTERNAL_SERVICES`: Whether to use mock external services

### Configuration Files

You can provide cloud provider-specific configuration through configuration files:

```yaml
# config/test-config.yaml
test:
  provider:
    name: "your-cloud-provider"
    region: "us-west-1"
    zone: "us-west-1a"
  
  cluster:
    name: "test-cluster"
    version: "1.24.0"
  
  resources:
    cleanup: true
    timeout: "5m"
  
  external_services:
    mock: false  # Use real cloud services for integration tests
```

## Integration with CI/CD

### GitHub Actions

Example GitHub Actions workflow for running tests:

```yaml
name: CCM Tests

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24'
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run unit tests
      run: go test -v ./pkg/testing/...
    
    - name: Run integration tests
      run: go test -v ./pkg/testing/integration_test.go
      env:
        CLOUD_PROVIDER_CONFIG: ${{ secrets.CLOUD_PROVIDER_CONFIG }}
    
    - name: Run test runner
      run: go run cmd/test-runner/main.go -suite=all -verbose
      env:
        CLOUD_PROVIDER_CONFIG: ${{ secrets.CLOUD_PROVIDER_CONFIG }}
```

### Jenkins

Example Jenkins pipeline:

```groovy
pipeline {
    agent any
    
    stages {
        stage('Test') {
            steps {
                sh 'go test -v ./pkg/testing/...'
                sh 'go run cmd/test-runner/main.go -suite=all -verbose'
            }
        }
    }
}
```

## Best Practices

### Test Design

1. **Keep tests independent**: Each test should be able to run independently
2. **Use descriptive names**: Test names should clearly describe what they test
3. **Handle cleanup properly**: Always clean up resources created during tests
4. **Use appropriate timeouts**: Set realistic timeouts for cloud operations
5. **Test error conditions**: Include tests for error scenarios and edge cases

### Resource Management

1. **Track created resources**: Keep track of all resources created during tests
2. **Implement proper cleanup**: Ensure resources are cleaned up even if tests fail
3. **Use unique names**: Use unique resource names to avoid conflicts
4. **Handle timeouts**: Implement proper timeout handling for long-running operations

### Configuration Management

1. **Use environment variables**: Store sensitive configuration in environment variables
2. **Provide defaults**: Provide sensible defaults for all configuration options
3. **Validate configuration**: Validate configuration before running tests
4. **Document configuration**: Document all configuration options and their effects

## Troubleshooting

### Common Issues

1. **Test timeouts**: Increase timeout values for slow cloud operations
2. **Resource conflicts**: Use unique resource names and proper cleanup
3. **Authentication issues**: Ensure proper cloud provider authentication
4. **Network issues**: Check network connectivity and firewall rules

### Debugging

Enable verbose logging to debug issues:

```bash
./bin/test-runner -verbose -log-level=debug
```

### Getting Help

- Check the [examples](pkg/testing/integration_test.go) for usage patterns
- Review the [test implementations](pkg/testing/test_suites.go) for best practices
- Open an issue in the [repository](https://github.com/kubernetes/ccm-cloudagnostic-tests/issues)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

### Code Style

- Follow Go coding standards
- Add comments for exported functions
- Include tests for new functionality
- Update documentation as needed

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built on top of the [cloud-provider-testing-interface](https://github.com/miyadav/cloud-provider-testing-interface)
- Inspired by the Kubernetes cloud provider architecture
- Community contributions and feedback
