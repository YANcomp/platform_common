package pgquery

import (
	"reflect"
	"strconv"
)

const (
	defaultLimit = 25
)

// Pagination query params
type PaginationQuery struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

func (q *PaginationQuery) IsEmpty() bool {
	return reflect.DeepEqual(q, &PaginationQuery{})
}

// Set limit
func (q *PaginationQuery) SetLimit(limitQuery string) error {
	if limitQuery == "" {
		q.Limit = defaultLimit
		return nil
	}
	n, err := strconv.Atoi(limitQuery)
	if err != nil {
		return err
	}
	q.Limit = n

	return nil
}

// Set offset
func (q *PaginationQuery) SetOffset(offsetQuery string) error {
	if offsetQuery == "" {
		q.Offset = 0
		return nil
	}
	n, err := strconv.Atoi(offsetQuery)
	if err != nil {
		return err
	}
	q.Offset = n

	return nil
}

// Get offset
func (q *PaginationQuery) GetOffset() int {
	if q.Offset == 0 {
		return 0
	}
	return q.Offset
}

// Get limit
func (q *PaginationQuery) GetLimit() int {
	return q.Limit
}
