package controller

import (
	"github.com/jinzhu/copier"
	"humpback/api/handle/models"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/internal/db"
	"humpback/pkg/utils"
	"humpback/types"
)

func NodeCreate(operator *types.User, nodeCh chan types.NodeSimpleInfo, nodes models.NodesCreateReqInfo) error {
	nodeList, err := Nodes()
	if err != nil {
		return err
	}
	addNodes := nodes.NewNodesInfo()
	for _, node := range addNodes {
		for _, existNode := range nodeList {
			if node.IpAddress == existNode.IpAddress {
				return response.NewBadRequestErr(locales.CodeNodesIpAddressAlreadyExist)
			}
		}
	}
	if err = db.NodesUpdate(addNodes); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	for _, node := range addNodes {
		sendNodeEvent(nodeCh, node.NodeId, "")
		InsertNodeActivity(&ActivityNodeInfo{
			NewNodeInfo:  node,
			Action:       types.ActivityActionAdd,
			OperatorInfo: operator,
			OperateAt:    node.UpdatedAt,
		})
	}
	InsertStatisticsCount(&StatisticalCountEvent{
		CreateAt: addNodes[0].UpdatedAt,
		Type:     types.CountTypeNode,
		Num:      len(addNodes),
		UserId:   operator.UserId,
	})
	return nil
}

func NodeUpdateLabel(operator *types.User, nodeCh chan types.NodeSimpleInfo, info *models.NodeUpdateLabelReqInfo) (string, error) {
	node, err := Node(info.NodeId)
	if err != nil {
		return "", err
	}
	newNode := new(types.Node)
	copier.Copy(newNode, node)
	newNode.Labels = info.Labels
	newNode.UpdatedAt = utils.NewActionTimestamp()
	if err = db.NodeUpdate(newNode); err != nil {
		return "", response.NewRespServerErr(err.Error())
	}
	sendNodeEvent(nodeCh, info.NodeId, "")
	InsertNodeActivity(&ActivityNodeInfo{
		OldNodeInfo:  node,
		NewNodeInfo:  newNode,
		Action:       types.ActivityActionUpdateLabel,
		OperatorInfo: operator,
		OperateAt:    node.UpdatedAt,
	})
	return info.NodeId, nil
}

func NodeUpdateSwitch(operator *types.User, nodeCh chan types.NodeSimpleInfo, nodeId string, enabled bool) (string, error) {
	node, err := Node(nodeId)
	if err != nil {
		return "", err
	}
	if node.IsEnable != enabled {
		node.IsEnable = enabled
		if err = db.NodeUpdate(node); err != nil {
			return "", response.NewRespServerErr(err.Error())
		}
		sendNodeEvent(nodeCh, nodeId, "")
		action := types.ActivityActionEnable
		if !node.IsEnable {
			action = types.ActivityActionDisable
		}
		InsertNodeActivity(&ActivityNodeInfo{
			NewNodeInfo:  node,
			Action:       action,
			OperatorInfo: operator,
			OperateAt:    0,
		})
	}
	return nodeId, nil
}

func NodeDelete(operator *types.User, nodeCh chan types.NodeSimpleInfo, id string) error {
	node, groups, err := db.NodeAndGroupsGetById(id)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil
		}
		return response.NewRespServerErr(err.Error())
	}
	for _, group := range groups {
		nodeIds := make([]string, 0)
		for _, nodeId := range group.Nodes {
			if nodeId != id {
				nodeIds = append(nodeIds, nodeId)
			}
		}
		group.Nodes = nodeIds
	}
	if err = db.NodeDeleteAndGroupsUpdate(node.NodeId, groups); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	if node.IsEnable {
		sendNodeEvent(nodeCh, id, "")
	}
	InsertNodeActivity(&ActivityNodeInfo{
		OldNodeInfo:  node,
		Action:       types.ActivityActionDelete,
		OperatorInfo: operator,
		OperateAt:    0,
	})
	return nil
}

func Node(id string) (*types.Node, error) {
	node, err := db.NodeGetById(id)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil, response.NewBadRequestErr(locales.CodeNodesNotExist)
		}
		return nil, response.NewRespServerErr(err.Error())
	}
	return node, nil
}

func Nodes() ([]*types.Node, error) {
	nodes, err := db.NodesGetAll()
	if err != nil {
		return nil, response.NewRespServerErr(err.Error())
	}
	return nodes, nil
}

func NodesQuery(queryInfo *models.NodeQueryReqInfo) (*response.QueryResult[types.Node], error) {
	nodes, err := Nodes()
	if err != nil {
		return nil, err
	}
	result := queryInfo.QueryFilter(nodes)
	return response.NewQueryResult[types.Node](
		len(result),
		types.QueryPagination[types.Node](queryInfo.PageInfo, result),
	), nil
}

func NodesGetByIds(ids []string, ignoreNotExist bool) ([]*types.Node, error) {
	list, err := db.NodesGetByIds(ids, ignoreNotExist)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil, response.NewBadRequestErr(locales.CodeNodesNotExist)
		}
		return nil, response.NewRespServerErr(err.Error())
	}
	return list, nil
}
