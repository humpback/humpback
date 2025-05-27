package response

import (
	"fmt"
	"net/http"
	"strings"

	"humpback/common/locales"
)

func NewBadRequestErr(code string, msg ...string) *ErrInfo {
	return &ErrInfo{
		StatusCode: http.StatusBadRequest,
		Code:       code,
		ErrMsg:     strings.Join(msg, ", "),
	}
}

func NewNoPermissionErr(code string, msg ...string) *ErrInfo {
	return &ErrInfo{
		StatusCode: http.StatusForbidden,
		Code:       code,
		ErrMsg:     strings.Join(msg, ", "),
	}
}

func NewRespUnauthorized(startup bool, msg ...string) *ErrInfo {
	if startup {
		return &ErrInfo{
			StatusCode: http.StatusUnauthorized,
			Code:       locales.CodeNotLogin,
			ErrMsg:     strings.Join(msg, ", "),
		}
	}
	return &ErrInfo{
		StatusCode: http.StatusUnauthorized,
		Code:       locales.CodeLoginExpired,
		ErrMsg:     strings.Join(msg, ", "),
	}
}

func NewRespServerErr(msg ...string) *ErrInfo {
	return &ErrInfo{
		StatusCode: http.StatusInternalServerError,
		Code:       locales.CodeServerErr,
		ErrMsg:     strings.Join(msg, ", "),
		BizCode:    locales.CodeServerErr,
		BizMsg:     strings.Join(msg, ", "),
	}
}

type ErrInfo struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	ErrMsg     string `json:"errMsg"`
	BizCode    string `json:"bizCode,omitempty"`
	BizMsg     string `json:"bizMsg,omitempty"`
}

func (e *ErrInfo) CopyToBizData() {
	e.BizCode = e.Code
	e.BizMsg = e.ErrMsg
}

func (e *ErrInfo) ClearBizData() {
	e.BizCode = ""
	e.BizMsg = ""
}

func (e *ErrInfo) ParseCodeMsg(language string) {
	msg := locales.GetMsg(language, e.Code)
	if msg != "" {
		e.ErrMsg = msg
	}
}

func (e *ErrInfo) ReplaceCode(code string) {
	e.Code = code
}

func (e *ErrInfo) String() string {
	if e.BizCode == "" {
		return fmt.Sprintf("StatusCode: %d, Code: %s, ErrMsg: %s", e.StatusCode, e.Code, e.ErrMsg)
	}
	return fmt.Sprintf("StatusCode: %d, Code: %s, BizCode: %s, ErrMsg: %s, BizMsg: %s", e.StatusCode, e.Code, e.BizCode, e.ErrMsg, e.BizMsg)
}

func (e *ErrInfo) Error() string {
	return fmt.Sprintf("StatusCode: %d, Code: %s, Msg: %s", e.StatusCode, e.Code, e.ErrMsg)
}
