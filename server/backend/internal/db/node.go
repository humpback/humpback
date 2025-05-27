package db

import (
    "encoding/json"
    "fmt"
    "slices"
    
    bolt "go.etcd.io/bbolt"
    
    "humpback/types"
)

func NodesGetTotalAndAbnormalByIds(nodes map[string]bool) (int, int, []*types.Node, error) {
    var (
        total         = 0
        enableTotal   = 0
        abnormalNodes = make([]*types.Node, 0)
    )
    if err := TransactionGet(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(BucketNodes))
        if bucket == nil {
            return ErrBucketNotExist
        }
        c := bucket.Cursor()
        for k, v := c.First(); k != nil; k, v = c.Next() {
            total++
            if nodes[string(k)] {
                tempNode := new(types.Node)
                if err := json.Unmarshal(v, tempNode); err != nil {
                    return err
                }
                
                if tempNode.IsEnable {
                    enableTotal++
                    if tempNode.Status == types.NodeStatusOffline {
                        abnormalNodes = append(abnormalNodes, tempNode)
                    }
                }
            }
        }
        return nil
    }); err != nil {
        return 0, 0, nil, err
    }
    return total, enableTotal, abnormalNodes, nil
}

func NodesGetAll() ([]*types.Node, error) {
    return GetDataAll[types.Node](BucketNodes)
}

func NodeUpdateStatus(nodeInfo *types.NodeSimpleInfo) error {
    node, err := GetDataById[types.Node](BucketNodes, nodeInfo.NodeId)
    if err != nil {
        return err
    }
    node.Name = nodeInfo.Name
    node.Port = nodeInfo.Port
    node.Status = nodeInfo.Status
    node.UpdatedAt = nodeInfo.LastHeartbeat
    node.CPU = nodeInfo.TotalCPU
    node.CPUUsage = nodeInfo.CPUUsage
    node.MemoryTotal = nodeInfo.TotalMemory
    node.MemoryUsed = nodeInfo.UsedMemory
    node.MemoryUsage = nodeInfo.MemoryUsage
    return SaveData(BucketNodes, nodeInfo.NodeId, node)
}

func NodeGetById(nodeId string) (*types.Node, error) {
    return GetDataById[types.Node](BucketNodes, nodeId)
}

func NodeAndGroupsGetById(nodeId string) (*types.Node, []*types.NodesGroups, error) {
    var (
        node   = &types.Node{}
        groups = make([]*types.NodesGroups, 0)
    )
    if err := TransactionGet(func(tx *bolt.Tx) error {
        nodeBucket := tx.Bucket([]byte(BucketNodes))
        if nodeBucket == nil {
            return ErrBucketNotExist
        }
        data := nodeBucket.Get([]byte(nodeId))
        if data == nil {
            return ErrKeyNotExist
        }
        if err := json.Unmarshal(data, &node); err != nil {
            return fmt.Errorf("failed to decode node: %s", err)
        }
        groupBucket := tx.Bucket([]byte(BucketNodesGroups))
        if groupBucket == nil {
            return ErrBucketNotExist
        }
        groupBucket.ForEach(func(k, v []byte) error {
            result := new(types.NodesGroups)
            if err := json.Unmarshal(v, result); err != nil {
                return err
            }
            if slices.Contains(result.Nodes, nodeId) {
                groups = append(groups, result)
            }
            return nil
        })
        return nil
    }); err != nil {
        return nil, nil, err
    }
    return node, groups, nil
}

func NodesGetAllEnabled() ([]*types.Node, error) {
    nodes, err := GetDataByQuery[types.Node](BucketNodes, func(key string, node interface{}) bool {
        return node.(*types.Node).IsEnable
    })
    return nodes, err
}

func NodesGetEnabledByGroupId(groupId string) ([]*types.Node, error) {
    ng, err := GroupGetById(groupId)
    if err != nil {
        return nil, err
    }
    nodes := make([]*types.Node, 0)
    for _, v := range ng.Nodes {
        node, err := GetDataById[types.Node](BucketNodes, v)
        if err == nil && node.IsEnable {
            nodes = append(nodes, node)
        }
    }
    return nodes, nil
}

func NodesGetByIds(ids []string, ignoreNotExist bool) ([]*types.Node, error) {
    return GetDataByIds[types.Node](BucketNodes, ids, ignoreNotExist)
}

func NodeUpdate(node *types.Node) error {
    return SaveData[*types.Node](BucketNodes, node.NodeId, node)
}

func NodesUpdate(nodes []*types.Node) error {
    return TransactionUpdates(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(BucketNodes))
        if bucket == nil {
            return ErrBucketNotExist
        }
        for _, node := range nodes {
            data, err := json.Marshal(node)
            if err != nil {
                return fmt.Errorf("failed to encode node data: %s", err)
            }
            if err = bucket.Put([]byte(node.NodeId), data); err != nil {
                return err
            }
        }
        return nil
    })
}

func NodeDeleteAndGroupsUpdate(nodeId string, groups []*types.NodesGroups) error {
    return TransactionUpdates(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(BucketNodes))
        if bucket == nil {
            return ErrBucketNotExist
        }
        if err := bucket.Delete([]byte(nodeId)); err != nil {
            return err
        }
        groupBucket := tx.Bucket([]byte(BucketNodesGroups))
        if groupBucket == nil {
            return ErrBucketNotExist
        }
        for _, group := range groups {
            v, err := json.Marshal(group)
            if err != nil {
                return fmt.Errorf("delete node(%s), failed to encode group(%s) data: %s", nodeId, group.GroupId, err)
            }
            if err = groupBucket.Put([]byte(group.GroupId), v); err != nil {
                return err
            }
        }
        return nil
    })
}
