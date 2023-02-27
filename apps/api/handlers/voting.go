package handlers

import (
	"net/http"

	"casper-dao-middleware/apps/api/serialization"
	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/voting"
	http_params "casper-dao-middleware/pkg/http-params"
	http_response "casper-dao-middleware/pkg/http-response"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/serialize"
)

type Voting struct {
	entityManager persistence.EntityManager
}

func NewVoting(entityManager persistence.EntityManager) *Voting {
	return &Voting{
		entityManager: entityManager,
	}
}

// HandleGetVotingVotes
//
//	@Summary	Return paginated list of votes for votingID
//
//	@Router		/votings/{voting_id}/votes [GET]
//
//	@Param		voting_id		query		string		false	"Comma-separated list of VotingIDs (number)"
//	@Param		includes		query		string		false	"Optional fields' schema (voting{})"
//	@Param		page			query		int			false	"Page number"													default(1)
//	@Param		page_size		query		string		false	"Number of items per page"										default(10)
//	@Param		order_direction	query		string		false	"Sorting direction"												Enums(ASC, DESC)		default(ASC)
//	@Param		order_by		query		[]string	false	"Comma-separated list of sorting fields (voting_id, address)"	collectionFormat(csv)	default(voting_id)
//
//	@Success	200				{object}	http_response.PaginatedResponse{data=entities.Vote}
//	@Failure	400,404,500		{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		Vote
func (h *Voting) HandleGetVotingVotes(w http.ResponseWriter, r *http.Request) {
	votingID, err := http_params.ParseUint32("voting_id", r)
	if err != nil {
		http_response.Error(w, r, err)
		return
	}

	includes, err := http_params.ParseOptionalData("includes", r)
	if err != nil {
		http_response.Error(w, r, err)
		return
	}

	paginationParams := pagination.NewParamsFromRequest(r)

	getVotes := voting.NewGetVotes()
	getVotes.SetVotingIDs([]uint32{votingID})
	getVotes.SetEntityManager(h.entityManager)
	getVotes.SetPaginationParams(paginationParams)

	paginatedVotes, err := getVotes.Execute()
	if err != nil {
		http_response.Error(w, r, err)
		return
	}

	votesJSON := serialize.ToRawJSONList(paginatedVotes.Data)

	if optionalVotingData, ok := includes.Contains("voting"); ok {
		votingIncluder := serialization.NewVotingIncluder(votesJSON, h.entityManager)
		votingIncluder.Include(optionalVotingData, "voting_id")
	}

	paginatedVotes.Data = votesJSON
	http_response.WriteJSON(w, http.StatusOK, paginatedVotes)
}

// HandleGetVotings
//
//	@Summary	Return paginated list of votings
//
//	@Router		/votings [GET]
//
//	@Param		includes		query		string		false	"Optional fields' schema (votes_number{}, account_vote(hash))"
//	@Param		page			query		int			false	"Page number"											default(1)
//	@Param		page_size		query		string		false	"Number of items per page"								default(10)
//	@Param		order_direction	query		string		false	"Sorting direction"										Enums(ASC, DESC)		default(ASC)
//	@Param		order_by		query		[]string	false	"Comma-separated list of sorting fields (voting_id)"	collectionFormat(csv)	default(voting_id)
//
//	@Success	200				{object}	http_response.PaginatedResponse{data=entities.Voting}
//	@Failure	400,404,500		{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		Voting
func (h *Voting) HandleGetVotings(w http.ResponseWriter, r *http.Request) {
	includes, err := http_params.ParseOptionalData("includes", r)
	if err != nil {
		http_response.Error(w, r, err)
		return
	}

	paginationParams := pagination.NewParamsFromRequest(r)

	getVotings := voting.NewGetVotings()
	getVotings.SetEntityManager(h.entityManager)
	getVotings.SetPaginationParams(paginationParams)

	paginatedVotings, err := getVotings.Execute()
	if err != nil {
		http_response.Error(w, r, err)
		return
	}

	votingsJSON := serialize.ToRawJSONList(paginatedVotings.Data)
	if _, ok := includes.Contains("votes_number"); ok {
		votesNumberIncluder := serialization.NewVotesNumberIncluder(votingsJSON, h.entityManager)
		votesNumberIncluder.Include("voting_id")
	}

	if arg, ok := includes.ContainsFunc("account_vote"); ok {
		voteIncluder := serialization.NewAccountVoteIncluder(votingsJSON, h.entityManager)
		voteIncluder.Include(arg, "voting_id")
	}

	paginatedVotings.Data = votingsJSON
	http_response.WriteJSON(w, http.StatusOK, paginatedVotings)
}
