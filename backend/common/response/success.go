package response

import (
	"net/http"
)

type SucceedMsg struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

func NewRespSucceed() *SucceedMsg {
	return &SucceedMsg{
		StatusCode: http.StatusOK,
		Msg:        "Succeed",
	}
}

type QueryResult[T any] struct {
	Total int  `json:"total"`
	List  []*T `json:"list"`
}

func NewQueryResult[T any](total int, list []*T) *QueryResult[T] {
	return &QueryResult[T]{Total: total, List: list}
}
