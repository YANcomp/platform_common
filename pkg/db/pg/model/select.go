package model

// Select query params
type SelectQuery struct {
	Fields []string `json:"fields,omitempty"`
}

// Set field to select
func (q *SelectQuery) SetField(filed string) error {
	if filed == "" {
		return nil
	}
	q.Fields = append(q.Fields, filed)

	return nil
}

// Get fields to select
func (q *SelectQuery) GetFields() []string {
	return q.Fields
}
