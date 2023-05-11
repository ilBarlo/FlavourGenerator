package flavourgenerator

import (
	"fmt"
	"strconv"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	NodeUID         string     `json:"nodeUid"`
	UID             string     `json:"uid"`
	Name            string     `json:"name"`
	Architecture    string     `json:"architecture"`
	OperatingSystem string     `json:"os"`
	CPUOffer        string     `json:"cpuOffer"`
	MemoryOffer     string     `json:"memoryOffer"`
	PodsOffer       []PodsPlan `json:"podsOffer"`
}

// PodsPlan represents a plan for which is possibile to have a specific amount of available pods
type PodsPlan struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
	Pods      int64  `json:"availablePods"`
}

// splitResources produces different Flavours with intelligent resource allocation
func splitResources(node NodeInfo) []Flavour {
	AvailCPU := node.ResourceMetrics.CPUAvailable
	AvailMemory := node.ResourceMetrics.MemoryAvailable

	// Define initial values for resource allocation
	cpuAllocation := 1.0
	memoryAllocation := 1.0

	flavours := []Flavour{}
	for AvailCPU > cpuAllocation && AvailMemory > memoryAllocation {

		// Create the flavour
		flavour := Flavour{
			NodeUID:         node.UID,
			UID:             node.UID + "-flavour-" + strconv.Itoa(len(flavours)+1),
			Name:            node.Name + "-flavour-" + strconv.Itoa(len(flavours)+1),
			Architecture:    node.Architecture,
			OperatingSystem: node.OperatingSystem,
			CPUOffer:        fmt.Sprintf("%.2f", cpuAllocation),
			MemoryOffer:     fmt.Sprintf("%.2fGi", memoryAllocation),
			PodsOffer: []PodsPlan{
				{Name: "Small", Available: true, Pods: 11},
				{Name: "Medium", Available: true, Pods: 33},
				{Name: "Large", Available: true, Pods: 66},
			},
		}
		flavours = append(flavours, flavour)

		// Increase the resource allocation for the next flavour
		cpuAllocation += float64(len(flavours) + 1)
		memoryAllocation += float64(len(flavours) + 1)

		// Check if the allocation exceeds the available resources
		if cpuAllocation > AvailCPU {
			cpuAllocation = AvailCPU
		}
		if memoryAllocation > AvailMemory {
			memoryAllocation = AvailMemory
		}
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
