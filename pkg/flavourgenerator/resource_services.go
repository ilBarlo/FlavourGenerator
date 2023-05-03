package flavourgenerator

import (
	"context"

	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	sharingv1alpha1 "github.com/liqotech/liqo/apis/sharing/v1alpha1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = discoveryv1alpha1.AddToScheme(scheme)
	_ = sharingv1alpha1.AddToScheme(scheme)
	_ = metricsv1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
}

// GetKClient creates a kubernetes API client and returns it.
func GetKClient(ctx context.Context) (client.Client, error) {
	config := ctrl.GetConfigOrDie()

	cl, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		klog.Fatalf("error creating manager: %", err)
	}

	return cl, nil
}

func GetNodesResources(ctx context.Context, cl client.Client) (*[]NodeInfo, error) {
	// Get a list of nodes
	nodes := &corev1.NodeList{}
	err := cl.List(ctx, nodes)
	if err != nil {
		return nil, err
	}

	// Get a list of nodes metrics
	nodesMetrics := &metricsv1beta1.NodeMetricsList{}
	err = cl.List(ctx, nodesMetrics)
	if err != nil {
		return nil, err
	}

	var nodesInfo []NodeInfo
	// Print the name of each node
	for _, node := range nodes.Items {
		for _, metrics := range nodesMetrics.Items {
			if node.Name != metrics.Name {
				// So that we can select just the nodes that we want
				continue
			}

			// Select just the workers (not the control-plane nodes)
			if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
				continue
			}

			metricsStruct := getNodeResourceMetrics(&metrics, &node)
			nodeInfo := getNodeInfo(&node, metricsStruct)
			nodesInfo = append(nodesInfo, *nodeInfo)
		}
	}

	return &nodesInfo, nil
}

func getNodeResourceMetrics(nodeMetrics *metricsv1beta1.NodeMetrics, node *corev1.Node) *ResourceMetrics {
	cpuTotal := node.Status.Allocatable.Cpu().MilliValue()
	cpuUsed := nodeMetrics.Usage.Cpu().MilliValue()
	memoryTotal := node.Status.Allocatable.Memory().Value()
	memoryUsed := nodeMetrics.Usage.Memory().Value()
	return fromResourceMetrics(cpuTotal, cpuUsed, memoryTotal, memoryUsed)
}

func getNodeInfo(node *corev1.Node, metrics *ResourceMetrics) *NodeInfo {
	return fromNodeInfo(string(node.UID), node.Name, node.Status.NodeInfo.Architecture, node.Status.NodeInfo.OperatingSystem, *metrics)
}
