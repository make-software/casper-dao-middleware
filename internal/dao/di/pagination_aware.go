package di

import "casper-dao-middleware/pkg/pagination"

type PaginationParamsAware struct {
	params *pagination.Params
}

func (p *PaginationParamsAware) SetPaginationParams(params *pagination.Params) {
	p.params = params
}

func (p *PaginationParamsAware) GetPaginationParams() *pagination.Params {
	return p.params
}
