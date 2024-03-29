package pgquery

import (
	"github.com/Masterminds/squirrel"
	"reflect"
)

const (
	// FilterTypeEQ Значение равно
	FilterTypeEQ FilterType = iota

	// FilterTypeNotEQ Значение не равно
	FilterTypeNotEQ

	// FilterTypeGTE Значение больше или равно
	FilterTypeGTE

	// FilterTypeGT Значение больше
	FilterTypeGT

	// FilterTypeLT Значение меньше
	FilterTypeLT

	// FilterTypeLTE Значение меньше или равно
	FilterTypeLTE

	// FilterTypeLike Значение может содержать в себе
	FilterTypeLike

	// FilterTypeNotLike Значение может не содержать в себе
	FilterTypeNotLike

	// FilterTypeILike Значение может содержать в себе (регистронезависимый)
	FilterTypeILike

	// FilterTypeNotILike Значение может не содержать в себе (регистронезависимый)
	FilterTypeNotILike
)

const (
	OperatorAnd Operator = iota
	OperatorOr
)

type (
	FilterType uint8

	Operator uint8

	Filter struct {
		column        string
		fType         FilterType
		value         any
		operator      Operator
		filters       []Filter
		join          bool
		joinMainTable string
		joinTable     string
		joinOn        string
		joinColumn    string
	}
)

// NewFilter Создание нового фильтра
func NewFilter(column string, ftype FilterType, value any) Filter {
	return Filter{
		column:   column,
		fType:    ftype,
		value:    value,
		operator: OperatorAnd,
		filters:  make([]Filter, 0),
		join:     false,
	}
}

// NewFilter Создание нового фильтра
func NewFilterWithJoin(column string, ftype FilterType, value any, joinMainTable string, joinTable string, joinOn string, joinColumn string) Filter {
	return Filter{
		column:        column,
		fType:         ftype,
		value:         value,
		operator:      OperatorAnd,
		filters:       make([]Filter, 0),
		join:          true,
		joinMainTable: joinMainTable,
		joinTable:     joinTable,
		joinOn:        joinOn,
		joinColumn:    joinColumn,
	}
}

func (f *Filter) IsEmpty() bool {
	return reflect.DeepEqual(f, &Filter{})
}

// SetOperator Установка оператора для связывания всех дополнительных фильтров
func (f *Filter) SetOperator(operator Operator) *Filter {
	f.operator = operator
	return f
}

// WithFilters Добавление дополнительных фильтров
func (f *Filter) WithFilters(filters ...Filter) *Filter {
	f.filters = append(f.filters, filters...)
	return f
}

func (f *Filter) condition() squirrel.Sqlizer {
	//if f.join {
	//	return squirrel.Eq{f.column: f.value}
	//}

	switch f.fType {
	case FilterTypeNotEQ:
		return squirrel.NotEq{f.column: f.value}
	case FilterTypeGTE:
		return squirrel.GtOrEq{f.column: f.value}
	case FilterTypeGT:
		return squirrel.Gt{f.column: f.value}
	case FilterTypeLT:
		return squirrel.Lt{f.column: f.value}
	case FilterTypeLTE:
		return squirrel.LtOrEq{f.column: f.value}
	case FilterTypeLike:
		return squirrel.Like{f.column: f.value}
	case FilterTypeNotLike:
		return squirrel.NotLike{f.column: f.value}
	case FilterTypeILike:
		return squirrel.ILike{f.column: f.value}
	case FilterTypeNotILike:
		return squirrel.NotILike{f.column: f.value}
	case FilterTypeEQ:
	default:
	}

	return squirrel.Eq{f.column: f.value}
}

func (f *Filter) getConditions() squirrel.Sqlizer {
	if len(f.filters) == 0 {
		return f.condition()
	}

	var conditions []squirrel.Sqlizer

	conditions = append(conditions, f.condition())

	for _, filter := range f.filters {
		conditions = append(conditions, filter.getConditions())
	}

	if f.operator == OperatorOr {
		return or(conditions)
	}

	return and(conditions)
}

// UseSelectBuilder Наполнение squirrel.SelectBuilder фильтрацией
func (f *Filter) UseSelectBuilder(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	query := builder

	if f.join {
		query = builder.Join(f.joinTable + " on " + f.joinTable + "." + f.joinOn + " = " + f.joinMainTable + "." + f.joinColumn)
	}

	for _, filter := range f.filters {
		if filter.join {
			query = query.Join(filter.joinTable + " on " + filter.joinTable + "." + filter.joinOn + " = " + filter.joinMainTable + "." + filter.joinColumn)
		}
	}

	query = query.Where(f.getConditions())

	return query
}

func and(conditions []squirrel.Sqlizer) squirrel.Sqlizer {
	result := squirrel.And{}
	for _, condition := range conditions {
		result = append(result, condition)
	}
	return result
}

func or(conditions []squirrel.Sqlizer) squirrel.Sqlizer {
	result := squirrel.Or{}
	for _, condition := range conditions {
		result = append(result, condition)
	}
	return result
}
