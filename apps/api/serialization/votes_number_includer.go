package serialization

import (
	"casper-dao-middleware/internal/crdao/persistence"
	"casper-dao-middleware/internal/crdao/services/voting"

	"go.uber.org/zap"
)

type VotesNumberIncluder struct {
	entityManager persistence.EntityManager
	entitiesJSON  []map[string]interface{}
}

func NewVotesNumberIncluder(entitiesJSON []map[string]interface{}, entityManager persistence.EntityManager) VotesNumberIncluder {
	return VotesNumberIncluder{
		entitiesJSON:  entitiesJSON,
		entityManager: entityManager,
	}
}

// Include map list of AccountInfo to target JSON
func (s *VotesNumberIncluder) Include(jsonMapKey string) {
	mapJSONCallback := func(entityJSON map[string]interface{}) uint32 {
		votingId, _ := entityJSON[jsonMapKey].(float64)
		return uint32(votingId)
	}

	votingIDs := make([]uint32, 0, len(s.entitiesJSON))
	for index := range s.entitiesJSON {
		mapJSONValue := mapJSONCallback(s.entitiesJSON[index])
		votingIDs = append(votingIDs, mapJSONValue)
	}

	getVotesNumber := voting.NewGetVotesNumber()
	getVotesNumber.SetEntityManager(s.entityManager)
	getVotesNumber.SetVotingIDs(votingIDs)
	votesNumberResult, err := getVotesNumber.Execute()
	if err != nil {
		zap.S().With(zap.Error(err)).Warn("Unable to find Votes number for including")
		return
	}

	for index := range s.entitiesJSON {
		mapJSONValue := mapJSONCallback(s.entitiesJSON[index])

		var votesNumber *uint32
		if number, ok := votesNumberResult[mapJSONValue]; ok {
			votesNumber = &number
		}

		s.entitiesJSON[index]["votes_number"] = votesNumber
	}
}
