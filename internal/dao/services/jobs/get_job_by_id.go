package jobs

import (
	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
)

type GetJobByID struct {
	di.EntityManagerAware

	jobID uint32
}

func NewGetJobById() *GetJobByID {
	return &GetJobByID{}
}

func (c *GetJobByID) SetJobID(jobID uint32) {
	c.jobID = jobID
}

func (c *GetJobByID) Execute() (*entities.Job, error) {
	return c.GetEntityManager().JobRepository().GetByID(c.jobID)
}
