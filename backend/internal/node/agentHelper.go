package node

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"humpback/internal/db"
	"humpback/pkg/httpx"
	"humpback/pkg/utils"
	"humpback/security"
	"humpback/types"
)

var (
	ErrNodeNotExist = errors.New("the node does not exist")
)

type AgentManager struct {
	NodeAgents map[string]httpx.HttpXClient
	sync.RWMutex
}

var agentManager *AgentManager

func NewAgentManager() {
	agentManager = &AgentManager{
		NodeAgents: make(map[string]httpx.HttpXClient),
	}
}

func getAgent(nodeId string) httpx.HttpXClient {
	var agent httpx.HttpXClient
	var ok bool
	agentManager.RLock()
	agent, ok = agentManager.NodeAgents[nodeId]
	agentManager.RUnlock()
	if ok {
		return agent
	} else {
		ip := utils.HostIP()
		serverBundle, err := security.CreateCertificateBundle("humpback-server", "", ip)
		if err != nil {
			return nil
		}
		agent = httpx.NewHttpXClient(serverBundle.CreateTLSConfig(false))
		agentManager.Lock()
		agentManager.NodeAgents[nodeId] = agent
		agentManager.Unlock()
		return agent
	}
}

func RemoveNodeContainer(nodeId string, containerId string, containerName string) error {
	// remove container
	node := GetNodeInfo(nodeId)
	if node != nil {
		url := fmt.Sprintf("https://%s:%d/api/v1/container/%s?force=true&containerName=%s", node.IpAddress, node.Port, containerId, containerName)
		slog.Info("[Agent Helper] Remove container", "url", url)
		agent := getAgent(nodeId)
		if agent == nil {
			return ErrNodeNotExist
		}
		err := agent.Delete(url, nil, nil, nil, node.RegisterInfo.AccessToken)
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
	url := fmt.Sprintf("https://%s:%d/api/v1/container/%s/%s", node.IpAddress, node.Port, containerId, strings.ToLower(action))
	slog.Info("[Agent Helper] Operate container", "url", url)
	agent := getAgent(nodeId)
	if agent == nil {
		return ErrNodeNotExist
	}
	return agent.Post(url, nil, nil, nil, nil, node.RegisterInfo.AccessToken)
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
		url := fmt.Sprintf("https://%s:%d/api/v1/container", node.IpAddress, node.Port)
		slog.Info("[Agent Helper] Create container", "url", url)
		agent := getAgent(nodeId)
		if agent == nil {
			return ErrNodeNotExist
		}
		err := agent.Post(url, nil, nil, task, nil, node.RegisterInfo.AccessToken)
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
	url := fmt.Sprintf("https://%s:%d/api/v1/container/%s/logs", node.IpAddress, node.Port, containerId)
	data := make([]string, 0)
	agent := getAgent(nodeId)
	if agent == nil {
		return nil, ErrNodeNotExist
	}
	err := agent.Get(url, querys, nil, &data, node.RegisterInfo.AccessToken)
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
	url := fmt.Sprintf("https://%s:%d/api/v1/container/%s/stats", node.IpAddress, node.Port, containerId)
	agent := getAgent(nodeId)
	if agent == nil {
		return nil, ErrNodeNotExist
	}
	err := agent.Get(url, nil, nil, stats, node.RegisterInfo.AccessToken)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
