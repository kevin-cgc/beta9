package agent

import (
	"context"
	"io"

	dto "github.com/prometheus/client_model/go"
	node "github.com/virtual-kubelet/virtual-kubelet/node"
	nodeapi "github.com/virtual-kubelet/virtual-kubelet/node/api"
	"github.com/virtual-kubelet/virtual-kubelet/node/api/statsv1alpha1"
	nodeutil "github.com/virtual-kubelet/virtual-kubelet/node/nodeutil"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes"
)

type RemoteProvider struct {
}

type RemoteNodeProvider struct {
}

func (p *RemoteNodeProvider) Ping(context.Context) error {
	return nil
}

func (p *RemoteNodeProvider) NotifyNodeStatus(ctx context.Context, cb func(*corev1.Node)) {
}

// Provider controller methods
func (p *RemoteProvider) AttachToContainer(context.Context, string, string, string, nodeapi.AttachIO) error {
	return nil
}

func (p *RemoteProvider) GetContainerLogs(context.Context, string, string, string, nodeapi.ContainerLogOpts) (io.ReadCloser, error) {
	return nil, nil
}

func (p *RemoteProvider) GetMetricsResource(context.Context) ([]*dto.MetricFamily, error) {
	return nil, nil
}

func (p *RemoteProvider) GetStatsSummary(context.Context) (*statsv1alpha1.Summary, error) {
	return nil, nil
}

func (p *RemoteProvider) RunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string, attach nodeapi.AttachIO) error {
	return nil
}

func (p *RemoteProvider) PortForward(ctx context.Context, namespace, pod string, port int32, stream io.ReadWriteCloser) error {
	return nil
}

// Pod controller

func (p *RemoteProvider) GetPod(context.Context, string, string) (*corev1.Pod, error) {
	return nil, nil
}

func (p *RemoteProvider) GetPods(ctx context.Context) ([]*corev1.Pod, error) {
	return nil, nil
}

func (p *RemoteProvider) GetPodStatus(context.Context, string, string) (*corev1.PodStatus, error) {
	return nil, nil
}

func (p *RemoteProvider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	return nil
}

func (p *RemoteProvider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	return nil
}

func (p *RemoteProvider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	return nil
}

func newNodeProvider() node.NodeProvider {
	return &RemoteNodeProvider{}
}

func NewProvider(clientset *kubernetes.Clientset, nodeName string) (nodeutil.Provider, node.NodeProvider, error) {
	provider := &RemoteProvider{}
	nodeProvider := newNodeProvider()
	return provider, nodeProvider, nil
}
