package handlers

import (
	"net/http"

	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/settings"
	http_response "casper-dao-middleware/pkg/http-response"
	"casper-dao-middleware/pkg/pagination"
)

type Setting struct {
	entityManager persistence.EntityManager
}

func NewSetting(entityManager persistence.EntityManager) *Setting {
	return &Setting{
		entityManager: entityManager,
	}
}

// HandleGetSettings
//
//	@Summary	Return paginated list of settings
//
//	@Router		/settings [GET]
//
//	@Param		page			query		int			false	"Page number"									default(1)
//	@Param		page_size		query		string		false	"Number of items per page"						default(10)
//	@Param		order_direction	query		string		false	"Sorting direction"								Enums(ASC, DESC)		default(ASC)
//	@Param		order_by		query		[]string	false	"Comma-separated list of sorting fields (name)"	collectionFormat(csv)	default(voting_id)
//
//	@Success	200				{object}	http_response.PaginatedResponse{data=entities.Setting}
//	@Failure	400,404,500		{object}	http_response.ErrorResponse{error=http_response.ErrorResult}
//
//	@tags		Setting
func (h *Setting) HandleGetSettings(w http.ResponseWriter, r *http.Request) {
	paginationParams := pagination.NewParamsFromRequest(r)

	getSettings := settings.NewGetPaginatedSettings()
	getSettings.SetEntityManager(h.entityManager)
	getSettings.SetPaginationParams(paginationParams)

	http_response.FromFunction(getSettings.Execute, w, r)
}
