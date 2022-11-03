package serialization

import (
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/internal/crdao/persistence"
	"casper-dao-middleware/internal/crdao/services/voting"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/pagination"

	"go.uber.org/zap"
)

type AccountVoteIncluder struct {
	entityManager persistence.EntityManager
	entitiesJSON  []map[string]interface{}
}

func NewAccountVoteIncluder(entitiesJSON []map[string]interface{}, entityManager persistence.EntityManager) AccountVoteIncluder {
	return AccountVoteIncluder{
		entitiesJSON:  entitiesJSON,
		entityManager: entityManager,
	}
}

// Include map list of AccountInfo to target JSON
func (s *AccountVoteIncluder) Include(args []string, jsonMapKey string) {
	if len(args) != 1 {
		zap.S().Info("account_vote function expected to have one parameter")
		return
	}

	addressHash, err := types.NewHashFromHexString(args[0])
	if err != nil {
		zap.S().With(zap.Error(err)).Info("invalid Hash arg provided provided for account_vote function")
		return
	}

	mapJSONCallback := func(entityJSON map[string]interface{}) uint32 {
		votingId, _ := entityJSON[jsonMapKey].(float64)
		return uint32(votingId)
	}

	votingIDs := make([]uint32, 0, len(s.entitiesJSON))
	for index := range s.entitiesJSON {
		mapJSONValue := mapJSONCallback(s.entitiesJSON[index])
		votingIDs = append(votingIDs, mapJSONValue)
	}

	getVotes := voting.NewGetVotes()
	getVotes.SetEntityManager(s.entityManager)
	getVotes.SetVotingIDs(votingIDs)
	getVotes.SetAddress(&addressHash)
	getVotes.SetPaginationParams(&pagination.Params{
		OrderDirection: pagination.OrderDirectionDESC,
		Page:           1,
		PageSize:       uint64(len(votingIDs)),
	})
	paginatedVotes, err := getVotes.Execute()
	if err != nil {
		zap.S().With(zap.Error(err)).Warn("Unable to find Votes for including")
		return
	}

	votes := paginatedVotes.Data.([]*entities.Vote)

	votesMap := make(map[uint32]*entities.Vote)
	for _, vote := range votes {
		votesMap[vote.VotingID] = vote
	}

	for index := range s.entitiesJSON {
		mapJSONValue := mapJSONCallback(s.entitiesJSON[index])

		var accountVote *entities.Vote
		if vote, ok := votesMap[mapJSONValue]; ok {
			accountVote = vote
		}

		s.entitiesJSON[index]["account_vote"] = accountVote
	}
}
