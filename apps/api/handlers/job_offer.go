package handlers

import (
	"net/http"

	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/bid"
	"casper-dao-middleware/internal/dao/services/job_offer"
	"casper-dao-middleware/internal/dao/services/jobs"
	http_params "casper-dao-middleware/pkg/http-params"
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

	http_response.FromFunction(getJobOffers.Execute, w, r)
}

// HandleGetJobOfferBids
//
//	@Summary	Return paginated list of bid for JobOffer
//
//	@Router		/job_offers/{job_offer_id}/bids [GET]
//
//	@Param		job_offer_id	path		uint		true	"JobOfferID uint"
//	@Param		page			query		int			false	"Page number"											default(1)
//	@Param		page_size		query		string		false	"Number of items per page"								default(10)
//	@Param		order_direction	query		string		false	"Sorting direction"										Enums(ASC, DESC)		default(ASC)
//	@Param		order_by		query		[]string	false	"Comma-separated list of sorting fields (job_offer_id)"	collectionFormat(csv)	default(voting_id)
//
//	@Success	200				{object}	http_response.PaginatedResponse{data=entities.Bid}
//	@Failure	400,404,500		{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		BidEscrow
func (h *JobOffer) HandleGetJobOfferBids(w http.ResponseWriter, r *http.Request) {
	jobOfferID, err := http_params.ParseUint32("job_offer_id", r)
	if err != nil {
		http_response.Error(w, r, err)
		return
	}

	paginationParams := pagination.NewParamsFromRequest(r)

	getBids := bid.NewGetBids()
	getBids.SetEntityManager(h.entityManager)
	getBids.SetPaginationParams(paginationParams)
	getBids.SetJobOfferID(jobOfferID)

	//TODO: make sense to add job offer including

	http_response.FromFunction(getBids.Execute, w, r)
}

// HandleGetBidJob
//
//	@Summary	Return Job by BidID
//
//	@Router		/bids/{bid_id}/job [GET]
//
//	@Param		bid_id		path		uint	true	"BidID uint"
//
//	@Success	200			{object}	http_response.SuccessResponse{data=entities.Job}
//	@Failure	400,404,500	{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		BidEscrow
func (h *JobOffer) HandleGetBidJob(w http.ResponseWriter, r *http.Request) {
	bidID, err := http_params.ParseUint32("bid_id", r)
	if err != nil {
		http_response.Error(w, r, err)
		return
	}

	getJob := jobs.NewGetJobByBid()
	getJob.SetEntityManager(h.entityManager)
	getJob.SetBidID(bidID)

	http_response.FromFunction(getJob.Execute, w, r)
}
