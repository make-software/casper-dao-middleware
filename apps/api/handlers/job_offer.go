package handlers

import (
	"net/http"

	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/job_offer"
	http_response "casper-dao-middleware/pkg/http-response"
	"casper-dao-middleware/pkg/pagination"
)

type JobOffer struct {
	entityManager persistence.EntityManager
}

func NewJobOffer(entityManager persistence.EntityManager) *JobOffer {
	return &JobOffer{
		entityManager: entityManager,
	}
}

// HandleGetJobOffers
//
//	@Summary	Return paginated list of votes for votingID
//
//	@Router		/job_offers [GET]
//
//	@Param		page			query		int			false	"Page number"											default(1)
//	@Param		page_size		query		string		false	"Number of items per page"								default(10)
//	@Param		order_direction	query		string		false	"Sorting direction"										Enums(ASC, DESC)		default(ASC)
//	@Param		order_by		query		[]string	false	"Comma-separated list of sorting fields (job_offer_id)"	collectionFormat(csv)	default(voting_id)
//
//	@Success	200				{object}	http_response.PaginatedResponse{data=entities.JobOffer}
//	@Failure	400,404,500		{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		BidEscrow
func (h *JobOffer) HandleGetJobOffers(w http.ResponseWriter, r *http.Request) {
	paginationParams := pagination.NewParamsFromRequest(r)

	getJobOffers := job_offer.NewGetJobOffers()
	getJobOffers.SetEntityManager(h.entityManager)
	getJobOffers.SetPaginationParams(paginationParams)

	//TODO: add includes of options bids

	http_response.FromFunction(getJobOffers.Execute, w, r)
}
