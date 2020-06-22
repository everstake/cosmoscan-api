package smodels

type PaginatableResponse struct {
	Items interface{} `json:"items"`
	Total uint64      `json:"total"`
}
