package di

import "casper-dao-middleware/internal/crdao/types"

type DAOContractsMetadataAware struct {
	metadata types.DAOContractsMetadata
}

func (a *DAOContractsMetadataAware) SetDAOContractsMetadata(hashes types.DAOContractsMetadata) {
	a.metadata = hashes
}

func (a *DAOContractsMetadataAware) GetDAOContractsMetadata() types.DAOContractsMetadata {
	return a.metadata
}
