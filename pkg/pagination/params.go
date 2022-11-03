package pagination

import (
	"net/http"
	"strconv"
	"strings"
)

type OrderDirection string

var (
	OrderDirectionASC  OrderDirection = "asc"
	OrderDirectionDESC OrderDirection = "desc"
)

func NewOrderDirection(direction string) OrderDirection {
	orderDirection := OrderDirectionDESC
	if strings.ToLower(direction) == string(OrderDirectionASC) {
		orderDirection = OrderDirectionASC
	}
	return orderDirection
}

type Params struct {
	OrderDirection OrderDirection
	OrderBy        []string
	Page           uint32
	PageSize       uint64
}

func (p Params) Offset() uint64 {
	return (uint64(p.Page) - 1) * p.PageSize
}

func (p *Params) SetDefaultOrder(orderBy string, direction OrderDirection) {
	if len(p.OrderBy) == 0 {
		p.OrderBy = append(p.OrderBy, orderBy)
		p.OrderDirection = direction
	}
}

func NewParamsFromRequest(r *http.Request) *Params {
	// TODO: consider situation to return error on parsing pagination params
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	return &Params{
		Page:           uint32(page),
		PageSize:       uint64(pageSize),
		OrderDirection: NewOrderDirection(r.URL.Query().Get("order_direction")),
		OrderBy:        strings.Split(r.URL.Query().Get("order_by"), ","),
	}
}
