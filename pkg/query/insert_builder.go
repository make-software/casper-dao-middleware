package query

import sq "github.com/Masterminds/squirrel"

type InsertBuilder struct {
	sq.InsertBuilder
}

func Insert(into string) InsertBuilder {
	return InsertBuilder{
		sq.Insert(into),
	}
}

func Replace(into string) InsertBuilder {
	return InsertBuilder{
		sq.Replace(into),
	}
}
