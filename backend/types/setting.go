package types

const (
	NodeStatusOnline  = "Online"
	NodeStatusOffline = "Offline"

	SwitchEnabled  = "Enabled"
	SwitchDisabled = "Disabled"
)

type QueryRegistry struct {
	HasAuth bool `json:"hasAuth"`
	*Registry
}

type Registry struct {
	RegistryId string `json:"registryId"`
	URL        string `json:"url"`
	IsDefault  bool   `json:"isDefault"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	CreatedAt  int64  `json:"createdAt"`
	UpdatedAt  int64  `json:"updatedAt"`
}

type ConfigType int

var (
	ConfigTypeStatic ConfigType = 1
	ConfigTypeVolume ConfigType = 2
)

var (
	ConfigTypesMap = map[ConfigType]string{
		ConfigTypeStatic: "Static",
		ConfigTypeVolume: "Volume",
	}
)

type Config struct {
	ConfigId    string     `json:"configId"`
	ConfigName  string     `json:"configName"`
	Description string     `json:"description"`
	ConfigType  ConfigType `json:"configType"`
	ConfigValue string     `json:"configValue"`
	CreatedAt   int64      `json:"createdAt"`
	UpdatedAt   int64      `json:"updatedAt"`
}

type Template struct {
	TemplateID   string      `json:"templateId"`
	TemplateName string      `json:"templateName"`
	Content      interface{} `json:"content"`
	CreateAt     int64       `json:"createAt"`
	UpdateAt     int64       `json:"updateAt"`
}
