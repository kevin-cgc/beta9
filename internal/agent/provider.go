package agent

import (
	"context"
	"io"
	"log"
	"time"

	dto "github.com/prometheus/client_model/go"
	node "github.com/virtual-kubelet/virtual-kubelet/node"
	nodeapi "github.com/virtual-kubelet/virtual-kubelet/node/api"
	"github.com/virtual-kubelet/virtual-kubelet/node/api/statsv1alpha1"
	nodeutil "github.com/virtual-kubelet/virtual-kubelet/node/nodeutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RemoteProvider struct {
	clients  *ProviderClients
	nodeName string
}

type RemoteNodeProvider struct {
	nodeName string
}

func (p *RemoteNodeProvider) Ping(context.Context) error {
	// TODO: add connectivity health check (maybe to tailscale)
	return nil
}

func (p *RemoteNodeProvider) NotifyNodeStatus(ctx context.Context, cb func(*corev1.Node)) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {

		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				nodeStatus := &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"type":                   "virtual-kubelet",
							"kubernetes.io/role":     "agent",
							"kubernetes.io/hostname": p.nodeName,
						},
					},
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{
							{
								Type:               corev1.NodeReady,
								Status:             corev1.ConditionTrue, // Or ConditionFalse based on actual health check
								LastHeartbeatTime:  metav1.Now(),
								LastTransitionTime: metav1.Now(),
								Reason:             "KubeletReady",
								Message:            "kubelet is posting ready status",
							},
						},
					},
				}

				// Call the callback with the updated node status
				log.Println("updating node status")
				cb(nodeStatus)
			}
		}
	}()
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

func (p *RemoteProvider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	return p.clients.LocalClient.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (p *RemoteProvider) GetPods(ctx context.Context) ([]*corev1.Pod, error) {
	podList, err := p.clients.LocalClient.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + p.nodeName,
	})
	if err != nil {
		return nil, err
	}
	pods := make([]*corev1.Pod, len(podList.Items))
	for i, pod := range podList.Items {
		pods[i] = &pod
	}
	return pods, nil
}

func (p *RemoteProvider) GetPodStatus(ctx context.Context, namespace, name string) (*corev1.PodStatus, error) {
	pod, err := p.GetPod(ctx, namespace, name)
	if err != nil {
		return nil, err
	}
	return &pod.Status, nil
}

func (p *RemoteProvider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
	_, err := p.clients.LocalClient.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	return err
}

func (p *RemoteProvider) UpdatePod(ctx context.Context, pod *corev1.Pod) error {
	_, err := p.clients.LocalClient.CoreV1().Pods(pod.Namespace).Update(ctx, pod, metav1.UpdateOptions{})
	return err
}

func (p *RemoteProvider) DeletePod(ctx context.Context, pod *corev1.Pod) error {
	return p.clients.LocalClient.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
}

func newNodeProvider(nodeName string) node.NodeProvider {
	return &RemoteNodeProvider{
		nodeName: nodeName,
	}
}

func NewProvider(clients *ProviderClients, nodeName string) (nodeutil.Provider, node.NodeProvider, error) {
	provider := &RemoteProvider{
		clients:  clients,
		nodeName: nodeName,
	}
	nodeProvider := newNodeProvider(nodeName)
	return provider, nodeProvider, nil
}
