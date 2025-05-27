package node

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"humpback/internal/db"
	"humpback/pkg/httpx"
	"humpback/pkg/utils"
	"humpback/types"
)

var (
	ErrNodeNotExist = errors.New("The node does not exist")
)

func RemoveNodeContainer(nodeId string, containerId string, containerName string) error {
	// remove container
	node := GetNodeInfo(nodeId)
	if node != nil {
		url := fmt.Sprintf("http://%s:%d/api/v1/container/%s?force=true&containerName=%s", node.IpAddress, node.Port, containerId, containerName)
		slog.Info("[Agent Helper] Remove container", "url", url)
		err := httpx.NewHttpXClient().Delete(url, nil, nil, nil)
		if err != nil {
			slog.Error("[Agent Helper] Remove container error", "error", err.Error())
			return err
		}
	}
	return ErrNodeNotExist
}

func OperateNodeContainer(nodeId string, containerId string, action string) error {
	// operate container
	node := GetNodeInfo(nodeId)
	if node == nil {
		return ErrNodeNotExist
	}
	url := fmt.Sprintf("http://%s:%d/api/v1/container/%s/%s", node.IpAddress, node.Port, containerId, strings.ToLower(action))
	slog.Info("[Agent Helper] Operate container", "url", url)
	return httpx.NewHttpXClient().Post(url, nil, nil, nil, nil)
}

func StartNewContainer(nodeId, containerName string, svc *types.Service) error {
	node := GetNodeInfo(nodeId)
	if node != nil {

		task := &types.AgentTask{
			ContainerName:     containerName,
			ServiceName:       svc.ServiceName,
			ServiceId:         svc.ServiceId,
			GroupId:           svc.GroupId,
			ManualExec:        svc.Deployment.ManualExec,
			ServiceMetaDocker: svc.Meta,
			ScheduleInfo:      svc.Deployment.Schedule,
		}

		if svc.Meta.RegistryId != "" {
			registry, _ := db.RegistryGetById(svc.Meta.RegistryId)
			if registry != nil {
				task.RegistryAuth = types.RegistryAuth{
					ServerAddress:    registry.URL,
					RegistryUsername: registry.Username,
					RegistryPassword: registry.Password,
				}
			}
		}

		utils.PrintJson(task)
		url := fmt.Sprintf("http://%s:%d/api/v1/container", node.IpAddress, node.Port)
		slog.Info("[Agent Helper] Create container", "url", url)
		err := httpx.NewHttpXClient().Post(url, nil, nil, task, nil)
		if err != nil {
			slog.Error("[Agent Helper] Start container error", "error", err.Error())
			return err
		}

		c := &types.ContainerStatus{
			ContainerName: task.ContainerName,
			State:         types.ContainerStatusPending,
			NodeId:        nodeId,
			Ip:            node.IpAddress,
		}
		svc.Containers = append(svc.Containers, c)
	}

	return nil
}

func QueryContainerLogs(nodeId string, containerId string, querys map[string]string) ([]string, error) {
	node := GetNodeInfo(nodeId)
	if node == nil {
		return nil, ErrNodeNotExist
	}
	url := fmt.Sprintf("http://%s:%d/api/v1/container/%s/logs", node.IpAddress, node.Port, containerId)
	data := make([]string, 0)
	err := httpx.NewHttpXClient().Get(url, querys, nil, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetContainerStats(nodeId string, containerId string) (*ContainerStats, error) {
	node := GetNodeInfo(nodeId)
	if node == nil {
		return nil, ErrNodeNotExist
	}
	var stats = new(ContainerStats)
	url := fmt.Sprintf("http://%s:%d/api/v1/container/%s/stats", node.IpAddress, node.Port, containerId)
	err := httpx.NewHttpXClient().Get(url, nil, nil, stats)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
