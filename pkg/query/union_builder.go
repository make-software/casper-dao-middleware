package query

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"casper-dao-middleware/pkg/pagination"

	sq "github.com/Masterminds/squirrel"
	"github.com/lann/builder"
)

type unionSelect struct {
	op       string // e.g. "UNION"
	selector sq.SelectBuilder
}

type unionData struct {
	Selects           []*unionSelect
	Limit             string
	Offset            uint64
	OrderBy           []string
	PlaceholderFormat sq.PlaceholderFormat
}

// UnionBuilder is a (rather hack) implementation of Unions for squirrel query builder. They
// currently don't offer this feature. When they do, this code should be trashed
type UnionBuilder builder.Builder

func Union(a sq.SelectBuilder, b sq.SelectBuilder) UnionBuilder {
	ub := UnionBuilder{}
	ub = ub.setFirstSelect(a)
	return ub.Union(b)
}

func (u UnionBuilder) ToSql() (sql string, args []interface{}, err error) {
	builderStruct := builder.GetStruct(u)

	data := builderStruct.(unionData)

	if len(data.Selects) == 0 {
		err = errors.New("require a minimum of 1 select clause in UnionBuilder")
		return
	}

	sqlBuf := &bytes.Buffer{}
	var selArgs []interface{}
	var selSql string

	for index, selector := range data.Selects {

		selSql, selArgs, err = selector.selector.ToSql()
		if err != nil {
			return
		}

		if index == 0 {
			sqlBuf.WriteString(selSql) // no operator for first select-clause
		} else {
			sqlBuf.WriteString(" " + selector.op + " ( " + selSql + " ) ")
		}

		args = append(args, selArgs...)
	}

	if len(data.OrderBy) > 0 {
		sqlBuf.WriteString(" ORDER BY ")
		sqlBuf.WriteString(strings.Join(data.OrderBy, ","))
	}

	if data.Limit != "" {
		sqlBuf.WriteString(" LIMIT ")
		sqlBuf.WriteString(data.Limit)
	}

	if data.Offset != 0 {
		sqlBuf.WriteString(" OFFSET ")
		sqlBuf.WriteString(strconv.FormatUint(data.Offset, 10))
	}

	sql = sqlBuf.String()
	return
}

func (u UnionBuilder) Union(selector sq.SelectBuilder) UnionBuilder {
	// use ? in children to prevent numbering issues
	selector = selector.PlaceholderFormat(sq.Question)

	return builder.Append(u, "Selects", &unionSelect{op: "UNION", selector: selector}).(UnionBuilder)
}

func (u UnionBuilder) setFirstSelect(selector sq.SelectBuilder) UnionBuilder {

	// copy the PlaceholderFormat value from children since we don't know what it should be
	value, _ := builder.Get(selector, "PlaceholderFormat")
	bld := u.setProp("PlaceholderFormat", value)

	// use ? in children to prevent numbering issues
	selector = selector.PlaceholderFormat(sq.Question)

	return builder.Append(bld, "Selects", &unionSelect{op: "", selector: selector}).(UnionBuilder)
}

func (u UnionBuilder) Limit(n uint64) UnionBuilder {
	return u.setProp("Limit", fmt.Sprintf("%d", n))
}

// Offset sets a OFFSET clause on the query.
func (u UnionBuilder) Offset(offset uint64) UnionBuilder {
	return u.setProp("Offset", offset)
}

func (u UnionBuilder) OrderBy(orderBys ...string) UnionBuilder {
	return u.setProp("OrderBy", orderBys)
}

func (u UnionBuilder) PlaceholderFormat(fmt sq.PlaceholderFormat) UnionBuilder {
	return u.setProp("PlaceholderFormat", fmt)
}

// Paginate apply pagination
func (u UnionBuilder) Paginate(params *pagination.Params, availableFields map[string]struct{}) UnionBuilder {
	u = u.Limit(params.PageSize).Offset(params.Offset())

	if len(params.OrderBy) > 0 {
		keys := make([]string, 0)
		for _, key := range params.OrderBy {
			if _, ok := availableFields[key]; ok {
				keys = append(keys, fmt.Sprintf("%s %s", key, params.OrderDirection))
			}
		}
		if len(keys) > 0 {
			u = u.OrderBy(keys...)
		}
	}

	return u
}

func init() {
	builder.Register(UnionBuilder{}, unionData{})
}

func (u UnionBuilder) setProp(key string, value interface{}) UnionBuilder {
	return builder.Set(u, key, value).(UnionBuilder)
}
