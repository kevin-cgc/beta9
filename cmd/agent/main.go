package main

import (
	"context"
	"fmt"
	"os"

	"github.com/beam-cloud/beta9/internal/agent"

	vk "github.com/virtual-kubelet/virtual-kubelet"
	node "github.com/virtual-kubelet/virtual-kubelet/node"
	nodeutil "github.com/virtual-kubelet/virtual-kubelet/node/nodeutil"
)

func NewProvider() (vk.Provider, node.NodeProvider, error) {

}

func main() {
	nodeName := "test-node"
	clusterName := "eks-stage-01"
	region := "us-east-1"

	// Initialize Kubernetes client for EKS
	clientset, err := agent.InitializeKubeClient(clusterName, region)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializing Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// Setup Virtual Kubelet provider
	provider, err := agent.SetupProvider(clientset, nodeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error setting up provider: %v\n", err)
		os.Exit(1)
	}

	// // Prepare node object for Virtual Kubelet
	// taints := []corev1.Taint{
	// 	{
	// 		Key:    "virtual-kubelet.io/provider",
	// 		Value:  "aws-eks",
	// 		Effect: corev1.TaintEffectNoSchedule,
	// 	},
	// }

	node, err := nodeutil.NewNode(nodeName, provider)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error setting up provider: %v\n", err)
		os.Exit(1)
	}

	// Run the node
	if err := node.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "error running node: %v\n", err)
		os.Exit(1)
	}
}
