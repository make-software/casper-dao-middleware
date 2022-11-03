package query

import sq "github.com/Masterminds/squirrel"

type DeleteBuilder struct {
	sq.DeleteBuilder
}

func Delete(from string) DeleteBuilder {
	return DeleteBuilder{
		sq.Delete(from),
	}
}
