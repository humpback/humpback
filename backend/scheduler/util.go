package scheduler

import (
	"regexp"
	"strings"

	"humpback/pkg/utils"
	"humpback/types"
)

const (
	InstanceNamePrefix = "Humpback"
)

func parseVersionByContainerId(containerId string) string {
	parts := strings.Split(containerId, "-")
	if len(parts) < 4 {
		return ""
	}
	return parts[2]
}

func ParseServiceIdByContainerId(containerId string) string {
	parts := strings.Split(containerId, "-")
	if len(parts) < 4 {
		return ""
	}
	return parts[1]
}

func GenerateContainerName(serviceId, version string) string {
	return InstanceNamePrefix + "-" + serviceId + "-" + version + "-" + utils.NewVersionId()
}

func isContainerRunning(status string) bool {
	return status == types.ContainerStatusRunning
}

func isContainerExited(status string) bool {
	return status == types.ContainerStatusExited ||
		status == types.ContainerStatusCreated
}

func isContainerRemoved(status string) bool {
	return status == types.ContainerStatusRemoved
}

func isContainerStarting(status string) bool {
	return status == types.ContainerStatusPending ||
		status == types.ContainerStatusStarting
}

func isContainerWarning(status string) bool {
	return status == types.ContainerStatusWarning
}

func isContainerFailed(status string) bool {
	return status == types.ContainerStatusFailed
}

func isPlacementMatched(node *types.Node, p *types.PlacementInfo) bool {
	if p.Mode == types.PlacementModeIP {
		if p.IsEqual {
			return node.IpAddress == p.Value
		} else {
			return node.IpAddress != p.Value
		}
	} else {
		label := p.Key
		if l, ok := node.Labels[label]; ok {
			if p.IsEqual {
				return l == p.Value
			} else {
				return l != p.Value
			}
		} else {
			return false
		}
	}
}

func hasConfigValue(value string) string {
	re := regexp.MustCompile(`^\{([^}]+)\}$`)

	// 查找匹配项
	match := re.FindStringSubmatch(value)

	// 判断是否匹配
	if len(match) == 2 {
		return match[1]
	} else {
		return ""
	}
}
