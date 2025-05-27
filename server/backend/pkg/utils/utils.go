package utils

import (
    "encoding/json"
    "fmt"
    "time"
)

func NewActionTimestamp() int64 {
    return time.Now().UnixMilli()
}

func PrintJson(data any) {
    value, _ := json.MarshalIndent(data, "", "    ")
    fmt.Printf("%s\n", value)
}

func NewGroupId() string {
    return GenerateRandomStringWithLength(8)
}

func NewServiceId(groupId string) string {
    return fmt.Sprintf("%s%s", groupId, GenerateRandomStringWithLength(8))
}

func NewVersionId() string {
    return GenerateRandomStringWithLength(5)
}
