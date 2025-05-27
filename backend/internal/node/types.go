package node

import (
	"humpback/types"
)

type ContainerStats struct {
	CpuPercent     float64                         `json:"cpuPercent"`
	MermoryUsed    uint64                          `json:"memoryUsageBytes"`
	MemoryLimit    uint64                          `json:"memoryLimitBytes"`
	DiskReadBytes  uint64                          `json:"diskReadBytes"`
	DiskWriteBytes uint64                          `json:"diskWriteBytes"`
	StatsTime      string                          `json:"statsTime"`
	Networks       []*types.ContainerNetWorkStatus `json:"networks"`
}
