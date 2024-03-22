package main

import (
	"context"
	"log"

	"github.com/beam-cloud/beta9/internal/agent"
	"github.com/virtual-kubelet/virtual-kubelet/node"
	nodeutil "github.com/virtual-kubelet/virtual-kubelet/node/nodeutil"
)

func main() {
	nodeName := "test-node"
	clusterName := "eks-stage-01"
	region := "us-east-1"

	clientset, err := agent.InitializeKubeClient(clusterName, region)
	if err != nil {
		log.Fatalf("error initializing Kubernetes client: %v\n", err)
	}

	// // Prepare node object for Virtual Kubelet
	// taints := []corev1.Taint{
	// 	{
	// 		Key:    "virtual-kubelet.io/provider",
	// 		Value:  "aws-eks",
	// 		Effect: corev1.TaintEffectNoSchedule,
	// 	},
	// }

	newProvider := func(cfg nodeutil.ProviderConfig) (nodeutil.Provider, node.NodeProvider, error) {
		provider, nodeProvider, err := agent.NewProvider(clientset, nodeName)
		if err != nil {
			return nil, nil, err
		}

		return provider, nodeProvider, nil
	}

	node, err := nodeutil.NewNode(nodeName, newProvider)
	if err != nil {
		log.Fatalf("error setting up provider: %v\n", err)
	}

	if err := node.Run(context.Background()); err != nil {
		log.Fatalf("error running node: %v\n", err)
	}
}
