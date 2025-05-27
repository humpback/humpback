package types

type Node struct {
	NodeId      string            `json:"nodeId"`
	Name        string            `json:"name"`
	IpAddress   string            `json:"ipAddress"`
	Port        int               `json:"port"`
	Status      string            `json:"status"`
	IsEnable    bool              `json:"isEnable"`
	CreatedAt   int64             `json:"createdAt"`
	UpdatedAt   int64             `json:"updatedAt"`
	CPUUsage    float32           `json:"cpuUsage"`
	CPU         uint              `json:"cpu"`
	MemoryUsage float32           `json:"memoryUsage"`
	MemoryTotal uint64            `json:"memoryTotal"`
	MemoryUsed  uint64            `json:"memoryUsed"`
	Labels      map[string]string `json:"labels"`
}

type NodesGroups struct {
	GroupId     string   `json:"groupId"`
	GroupName   string   `json:"groupName"`
	Description string   `json:"description"`
	CreatedAt   int64    `json:"createdAt"`
	UpdatedAt   int64    `json:"updatedAt"`
	Users       []string `json:"users"`
	Teams       []string `json:"teams"`
	Nodes       []string `json:"nodes"`
}

type NodeSimpleInfo struct {
	NodeId          string
	Name            string
	IpAddress       string
	Port            int
	Status          string
	LastHeartbeat   int64
	OnlineThreshold int
	CPUUsage        float32
	MemoryUsage     float32
	TotalCPU        uint
	TotalMemory     uint64
	UsedMemory      uint64
}
