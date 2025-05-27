package models

import (
	"slices"
	"strings"

	"github.com/jinzhu/copier"
	"golang.org/x/exp/maps"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/pkg/utils"
	"humpback/types"
)

type GroupUpdateNodesReqInfo struct {
	GroupId  string   `json:"-"`
	IsDelete bool     `json:"isDelete"`
	Nodes    []string `json:"nodes"`
}

func (g *GroupUpdateNodesReqInfo) Check() error {
	nodeList := make([]string, 0)
	for _, node := range g.Nodes {
		if strings.TrimSpace(node) != "" {
			nodeList = append(nodeList, node)
		}
	}
	g.Nodes = nodeList
	if len(g.Nodes) == 0 {
		return response.NewBadRequestErr(locales.CodeNodesNotEmpty)
	}
	return nil
}

func (g *GroupUpdateNodesReqInfo) NewGroupInfo(oldInfo *types.NodesGroups) *types.NodesGroups {
	tempNodes := make(map[string]bool)
	for _, id := range oldInfo.Nodes {
		tempNodes[id] = true
	}
	for _, id := range g.Nodes {
		if g.IsDelete {
			delete(tempNodes, id)
		} else {
			tempNodes[id] = true
		}
	}
	newInfo := &types.NodesGroups{}
	copier.Copy(newInfo, oldInfo)
	newInfo.Nodes = maps.Keys(tempNodes)
	newInfo.UpdatedAt = utils.NewActionTimestamp()
	return newInfo
}

type GroupQueryNodesReqInfo struct {
	types.QueryInfo
}

func (g *GroupQueryNodesReqInfo) Check() error {
	g.CheckBase()
	if g.Keywords != "" && g.Mode != "keywords" {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}

func (g *GroupQueryNodesReqInfo) QueryFilter(nodes []*types.Node) []*types.Node {
	result := make([]*types.Node, 0)
	for _, node := range nodes {
		if strings.Contains(node.IpAddress, g.Keywords) || strings.Contains(strings.ToLower(node.Name), strings.ToLower(g.Keywords)) {
			result = append(result, node)
		}
	}
	g.sort(result)
	return result
}

func (g *GroupQueryNodesReqInfo) sort(list []*types.Node) []*types.Node {
	var sortField = []string{"ipAddress", "name", "updatedAt", "createdAt"}
	if g.SortInfo == nil || g.SortInfo.Field == "" || slices.Index(sortField, g.SortInfo.Field) == -1 {
		return list
	}
	slices.SortFunc(list, func(a, b *types.Node) int {
		switch g.SortInfo.Field {
		case "ipAddress":
			return types.QuerySortOrder(g.SortInfo.Order, strings.ToLower(a.IpAddress), strings.ToLower(b.IpAddress))
		case "name":
			return types.QuerySortOrder(g.SortInfo.Order, strings.ToLower(a.Name), strings.ToLower(b.Name))
		case "updatedAt":
			return types.QuerySortOrder(g.SortInfo.Order, a.UpdatedAt, b.UpdatedAt)
		case "createdAt":
			return types.QuerySortOrder(g.SortInfo.Order, a.CreatedAt, b.CreatedAt)
		}
		return 1
	})
	return list
}
