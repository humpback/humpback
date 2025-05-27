package models

import (
	"slices"
	"strings"

	"humpback/common/enum"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/common/verify"
	"humpback/pkg/utils"
	"humpback/types"
)

type GroupCreateReqInfo struct {
	GroupName   string   `json:"groupName"`
	Description string   `json:"description"`
	Users       []string `json:"users"`
	Teams       []string `json:"teams"`
}

func (g *GroupCreateReqInfo) Check() error {
	if err := verify.CheckRequiredAndLengthLimit(g.GroupName, enum.LimitGroupName.Min, enum.LimitGroupName.Max, locales.CodeGroupNameNotEmpty, locales.CodeGroupNameLimitLength); err != nil {
		return err
	}
	if err := verify.CheckLengthLimit(g.Description, enum.LimitDescription.Min, enum.LimitDescription.Max, locales.CodeDescriptionLimitLength); err != nil {
		return err
	}
	if g.Users == nil {
		g.Users = make([]string, 0)
	}
	if g.Teams == nil {
		g.Teams = make([]string, 0)
	}
	return nil
}

func (g *GroupCreateReqInfo) NewGroupInfo() *types.NodesGroups {
	nowT := utils.NewActionTimestamp()
	return &types.NodesGroups{
		GroupId:     utils.NewGroupId(),
		GroupName:   g.GroupName,
		Description: g.Description,
		Users:       g.Users,
		Teams:       g.Teams,
		Nodes:       make([]string, 0),
		CreatedAt:   nowT,
		UpdatedAt:   nowT,
	}
}

type GroupUpdateReqInfo struct {
	GroupId string `json:"groupId"`
	GroupCreateReqInfo
}

func (g *GroupUpdateReqInfo) Check() error {
	if err := verify.CheckIsEmpty(g.GroupId, locales.CodeGroupIdNotEmpty); err != nil {
		return err
	}
	return g.GroupCreateReqInfo.Check()
}

func (g *GroupUpdateReqInfo) NewGroupInfo(oldInfo *types.NodesGroups) *types.NodesGroups {
	return &types.NodesGroups{
		GroupId:     oldInfo.GroupId,
		GroupName:   g.GroupName,
		Description: g.Description,
		Users:       g.Users,
		Teams:       g.Teams,
		Nodes:       oldInfo.Nodes,
		CreatedAt:   oldInfo.CreatedAt,
		UpdatedAt:   utils.NewActionTimestamp(),
	}
}

type GroupQueryReqInfo struct {
	UserInfo *types.User `json:"-"`
	types.QueryInfo
}

func (g *GroupQueryReqInfo) Check() error {
	g.QueryInfo.CheckBase()
	if g.Keywords != "" && g.Mode != "groupName" {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}

func (g *GroupQueryReqInfo) QueryFilter(groups []*types.NodesGroups) []*types.NodesGroups {
	result := make([]*types.NodesGroups, 0)
	for _, group := range groups {
		if g.UserInfo.InGroup(group) && strings.Contains(strings.ToLower(group.GroupName), strings.ToLower(g.Keywords)) {
			result = append(result, group)
		}
	}
	g.sort(result)
	return result
}

func (g *GroupQueryReqInfo) sort(list []*types.NodesGroups) []*types.NodesGroups {
	var sortField = []string{"groupName", "updatedAt", "createdAt"}
	if g.SortInfo == nil || g.SortInfo.Field == "" || slices.Index(sortField, g.SortInfo.Field) == -1 {
		return list
	}
	slices.SortFunc(list, func(a, b *types.NodesGroups) int {
		switch g.SortInfo.Field {
		case "groupName":
			return types.QuerySortOrder(g.SortInfo.Order, strings.ToLower(a.GroupName), strings.ToLower(b.GroupName))
		case "updatedAt":
			return types.QuerySortOrder(g.SortInfo.Order, a.UpdatedAt, b.UpdatedAt)
		case "createdAt":
			return types.QuerySortOrder(g.SortInfo.Order, a.CreatedAt, b.CreatedAt)
		}
		return 1
	})
	return list
}
