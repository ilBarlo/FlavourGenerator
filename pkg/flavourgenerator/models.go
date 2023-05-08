package flavourgenerator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Small  Plan = "Small: 11"
	Medium Plan = "Medium: 33"
	Large  Plan = "Large: 66"
)

// NodeInfo represents a node and its resources
type NodeInfo struct {
	UID             string          `json:"uid"`
	Name            string          `json:"name"`
	Architecture    string          `json:"architecture"`
	OperatingSystem string          `json:"os"`
	ResourceMetrics ResourceMetrics `json:"resources"`
}

// ResourceMetrics represents resources of a certain node
type ResourceMetrics struct {
	CPUTotal        float64 `json:"totalCPU"`
	CPUAvailable    float64 `json:"availableCPU"`
	MemoryTotal     float64 `json:"totalMemory"`
	MemoryAvailable float64 `json:"availableMemory"`
}

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Flavour represents a subset of a node's resources
type Flavour struct {
	UID             string `json:"uid"`
	Name            string `json:"name"`
	Architecture    string `json:"architecture"`
	OperatingSystem string `json:"os"`
	CPUOffer        string `json:"cpuOffer"`
	MemoryOffer     string `json:"memoryOffer"`
	PodsOffer       []Plan `json:"podsOffer"`
}

// Plan represents a specific Plan for each flavour depending the number of pods
type Plan string

// splitResources produces different Flavours (60% - 30% - 10% of available resources)
func splitResources(node NodeInfo) []Flavour {

	AvailCPU := node.ResourceMetrics.CPUAvailable
	AvailMemory := node.ResourceMetrics.MemoryAvailable

	flavours := []Flavour{
		{
			UID:             node.UID + "flavour-1",
			Name:            node.Name + "flavour-small",
			Architecture:    node.Architecture,
			OperatingSystem: node.OperatingSystem,
			CPUOffer:        fmt.Sprintf("%.0f", float64(AvailCPU)*0.6),
			MemoryOffer:     fmt.Sprintf("%.0fGi", float64(AvailMemory)*0.6),
			PodsOffer:       []Plan{Small, Medium, Large},
		},
		{
			UID:             node.UID + "flavour-2",
			Name:            node.Name + "flavour-medium",
			Architecture:    node.Architecture,
			OperatingSystem: node.OperatingSystem,
			CPUOffer:        fmt.Sprintf("%.0f", float64(AvailCPU)*0.3),
			MemoryOffer:     fmt.Sprintf("%.0fGi", float64(AvailMemory)*0.3),
			PodsOffer:       []Plan{Small, Medium, Large},
		},
		{
			UID:             node.UID + "flavour-3",
			Name:            node.Name + "flavour-large",
			Architecture:    node.Architecture,
			OperatingSystem: node.OperatingSystem,
			CPUOffer:        fmt.Sprintf("%.0f", float64(AvailCPU)*0.1),
			MemoryOffer:     fmt.Sprintf("%.0fGi", float64(AvailMemory)*0.1),
			PodsOffer:       []Plan{Small, Medium, Large},
		},
	}
	return flavours
}

// fromNodeInfo creates from params a new NodeInfo Struct
func fromNodeInfo(uid, name, arch, os string, metrics ResourceMetrics) *NodeInfo {
	return &NodeInfo{
		UID:             uid,
		Name:            name,
		Architecture:    arch,
		OperatingSystem: os,
		ResourceMetrics: metrics,
	}
}

// fromResourceMetrics creates from params a new ResourceMetrics Struct
func fromResourceMetrics(cpuTotal int64, cpuUsed int64, memoryTotal int64, memoryUsed int64) *ResourceMetrics {
	return &ResourceMetrics{
		CPUTotal:        float64(cpuTotal) / 1000,
		CPUAvailable:    float64(cpuTotal-cpuUsed) / 1000,
		MemoryTotal:     float64(memoryTotal) / (1024 * 1024 * 1024),
		MemoryAvailable: float64(memoryTotal-memoryUsed) / (1024 * 1024 * 1024),
	}
}
