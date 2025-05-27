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

type TeamCreateReqInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Users       []string `json:"users"`
}

func (t *TeamCreateReqInfo) Check() error {
	if err := verify.CheckRequiredAndLengthLimit(t.Name, enum.LimitTeamName.Min, enum.LimitTeamName.Max, locales.CodeTeamNameNotEmpty, locales.CodeTeamNameLimitLength); err != nil {
		return err
	}
	if err := verify.CheckLengthLimit(t.Description, 0, enum.LimitDescription.Max, locales.CodeDescriptionLimitLength); err != nil {
		return err
	}
	if len(t.Users) == 0 {
		t.Users = make([]string, 0)
	}
	return nil
}

func (t *TeamCreateReqInfo) NewTeamInfo() *types.Team {
	nowT := utils.NewActionTimestamp()
	return &types.Team{
		Name:        t.Name,
		Description: t.Description,
		Users:       t.Users,
		CreatedAt:   nowT,
		UpdatedAt:   nowT,
		TeamId:      utils.NewGuidStr(),
	}
}

type TeamUpdateReqInfo struct {
	TeamId string `json:"teamId"`
	TeamCreateReqInfo
}

func (t *TeamUpdateReqInfo) Check() error {
	if err := verify.CheckIsEmpty(t.TeamId, locales.CodeTeamIdNotEmpty); err != nil {
		return err
	}
	return t.TeamCreateReqInfo.Check()
}

func (t *TeamUpdateReqInfo) NewTeamInfo(oldTeamInfo *types.Team) *types.Team {
	return &types.Team{
		Name:        t.Name,
		Description: t.Description,
		Users:       t.Users,
		CreatedAt:   oldTeamInfo.CreatedAt,
		UpdatedAt:   utils.NewActionTimestamp(),
		TeamId:      t.TeamId,
	}
}

type TeamQueryReqInfo struct {
	types.QueryInfo
}

func (t *TeamQueryReqInfo) Check() error {
	t.QueryInfo.CheckBase()
	if t.Keywords != "" && t.Mode != "name" {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}

func (t *TeamQueryReqInfo) QueryFilter(teams []*types.Team) []*types.Team {
	result := make([]*types.Team, 0)
	for _, team := range teams {
		if strings.Contains(strings.ToLower(team.Name), strings.ToLower(t.Keywords)) {
			result = append(result, team)
		}
	}
	t.sort(result)
	return result
}

func (t *TeamQueryReqInfo) sort(list []*types.Team) []*types.Team {
	var sortField = []string{"name", "updatedAt", "createdAt"}
	if t.SortInfo == nil || t.SortInfo.Field == "" || slices.Index(sortField, t.SortInfo.Field) == -1 {
		return list
	}
	slices.SortFunc(list, func(a, b *types.Team) int {
		switch t.SortInfo.Field {
		case "name":
			return types.QuerySortOrder(t.SortInfo.Order, strings.ToLower(a.Name), strings.ToLower(b.Name))
		case "updatedAt":
			return types.QuerySortOrder(t.SortInfo.Order, a.UpdatedAt, b.UpdatedAt)
		case "createdAt":
			return types.QuerySortOrder(t.SortInfo.Order, a.CreatedAt, b.CreatedAt)
		}
		return 1
	})
	return list
}
