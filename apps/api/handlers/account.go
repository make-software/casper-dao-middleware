package handlers

import (
	"net/http"

	"casper-dao-middleware/apps/api/serialization"
	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/account"
	"casper-dao-middleware/internal/dao/services/votes"
	"casper-dao-middleware/pkg/errors"
	http_params "casper-dao-middleware/pkg/http-params"
	http_response "casper-dao-middleware/pkg/http-response"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/serialize"
)

type Account struct {
	entityManager persistence.EntityManager
}

func NewAccount(entityManager persistence.EntityManager) *Account {
	return &Account{
		entityManager: entityManager,
	}
}

// HandleGetAccountVotes
//
//	@Summary	Return paginated list of votes for address
//
//	@Router		/accounts/{address}/votes [GET]
//
//	@Param		address			path		string		true	"Hash or PublicKey"	maxlength(66)
//	@Param		includes		query		string		false	"Optional fields' schema (voting{})"
//	@Param		page			query		int			false	"Page number"													default(1)
//	@Param		page_size		query		string		false	"Number of items per page"										default(10)
//	@Param		order_direction	query		string		false	"Sorting direction"												Enums(ASC, DESC)		default(ASC)
//	@Param		order_by		query		[]string	false	"Comma-separated list of sorting fields (voting_id,address)"	collectionFormat(csv)	default(voting_id)
//
//	@Success	200				{object}	http_response.PaginatedResponse{data=entities.Vote}
//	@Failure	400,404,500		{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		Vote
func (h *Account) HandleGetAccountVotes(w http.ResponseWriter, r *http.Request) {
	addressHash, err := http_params.ParseOptionalHash("address", r)
	if err != nil {
		accountPubKey, err := http_params.ParseOptionalPublicKey("address", r)
		if err != nil {
			http_response.Error(w, r, errors.NewInvalidInputError("Account address is not a valid account hash or public key"))
			return
		}
		addressHash = accountPubKey.AccountHash()
	}

	includes, err := http_params.ParseOptionalData("includes", r)
	if err != nil {
		http_response.Error(w, r, err)
		return
	}

	paginationParams := pagination.NewParamsFromRequest(r)

	getVotes := votes.NewGetVotes()
	getVotes.SetAddress(addressHash)
	getVotes.SetEntityManager(h.entityManager)
	getVotes.SetPaginationParams(paginationParams)

	paginatedVotes, err := getVotes.Execute()
	if err != nil {
		http_response.Error(w, r, err)
		return
	}

	votesJSON := serialize.ToRawJSONList(paginatedVotes.Data)

	if optionalVotingData, ok := includes.Contains("voting"); ok {
		votingsIncluder := serialization.NewVotingIncluder(votesJSON, h.entityManager)
		votingsIncluder.Include(optionalVotingData, "voting_id")
	}

	paginatedVotes.Data = votesJSON
	http_response.WriteJSON(w, http.StatusOK, paginatedVotes)
}

// HandleGetAccounts
//
//	@Summary	Return paginated list of accounts
//
//	@Router		/accounts [GET]
//
//	@Param		page		query		int		false	"Page number"				default(1)
//	@Param		page_size	query		string	false	"Number of items per page"	default(10)
//
//	@Success	200			{object}	http_response.PaginatedResponse{data=entities.Account}
//	@Failure	400,404,500	{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		Vote
func (h *Account) HandleGetAccounts(w http.ResponseWriter, r *http.Request) {
	paginationParams := pagination.NewParamsFromRequest(r)

	getAccounts := account.NewGetAccounts()
	getAccounts.SetEntityManager(h.entityManager)
	getAccounts.SetPaginationParams(paginationParams)

	http_response.FromFunction(getAccounts.Execute, w, r)
}

// HandleGetAccountsByAddress
//
//	@Summary	Return account by its address
//
//	@Router		/accounts/{address}  [GET]
//
//	@Param		address		path		string	true	"Hash or PublicKey"	maxlength(66)
//
//	@Success	200			{object}	http_response.SuccessResponse{data=entities.Account}
//	@Failure	400,404,500	{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		Vote
func (h *Account) HandleGetAccountsByAddress(w http.ResponseWriter, r *http.Request) {
	addressHash, err := http_params.ParseOptionalHash("address", r)
	if err != nil {
		accountPubKey, err := http_params.ParseOptionalPublicKey("address", r)
		if err != nil {
			http_response.Error(w, r, errors.NewInvalidInputError("Account address is not a valid account hash or public key"))
			return
		}
		addressHash = accountPubKey.AccountHash()
	}

	getAccount := account.NewGetAccountByHash()
	getAccount.SetHash(*addressHash)
	getAccount.SetEntityManager(h.entityManager)

	http_response.FromFunction(getAccount.Execute, w, r)
}
