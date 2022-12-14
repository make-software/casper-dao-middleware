package di

import "casper-dao-middleware/internal/crdao/dao_event_parser"

type DAOContractsMetadataAware struct {
	metadata dao_event_parser.DAOContractsMetadata
}

func (a *DAOContractsMetadataAware) SetDAOContractsMetadata(hashes dao_event_parser.DAOContractsMetadata) {
	a.metadata = hashes
}

func (a *DAOContractsMetadataAware) GetDAOContractsMetadata() dao_event_parser.DAOContractsMetadata {
	return a.metadata
}
