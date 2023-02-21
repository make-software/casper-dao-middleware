package di

import "casper-dao-middleware/internal/dao/utils"

type DAOContractsMetadataAware struct {
	metadata utils.DAOContractsMetadata
}

func (a *DAOContractsMetadataAware) SetDAOContractsMetadata(hashes utils.DAOContractsMetadata) {
	a.metadata = hashes
}

func (a *DAOContractsMetadataAware) GetDAOContractsMetadata() utils.DAOContractsMetadata {
	return a.metadata
}
