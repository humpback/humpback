package db

import (
	"humpback/pkg/utils"
	"humpback/types"
)

func SessionsGetAll() ([]*types.Session, error) {
	return GetDataAll[types.Session](BucketSessions)
}

func SessionGCByIds(ids []string) error {
	return DeleteDataByIds(BucketSessions, ids)
}

func SessionGetById(sessionId string) (*types.Session, bool, error) {
	sessionInfo, err := GetDataById[types.Session](BucketSessions, sessionId)
	if err != nil {
		if err == ErrKeyNotExist {
			return nil, true, nil
		}
		return nil, false, err
	}
	return sessionInfo, sessionInfo.ExpiredAt < utils.NewActionTimestamp(), nil
}

func SessionUpdate(data *types.Session) error {
	return SaveData[*types.Session](BucketSessions, data.SessionId, data)
}

func SessionDelete(sessionId string) error {
	return DeleteData(BucketSessions, sessionId)
}

func SessionBatchDeleteByUserId(userId string) error {
	sessions, err := SessionsGetAll()
	if err != nil {
		return err
	}
	var ids []string
	for _, session := range sessions {
		if session.UserId == userId {
			ids = append(ids, session.SessionId)
		}
	}
	if len(ids) > 0 {
		if err = DeleteDataByIds(BucketSessions, ids); err != nil {
			return err
		}
	}
	return nil
}
