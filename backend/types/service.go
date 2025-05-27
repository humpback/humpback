package types

const (
	ServiceStatusNotReady = "NotReady"
	ServiceStatusRunning  = "Running"
	ServiceStatusFailed   = "Failed"

	ContainerStatusPending  = "Pending"
	ContainerStatusStarting = "Starting"
	ContainerStatusCreated  = "Created"
	ContainerStatusRunning  = "Running"
	ContainerStatusFailed   = "Failed"
	ContainerStatusExited   = "Exited"
	ContainerStatusRemoved  = "Removed"
	ContainerStatusWarning  = "Warning"

	ServiceActionStart    = "Start"
	ServiceActionStop     = "Stop"
	ServiceActionRestart  = "Restart"
	ServiceActionEnable   = "Enable"
	ServiceActionDisable  = "Disable"
	ServiceActionDelete   = "Delete"
	ServiceActionDispatch = "Dispatch"

	ContainerLabelServiceId   = "Humpback-ServiceId"
	ContainerLabelServiceName = "Humpback-ServiceName"
	ContainerLabelGroupId     = "Humpback-GroupId"
)

type ServiceChangeInfo struct {
	ServiceId string
	Action    string
	Version   string
}

type Service struct {
	ServiceId   string             `json:"serviceId"`
	GroupId     string             `json:"groupId"`
	ServiceName string             `json:"serviceName"`
	Description string             `json:"description"`
	Version     string             `json:"version"`
	Action      string             `json:"action"`
	IsEnabled   bool               `json:"isEnabled"`
	IsDelete    bool               `json:"isDelete"`
	Status      string             `json:"status"`
	Memo        string             `json:"memo"`
	Meta        *ServiceMetaDocker `json:"meta"`
	Deployment  *Deployment        `json:"deployment"`
	Containers  []*ContainerStatus `json:"containers"`
	CreatedAt   int64              `json:"createdAt"`
	UpdatedAt   int64              `json:"updatedAt"`
}

type AgentTask struct {
	ContainerName string       `json:"containerName"`
	ServiceName   string       `json:"serviceName"`
	ServiceId     string       `json:"serviceId"`
	GroupId       string       `json:"groupId"`
	ManualExec    bool         `json:"manualExec"`
	RegistryAuth  RegistryAuth `json:"registryAuth"`
	*ServiceMetaDocker
	*ScheduleInfo
}

type RegistryAuth struct {
	ServerAddress    string `json:"serverAddress"`
	RegistryUsername string `json:"registryUsername"`
	RegistryPassword string `json:"registryPassword"`
}

type MountInfo struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type ContainerPort struct {
	BindIP      string `json:"bindIP"`
	PrivatePort int    `json:"privatePort"`
	PublicPort  int    `json:"publicPort"`
	Type        string `json:"type"`
}

type ContainerStatus struct {
	ContainerId   string            `json:"containerId"`
	ContainerName string            `json:"containerName"`
	NodeId        string            `json:"nodeId"`
	Ip            string            `json:"ip"`
	State         string            `json:"state"`
	Status        string            `json:"status"`
	ErrorMsg      string            `json:"errorMsg"`
	Image         string            `json:"image"`
	Command       string            `json:"command"`
	Network       string            `json:"network"`
	CreateAt      int64             `json:"created"`
	StartAt       int64             `json:"started"`
	NextAt        int64             `json:"nextAt"`
	LastHeartbeat int64             `json:"lastHeartbeat"`
	Labels        map[string]string `json:"labels"`
	Env           []string          `json:"env"`
	Mounts        []MountInfo       `json:"mounts"`
	Ports         []ContainerPort   `json:"ports"`
}

type Deployment struct {
	Type       DeployType       `json:"type"`
	Mode       DeployMode       `json:"mode"`
	Replicas   int              `json:"replicas"`
	Placements []*PlacementInfo `json:"placements"`
	Schedule   *ScheduleInfo    `json:"schedule"`
	ManualExec bool             `json:"manualExec"`
}

type DeployMode string
type DeployType string

var (
	DeployModeGlobal    DeployMode = "global"
	DeployModeReplicate DeployMode = "replicated"
)

var (
	DeployTypeSchedule   DeployType = "schedule"
	DeployTypeBackground DeployType = "background"
)

type PlacementMode string

var (
	PlacementModeLabel PlacementMode = "label"
	PlacementModeIP    PlacementMode = "ip"
)

type PlacementInfo struct {
	Mode    PlacementMode `json:"mode"`
	Key     string        `json:"key"`
	Value   string        `json:"value"`
	IsEqual bool          `json:"isEqual"`
}

// Rule is cron string
type ScheduleInfo struct {
	Timeout string   `json:"timeout"`
	Rules   []string `json:"rules"`
}

type ServiceMetaDocker struct {
	RegistryDomain string            `json:"registryDomain"`
	Image          string            `json:"image"`
	RegistryId     string            `json:"registryId"`
	AlwaysPull     bool              `json:"alwaysPull"`
	Command        string            `json:"command"`
	EnvConfig      []string          `json:"envConfig"`
	Envs           []string          `json:"env"`
	Labels         map[string]string `json:"labels"`
	Privileged     bool              `json:"privileged"`
	Capabilities   *Capabilities     `json:"capabilities"`
	LogConfig      *ServiceLogConfig `json:"logConfig"`
	Resources      *ServiceResources `json:"resources"`
	Volumes        []*ServiceVolume  `json:"volumes"`
	Network        *NetworkInfo      `json:"network"`
	RestartPolicy  *RestartPolicy    `json:"restartPolicy"`
}

type Capabilities struct {
	CapAdd  []string `json:"capAdd"`
	CapDrop []string `json:"capDrop"`
}

type ServiceLogConfig struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

type ServiceResources struct {
	Memory            uint64 `json:"memory"`
	MemoryReservation uint64 `json:"memoryReservation"`
	MaxCpuUsage       uint64 `json:"maxCpuUsage"`
}

type ServiceVolumeType string

var (
	ServiceVolumeTypeBind   ServiceVolumeType = "bind"
	ServiceVolumeTypeVolume ServiceVolumeType = "volume"
)

type ServiceVolume struct {
	Type     ServiceVolumeType `json:"type"`
	Target   string            `json:"target"`
	Source   string            `json:"source"`
	Readonly bool              `json:"readonly"`
}

type NetworkMode string

var (
	NetworkModeHost   NetworkMode = "host"
	NetworkModeBridge NetworkMode = "bridge"
	NetworkModeCustom NetworkMode = "custom"
)

type NetworkInfo struct {
	Mode               NetworkMode `json:"mode"`        // custom模式需要创建网络
	Hostname           string      `json:"hostname"`    // bridge及custom模式时可设置，用户容器的hostname
	NetworkName        string      `json:"networkName"` //custom模式使用
	UseMachineHostname bool        `json:"useMachineHostname"`
	Ports              []*PortInfo `json:"ports"`
}

type PortInfo struct {
	HostPort      uint64 `json:"hostPort"`
	ContainerPort uint64 `json:"containerPort"`
	Protocol      string `json:"protocol"`
}

type RestartPolicyMode string

var (
	RestartPolicyModeNo            RestartPolicyMode = "no"
	RestartPolicyModeAlways        RestartPolicyMode = "always"
	RestartPolicyModeOnFail        RestartPolicyMode = "on-failure"
	RestartPolicyModeUnlessStopped RestartPolicyMode = "unless-stopped"
)

type RestartPolicy struct {
	Mode          RestartPolicyMode `json:"mode"`
	MaxRetryCount uint64            `json:"maxRetryCount"`
}
