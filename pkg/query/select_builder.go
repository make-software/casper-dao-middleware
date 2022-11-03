package query

import (
	"fmt"

	"casper-dao-middleware/pkg/pagination"

	sq "github.com/Masterminds/squirrel"
)

type SelectBuilder struct {
	inner sq.SelectBuilder
}

func Select(columns ...string) *SelectBuilder {
	return &SelectBuilder{
		inner: sq.Select(columns...),
	}
}

func (b *SelectBuilder) Where(pred interface{}, args ...interface{}) *SelectBuilder {
	b.inner = b.inner.Where(pred, args)
	return b
}

// From sets the FROM clause of the query.
func (b *SelectBuilder) From(from string) *SelectBuilder {
	b.inner = b.inner.From(from)
	return b
}

// Limit sets a LIMIT clause on the query.
func (b *SelectBuilder) Limit(limit uint64) *SelectBuilder {
	b.inner = b.inner.Limit(limit)
	return b
}

// Offset sets a OFFSET clause on the query.
func (b *SelectBuilder) Offset(offset uint64) *SelectBuilder {
	b.inner = b.inner.Offset(offset)
	return b
}

// OrderBy adds ORDER BY expressions to the query.
func (b *SelectBuilder) OrderBy(orderBys ...string) *SelectBuilder {
	b.inner = b.inner.OrderBy(orderBys...)
	return b
}

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b *SelectBuilder) PlaceholderFormat(f sq.PlaceholderFormat) *SelectBuilder {
	b.inner = b.inner.PlaceholderFormat(f)
	return b
}

// Paginate apply pagination
func (b *SelectBuilder) Paginate(params *pagination.Params, availableFields map[string]struct{}) *SelectBuilder {
	b.inner = b.inner.Limit(params.PageSize)
	b.inner = b.inner.Offset(params.Offset())

	if len(params.OrderBy) > 0 {
		keys := make([]string, 0)
		for _, key := range params.OrderBy {
			if _, ok := availableFields[key]; ok {
				keys = append(keys, fmt.Sprintf("%s %s", key, params.OrderDirection))
			}
		}
		if len(keys) > 0 {
			b.inner = b.inner.OrderBy(keys...)
		}
	}

	return b
}

// FilterBy apply filters to select builder.
func (b *SelectBuilder) FilterBy(filters map[string]interface{}, availableFields map[string]struct{}) *SelectBuilder {
	for key := range availableFields {
		if filter, ok := filters[key]; ok {
			b.inner = b.inner.Where(sq.Eq{key: filter})
		}
	}
	return b
}

// GroupBy adds GROUP BY expressions to the query.
func (b *SelectBuilder) GroupBy(groupBys ...string) *SelectBuilder {
	b.inner = b.inner.GroupBy(groupBys...)
	return b
}

// ToSql builds the query into a SQL string and bound args.
// nolint
func (b *SelectBuilder) ToSql() (string, []interface{}, error) {
	return b.inner.ToSql()
}
