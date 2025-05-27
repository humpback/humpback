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

type RegistryCreateReqInfo struct {
	URL       string `json:"url"`
	IsDefault bool   `json:"isDefault"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (r *RegistryCreateReqInfo) Check() error {
	if r.Username != "" {
		r.Username = utils.RSADecrypt(r.Username)
	}
	if r.Password != "" {
		r.Password = utils.RSADecrypt(r.Password)
	}

	re := regexp.MustCompile(`/+`)
	r.URL = strings.Trim(re.ReplaceAllString(r.URL, "/"), "/")

	if err := verify.CheckRequiredAndLengthLimit(r.URL, 0, enum.LimitRegistryUrl.Max, locales.CodeRegistryUrlNotEmpty, locales.CodeRegistryUrlLimitLength); err != nil {
		return err
	}
	if err := verify.CheckLengthLimit(r.Username, 0, enum.LimitRegistryUsername.Max, locales.CodeRegistryUsernameLimitLength); err != nil {
		return err
	}
	if err := verify.CheckLengthLimit(r.Password, 0, enum.LimitRegistryPassword.Max, locales.CodeRegistryPasswordLimitLength); err != nil {
		return err
	}
	return nil
}

func (r *RegistryCreateReqInfo) NewRegistryInfo() *types.Registry {
	nowT := utils.NewActionTimestamp()
	return &types.Registry{
		RegistryId: utils.NewGuidStr(),
		URL:        r.URL,
		IsDefault:  r.IsDefault,
		Username:   r.Username,
		Password:   r.Password,
		CreatedAt:  nowT,
		UpdatedAt:  nowT,
	}
}

type RegistryUpdateReqInfo struct {
	RegistryId string `json:"registryId"`
	RegistryCreateReqInfo
}

func (r *RegistryUpdateReqInfo) Check() error {
	if err := verify.CheckIsEmpty(r.RegistryId, locales.CodeRegistryIdNotEmpty); err != nil {
		return err
	}
	return r.RegistryCreateReqInfo.Check()
}

func (r *RegistryUpdateReqInfo) NewRegistryInfo(oldInfo *types.Registry) *types.Registry {
	return &types.Registry{
		RegistryId: oldInfo.RegistryId,
		URL:        r.URL,
		IsDefault:  r.IsDefault,
		Username:   r.Username,
		Password:   r.Password,
		CreatedAt:  oldInfo.CreatedAt,
		UpdatedAt:  utils.NewActionTimestamp(),
	}
}

type RegistryQueryReqInfo struct {
	types.QueryInfo
}

func (r *RegistryQueryReqInfo) Check() error {
	r.QueryInfo.CheckBase()
	if r.Keywords != "" && r.Mode != "url" {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}

func (r *RegistryQueryReqInfo) QueryFilter(registries []*types.Registry) []*types.Registry {
	result := make([]*types.Registry, 0)
	for _, registry := range registries {
		if strings.Contains(strings.ToLower(registry.URL), strings.ToLower(r.Keywords)) {
			result = append(result, registry)
		}
	}
	r.sort(result)
	return result
}

func (r *RegistryQueryReqInfo) sort(list []*types.Registry) []*types.Registry {
	var sortField = []string{"url", "updatedAt", "createdAt"}
	if r.SortInfo == nil || r.SortInfo.Field == "" || slices.Index(sortField, r.SortInfo.Field) == -1 {
		return list
	}
	slices.SortFunc(list, func(a, b *types.Registry) int {
		switch r.SortInfo.Field {
		case "url":
			return types.QuerySortOrder(r.SortInfo.Order, strings.ToLower(a.URL), strings.ToLower(b.URL))
		case "updatedAt":
			return types.QuerySortOrder(r.SortInfo.Order, a.UpdatedAt, b.UpdatedAt)
		case "createdAt":
			return types.QuerySortOrder(r.SortInfo.Order, a.CreatedAt, b.CreatedAt)
		}
		return 1
	})
	return list
}
