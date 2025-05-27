package types

import (
	"cmp"
)

var (
	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

type QueryInfo struct {
	Keywords string         `json:"keywords"`
	Mode     string         `json:"mode"`
	Filter   map[string]any `json:"filter"`
	PageInfo *PageInfo      `json:"pageInfo"`
	SortInfo *SortInfo      `json:"sortInfo"`
}

func (q *QueryInfo) CheckBase() {
	if q.Filter == nil {
		q.Filter = make(map[string]any)
	}
	if q.PageInfo != nil {
		if q.PageInfo.Index < 1 {
			q.PageInfo.Index = 1
		}
		if q.PageInfo.Size < 1 {
			q.PageInfo.Size = 20
		}
	}
	if q.SortInfo != nil {
		if q.SortInfo.Field == "" {
			q.SortInfo = nil
		} else if q.SortInfo.Order != SortOrderAsc && q.SortInfo.Order != SortOrderDesc {
			q.SortInfo.Order = SortOrderAsc
		}
	}
}

type SortInfo struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

type PageInfo struct {
	Index int `json:"index"`
	Size  int `json:"size"`
}

func QueryPagination[T any](pageInfo *PageInfo, list []*T) []*T {
	if pageInfo == nil {
		return list
	}
	start := (pageInfo.Index - 1) * pageInfo.Size
	if start >= len(list) {
		return []*T{}
	}
	end := start + pageInfo.Size
	if end > len(list) {
		end = len(list)
	}
	return list[start:end]
}

func QuerySortOrder[T cmp.Ordered](order string, a, b T) int {
	if order == SortOrderAsc {
		return cmp.Compare(a, b)
	}
	return cmp.Compare(b, a)
}
