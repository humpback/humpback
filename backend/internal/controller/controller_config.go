package controller

import (
	"humpback/api/handle/models"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/internal/db"
	"humpback/types"
)

func ConfigCreate(operator *types.User, reqInfo *models.ConfigCreateReqInfo) (string, error) {
	if err := configCheckNameExist(reqInfo.ConfigName, ""); err != nil {
		return "", err
	}
	newInfo := reqInfo.NewConfigInfo()
	if err := db.ConfigUpdate(newInfo); err != nil {
		return "", response.NewRespServerErr(err.Error())
	}
	InsertConfigActivity(&ActivityConfigInfo{
		NewConfigInfo: newInfo,
		Action:        types.ActivityActionAdd,
		OperatorInfo:  operator,
		OperateAt:     newInfo.UpdatedAt,
	})
	return newInfo.ConfigId, nil
}

func ConfigUpdate(operator *types.User, reqInfo *models.ConfigUpdateReqInfo) (string, error) {
	if err := configCheckNameExist(reqInfo.ConfigName, reqInfo.ConfigId); err != nil {
		return "", err
	}
	oldInfo, err := Config(reqInfo.ConfigId)
	if err != nil {
		return "", err
	}
	newInfo := reqInfo.NewConfigInfo(oldInfo)
	if err = db.ConfigUpdate(newInfo); err != nil {
		return "", response.NewRespServerErr(err.Error())
	}
	InsertConfigActivity(&ActivityConfigInfo{
		OldConfigInfo: oldInfo,
		NewConfigInfo: newInfo,
		Action:        types.ActivityActionUpdate,
		OperatorInfo:  operator,
		OperateAt:     newInfo.UpdatedAt,
	})
	return newInfo.ConfigId, err
}

func configCheckNameExist(name, id string) error {
	sameNames, err := db.ConfigsGetByName(name, true)
	if err != nil {
		return response.NewRespServerErr(err.Error())
	}
	for _, info := range sameNames {
		if info.ConfigId != id {
			return response.NewBadRequestErr(locales.CodeConfigNameAlreadyExist)
		}
	}
	return nil
}

func Config(id string) (*types.Config, error) {
	info, err := db.ConfigGetById(id)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil, response.NewBadRequestErr(locales.CodeConfigNotExist)
		}
		return nil, response.NewRespServerErr(err.Error())
	}
	return info, nil
}

func ConfigQuery(queryInfo *models.ConfigQueryReqInfo) (*response.QueryResult[types.Config], error) {
	configs, err := db.ConfigsGetAll()
	if err != nil {
		return nil, response.NewRespServerErr(err.Error())
	}
	result := queryInfo.QueryFilter(configs)
	return response.NewQueryResult[types.Config](
		len(result),
		types.QueryPagination[types.Config](queryInfo.PageInfo, result),
	), nil
}

func ConfigDelete(operator *types.User, id string) error {
	info, err := db.ConfigGetById(id)
	if err != nil {
		if err == db.ErrKeyNotExist {
			return nil
		}
		return response.NewRespServerErr(err.Error())
	}
	if err := db.ConfigDelete(id); err != nil {
		return response.NewRespServerErr(err.Error())
	}
	InsertConfigActivity(&ActivityConfigInfo{
		OldConfigInfo: info,
		Action:        types.ActivityActionDelete,
		OperatorInfo:  operator,
		OperateAt:     0,
	})
	return nil
}
