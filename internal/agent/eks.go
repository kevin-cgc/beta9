package agent

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	vknode "github.com/virtual-kubelet/virtual-kubelet/node"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// initializeKubeClient sets up the Kubernetes client using the AWS SDK
func InitializeKubeClient(clusterName, region string) (*kubernetes.Clientset, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	// Create an EKS client
	eksClient := eks.NewFromConfig(cfg)

	// Retrieve the cluster information
	clusterInfo, err := eksClient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
		Name: &clusterName,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to describe cluster: %w", err)
	}

	// Decode the base64-encoded certificate authority data
	caData, err := base64.StdEncoding.DecodeString(*clusterInfo.Cluster.CertificateAuthority.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CA data: %w", err)
	}

	kubeConfig := api.Config{
		Clusters: map[string]*api.Cluster{
			"eks-cluster": {
				Server:                   *clusterInfo.Cluster.Endpoint,
				CertificateAuthorityData: caData,
			},
		},
		Contexts: map[string]*api.Context{
			"eks-context": {
				Cluster:  "eks-cluster",
				AuthInfo: "eks-auth",
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			"eks-auth": {
				Exec: &api.ExecConfig{
					APIVersion: "client.authentication.k8s.io/v1beta1",
					Command:    "aws",
					Args:       []string{"eks", "get-token", "--cluster-name", clusterName},
					Env:        nil,
				},
			},
		},
		CurrentContext: "eks-context",
	}

	// Generate a REST config from the kubeconfig
	clientConfig := clientcmd.NewNonInteractiveClientConfig(kubeConfig, "eks-context", &clientcmd.ConfigOverrides{}, nil)
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes REST config: %w", err)
	}

	// Create the Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes clientset: %w", err)
	}

	return clientset, nil
}

// setupProvider initializes your custom Virtual Kubelet provider.
// This example provides a stub implementation. Replace it with your provider's setup logic.
func SetupProvider(clientset *kubernetes.Clientset, nodeName string) (vknode.NodeProvider, error) {
	provider := &vknode.NaiveNodeProviderV2{}

	// Example: Use NotifyNodeStatus to respond to node status updates
	provider.NotifyNodeStatus(context.Background(), func(node *corev1.Node) {
		// Example: Update node status based on external conditions
	})

	return provider, nil
}
