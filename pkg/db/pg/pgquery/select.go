package pgquery

import (
	"reflect"
	"strings"
)

// Select query params
type SelectQuery struct {
	Fields []string `json:"fields,omitempty"`
}

func (s *SelectQuery) IsEmpty() bool {
	return reflect.DeepEqual(s, &SelectQuery{})
}

// Set field to select
func (s *SelectQuery) SetField(filed string) error {
	if filed == "" {
		return nil
	}
	s.Fields = append(s.Fields, strings.ToLower(filed))

	return nil
}

// Get fields to select
func (s *SelectQuery) GetFields() []string {
	return s.Fields
}
