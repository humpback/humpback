package db

import (
    "encoding/json"
    "fmt"
    "slices"
    "strings"
    
    bolt "go.etcd.io/bbolt"
    "humpback/types"
)

func GroupsGetAll() ([]*types.NodesGroups, error) {
    return GetDataAll[types.NodesGroups](BucketNodesGroups)
}

func GroupGetById(id string) (*types.NodesGroups, error) {
    return GetDataById[types.NodesGroups](BucketNodesGroups, id)
}

func GroupsGetByNodeId(nodeId string) ([]*types.NodesGroups, error) {
    groups, err := GetDataByQuery[types.NodesGroups](BucketNodesGroups, func(key string, value interface{}) bool {
        ngp := value.(*types.NodesGroups)
        return slices.Contains(ngp.Nodes, nodeId)
    })
    if err != nil {
        return nil, err
    }
    return groups, nil
}

func GroupsGetByUserId(userId string) ([]*types.NodesGroups, error) {
    groups, err := GetDataByQuery[types.NodesGroups](BucketNodesGroups, func(key string, value interface{}) bool {
        ngp := value.(*types.NodesGroups)
        return slices.Contains(ngp.Users, userId)
    })
    if err != nil {
        return nil, err
    }
    return groups, nil
}

func GroupsGetByTeamId(TeamId string) ([]*types.NodesGroups, error) {
    groups, err := GetDataByQuery[types.NodesGroups](BucketNodesGroups, func(key string, value interface{}) bool {
        ngp := value.(*types.NodesGroups)
        return slices.Contains(ngp.Teams, TeamId)
    })
    if err != nil {
        return nil, err
    }
    return groups, nil
}

func GroupsGetByName(name string, isLower bool) ([]*types.NodesGroups, error) {
    groups, err := GroupsGetAll()
    if err != nil {
        return nil, err
    }
    var result []*types.NodesGroups
    for _, group := range groups {
        if isLower && strings.EqualFold(group.GroupName, name) {
            result = append(result, group)
        }
        if !isLower && group.GroupName == name {
            result = append(result, group)
        }
    }
    return result, nil
}

func GroupUpdate(info *types.NodesGroups) error {
    return SaveData(BucketNodesGroups, info.GroupId, info)
}

func GroupDelete(id string) error {
    return DeleteData(BucketNodesGroups, id)
}

func GroupDeleteAndServiceSoftDelete(id string, services []*types.Service) error {
    if err := TransactionUpdates(func(tx *bolt.Tx) error {
        var (
            groupBucket   *bolt.Bucket
            serviceBucket *bolt.Bucket
        )
        groupBucket = tx.Bucket([]byte(BucketNodesGroups))
        if groupBucket == nil {
            return ErrBucketNotExist
        }
        serviceBucket = tx.Bucket([]byte(BucketServices))
        if serviceBucket == nil {
            return ErrBucketNotExist
        }
        if err := groupBucket.Delete([]byte(id)); err != nil {
            return fmt.Errorf("failed to delete group(%s): %s", id, err)
        }
        
        for _, service := range services {
            service.IsDelete = true
            v, err := json.Marshal(service)
            if err != nil {
                return fmt.Errorf("failed to delete group(%s), failed to encode servie(%s) data: %s", id, service.ServiceId, err)
            }
            if err = serviceBucket.Put([]byte(service.ServiceId), v); err != nil {
                return fmt.Errorf("failed to delete group(%s), failed to delete service(%s): %s", id, service.ServiceId, err)
            }
        }
        return nil
    }); err != nil {
        return err
    }
    return nil
}
