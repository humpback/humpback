package models

import (
	"strconv"
	"strings"

	"humpback/common/locales"
	"humpback/common/response"
	"humpback/types"
)

type StatisticsCountQueryReqInfo struct {
	UserInfo *types.User `json:"-"`
	StartAt  int64       `json:"startAt"`
	EndAt    int64       `json:"endAt"`
	OnlyMe   bool        `json:"onlyMe"`
}

func (s *StatisticsCountQueryReqInfo) Check() error {
	if s.StartAt > 0 && s.EndAt > 0 && s.StartAt > s.EndAt {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}

func (s *StatisticsCountQueryReqInfo) IsValid(key string) (bool, error) {
	keys := strings.Split(key, "-")
	if len(keys) != 4 {
		return false, nil
	}
	if s.OnlyMe && s.UserInfo.UserId != keys[1] {
		return false, nil
	}
	timestamp, err := strconv.ParseInt(keys[0], 10, 64)
	if err != nil {
		return false, err
	}
	if timestamp >= s.StartAt && (s.EndAt == 0 || timestamp < s.EndAt) {
		return true, nil
	}
	return false, nil
}
