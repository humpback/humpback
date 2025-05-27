package models

import (
	"encoding/json"

	"humpback/common/locales"
	"humpback/common/response"
)

func ParseMapToStructConvert(m map[string]any, obj interface{}) error {
	if len(m) == 0 {
		return nil
	}
	v, err := json.Marshal(m)
	if err != nil {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}

	if err = json.Unmarshal(v, obj); err != nil {
		return response.NewBadRequestErr(locales.CodeRequestParamsInvalid)
	}
	return nil
}

func removeEmptyStrings(arr []string) []string {
	result := make([]string, 0)
	for _, str := range arr {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}
