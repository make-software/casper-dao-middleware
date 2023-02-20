package di

import (
	"casper-dao-middleware/internal/dao/config"
)

type DAOContractsMetadataAware struct {
	metadata config.DAOContractsMetadata
}

func (a *DAOContractsMetadataAware) SetDAOContractsMetadata(hashes config.DAOContractsMetadata) {
	a.metadata = hashes
}

func (a *DAOContractsMetadataAware) GetDAOContractsMetadata() config.DAOContractsMetadata {
	return a.metadata
}
