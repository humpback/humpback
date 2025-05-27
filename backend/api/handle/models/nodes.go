package models

import (
	"regexp"
	"slices"
	"strings"

	"humpback/common/enum"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/common/verify"
	"humpback/pkg/utils"
	"humpback/types"
)

type NodesCreateReqInfo []string

func (n *NodesCreateReqInfo) Check() error {
	if len(*n) == 0 {
		return response.NewBadRequestErr(locales.CodeNodesNotEmpty)
	}
	ipRule := regexp.MustCompile(enum.RegularIpAddress)
	for _, ipAddress := range *n {
		if !ipRule.MatchString(ipAddress) {
			return response.NewBadRequestErr(locales.CodeNodesIpAddressInvalid)
		}
	}
	return nil
}

func (n *NodesCreateReqInfo) NewNodesInfo() []*types.Node {
	result := make([]*types.Node, 0)
	nowT := utils.NewActionTimestamp()
	for _, ip := range *n {
		result = append(result, &types.Node{
			NodeId:      utils.NewGuidStr(),
			Name:        "",
			IpAddress:   ip,
			Port:        8677,
			Status:      types.NodeStatusOffline,
			IsEnable:    true,
			CreatedAt:   nowT,
			UpdatedAt:   nowT,
			CPUUsage:    0,
			CPU:         0,
			MemoryUsage: 0,
			MemoryTotal: 0,
			MemoryUsed:  0,
			Labels:      make(map[string]string),
		})
	}
	return result
}

type NodeUpdateLabelReqInfo struct {
	NodeId string            `json:"nodeId"`
	Labels map[string]string `json:"labels"`
}

func (n *NodeUpdateLabelReqInfo) Check() error {
	if err := verify.CheckIsEmpty(n.NodeId, locales.CodeNodesIdNotEmpty); err != nil {
		return err
	}
	if n.Labels == nil {
		n.Labels = make(map[string]string)
	}
	return nil
}

type NodeUpdateSwitchReqInfo struct {
	NodeId string `json:"nodeId"`
	Enable bool   `json:"enable"`
}

func (n *NodeUpdateSwitchReqInfo) Check() error {
	if err := verify.CheckIsEmpty(n.NodeId, locales.CodeNodesIdNotEmpty); err != nil {
		return err
	}
	return nil
}

type NodeQueryFilterInfo struct {
	Status string `json:"status"`
}

type NodeQueryReqInfo struct {
	types.QueryInfo
	FilterInfo *NodeQueryFilterInfo `json:"-"`
}

func (n *NodeQueryReqInfo) Check() error {
	n.CheckBase()
	n.FilterInfo = new(NodeQueryFilterInfo)
	if err := ParseMapToStructConvert(n.Filter, n.FilterInfo); err != nil {
		return err
	}

	if n.FilterInfo != nil && n.FilterInfo.Status != "" && slices.Index([]string{
		types.NodeStatusOnline,
		types.NodeStatusOffline,
		types.SwitchEnabled,
		types.SwitchDisabled,
	}, n.FilterInfo.Status) == -1 {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}

	if n.Keywords != "" && n.Mode != "keywords" {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}

func (n *NodeQueryReqInfo) QueryFilter(nodes []*types.Node) []*types.Node {
	result := make([]*types.Node, 0)
	for _, node := range nodes {
		if n.filter(node) {
			result = append(result, node)
		}
	}
	n.sort(result)
	return result
}

func (n *NodeQueryReqInfo) filter(info *types.Node) bool {
	if n.FilterInfo != nil {
		switch n.FilterInfo.Status {
		case types.SwitchDisabled:
			if info.IsEnable {
				return false
			}
		case types.SwitchEnabled, types.NodeStatusOffline, types.NodeStatusOnline:
			if !info.IsEnable {
				return false
			}
			if n.FilterInfo.Status != types.SwitchEnabled && info.Status != n.FilterInfo.Status {
				return false
			}
		}
	}

	if strings.Contains(info.IpAddress, n.Keywords) || strings.Contains(strings.ToLower(info.Name), strings.ToLower(n.Keywords)) {
		return true
	}
	for key := range info.Labels {
		if strings.Contains(strings.ToLower(key), strings.ToLower(n.Keywords)) {
			return true
		}
	}
	return false
}

func (n *NodeQueryReqInfo) sort(list []*types.Node) []*types.Node {
	var sortField = []string{"ipAddress", "name", "updatedAt", "createdAt"}
	if n.SortInfo == nil || n.SortInfo.Field == "" || slices.Index(sortField, n.SortInfo.Field) == -1 {
		return list
	}
	slices.SortFunc(list, func(a, b *types.Node) int {
		switch n.SortInfo.Field {
		case "ipAddress":
			return types.QuerySortOrder(n.SortInfo.Order, strings.ToLower(a.IpAddress), strings.ToLower(b.IpAddress))
		case "name":
			return types.QuerySortOrder(n.SortInfo.Order, strings.ToLower(a.Name), strings.ToLower(b.Name))
		case "updatedAt":
			return types.QuerySortOrder(n.SortInfo.Order, a.UpdatedAt, b.UpdatedAt)
		case "createdAt":
			return types.QuerySortOrder(n.SortInfo.Order, a.CreatedAt, b.CreatedAt)
		}
		return 1
	})
	return list
}
