package jobs

import (
	"casper-dao-middleware/internal/dao/entities"
)

type GetJobByStatuses struct{}

func NewGetJobStatuses() *GetJobByStatuses {
	return &GetJobByStatuses{}
}

func (c *GetJobByStatuses) Execute() ([]entities.JobStatus, error) {
	return []entities.JobStatus{
		{
			ID:   entities.JobStatusIDCreated,
			Name: "created",
		},
		{
			ID:   entities.JobStatusIDSubmitted,
			Name: "submitted",
		},
		{
			ID:   entities.JobStatusIDCancelled,
			Name: "cancelled",
		},
		{
			ID:   entities.JobStatusIDDone,
			Name: "done",
		}, {
			ID:   entities.JobStatusIDRejected,
			Name: "rejected",
		},
	}, nil
}
