package enum

type LengthLimit struct {
	Min int `json:"Min"`
	Max int `json:"Max"`
}

var (
	LimitUsername          = LengthLimit{Min: 2, Max: 100}
	LimitTeamName          = LengthLimit{Min: 2, Max: 100}
	LimitEmail             = LengthLimit{Min: 0, Max: 200}
	LimitPassword          = LengthLimit{Min: 8, Max: 20}
	LimitPhone             = LengthLimit{Min: 0, Max: 11}
	LimitDescription       = LengthLimit{Min: 0, Max: 500}
	LimitConfigName        = LengthLimit{Min: 1, Max: 100}
	LimitConfigValue       = LengthLimit{Min: 1, Max: 10000}
	LimitRegistryUrl       = LengthLimit{Min: 1, Max: 200}
	LimitRegistryUsername  = LengthLimit{Min: 1, Max: 100}
	LimitRegistryPassword  = LengthLimit{Min: 1, Max: 100}
	LimitGroupName         = LengthLimit{Min: 1, Max: 100}
	LimitServiceName       = LengthLimit{Min: 1, Max: 100}
	LimitImageName         = LengthLimit{Min: 1, Max: 200}
	LimitMemoryLimit       = LengthLimit{Min: 0, Max: 20480}
	LimitMemoryReservation = LengthLimit{Min: 0, Max: 20480}
	LimitMaxCpuUsage       = LengthLimit{Min: 0, Max: 100}
	LimitInstanceNum       = LengthLimit{Min: 1, Max: 100}
	LimitLogsLine          = LengthLimit{Min: 0, Max: 10000}
)

var RuleLength = map[string]LengthLimit{
	"Username":          LimitUsername,
	"TeamName":          LimitTeamName,
	"Email":             LimitEmail,
	"Password":          LimitPassword,
	"Phone":             LimitPhone,
	"Description":       LimitDescription,
	"ConfigName":        LimitConfigName,
	"ConfigValue":       LimitConfigValue,
	"RegistryUrl":       LimitRegistryUrl,
	"RegistryUsername":  LimitRegistryUsername,
	"RegistryPassword":  LimitRegistryPassword,
	"GroupName":         LimitGroupName,
	"ServiceName":       LimitServiceName,
	"ImageName":         LimitImageName,
	"MemoryLimit":       LimitMemoryLimit,
	"MemoryReservation": LimitMemoryReservation,
	"MaxCpuUsage":       LimitMaxCpuUsage,
	"InstanceNum":       LimitInstanceNum,
	"LogsLine":          LimitLogsLine,
}

var (
	RegularEmail     = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	RegularPhone     = `^\d+$`
	RegularIpAddress = `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
)

var RuleFormat = map[string]string{
	"Email":     RegularEmail,
	"Phone":     RegularPhone,
	"IPAddress": RegularIpAddress,
}
