package main

import (
	"context"
	"log"
	"runtime"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/beam-cloud/beta9/internal/agent"
	"github.com/virtual-kubelet/virtual-kubelet/node"
	nodeutil "github.com/virtual-kubelet/virtual-kubelet/node/nodeutil"
	v1 "k8s.io/api/core/v1"
)

func main() {
	nodeName := "test-node5"
	clusterName := "eks-stage-01"
	region := "us-east-1"

	clients, err := agent.InitializeClients(clusterName, region)
	if err != nil {
		log.Fatalf("error initializing Kubernetes client: %v\n", err)
	}

	newProvider := func(cfg nodeutil.ProviderConfig) (nodeutil.Provider, node.NodeProvider, error) {
		provider, nodeProvider, err := agent.NewProvider(clients, nodeName)
		if err != nil {
			return nil, nil, err
		}

		return provider, nodeProvider, nil
	}

	c := nodeutil.NodeConfig{
		NumWorkers:           runtime.NumCPU(),
		InformerResyncPeriod: time.Minute,
		HTTPListenAddr:       ":10250",
		Client:               clients.LocalClient,
		NodeSpec: v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: nodeName,
				Labels: map[string]string{
					"type":                   "virtual-kubelet",
					"kubernetes.io/role":     "agent",
					"kubernetes.io/hostname": nodeName,
				},
			},
			Status: v1.NodeStatus{
				Phase: v1.NodePending,
				Conditions: []v1.NodeCondition{
					{Type: v1.NodeReady},
					{Type: v1.NodeDiskPressure},
					{Type: v1.NodeMemoryPressure},
					{Type: v1.NodePIDPressure},
					{Type: v1.NodeNetworkUnavailable},
				},
			},
		}}

	node, err := nodeutil.NewNode(nodeName, newProvider, nodeutil.WithNodeConfig(c))
	if err != nil {
		log.Fatalf("error setting up provider: %v\n", err)
	}

	if err := node.Run(context.Background()); err != nil {
		log.Fatalf("error running node: %v\n", err)
	}
}
