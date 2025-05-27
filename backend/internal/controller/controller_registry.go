package controller

import (
	"strings"

	"humpback/api/handle/models"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/internal/db"
	"humpback/types"
)

func RegistryCreate(operator *types.User, reqInfo *models.RegistryCreateReqInfo) (string, error) {
	_, defaultRegistry, err := registryCheck(reqInfo, "")
	if err != nil {
		return "", err
	}
	newInfo := reqInfo.NewRegistryInfo()
	updateList := []*types.Registry{newInfo}
	if newInfo.IsDefault && defaultRegistry != nil {
		defaultRegistry.IsDefault = false
		updateList = append(updateList, defaultRegistry)
	}
	if err = db.RegistryUpdate(updateList); err != nil {
		return "", response.NewRespServerErr(err.Error())
	}
	InsertRegistryActivity(&ActivityRegistryInfo{
		NewRegistryInfo: newInfo,
		Action:          types.ActivityActionAdd,
		OperatorInfo:    operator,
		OperateAt:       newInfo.UpdatedAt,
	})
	return newInfo.RegistryId, err
}

func RegistryUpdate(operator *types.User, reqInfo *models.RegistryUpdateReqInfo) (string, error) {
	oldRegistry, defaultRegistry, err := registryCheck(&reqInfo.RegistryCreateReqInfo, reqInfo.RegistryId)
	if err != nil {
		return "", err
	}
	newInfo := reqInfo.NewRegistryInfo(oldRegistry)
	updateList := []*types.Registry{newInfo}
	if newInfo.IsDefault && defaultRegistry != nil {
		defaultRegistry.IsDefault = false
		updateList = append(updateList, defaultRegistry)
	}
	if err = db.RegistryUpdate(updateList); err != nil {
		return "", response.NewRespServerErr(err.Error())
	}
	InsertRegistryActivity(&ActivityRegistryInfo{
		OldRegistryInfo: oldRegistry,
		NewRegistryInfo: newInfo,
		Action:          types.ActivityActionUpdate,
		OperatorInfo:    operator,
		OperateAt:       newInfo.UpdatedAt,
	})
	return newInfo.RegistryId, err
}

func registryCheck(reqInfo *models.RegistryCreateReqInfo, id string) (*types.Registry, *types.Registry, error) {
	list, err := Registries()
	if err != nil {
		return nil, nil, err
	}
	var (
		defaultInfo *types.Registry
		currentInfo *types.Registry
	)
	for _, info := range list {
		if info.RegistryId == id {
			currentInfo = info
			continue
		}
		if strings.ToLower(info.URL) == strings.ToLower(reqInfo.URL) {
			return nil, nil, response.NewBadRequestErr(locales.CodeRegistryUrlAlreadyExist)
		}
		if info.IsDefault {
			defaultInfo = info
		}
	}
	return currentInfo, defaultInfo, nil
}

func Registry(id string) (*types.Registry, error) {
	info, err := db.RegistryGetById(id)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil, response.NewBadRequestErr(locales.CodeRegistryNotExist)
		}
		return nil, response.NewRespServerErr(err.Error())
	}
	return info, nil
}

func Registries() ([]*types.Registry, error) {
	list, err := db.RegistryGetAll()
	if err != nil {
		return nil, response.NewRespServerErr(err.Error())
	}
	return list, nil
}

func RegistryQuery(queryInfo *models.RegistryQueryReqInfo) (*response.QueryResult[types.Registry], error) {
	registries, err := Registries()
	if err != nil {
		return nil, err
	}
	result := queryInfo.QueryFilter(registries)
	return response.NewQueryResult[types.Registry](
		len(result),
		types.QueryPagination[types.Registry](queryInfo.PageInfo, result),
	), nil
}

func RegistryDelete(operator *types.User, id string) error {
	info, err := Registry(id)
	if err != nil {
		return err
	}
	if strings.ToLower(info.URL) == "docker.io" {
		return response.NewBadRequestErr(locales.CodeRegistryDefaultNotDelete)
	}
	if err := db.RegistryDelete(id); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	InsertRegistryActivity(&ActivityRegistryInfo{
		OldRegistryInfo: info,
		Action:          types.ActivityActionDelete,
		OperatorInfo:    operator,
		OperateAt:       0,
	})
	return nil
}
