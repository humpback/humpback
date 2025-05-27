package locales

import (
	"strings"
)

func GetMsg(language string, key string) string {
	msg := enUSMsg[key]
	switch strings.ToLower(language) {
	case "zh-cn":
		if v := zhCnMsg[key]; v != "" {
			msg = v
		}
	}
	return msg
}

const (
	CodeSucceed              = "S20000"
	CodeServerErr            = "R50000"
	CodeRequestParamsInvalid = "R40000"
	CodeLoginExpired         = "R40101"
	CodeNotLogin             = "R40102"
	CodeNoPermission         = "RNoPermission-001"

	CodeNameNotEmpty = "R4Name-001"

	CodeUserNameNotEmpty     = "R4User-001"
	CodeUserNameIsInvalid    = "R4User-002"
	CodeUserNameLimitLength  = "R4User-003"
	CodeUserNotExist         = "R4User-004"
	CodeUserIdNotEmpty       = "R4User-005"
	CodeUserAlreadyExist     = "R4User-006"
	CodeUserIsOwner          = "R4User-007"
	CodeUserNameAlreadyExist = "R4User-008"

	CodeUserRoleIsInvalid = "R4Role-001"

	CodeTeamNameNotEmpty     = "R4Team-001"
	CodeTeamNameLimitLength  = "R4Team-002"
	CodeTeamIdNotEmpty       = "R4Team-003"
	CodeTeamNotExist         = "R4Team-004"
	CodeTeamAlreadyExist     = "R4Team-005"
	CodeTeamNameAlreadyExist = "R4Team-006"

	CodeEmailIsInvalid   = "R4Email-001"
	CodeEmailLimitLength = "R4Email-002"

	CodePhoneLimitLength = "R4Phone-001"
	CodePhoneIsInvalid   = "R4Phone-002"

	CodeDescriptionLimitLength = "R4Desc-001"

	CodePasswordLimitLength    = "R4Psd-001"
	CodePasswordIsWrong        = "R4Psd-002"
	CodePasswordNotEmpty       = "R4Psd-003"
	CodeOldPasswordNotEmpty    = "R4Psd-004"
	CodeOldPasswordLimitLength = "R4Psd-005"
	CodeOldPasswordIsWrong     = "R4Psd-006"
	CodeNewPasswordNotEmpty    = "R4Psd-007"
	CodeNewPasswordLimitLength = "R4Psd-008"

	CodeConfigNameNotEmpty           = "R4Cfg-001"
	CodeConfigNameLimitLength        = "R4Cfg-002"
	CodeConfigNameAlreadyExist       = "R4Cfg-003"
	CodeConfigNotExist               = "R4Cfg-004"
	CodeConfigValueNotEmpty          = "R4Cfg-005"
	CodeConfigStaticValueLimitLength = "R4Cfg-006"
	CodeConfigVolumeValueLimitLength = "R4Cfg-007"
	CodeConfigTypeIsInvalid          = "R4Cfg-008"
	CodeConfigIdNotEmpty             = "R4Cfg-009"

	CodeRegistryNotExist            = "R4Registry-001"
	CodeRegistryIdNotEmpty          = "R4Registry-002"
	CodeRegistryUrlNotEmpty         = "R4Registry-003"
	CodeRegistryUrlLimitLength      = "R4Registry-004"
	CodeRegistryUrlAlreadyExist     = "R4Registry-005"
	CodeRegistryUsernameLimitLength = "R4Registry-006"
	CodeRegistryPasswordLimitLength = "R4Registry-007"
	CodeRegistryDefaultNotDelete    = "R4Registry-008"
	CodeRegistryDomainNotEmpty      = "R4Registry-009"

	CodeNodesNotEmpty              = "R4Nodes-001"
	CodeNodesIdNotEmpty            = "R4Nodes-002"
	CodeNodesNotExist              = "R4Nodes-003"
	CodeNodesIpAddressInvalid      = "R4Nodes-004"
	CodeNodesLabelKeyIsEmpty       = "R4Nodes-005"
	CodeNodesLabelValueIsEmpty     = "R4Nodes-006"
	CodeNodesLabelKeyIsDuplicated  = "R4Nodes-007"
	CodeNodesIpAddressAlreadyExist = "R4Nodes-008"

	CodeGroupIdNotEmpty       = "R4Group-001"
	CodeGroupNotExist         = "R4Group-NotExist"
	CodeGroupNameAlreadyExist = "R4Group-003"
	CodeGroupNameLimitLength  = "R4Group-004"
	CodeGroupNameNotEmpty     = "R4Group-005"
	CodeGroupNoPermission     = "R4Group-006"

	CodeServiceNotExist                          = "R4Service-NotExist"
	CodeServiceNameNotEmpty                      = "R4Service-002"
	CodeServiceNameLimitLength                   = "R4Service-003"
	CodeServiceNameAlreadyExist                  = "R4Service-004"
	CodeServiceIdNotEmpty                        = "R4Service-005"
	CodeServiceImageNotEmpty                     = "R4Service-006"
	CodeServiceImageLimitLength                  = "R4Service-007"
	CodeServiceNetworkModeInvalid                = "R4Service-008"
	CodeServiceRestartPolicyInvalid              = "R4Service-009"
	CodeServiceContainerPortNotEmpty             = "R4Service-010"
	CodeServiceContainerPortIsDuplicated         = "R4Service-011"
	CodeServiceHostPortIsDuplicated              = "R4Service-012"
	CodeServiceContainerVolumeNotEmpty           = "R4Service-013"
	CodeServiceContainerVolumeIsDuplicated       = "R4Service-014"
	CodeServiceContainerVolumeTypeInvlaid        = "R4Service-015"
	CodeServiceHostVolumeNotEmpty                = "R4Service-016"
	CodeServiceLabelNameNotEmpty                 = "R4Service-017"
	CodeServiceLabelValueNotEmpty                = "R4Service-018"
	CodeServiceResourceMemoryLimitLimitMax       = "R4Service-019"
	CodeServiceResourceMemoryReservationLimitMax = "R4Service-020"
	CodeServiceResourceMaxCpuUsageInvalid        = "R4Service-021"
	CodeServiceDispatchModeInvalid               = "R4Service-022"
	CodeServiceInstanceNumInvalid                = "R4Service-023"
	CodeServicePlacementModeInvalid              = "R4Service-024"
	CodeServicePlacementLabelNotEmpty            = "R4Service-025"
	CodeServicePlacementValueNotEmpty            = "R4Service-026"
	CodeServiceScheduleCronNotEmpty              = "R4Service-027"
	CodeServiceScheduleCronInvalid               = "R4Service-028"
	CodeServiceScheduleTimeoutInvalid            = "R4Service-029"
	CodeServiceProtocolInvalid                   = "R4Service-030"
	CodeServiceOperateInvalid                    = "R4Service-031"
	CodeServiceIsNotEnable                       = "R4Service-032"
	CodeServiceIsEnabled                         = "R4Service-033"
	CodeServiceIsDisabled                        = "R4Service-034"

	CodeContainerIdNotEmpty     = "R4Container-001"
	CodeContainerActionInvalid  = "R4Container-002"
	CodeContainerLogTimeInvalid = "R4Container-003"
	CodeContainerNotExist       = "R4Container-004"

	CodeActivityTypeInvalid  = "R4Activity-001"
	CodeActivityTypeNotEmpty = "R4Activity-002"
)
