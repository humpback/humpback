package types

var (
	MemoNoAvailableNode       = "NoAvailableNode"
	MemoCreateContainerFailed = "CreateContainerFailed"
)

type SearchGroupServiceSimpleInfo struct {
	GroupName   string `json:"groupName"`
	GroupId     string `json:"groupId"`
	ServiceName string `json:"serviceName,omitempty"`
	ServiceId   string `json:"serviceId,omitempty"`
}
