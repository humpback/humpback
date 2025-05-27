package verify

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"humpback/common/enum"
	"humpback/common/locales"
	"humpback/common/response"
)

func IsValidEmail(email string) bool {
	return regexp.MustCompile(enum.RegularEmail).MatchString(email)
}

func CheckIsEmpty(value string, code string) error {
	if strings.TrimSpace(value) == "" {
		return response.NewBadRequestErr(code)
	}
	return nil
}

func CheckRequiredAndLengthLimit(v string, min, max int, requiredCode, lengthCode string) error {
	if err := CheckIsEmpty(v, requiredCode); err != nil {
		return err
	}
	return CheckLengthLimit(v, min, max, lengthCode)
}

func CheckLengthLimit(v string, min, max int, lengthCode string) error {
	if (min > 0 && utf8.RuneCountInString(v) < min) || (max > 0 && utf8.RuneCountInString(v) > max) {
		return response.NewBadRequestErr(lengthCode)
	}
	return nil
}

func CheckEmail(email string) error {
	if !IsValidEmail(email) {
		return response.NewBadRequestErr(locales.CodeEmailIsInvalid)
	}
	return nil
}

func CheckPhone(phone string) error {
	if !regexp.MustCompile(enum.RegularPhone).MatchString(phone) {
		return response.NewBadRequestErr(locales.CodePhoneIsInvalid)
	}
	return nil
}
