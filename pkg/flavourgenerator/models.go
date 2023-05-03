package flavourgenerator

import "fmt"

// Info of the node
type NodeInfo struct {
	uid             string          /* UID - Unique Identifier: used in case we have two nodes from different clusters with the same name */
	Name            string          `json:"name"`
	Architecture    string          `json:"architecture"`
	OperatingSystem string          `json:"os"`
	ResourceMetrics ResourceMetrics `json:"resources"`
}

// Metrics of the resource of a certain node
type ResourceMetrics struct {
	CPUTotal        string `json:"totalCPU"`
	CPUAvailable    string `json:"availableCPU"`
	MemoryTotal     string `json:"totalMemory"`
	MemoryAvailable string `json:"availableMemory"`
}

func fromNodeInfo(uid, name, arch, os string, metrics ResourceMetrics) *NodeInfo {
	return &NodeInfo{
		uid:             uid,
		Name:            name,
		Architecture:    arch,
		OperatingSystem: os,
		ResourceMetrics: metrics,
	}
}

func fromResourceMetrics(cpuTotal int64, cpuUsed int64, memoryTotal int64, memoryUsed int64) *ResourceMetrics {
	return &ResourceMetrics{
		CPUTotal:        fmt.Sprintf("%.2f", float64(cpuTotal)/1000),
		CPUAvailable:    fmt.Sprintf("%.2f", float64(cpuTotal-cpuUsed)/1000),
		MemoryTotal:     fmt.Sprintf("%.2fGi", float64(memoryTotal)/(1024*1024*1024)),
		MemoryAvailable: fmt.Sprintf("%.2fGi", float64(memoryTotal-memoryUsed)/(1024*1024*1024)),
	}
}
