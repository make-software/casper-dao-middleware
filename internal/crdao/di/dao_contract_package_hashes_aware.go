package di

import "casper-dao-middleware/internal/crdao/dao_event_parser"

type DAOContractPackageHashesAware struct {
	hashes dao_event_parser.DAOContractPackageHashes
}

func (a *DAOContractPackageHashesAware) SetDAOContractPackageHashes(hashes dao_event_parser.DAOContractPackageHashes) {
	a.hashes = hashes
}

func (a *DAOContractPackageHashesAware) GetDAOContractPackageHashes() dao_event_parser.DAOContractPackageHashes {
	return a.hashes
}
