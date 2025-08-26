# IBM Cloud Provider Implementation Guide

This guide explains how to complete the IBM Cloud provider implementation for the CCM Cloud-Agnostic Testing Framework.

## Current Status

The IBM Cloud provider is now supported in the framework with:
- ✅ Provider recognition in the e2e test runner
- ✅ IBM Cloud provider adapter structure
- ✅ Configuration handling
- ❌ Actual cloud provider implementation (placeholder)

## Implementation Steps

### 1. Install IBM Cloud Dependencies

Add the IBM Cloud Kubernetes Service SDK to your `go.mod`:

```bash
go get github.com/IBM-Cloud/ibm-cloud-cli-sdk
go get github.com/IBM-Cloud/ibm-cloud-kubernetes-service
```

### 2. Complete the IBM Cloud Provider Implementation

Replace the placeholder function in `pkg/testing/real_cloud_provider.go`:

```go
func initializeIBMCloudProvider(config *RealCloudProviderConfig) (cloudprovider.Interface, error) {
    // Initialize IBM Cloud authentication
    apiKey := config.Credentials["api-key"]
    if apiKey == "" {
        return nil, fmt.Errorf("IBM Cloud API key is required")
    }

    // Set up IBM Cloud client
    authenticator := &core.IamAuthenticator{
        ApiKey: apiKey,
    }

    // Create IBM Cloud Kubernetes Service client
    iksClient, err := iks.NewIksClient(&iks.IksClientOptions{
        Authenticator: authenticator,
        URL:           "https://containers.cloud.ibm.com/global",
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create IBM Cloud client: %w", err)
    }

    // Create cloud provider interface
    ibmProvider := &IBMCloudProvider{
        client:     iksClient,
        region:     config.Region,
        clusterID:  config.Credentials["cluster-id"],
        resourceGroup: config.Credentials["resource-group"],
    }

    return ibmProvider, nil
}
```

### 3. Implement the IBM Cloud Provider Interface

Create a new file `pkg/testing/ibmcloud_provider.go`:

```go
package testing

import (
    "context"
    "fmt"

    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/types"
    cloudprovider "k8s.io/cloud-provider"
    "k8s.io/klog/v2"

    "github.com/IBM-Cloud/ibm-cloud-kubernetes-service/iks"
)

// IBMCloudProvider implements the cloudprovider.Interface for IBM Cloud
type IBMCloudProvider struct {
    client       *iks.IksClient
    region       string
    clusterID    string
    resourceGroup string
}

// Initialize initializes the IBM Cloud provider
func (p *IBMCloudProvider) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
    klog.Info("Initializing IBM Cloud provider")
}

// LoadBalancer returns the load balancer interface
func (p *IBMCloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
    return &IBMLoadBalancer{
        client:    p.client,
        region:    p.region,
        clusterID: p.clusterID,
    }, true
}

// Instances returns the instances interface
func (p *IBMCloudProvider) Instances() (cloudprovider.Instances, bool) {
    return &IBMInstances{
        client:    p.client,
        region:    p.region,
        clusterID: p.clusterID,
    }, true
}

// InstancesV2 returns the instances v2 interface
func (p *IBMCloudProvider) InstancesV2() (cloudprovider.InstancesV2, bool) {
    return nil, false
}

// Zones returns the zones interface
func (p *IBMCloudProvider) Zones() (cloudprovider.Zones, bool) {
    return &IBMZones{
        client:    p.client,
        region:    p.region,
        clusterID: p.clusterID,
    }, true
}

// Clusters returns the clusters interface
func (p *IBMCloudProvider) Clusters() (cloudprovider.Clusters, bool) {
    return &IBMClusters{
        client:    p.client,
        region:    p.region,
        clusterID: p.clusterID,
    }, true
}

// Routes returns the routes interface
func (p *IBMCloudProvider) Routes() (cloudprovider.Routes, bool) {
    return &IBMRoutes{
        client:    p.client,
        region:    p.region,
        clusterID: p.clusterID,
    }, true
}

// ProviderName returns the cloud provider name
func (p *IBMCloudProvider) ProviderName() string {
    return "ibmcloud"
}

// HasClusterID returns true if the cluster has a clusterID
func (p *IBMCloudProvider) HasClusterID() bool {
    return p.clusterID != ""
}

// IBMLoadBalancer implements cloudprovider.LoadBalancer for IBM Cloud
type IBMLoadBalancer struct {
    client    *iks.IksClient
    region    string
    clusterID string
}

// EnsureLoadBalancer creates or updates a load balancer
func (lb *IBMLoadBalancer) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
    // Implementation for creating/updating IBM Cloud load balancer
    klog.Infof("Ensuring load balancer for service %s in cluster %s", service.Name, clusterName)
    
    // TODO: Implement actual IBM Cloud load balancer creation
    // This would involve calling the IBM Cloud API to create/update a load balancer
    
    return &v1.LoadBalancerStatus{
        Ingress: []v1.LoadBalancerIngress{
            {
                IP: "192.168.1.100", // Example IP
            },
        },
    }, nil
}

// UpdateLoadBalancer updates an existing load balancer
func (lb *IBMLoadBalancer) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
    klog.Infof("Updating load balancer for service %s in cluster %s", service.Name, clusterName)
    
    // TODO: Implement actual IBM Cloud load balancer update
    return nil
}

// EnsureLoadBalancerDeleted deletes a load balancer
func (lb *IBMLoadBalancer) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
    klog.Infof("Deleting load balancer for service %s in cluster %s", service.Name, clusterName)
    
    // TODO: Implement actual IBM Cloud load balancer deletion
    return nil
}

// GetLoadBalancer returns the status of the specified load balancer
func (lb *IBMLoadBalancer) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (*v1.LoadBalancerStatus, bool, error) {
    klog.Infof("Getting load balancer for service %s in cluster %s", service.Name, clusterName)
    
    // TODO: Implement actual IBM Cloud load balancer status retrieval
    return nil, false, nil
}

// GetLoadBalancerName returns the name of the load balancer
func (lb *IBMLoadBalancer) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
    return fmt.Sprintf("ibm-lb-%s-%s", clusterName, service.Name)
}

// Similar implementations for IBMInstances, IBMZones, IBMClusters, and IBMRoutes...
```

### 4. Configure IBM Cloud Credentials

Create a credentials file `ibmcloud-credentials.json`:

```json
{
    "api-key": "your-ibm-cloud-api-key",
    "cluster-id": "your-iks-cluster-id",
    "resource-group": "your-resource-group",
    "region": "us-south"
}
```

### 5. Test the Implementation

```bash
# Test with IBM Cloud provider
./bin/e2e-test-runner \
  --provider ibmcloud \
  --kubeconfig ~/.kube/config \
  --credentials ibmcloud-credentials.json \
  --region us-south \
  --zone us-south-1 \
  --cluster my-iks-cluster \
  --suite loadbalancer \
  --verbose
```

## IBM Cloud Specific Considerations

### 1. Authentication
- IBM Cloud uses API keys for authentication
- Service IDs can be used for automated access
- IAM tokens are required for API access

### 2. Load Balancer Types
- IBM Cloud offers different load balancer types:
  - Application Load Balancer (ALB)
  - Network Load Balancer (NLB)
  - Classic Load Balancer (deprecated)

### 3. Regions and Zones
- IBM Cloud regions: `us-south`, `us-east`, `eu-gb`, `eu-de`, `jp-tok`, etc.
- Zones within regions: `us-south-1`, `us-south-2`, `us-south-3`

### 4. Resource Groups
- IBM Cloud uses resource groups for resource organization
- Resources must be created in a specific resource group

### 5. Cluster ID
- IBM Cloud Kubernetes Service clusters have unique IDs
- Cluster ID is required for API operations

## Example Implementation

Here's a more complete example of the load balancer implementation:

```go
func (lb *IBMLoadBalancer) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
    // Get service annotations for IBM Cloud specific configuration
    lbType := service.Annotations["service.kubernetes.io/ibm-load-balancer-cloud-provider-type"]
    if lbType == "" {
        lbType = "public" // Default to public load balancer
    }

    // Create load balancer request
    lbRequest := &iks.CreateLoadBalancerRequest{
        ClusterID: lb.clusterID,
        ServiceName: service.Name,
        ServiceNamespace: service.Namespace,
        Type: lbType,
        // Add other required fields
    }

    // Call IBM Cloud API
    lbResponse, err := lb.client.CreateLoadBalancer(ctx, lbRequest)
    if err != nil {
        return nil, fmt.Errorf("failed to create load balancer: %w", err)
    }

    // Return load balancer status
    return &v1.LoadBalancerStatus{
        Ingress: []v1.LoadBalancerIngress{
            {
                IP: lbResponse.IP,
                Hostname: lbResponse.Hostname,
            },
        },
    }, nil
}
```

## Testing

### 1. Unit Tests
Create unit tests for the IBM Cloud provider:

```go
func TestIBMCloudProvider_LoadBalancer(t *testing.T) {
    // Test load balancer operations
}

func TestIBMCloudProvider_Instances(t *testing.T) {
    // Test instance operations
}
```

### 2. Integration Tests
Test with a real IBM Cloud cluster:

```bash
# Set up test environment
ibmcloud login --apikey $IBM_CLOUD_API_KEY
ibmcloud target --resource-group my-test-rg
ibmcloud ks cluster config --cluster my-test-cluster

# Run tests
./bin/e2e-test-runner \
  --provider ibmcloud \
  --kubeconfig ~/.kube/config \
  --region us-south \
  --zone us-south-1 \
  --cluster my-test-cluster \
  --suite all \
  --verbose
```

## Troubleshooting

### Common Issues

1. **Authentication Failed**
   - Verify API key is valid and has proper permissions
   - Check if the API key has access to the resource group

2. **Cluster Not Found**
   - Verify cluster ID is correct
   - Ensure cluster is in the specified region

3. **Resource Group Access**
   - Verify API key has access to the resource group
   - Check resource group permissions

4. **Load Balancer Creation Failed**
   - Check load balancer quotas
   - Verify network configuration
   - Ensure proper IAM permissions

## Next Steps

1. **Complete Implementation**: Implement all cloud provider interfaces
2. **Add Tests**: Create comprehensive test suite
3. **Documentation**: Update documentation with IBM Cloud specifics
4. **CI/CD Integration**: Add IBM Cloud testing to CI/CD pipeline
5. **Performance Testing**: Add performance benchmarks
6. **Security Testing**: Add security-focused test cases

## Resources

- [IBM Cloud Kubernetes Service Documentation](https://cloud.ibm.com/docs/containers)
- [IBM Cloud API Documentation](https://cloud.ibm.com/apis)
- [IBM Cloud CLI Documentation](https://cloud.ibm.com/docs/cli)
- [Kubernetes Cloud Controller Manager](https://kubernetes.io/docs/concepts/architecture/cloud-controller/)
