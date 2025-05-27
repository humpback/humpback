package controller

import (
    "sync"
    
    "humpback/common/response"
    "humpback/internal/db"
    "humpback/types"
)

func DashboardResourceStatistics(userInfo *types.User) (*types.ResourceTotalInfo, error) {
    var (
        result = &types.ResourceTotalInfo{
            Services:          0,
            Nodes:             0,
            Users:             0,
            Groups:            0,
            OwnGroups:         0,
            OwnServices:       0,
            EnableOwnServices: 0,
            EnableOwnNodes:    0,
            ExceptionServices: make([]*types.ResourceExceptionServiceInfo, 0),
            AbnormalNodes:     make([]*types.Node, 0),
        }
        groups  = make(map[string]string)
        nodeIds = make(map[string]bool)
        wg      = &sync.WaitGroup{}
        l       = &sync.Mutex{}
        errMap  = map[string]error{}
    )
    
    allGroups, err := db.GroupsGetAll()
    if err != nil {
        return nil, response.NewRespServerErr(err.Error())
    }
    result.Groups = len(allGroups)
    for _, group := range allGroups {
        if userInfo.InGroup(group) {
            result.OwnGroups++
            groups[group.GroupId] = group.GroupName
            for _, node := range group.Nodes {
                nodeIds[node] = true
            }
        }
    }
    wg.Add(3)
    go func() {
        defer wg.Done()
        userTotal, err := db.UsersGetTotal()
        if err != nil {
            l.Lock()
            errMap["user"] = err
            l.Unlock()
            return
        }
        result.Users = userTotal
    }()
    go func() {
        defer wg.Done()
        nodes, enableNodes, abnormalNodes, err := db.NodesGetTotalAndAbnormalByIds(nodeIds)
        if err != nil {
            l.Lock()
            errMap["node"] = err
            l.Unlock()
            return
        }
        result.Nodes = nodes
        result.AbnormalNodes = abnormalNodes
        result.EnableOwnNodes = enableNodes
    }()
    go func() {
        defer wg.Done()
        services, ownServices, enableOwnServices, exceptionServices, err := db.ServicesGetTotalAndExceptionByGroups(groups)
        if err != nil {
            l.Lock()
            errMap["service"] = err
            l.Unlock()
        }
        result.Services = services
        result.OwnServices = ownServices
        result.ExceptionServices = exceptionServices
        result.EnableOwnServices = enableOwnServices
    }()
    wg.Wait()
    for _, err := range errMap {
        return nil, response.NewRespServerErr(err.Error())
    }
    return result, nil
}
