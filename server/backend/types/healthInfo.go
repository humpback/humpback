package types

// HostInfo 主机基础信息
type HostInfo struct {
	Hostname      string   `json:"hostname"`
	IpAddress     []string `json:"hostIPs"`
	Port          int      `json:"hostPort"`
	OsInformation string   `json:"osInformation"`
	KernelVersion string   `json:"kernelVersion"`
	TotalCPU      uint     `json:"totalCPU"`
	UsedCPU       float32  `json:"usedCPU"`
	CPUUsage      float32  `json:"cpuUsage"`
	TotalMemory   uint64   `json:"totalMemory"`
	UsedMemory    uint64   `json:"usedMemory"`
	MemoryUsage   float32  `json:"memoryUsage"`
}

// DockerEngine Docker 基础信息
type DockerEngine struct {
	Version        string   `json:"version"`
	ApiVersion     string   `json:"apiVersion"`
	RootDirectory  string   `json:"rootDirectory"`
	StorageDriver  string   `json:"storageDriver"`
	LoggingDriver  string   `json:"loggingDriver"`
	VolumePlugins  []string `json:"volumePlugins"`
	NetworkPlugins []string `json:"networkPlugins"`
}

// HealthInfo 整体信息，包含主机信息、Docker 引擎信息和容器列表
type HealthInfo struct {
	NodeId        string
	IpAddress     string
	HostInfo      HostInfo          `json:"hostInfo"`
	DockerEngine  DockerEngine      `json:"dockerEngine"`
	ContainerList []ContainerStatus `json:"containers"`
}
