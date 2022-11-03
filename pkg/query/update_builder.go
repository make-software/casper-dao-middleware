package query

import sq "github.com/Masterminds/squirrel"

type UpdateBuilder struct {
	sq.UpdateBuilder
}

func Update(table string) UpdateBuilder {
	return UpdateBuilder{
		sq.Update(table),
	}
}
