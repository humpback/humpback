package types

type CountType string

var (
    CountTypeGroup   CountType = "group"
    CountTypeService CountType = "service"
    CountTypeNode    CountType = "node"
    CountTypeDeploy  CountType = "deploy"
)

type StatisticalCountInfo struct {
    Id       string    `json:"id"`
    CreateAt int64     `json:"createAt"`
    Type     CountType `json:"type"`
    Num      int       `json:"num"`
    UserId   string    `json:"userId"`
}
