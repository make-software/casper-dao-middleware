package pagination

import (
	"sort"
)

type Result struct {
	ItemCount uint64      `json:"item_count"`
	PageCount uint64      `json:"page_count"`
	Data      interface{} `json:"data"`
}

func NewResult(itemCount uint64, limit uint64, data interface{}) *Result {
	pageCount := itemCount / limit
	if rem := itemCount % limit; rem != 0 {
		pageCount++
	}
	return &Result{
		ItemCount: itemCount,
		PageCount: pageCount,
		Data:      data,
	}
}

func SortAndPaginate[T JSONTagValueExtractor](collection []T, params *Params) *Result {
	jsonSorter := NewSortableJSONCollection(&collection, OrderCriteria{
		orderBy:        params.OrderBy[0],
		orderDirection: params.OrderDirection,
	}, nil)
	sort.Sort(jsonSorter)

	return Paginate(collection, params)
}

func Paginate[T any](collection []T, params *Params) *Result {
	collectionLimit := params.Offset() + params.PageSize
	if collectionLimit > uint64(len(collection)) {
		collectionLimit = uint64(len(collection))
	}

	result := make([]T, 0, collectionLimit)
	for i := params.Offset(); i < collectionLimit; i++ {
		result = append(result, collection[i])
	}
	return NewResult(uint64(len(collection)), params.PageSize, result)
}
