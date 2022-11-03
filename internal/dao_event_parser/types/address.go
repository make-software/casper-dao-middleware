package types

import "casper-dao-middleware/pkg/casper/types"

type Address struct {
	AccountHash         *types.Hash
	ContractPackageHash *types.Hash
}
