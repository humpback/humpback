package models

import (
	"slices"
	"strings"

	"humpback/common/enum"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/common/verify"
	"humpback/types"
)

type InstanceOperateReqInfo struct {
	GroupId     string `json:"-"`
	ServiceId   string `json:"-"`
	ContainerId string `json:"containerId"`
	NodeId      string `json:"nodeId"`
	Action      string `json:"action"`
}

func (g *InstanceOperateReqInfo) Check() error {
	if err := verify.CheckIsEmpty(g.NodeId, locales.CodeNodesIdNotEmpty); err != nil {
		return err
	}
	if err := verify.CheckIsEmpty(g.ContainerId, locales.CodeContainerIdNotEmpty); err != nil {
		return err
	}
	var actions = []string{strings.ToLower(types.ServiceActionStart), strings.ToLower(types.ServiceActionRestart), strings.ToLower(types.ServiceActionStop)}
	if slices.Index(actions, strings.ToLower(g.Action)) == -1 {
		return response.NewBadRequestErr(locales.CodeContainerActionInvalid)
	}
	return nil
}

type InstanceLogsReqInfo struct {
	GroupId       string `json:"-"`
	ServiceId     string `json:"-"`
	ContainerId   string `json:"containerId"`
	NodeId        string `json:"nodeId"`
	Line          uint   `json:"line"`
	StartAt       int64  `json:"startAt"`
	EndAt         int64  `json:"endAt"`
	ShowTimestamp bool   `json:"showTimestamp"`
}

func (g *InstanceLogsReqInfo) Check() error {
	if err := verify.CheckIsEmpty(g.ContainerId, locales.CodeContainerIdNotEmpty); err != nil {
		return err
	}
	if err := verify.CheckIsEmpty(g.NodeId, locales.CodeNodesIdNotEmpty); err != nil {
		return err
	}
	if g.StartAt > g.EndAt {
		return response.NewBadRequestErr(locales.CodeContainerLogTimeInvalid)
	}
	if g.Line > uint(enum.LimitLogsLine.Max) {
		g.Line = uint(enum.LimitLogsLine.Max)
	}
	return nil
}

type InstanceStatsReqInfo struct {
	ContainerId string `json:"containerId"`
	NodeId      string `json:"nodeId"`
}

type InstancesPerformanceReqInfo struct {
	GroupId    string                  `json:"-"`
	ServiceId  string                  `json:"-"`
	Containers []*InstanceStatsReqInfo `json:"containers"`
}

func (i InstancesPerformanceReqInfo) Check() error {
	for _, info := range i.Containers {
		if info.ContainerId == "" {
			return response.NewBadRequestErr(locales.CodeContainerIdNotEmpty)
		}
		if info.NodeId == "" {
			return response.NewBadRequestErr(locales.CodeNodesIdNotEmpty)
		}
	}
	if len(i.Containers) == 0 {
		response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}
