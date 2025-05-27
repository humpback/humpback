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

type ConfigCreateReqInfo struct {
    ConfigName  string           `json:"configName"`
    Description string           `json:"description"`
    ConfigType  types.ConfigType `json:"configType"`
    ConfigValue string           `json:"configValue"`
}

func (c *ConfigCreateReqInfo) Check() error {
    if err := verify.CheckRequiredAndLengthLimit(c.ConfigName, enum.LimitConfigName.Min, enum.LimitConfigName.Max, locales.CodeConfigNameNotEmpty, locales.CodeConfigNameLimitLength); err != nil {
        return err
    }
    if err := verify.CheckLengthLimit(c.Description, 0, enum.LimitDescription.Max, locales.CodeDescriptionLimitLength); err != nil {
        return err
    }
    if c.ConfigType != types.ConfigTypeStatic && c.ConfigType != types.ConfigTypeVolume {
        return response.NewBadRequestErr(locales.CodeConfigTypeIsInvalid)
    }
    if c.ConfigType == types.ConfigTypeStatic {
        if err := verify.CheckRequiredAndLengthLimit(c.ConfigValue, enum.LimitConfigValue.Min, enum.LimitConfigValue.Max/2, locales.CodeConfigValueNotEmpty, locales.CodeConfigStaticValueLimitLength); err != nil {
            return err
        }
    } else {
        if err := verify.CheckRequiredAndLengthLimit(c.ConfigValue, enum.LimitConfigValue.Min, enum.LimitConfigValue.Max, locales.CodeConfigValueNotEmpty, locales.CodeConfigVolumeValueLimitLength); err != nil {
            return err
        }
    }
    return nil
}

func (c *ConfigCreateReqInfo) NewConfigInfo() *types.Config {
    nowT := utils.NewActionTimestamp()
    return &types.Config{
        ConfigName:  c.ConfigName,
        Description: c.Description,
        ConfigType:  c.ConfigType,
        ConfigValue: c.ConfigValue,
        CreatedAt:   nowT,
        UpdatedAt:   nowT,
        ConfigId:    utils.NewGuidStr(),
    }
}

type ConfigUpdateReqInfo struct {
    ConfigId string `json:"configId"`
    ConfigCreateReqInfo
}

func (c *ConfigUpdateReqInfo) Check() error {
    if err := verify.CheckIsEmpty(c.ConfigId, locales.CodeConfigIdNotEmpty); err != nil {
        return err
    }
    return c.ConfigCreateReqInfo.Check()
}

func (c *ConfigUpdateReqInfo) NewConfigInfo(oldInfo *types.Config) *types.Config {
    return &types.Config{
        ConfigName:  c.ConfigName,
        Description: c.Description,
        ConfigType:  c.ConfigType,
        ConfigValue: c.ConfigValue,
        CreatedAt:   oldInfo.CreatedAt,
        UpdatedAt:   utils.NewActionTimestamp(),
        ConfigId:    oldInfo.ConfigId,
    }
}

type ConfigQueryFilterInfo struct {
    ConfigType int `json:"configType"`
}

type ConfigQueryReqInfo struct {
    types.QueryInfo
    FilterInfo *ConfigQueryFilterInfo `json:"-"`
}

func (c *ConfigQueryReqInfo) Check() error {
    c.QueryInfo.CheckBase()
    if c.Keywords != "" && c.Mode != "configName" {
        return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
    }
    c.FilterInfo = new(ConfigQueryFilterInfo)
    if err := ParseMapToStructConvert(c.Filter, c.FilterInfo); err != nil {
        return err
    }
    return nil
}

func (c *ConfigQueryReqInfo) QueryFilter(configs []*types.Config) []*types.Config {
    result := make([]*types.Config, 0)
    for _, config := range configs {
        if c.filter(config) {
            result = append(result, config)
        }
    }
    c.sort(result)
    return result
}

func (c *ConfigQueryReqInfo) filter(info *types.Config) bool {
    if c.FilterInfo != nil && c.FilterInfo.ConfigType != 0 && int(info.ConfigType) != c.FilterInfo.ConfigType {
        return false
    }
    return strings.Contains(strings.ToLower(info.ConfigName), strings.ToLower(c.Keywords))
}

func (c *ConfigQueryReqInfo) sort(list []*types.Config) []*types.Config {
    var sortField = []string{"configName", "updatedAt", "createdAt"}
    if c.SortInfo == nil || c.SortInfo.Field == "" || slices.Index(sortField, c.SortInfo.Field) == -1 {
        return list
    }
    slices.SortFunc(list, func(a, b *types.Config) int {
        switch c.SortInfo.Field {
        case "configName":
            return types.QuerySortOrder(c.SortInfo.Order, strings.ToLower(a.ConfigName), strings.ToLower(b.ConfigName))
        case "updatedAt":
            return types.QuerySortOrder(c.SortInfo.Order, a.UpdatedAt, b.UpdatedAt)
        case "createdAt":
            return types.QuerySortOrder(c.SortInfo.Order, a.CreatedAt, b.CreatedAt)
        }
        return 1
    })
    return list
}
