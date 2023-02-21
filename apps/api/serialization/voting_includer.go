package serialization

import (
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/voting"
	"casper-dao-middleware/pkg/pagination"
	"casper-dao-middleware/pkg/serialize"
	pkgTypes "casper-dao-middleware/pkg/types"
	optional_data "casper-dao-middleware/pkg/types/optional-data"

	"go.uber.org/zap"
)

type VotingIncluder struct {
	entityManager persistence.EntityManager
	entitiesJSON  []map[string]interface{}
}

func NewVotingIncluder(entitiesJSON []map[string]interface{}, entityManager persistence.EntityManager) VotingIncluder {
	return VotingIncluder{
		entitiesJSON:  entitiesJSON,
		entityManager: entityManager,
	}
}

// Include map list of AccountInfo to target JSON
func (s *VotingIncluder) Include(optData *pkgTypes.OptionalData, jsonMapKey string) {
	mapJSONCallback := func(entityJSON map[string]interface{}) uint32 {
		votingId, _ := entityJSON[jsonMapKey].(float64)
		return uint32(votingId)
	}

	votingIDsMap := make(map[uint32]struct{}, len(s.entitiesJSON))
	for index := range s.entitiesJSON {
		mapJSONValue := mapJSONCallback(s.entitiesJSON[index])

		if _, ok := votingIDsMap[mapJSONValue]; !ok {
			votingIDsMap[mapJSONValue] = struct{}{}
		}
	}

	votingIDs := make([]uint32, 0, len(votingIDsMap))
	for votingID := range votingIDsMap {
		votingIDs = append(votingIDs, votingID)
	}

	getVotings := voting.NewGetVotings()
	getVotings.SetEntityManager(s.entityManager)
	getVotings.SetVotingIDs(votingIDs)
	getVotings.SetPaginationParams(&pagination.Params{
		OrderDirection: pagination.OrderDirectionDESC,
		Page:           1,
		PageSize:       uint64(len(votingIDs)),
	})

	paginatedVotings, err := getVotings.Execute()
	if err != nil {
		zap.S().With(zap.Error(err)).Warn("Unable to find Votings for including")
		return
	}
	votings := paginatedVotings.Data.([]*entities.Voting)

	votingHashMap := make(map[uint32]*entities.Voting)
	for _, voting := range votings {
		votingHashMap[voting.VotingID] = voting
	}

	for index := range s.entitiesJSON {
		mapJSONValue := mapJSONCallback(s.entitiesJSON[index])

		var optionalData map[string]interface{}
		if voting, ok := votingHashMap[mapJSONValue]; ok {
			optionalData = serialize.ToRawJSON(voting)
		}

		err := optional_data.MapJSON(optData, s.entitiesJSON[index], optionalData)
		if err != nil {
			zap.S().With(zap.Error(err)).Info("Error on mapping optional Votings data")
		}
	}
}
