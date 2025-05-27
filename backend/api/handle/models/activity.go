package models

import (
	"strconv"
	"strings"

	"humpback/common/locales"
	"humpback/common/response"
	"humpback/common/verify"
	"humpback/types"
)

type ActivityQueryFilterInfo struct {
	StartAt   int64                `json:"startAt"`
	EndAt     int64                `json:"endAt"`
	GroupId   string               `json:"groupId"`
	Action    types.ActivityAction `json:"action"`
	Operator  string               `json:"operator"`
	Type      string               `json:"type"`
	ServiceId string               `json:"serviceId"`
}

type ActivityQueryReqInfo struct {
	UserInfo *types.User `json:"-"`
	types.QueryInfo
	FilterInfo *ActivityQueryFilterInfo `json:"-"`
}

func (a *ActivityQueryReqInfo) Check() error {
	a.CheckBase()
	a.FilterInfo = new(ActivityQueryFilterInfo)
	if len(a.Filter) == 0 {
		return response.NewBadRequestErr(locales.CodeActivityTypeNotEmpty)
	}
	if err := ParseMapToStructConvert(a.Filter, a.FilterInfo); err != nil {
		return err
	}
	if err := verify.CheckIsEmpty(a.FilterInfo.Type, locales.CodeActivityTypeNotEmpty); err != nil {
		return err
	}
	if a.FilterInfo.StartAt > 0 && a.FilterInfo.EndAt < a.FilterInfo.StartAt {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}

func (a *ActivityQueryReqInfo) IsValidTimeRange(timestampStr string) (bool, bool, error) {
	if a.FilterInfo.StartAt > 0 || a.FilterInfo.EndAt > 0 {
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			return false, false, err
		}
		if a.FilterInfo.StartAt > 0 && timestamp < a.FilterInfo.StartAt {
			return true, false, nil
		}
		if a.FilterInfo.EndAt > 0 && timestamp > a.FilterInfo.EndAt {
			return false, false, nil
		}
	}
	return false, true, nil
}

func (a *ActivityQueryReqInfo) IsValid(keys []string) bool {
	if len(keys) != 4 {
		return false
	}

	if a.FilterInfo.Operator != "" && a.FilterInfo.Operator != keys[1] {
		return false
	}

	if a.FilterInfo.Action != "" && a.FilterInfo.Action != types.ActivityAction(keys[2]) {
		return false
	}

	if a.FilterInfo.Type == "account" {
		if a.UserInfo.UserId != keys[3] {
			return false
		}
	}

	if a.FilterInfo.Type == "groups" {
		if a.FilterInfo.GroupId != "" && a.FilterInfo.GroupId != keys[3] {
			return false
		}
	}
	if a.FilterInfo.Type == "services" {
		if a.FilterInfo.GroupId != "" && !strings.HasPrefix(keys[3], a.FilterInfo.GroupId) {
			return false
		}
		if a.FilterInfo.ServiceId != "" && a.FilterInfo.ServiceId != keys[3] {
			return false
		}
	}

	return true
}

type ActivityAllQueryReqInfo struct {
	UserInfo *types.User `json:"-"`
	StartAt  int64       `json:"startAt"`
}

func (s *ActivityAllQueryReqInfo) Check() error {
	if s.StartAt <= 0 {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}

func (s *ActivityAllQueryReqInfo) IsValid(timestampStr, userId string) (bool, error) {
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false, err
	}
	if timestamp < s.StartAt {
		return false, nil
	}
	if types.IsUser(s.UserInfo.Role) && userId != s.UserInfo.UserId {
		return false, nil
	}
	return true, nil
}
