package settings

import (
	"time"

	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/pkg/pagination"
)

type GetSettings struct {
	di.PaginationParamsAware
	di.EntityManagerAware
}

func NewGetSettings() *GetSettings {
	return &GetSettings{}
}

func (c *GetSettings) Execute() (*pagination.Result, error) {
	filters := map[string]interface{}{}

	count, err := c.GetEntityManager().SettingRepository().Count(filters)
	if err != nil {
		return nil, err
	}

	settings, err := c.GetEntityManager().SettingRepository().Find(c.GetPaginationParams(), filters)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now().UTC()
	for i := range settings {
		activationTime := settings[i].ActivationTime
		if activationTime == nil || settings[i].NextValue == nil {
			continue
		}

		if activationTime.Before(currentTime) {
			settings[i].Value = *settings[i].NextValue
		}
	}

	return pagination.NewResult(count, c.GetPaginationParams().PageSize, settings), nil
}
