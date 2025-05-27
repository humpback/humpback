package scheduler

import (
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"

	"humpback/config"
	"humpback/internal/db"
	"humpback/types"
)

type NodeController struct {
	NodesInfo           map[string]*types.NodeSimpleInfo
	NodeHeartbeatChan   chan types.NodeSimpleInfo
	ContainerChangeChan chan types.ContainerStatus
	CheckInterval       int64
	CheckThreshold      int
	ThresholdInvterval  int64
	sync.RWMutex
}

func NewNodeController(nodeChan chan types.NodeSimpleInfo, containerChan chan types.ContainerStatus) *NodeController {
	nc := &NodeController{
		NodesInfo:           make(map[string]*types.NodeSimpleInfo),
		NodeHeartbeatChan:   nodeChan,
		ContainerChangeChan: containerChan,
		CheckInterval:       int64(config.BackendArgs().CheckInterval),
		CheckThreshold:      config.BackendArgs().CheckThreshold,
		ThresholdInvterval:  int64(config.BackendArgs().CheckInterval) * int64(config.BackendArgs().CheckThreshold),
	}

	go nc.CheckNodes()

	return nc
}

func (nc *NodeController) RestoreNodes() {
	nc.Lock()
	defer nc.Unlock()

	nodes, err := db.NodesGetAllEnabled()
	if err != nil {
		slog.Info("[Node Controller] restore nodes failed", "error", err)
		return
	}

	for _, node := range nodes {
		nc.NodesInfo[node.NodeId] = &types.NodeSimpleInfo{
			NodeId:          node.NodeId,
			IpAddress:       node.IpAddress,
			Status:          node.Status,
			LastHeartbeat:   node.UpdatedAt,
			OnlineThreshold: 1,
			CPUUsage:        node.CPUUsage,
			MemoryUsage:     node.MemoryUsage,
		}
	}
}

func (nc *NodeController) CheckNodes() {
	interval := nc.CheckInterval
	time.Sleep(time.Duration(rand.Int32N(int32(interval))) * time.Second)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	for range ticker.C {
		slog.Info("[Node Controller] Check nodes......")
		nc.CheckNodesCore()
	}
}

// 机器上下线时需要通知该机器所属的Group，去检查Group中所有service的状态
func (nc *NodeController) CheckNodesCore() {
	nc.Lock()
	defer nc.Unlock()

	currentTime := time.Now().Unix()
	statusChanged := false
	for nodeId, nodeInfo := range nc.NodesInfo {
		statusChanged = false
		if nodeInfo.Status == types.NodeStatusOnline {
			if currentTime-nodeInfo.LastHeartbeat > nc.ThresholdInvterval {
				slog.Info("[Node Controller] Node is not responding.", "nodeId", nodeId, "Last heartbeat", nodeInfo.LastHeartbeat)
				nodeInfo.Status = types.NodeStatusOffline
				statusChanged = true
			}
		}

		if nodeInfo.Status == types.NodeStatusOffline {
			if currentTime-nodeInfo.LastHeartbeat < nc.ThresholdInvterval &&
				nodeInfo.OnlineThreshold >= nc.CheckThreshold {
				slog.Info("[Node Controller] need report online node", "nodeId", nodeId)
				nodeInfo.Status = types.NodeStatusOnline
				statusChanged = true
			}
		}

		err := db.NodeUpdateStatus(nodeInfo)
		if err != nil {
			slog.Info("[Node Controller] update node status to DB failed", "error", err)
		}

		if statusChanged {
			nc.NodeHeartbeatChan <- *nodeInfo
		}
	}
}

func (nc *NodeController) HeartBeat(healthInfo types.HealthInfo) {
	nc.Lock()
	defer nc.Unlock()

	nodeId := healthInfo.NodeId
	ts := time.Now().Unix()
	if n, ok := nc.NodesInfo[nodeId]; ok {
		n.Name = healthInfo.HostInfo.Hostname
		n.Port = healthInfo.HostInfo.Port
		n.CPUUsage = healthInfo.HostInfo.CPUUsage
		n.MemoryUsage = healthInfo.HostInfo.MemoryUsage
		n.TotalCPU = healthInfo.HostInfo.TotalCPU
		n.UsedMemory = healthInfo.HostInfo.UsedMemory
		n.TotalMemory = healthInfo.HostInfo.TotalMemory
		if n.Status == types.NodeStatusOffline && ts-n.LastHeartbeat < nc.ThresholdInvterval {
			n.OnlineThreshold++
		} else {
			n.OnlineThreshold = 1
			nc.CheckContainers(healthInfo)
		}
		n.LastHeartbeat = ts
	} else {
		n, err := db.NodeGetById(nodeId)
		if err == nil {
			nc.NodesInfo[nodeId] = &types.NodeSimpleInfo{
				NodeId:          n.NodeId,
				IpAddress:       n.IpAddress,
				Name:            healthInfo.HostInfo.Hostname,
				Port:            healthInfo.HostInfo.Port,
				Status:          types.NodeStatusOffline,
				LastHeartbeat:   ts,
				OnlineThreshold: 1,
				CPUUsage:        healthInfo.HostInfo.CPUUsage,
				TotalCPU:        healthInfo.HostInfo.TotalCPU,
				MemoryUsage:     healthInfo.HostInfo.MemoryUsage,
				TotalMemory:     healthInfo.HostInfo.TotalMemory,
				UsedMemory:      healthInfo.HostInfo.UsedMemory,
			}
		}
	}
}

// 对于在线的机器，检查容器状态
func (nc *NodeController) CheckContainers(healthInfo types.HealthInfo) {
	for _, container := range healthInfo.ContainerList {
		container.NodeId = healthInfo.NodeId
		container.Ip = healthInfo.IpAddress
		nc.ContainerChangeChan <- container
	}
}
