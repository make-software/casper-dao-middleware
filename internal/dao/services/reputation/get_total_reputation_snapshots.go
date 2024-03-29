package reputation

import (
	"github.com/make-software/casper-go-sdk/casper"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/pkg/pagination"
)

type GetTotalReputationSnapshots struct {
	di.EntityManagerAware
	di.PaginationParamsAware

	address *casper.Hash
}

func NewGetTotalReputationSnapshots() *GetTotalReputationSnapshots {
	return &GetTotalReputationSnapshots{}
}

func (s *GetTotalReputationSnapshots) SetAddress(address *casper.Hash) {
	s.address = address
}

func (s *GetTotalReputationSnapshots) Execute() (*pagination.Result, error) {
	filters := map[string]interface{}{}

	if s.address != nil {
		filters["address"] = *s.address
	}

	count, err := s.GetEntityManager().TotalReputationSnapshotRepository().Count(filters)
	if err != nil {
		return nil, err
	}

	snapshots, err := s.GetEntityManager().TotalReputationSnapshotRepository().Find(s.GetPaginationParams(), filters)
	if err != nil {
		return nil, err
	}

	return pagination.NewResult(count, s.GetPaginationParams().PageSize, snapshots), nil
}
