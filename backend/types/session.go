package types

type Session struct {
	SessionId string `json:"sessionId"`
	UserId    string `json:"userId"`
	ExpiredAt int64  `json:"expiredAt"`
}
