package controller

import (
    "fmt"
    "slices"
    "strings"
    "sync"
    "time"
    
    "humpback/api/handle/models"
    "humpback/common/locales"
    "humpback/common/response"
    "humpback/internal/node"
    "humpback/types"
)

func InstanceOperate(operator *types.User, groupInfo *types.NodesGroups, info *models.InstanceOperateReqInfo) error {
    if slices.Index(groupInfo.Nodes, info.NodeId) == -1 {
        return response.NewBadRequestErr(locales.CodeNodesNotExist)
    }
    service, container, err := getServcieInstanceInfo(info.GroupId, info.ServiceId, info.ContainerId)
    if err != nil {
        return err
    }
    if err = node.OperateNodeContainer(info.NodeId, info.ContainerId, info.Action); err != nil {
        return err
    }
    action := types.ActivityActionStart
    if strings.ToLower(info.Action) == strings.ToLower(types.ServiceActionRestart) {
        action = types.ActivityActionRestart
    }
    if strings.ToLower(info.Action) == strings.ToLower(types.ServiceActionStop) {
        action = types.ActivityActionStop
    }
    InsertServiceActivity(&ActivityServiceInfo{
        NewServiceInfo: service,
        Action:         action,
        InstanceName:   container.ContainerName,
        OperatorInfo:   operator,
        OperateAt:      0,
        GroupName:      groupInfo.GroupName,
    })
    return nil
}

func InstanceQueryLogs(groupInfo *types.NodesGroups, info *models.InstanceLogsReqInfo) ([]string, error) {
    if slices.Index(groupInfo.Nodes, info.NodeId) == -1 {
        return nil, response.NewBadRequestErr(locales.CodeNodesNotExist)
    }
    _, _, err := getServcieInstanceInfo(info.GroupId, info.ServiceId, info.ContainerId)
    if err != nil {
        return nil, err
    }
    query := map[string]string{
        "containerId": info.ContainerId,
        "timestamps":  fmt.Sprintf("%v", info.ShowTimestamp),
    }
    if info.Line > 0 {
        query["tail"] = fmt.Sprintf("%d", info.Line)
    }
    if info.StartAt > 0 {
        query["since"] = time.UnixMilli(info.StartAt).Format(time.RFC3339Nano)
    }
    if info.EndAt > 0 {
        query["until"] = time.UnixMilli(info.EndAt).Format(time.RFC3339Nano)
    }
    logs, err := node.QueryContainerLogs(info.NodeId, info.ContainerId, query)
    if err != nil {
        return nil, err
    }
    return logs, nil
}

func getServcieInstanceInfo(groupId, serviceId, containerId string) (*types.Service, *types.ContainerStatus, error) {
    service, err := Service(groupId, serviceId)
    if err != nil {
        return nil, nil, err
    }
    index := slices.IndexFunc(service.Containers, func(item *types.ContainerStatus) bool {
        return containerId == item.ContainerId
    })
    if index == -1 {
        return nil, nil, response.NewBadRequestErr(locales.CodeContainerNotExist)
    }
    return service, service.Containers[index], nil
}

func InsatncePerformances(groupInfo *types.NodesGroups, info *models.InstancesPerformanceReqInfo) (*types.InstancesPerformance, error) {
    service, err := Service(info.GroupId, info.ServiceId)
    if err != nil {
        return nil, err
    }
    
    type tempResult struct {
        NodeId      string
        ContaienrId string
        Stats       *node.ContainerStats
        Err         error
    }
    var (
        l       sync.Mutex
        wg      = &sync.WaitGroup{}
        tempMap = map[string]*tempResult{}
        result  = &types.InstancesPerformance{
            StatsAt:    time.Now().UnixMilli(),
            Containers: make([]*types.GroupContainerPerformance, 0),
        }
        validContainer = make([]*models.InstanceStatsReqInfo, 0)
    )
    for _, container := range info.Containers {
        if slices.IndexFunc(groupInfo.Nodes, func(item string) bool {
            return item == container.NodeId
        }) != -1 &&
            slices.IndexFunc(service.Containers, func(item *types.ContainerStatus) bool {
                return container.ContainerId != "" && container.ContainerId == item.ContainerId
            }) != -1 {
            validContainer = append(validContainer, container)
        }
    }
    
    for _, info := range validContainer {
        wg.Add(1)
        go func(info *models.InstanceStatsReqInfo) {
            stats, err := node.GetContainerStats(info.NodeId, info.ContainerId)
            l.Lock()
            tempMap[info.ContainerId] = &tempResult{NodeId: info.NodeId, ContaienrId: info.ContainerId, Stats: stats, Err: err}
            l.Unlock()
            wg.Done()
        }(info)
    }
    wg.Wait()
    for _, p := range tempMap {
        t := &types.GroupContainerPerformance{
            ContainerId: p.ContaienrId,
            NodeId:      p.NodeId,
            IsSuccess:   p.Err == nil,
            Error:       "",
            Stats:       nil,
        }
        if p.Err == nil {
            t.Stats = &types.GroupContainerStats{
                CpuPercent:  p.Stats.CpuPercent,
                MemoryUsed:  p.Stats.MermoryUsed,
                MemoryLimit: p.Stats.MemoryLimit,
                IORead:      p.Stats.DiskReadBytes,
                IOWrite:     p.Stats.DiskWriteBytes,
                Networks:    p.Stats.Networks,
            }
        } else {
            t.Error = p.Err.Error()
        }
        result.Containers = append(result.Containers, t)
    }
    return result, nil
}
