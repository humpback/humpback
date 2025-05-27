package utils

import (
	"regexp"
	"strings"
)

// MaskPhoneNumber 对电话进行掩码
func MaskPhoneNumber(phone string) string {
	if strings.TrimSpace(phone) == "" {
		return phone
	}
	re := regexp.MustCompile(`(\d{3})\d{4}(\d{4})`)
	return re.ReplaceAllString(phone, `$1****$2`)
}

// MaskEmail 对邮箱进行掩码
func MaskEmail(email string) string {
	if strings.TrimSpace(email) == "" {
		return email
	}
	parts := strings.Split(email, "@")
	if len(parts[0]) > 3 {
		masked := parts[0][:2] + strings.Repeat("*", len(parts[0])-3) + parts[0][len(parts[0])-1:]
		return masked + "@" + parts[1]
	}
	return email // 如果用户名部分长度不超过3，直接返回原始邮件地址
}
