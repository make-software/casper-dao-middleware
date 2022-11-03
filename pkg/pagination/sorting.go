package pagination

import "log"

type JSONTagValueExtractor interface {
	GetValueByJSONTag(tag string) interface{}
}

type OrderCriteria struct {
	orderBy        string
	orderDirection OrderDirection
}

func NewOrderCriteria(orderBy string, direction OrderDirection) OrderCriteria {
	return OrderCriteria{
		orderBy:        orderBy,
		orderDirection: direction,
	}
}

type SortableJSONCollection[T JSONTagValueExtractor] struct {
	data          *[]T
	mainCriteria  OrderCriteria
	extraCriteria *OrderCriteria
}

func NewSortableJSONCollection[T JSONTagValueExtractor](data *[]T, main OrderCriteria, extra *OrderCriteria) *SortableJSONCollection[T] {
	return &SortableJSONCollection[T]{
		data:          data,
		mainCriteria:  main,
		extraCriteria: extra,
	}
}

func (s *SortableJSONCollection[T]) Len() int {
	return len(*s.data)
}

func (s *SortableJSONCollection[T]) Less(i, j int) bool {
	sortable := *s.data

	sortablei := s.getComparableValue(sortable[i], s.mainCriteria.orderBy)
	sortablej := s.getComparableValue(sortable[j], s.mainCriteria.orderBy)

	if s.extraCriteria != nil && sortablei == sortablej {
		sortablei = s.getComparableValue(sortable[i], s.extraCriteria.orderBy)
		sortablej = s.getComparableValue(sortable[j], s.extraCriteria.orderBy)

		if s.extraCriteria.orderDirection == OrderDirectionDESC {
			return sortablei > sortablej
		}
		return sortablei < sortablej
	}

	if s.mainCriteria.orderDirection == OrderDirectionDESC {
		return sortablei > sortablej
	}
	return sortablei < sortablej
}

func (s *SortableJSONCollection[T]) Swap(i, j int) {
	sortable := *s.data
	sortable[i], sortable[j] = sortable[j], sortable[i]
}

func (s *SortableJSONCollection[T]) getComparableValue(item T, field string) float64 {
	switch casted := item.GetValueByJSONTag(field).(type) {
	case float64:
		return casted
	case float32:
		return float64(casted)
	case int64:
		return float64(casted)
	case uint64:
		return float64(casted)
	case int32:
		return float64(casted)
	case uint32:
		return float64(casted)
	case int:
		return float64(casted)
	case uint:
		return float64(casted)
	case uint8:
		return float64(casted)
	case int8:
		return float64(casted)
	case bool:
		if casted {
			return 1
		}
		return 0
	case string:
		// TODO: extend to sort alphabetically
		log.Fatal("Can support sorting by string value")
		return 0
	default:
		return 0
	}
}
