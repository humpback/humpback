package models

import (
	"regexp"
	"slices"
	"strings"
	"time"
	
	"humpback/common/enum"
	"humpback/common/locales"
	"humpback/common/response"
	"humpback/common/verify"
	"humpback/pkg/utils"
	"humpback/types"
	
	cronv3 "github.com/robfig/cron/v3"
)

type ServiceCreateReqInfo struct {
	GroupId     string `json:"-"`
	ServiceName string `json:"serviceName"`
	Description string `json:"description"`
}

func (s *ServiceCreateReqInfo) Check() error {
	if err := verify.CheckRequiredAndLengthLimit(s.ServiceName, enum.LimitServiceName.Min, enum.LimitServiceName.Max, locales.CodeServiceNameNotEmpty, locales.CodeServiceNameLimitLength); err != nil {
		return err
	}
	if err := verify.CheckLengthLimit(s.Description, enum.LimitDescription.Min, enum.LimitDescription.Max, locales.CodeDescriptionLimitLength); err != nil {
		return err
	}
	return nil
}

func (s *ServiceCreateReqInfo) NewServiceInfo() *types.Service {
	nowT := utils.NewActionTimestamp()
	return &types.Service{
		ServiceId:   utils.NewServiceId(s.GroupId),
		GroupId:     s.GroupId,
		ServiceName: s.ServiceName,
		Description: s.Description,
		Version:     utils.NewVersionId(),
		Action:      "",
		IsEnabled:   false,
		IsDelete:    false,
		Status:      types.ServiceStatusNotReady,
		Meta:        nil,
		Deployment:  nil,
		Containers:  make([]*types.ContainerStatus, 0),
		CreatedAt:   nowT,
		UpdatedAt:   nowT,
	}
}

type ServiceQueryFilterInfo struct {
	Status   string `json:"status"`
	Schedule string `json:"schedule"`
}

type ServiceQueryReqInfo struct {
	types.QueryInfo
	UserInfo   *types.User             `json:"-"`
	FilterInfo *ServiceQueryFilterInfo `json:"-"`
}

func (s *ServiceQueryReqInfo) Check() error {
	s.CheckBase()
	s.FilterInfo = new(ServiceQueryFilterInfo)
	if err := ParseMapToStructConvert(s.Filter, s.FilterInfo); err != nil {
		return err
	}

	if s.FilterInfo != nil {
		if s.FilterInfo.Status != "" && slices.Index([]string{
			types.SwitchEnabled,
			types.SwitchDisabled,
			types.ServiceStatusRunning,
			types.ServiceStatusNotReady,
			types.ServiceStatusFailed,
		}, s.FilterInfo.Status) == -1 {
			return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
		}
		if s.FilterInfo.Schedule != "" && slices.Index([]string{"Yes", "No"}, s.FilterInfo.Schedule) == -1 {
			return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
		}
	}

	if s.Keywords != "" && s.Mode != "keywords" {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}

func (s *ServiceQueryReqInfo) QueryFilter(services []*types.Service) []*types.Service {
	result := make([]*types.Service, 0)
	for _, service := range services {
		if s.filter(service) {
			result = append(result, service)
		}
	}
	s.sort(result)
	return result
}

func (s *ServiceQueryReqInfo) filter(service *types.Service) bool {
	if s.FilterInfo != nil {
		switch s.FilterInfo.Status {
		case types.SwitchEnabled:
			if !service.IsEnabled {
				return false
			}
		case types.SwitchDisabled:
			if service.IsEnabled {
				return false
			}
		case types.ServiceStatusRunning, types.ServiceStatusNotReady, types.ServiceStatusFailed:
			if !service.IsEnabled || service.Status != s.FilterInfo.Status {
				return false
			}
		}

		if s.FilterInfo.Schedule == "Yes" && (service.Deployment == nil || service.Deployment.Type != types.DeployTypeSchedule) {
			return false
		}

		if s.FilterInfo.Schedule == "No" && (service.Deployment == nil || service.Deployment.Type != types.DeployTypeBackground) {
			return false
		}
	}
	if strings.Contains(strings.ToLower(service.ServiceName), strings.ToLower(s.Keywords)) {
		return true
	}
	if service.Meta != nil && strings.Contains(strings.ToLower(strings.Join([]string{service.Meta.RegistryDomain, service.Meta.Image}, "/")), strings.ToLower(s.Keywords)) {
		return true
	}
	return false
}

func (s *ServiceQueryReqInfo) sort(list []*types.Service) []*types.Service {
	var sortField = []string{"serviceName", "updatedAt", "createdAt"}
	if s.SortInfo == nil || s.SortInfo.Field == "" || slices.Index(sortField, s.SortInfo.Field) == -1 {
		return list
	}
	slices.SortFunc(list, func(a, b *types.Service) int {
		switch s.SortInfo.Field {
		case "serviceName":
			return types.QuerySortOrder(s.SortInfo.Order, strings.ToLower(a.ServiceName), strings.ToLower(b.ServiceName))
		case "updatedAt":
			return types.QuerySortOrder(s.SortInfo.Order, a.UpdatedAt, b.UpdatedAt)
		case "createdAt":
			return types.QuerySortOrder(s.SortInfo.Order, a.CreatedAt, b.CreatedAt)
		}
		return 1
	})
	return list
}

const (
	ServiceUpdateBasicInfo   = "basic-info"
	ServiceUpdateApplication = "application"
	ServiceUpdateDeployment  = "deployment"
)

type ServiceUpdateReqInfo struct {
	Type           string                   `json:"type"`
	ServiceId      string                   `json:"serviceId"`
	GroupId        string                   `json:"groupId"`
	Data           any                      `json:"data"`
	Description    string                   `json:"-"`
	MetaInfo       *types.ServiceMetaDocker `json:"-"`
	DeploymentInfo *types.Deployment        `json:"-"`
}

func (s *ServiceUpdateReqInfo) Check() error {
	if err := verify.CheckIsEmpty(s.ServiceId, locales.CodeServiceIdNotEmpty); err != nil {
		return err
	}
	switch strings.ToLower(s.Type) {
	case ServiceUpdateBasicInfo:
		description, ok := s.Data.(string)
		if !ok {
			return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
		}
		s.Description = description
		if err := verify.CheckLengthLimit(s.Description, 0, enum.LimitDescription.Max, locales.CodeDescriptionLimitLength); err != nil {
			return err
		}
	case ServiceUpdateApplication:
		data, ok := s.Data.(map[string]any)
		if !ok {
			return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
		}
		metaInfo := new(types.ServiceMetaDocker)
		if err := ParseMapToStructConvert(data, metaInfo); err != nil {
			return err
		}
		s.MetaInfo = metaInfo
		if err := s.checkMetaInfo(); err != nil {
			return err
		}
	case ServiceUpdateDeployment:
		data, ok := s.Data.(map[string]any)
		if !ok {
			return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
		}
		deploymentInfo := new(types.Deployment)
		if err := ParseMapToStructConvert(data, deploymentInfo); err != nil {
			return err
		}
		s.DeploymentInfo = deploymentInfo
		if err := s.checkDeployment(); err != nil {
			return err
		}
	default:
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}

func (s *ServiceUpdateReqInfo) checkMetaInfo() error {
	if s.MetaInfo == nil {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	info := s.MetaInfo
	re := regexp.MustCompile(`/+`)
	info.Image = strings.Trim(re.ReplaceAllString(info.Image, "/"), "/")
	if err := verify.CheckIsEmpty(info.Image, locales.CodeServiceImageNotEmpty); err != nil {
		return err
	}
	if err := verify.CheckIsEmpty(info.RegistryDomain, locales.CodeRegistryDomainNotEmpty); err != nil {
		return err
	}
	if err := verify.CheckIsEmpty(info.RegistryId, locales.CodeRegistryIdNotEmpty); err != nil {
		return err
	}
	if err := verify.CheckRequiredAndLengthLimit(info.Image, enum.LimitImageName.Min, enum.LimitImageName.Max, locales.CodeServiceImageNotEmpty, locales.CodeServiceImageLimitLength); err != nil {
		return err
	}
	info.EnvConfig = removeEmptyStrings(info.EnvConfig)
	if len(info.Labels) == 0 {
		info.Labels = make(map[string]string)
	}
	for k, v := range info.Labels {
		if strings.TrimSpace(k) == "" {
			return response.NewBadRequestErr(locales.CodeServiceLabelNameNotEmpty)
		}
		if strings.TrimSpace(v) == "" {
			return response.NewBadRequestErr(locales.CodeServiceLabelValueNotEmpty)
		}
	}
	if info.Capabilities == nil {
		info.Capabilities = &types.Capabilities{
			CapAdd:  make([]string, 0),
			CapDrop: make([]string, 0),
		}
	}
	info.Capabilities.CapAdd = removeEmptyStrings(info.Capabilities.CapAdd)
	info.Capabilities.CapDrop = removeEmptyStrings(info.Capabilities.CapDrop)
	if info.LogConfig == nil {
		info.LogConfig = &types.ServiceLogConfig{
			Type:   "",
			Config: make(map[string]string),
		}
	}
	if info.LogConfig.Config == nil {
		info.LogConfig.Config = make(map[string]string)
	}

	if info.Resources == nil {
		info.Resources = &types.ServiceResources{
			Memory:            0,
			MemoryReservation: 0,
			MaxCpuUsage:       0,
		}
	}
	if info.Resources.Memory < uint64(enum.LimitMemoryLimit.Min) || info.Resources.Memory > uint64(enum.LimitMemoryLimit.Max) {
		return response.NewBadRequestErr(locales.CodeServiceResourceMemoryLimitLimitMax)
	}
	if info.Resources.MemoryReservation < uint64(enum.LimitMemoryReservation.Min) || info.Resources.MemoryReservation > uint64(enum.LimitMemoryReservation.Max) {
		return response.NewBadRequestErr(locales.CodeServiceResourceMemoryReservationLimitMax)
	}
	if info.Resources.MaxCpuUsage < uint64(enum.LimitMaxCpuUsage.Min) || info.Resources.MaxCpuUsage > uint64(enum.LimitMaxCpuUsage.Max) {
		return response.NewBadRequestErr(locales.CodeServiceResourceMaxCpuUsageInvalid)
	}
	if len(info.Volumes) == 0 {
		info.Volumes = make([]*types.ServiceVolume, 0)
	}
	volumeMap := make(map[string]bool)
	for _, volume := range info.Volumes {
		if volume.Type != types.ServiceVolumeTypeBind && volume.Type != types.ServiceVolumeTypeVolume {
			return response.NewBadRequestErr(locales.CodeServiceContainerVolumeTypeInvlaid)
		}
		if strings.TrimSpace(volume.Target) == "" {
			return response.NewBadRequestErr(locales.CodeServiceContainerVolumeNotEmpty)
		}
		if strings.TrimSpace(volume.Source) == "" {
			return response.NewBadRequestErr(locales.CodeServiceHostVolumeNotEmpty)
		}
		if volumeMap[volume.Target] {
			return response.NewBadRequestErr(locales.CodeServiceContainerVolumeIsDuplicated)
		}
		volumeMap[volume.Target] = true
	}

	if info.Network == nil {
		info.Network = &types.NetworkInfo{
			Mode:               types.NetworkModeHost,
			Hostname:           "",
			NetworkName:        "",
			UseMachineHostname: false,
			Ports:              make([]*types.PortInfo, 0),
		}
	}
	if info.Network.Ports == nil {
		info.Network.Ports = make([]*types.PortInfo, 0)
	}

	if info.Network.Mode != types.NetworkModeHost && info.Network.Mode != types.NetworkModeBridge && info.Network.Mode != types.NetworkModeCustom {
		return response.NewBadRequestErr(locales.CodeServiceNetworkModeInvalid)
	}
	containerPorts := make(map[uint64]bool)
	hostPorts := make(map[uint64]bool)
	for _, p := range info.Network.Ports {
		if p.Protocol != "TCP" && p.Protocol != "UDP" {
			return response.NewBadRequestErr(locales.CodeServiceProtocolInvalid)
		}
		if p.ContainerPort <= 0 {
			return response.NewBadRequestErr(locales.CodeServiceContainerPortNotEmpty)
		}
		if p.HostPort < 0 {
			p.HostPort = 0
		}
		if containerPorts[p.ContainerPort] {
			return response.NewBadRequestErr(locales.CodeServiceContainerPortIsDuplicated)
		} else {
			containerPorts[p.ContainerPort] = true
		}
		if p.HostPort > 0 {
			if hostPorts[p.ContainerPort] {
				return response.NewBadRequestErr(locales.CodeServiceHostPortIsDuplicated)
			}
			hostPorts[p.ContainerPort] = true
		}
	}

	if info.RestartPolicy == nil {
		info.RestartPolicy = &types.RestartPolicy{
			Mode:          types.RestartPolicyModeNo,
			MaxRetryCount: 0,
		}
	}
	if info.RestartPolicy.Mode != types.RestartPolicyModeNo &&
		info.RestartPolicy.Mode != types.RestartPolicyModeAlways &&
		info.RestartPolicy.Mode != types.RestartPolicyModeOnFail &&
		info.RestartPolicy.Mode != types.RestartPolicyModeUnlessStopped {
		return response.NewBadRequestErr(locales.CodeServiceRestartPolicyInvalid)
	}
	if info.RestartPolicy.Mode != types.RestartPolicyModeOnFail {
		info.RestartPolicy.MaxRetryCount = 0
	}
	s.MetaInfo = info
	return nil
}

func (s *ServiceUpdateReqInfo) checkDeployment() error {
	if s.DeploymentInfo == nil {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	info := s.DeploymentInfo
	if info.Mode != types.DeployModeGlobal && info.Mode != types.DeployModeReplicate {
		return response.NewBadRequestErr(locales.CodeServiceDispatchModeInvalid)
	}
	if info.Mode == types.DeployModeGlobal {
		info.Replicas = 1
	}
	if info.Replicas < enum.LimitInstanceNum.Min || info.Replicas > enum.LimitInstanceNum.Max {
		return response.NewBadRequestErr(locales.CodeServiceInstanceNumInvalid)
	}
	if info.Placements == nil {
		info.Placements = make([]*types.PlacementInfo, 0)
	}
	for _, placement := range info.Placements {
		if placement.Mode != types.PlacementModeIP && placement.Mode != types.PlacementModeLabel {
			return response.NewBadRequestErr(locales.CodeServicePlacementModeInvalid)
		}
		if placement.Mode == types.PlacementModeLabel && strings.TrimSpace(placement.Key) == "" {
			return response.NewBadRequestErr(locales.CodeServicePlacementLabelNotEmpty)
		}
		if strings.TrimSpace(placement.Value) == "" {
			return response.NewBadRequestErr(locales.CodeServicePlacementValueNotEmpty)
		}
	}
	if info.Schedule == nil || info.Schedule.Rules == nil {
		info.Schedule = &types.ScheduleInfo{
			Timeout: "",
			Rules:   make([]string, 0),
		}
	}
	for _, cron := range info.Schedule.Rules {
		if strings.TrimSpace(cron) == "" {
			return response.NewBadRequestErr(locales.CodeServiceScheduleCronNotEmpty)
		}
		if _, err := cronv3.ParseStandard(cron); err != nil {
			return response.NewBadRequestErr(locales.CodeServiceScheduleCronInvalid)
		}
	}
	if len(info.Schedule.Rules) > 0 && info.Schedule.Timeout != "" {
		_, err := time.ParseDuration(info.Schedule.Timeout)
		if err != nil {
			return response.NewBadRequestErr(locales.CodeServiceScheduleTimeoutInvalid)
		}
	} else {
		info.Schedule.Timeout = ""
	}
	if len(info.Schedule.Rules) > 0 {
		info.ManualExec = false
	}
	s.DeploymentInfo = info
	return nil
}

type ServiceOperateReqInfo struct {
	ServiceId string `json:"serviceId"`
	GroupId   string `json:"-"`
	Action    string `json:"action"`
}

func (s *ServiceOperateReqInfo) Check() error {
	if err := verify.CheckIsEmpty(s.ServiceId, locales.CodeServiceIdNotEmpty); err != nil {
		return err
	}
	var actions = []string{
		types.ServiceActionStart,
		types.ServiceActionStop,
		types.ServiceActionRestart,
		types.ServiceActionEnable,
		types.ServiceActionDisable,
	}
	if slices.Index(actions, s.Action) == -1 {
		return response.NewBadRequestErr(locales.CodeServiceOperateInvalid)
	}
	return nil
}

type ServiceCloneReqInfo struct {
	ServiceId   string `json:"serviceId"`
	GroupId     string `json:"-"`
	NewGroupId  string `json:"newGroupId"`
	ServiceName string `json:"serviceName"`
	Description string `json:"description"`
}

func (s *ServiceCloneReqInfo) Check() error {
	if err := verify.CheckIsEmpty(s.ServiceId, locales.CodeServiceIdNotEmpty); err != nil {
		return err
	}
	if err := verify.CheckRequiredAndLengthLimit(s.ServiceName, enum.LimitServiceName.Min, enum.LimitServiceName.Max, locales.CodeServiceNameNotEmpty, locales.CodeServiceNameLimitLength); err != nil {
		return err
	}
	if err := verify.CheckIsEmpty(s.NewGroupId, locales.CodeGroupIdNotEmpty); err != nil {
		return err
	}
	if err := verify.CheckLengthLimit(s.Description, enum.LimitDescription.Min, enum.LimitDescription.Max, locales.CodeDescriptionLimitLength); err != nil {
		return err
	}
	return nil
}
