package stakingseq

import (
	"github.com/figment-networks/oasishub-indexer/models/shared"
	"github.com/figment-networks/oasishub-indexer/types"
)

type Model struct {
	*shared.Model
	*shared.Sequence

	TotalSupply         types.Quantity `json:"total_supply"`
	CommonPool          types.Quantity `json:"common_pool"`
	DebondingInterval   uint64         `json:"debonding_interval"`
	MinDelegationAmount types.Quantity `json:"min_delegation_amount"`
}

// - Methods
func (Model) TableName() string {
	return "staking_sequences"
}

func (ss *Model) ValidOwn() bool {
	return ss.TotalSupply.Valid() &&
		ss.CommonPool.Valid()
}

func (ss *Model) EqualOwn(m Model) bool {
	return true
}

func (ss *Model) Valid() bool {
	return ss.Model.Valid() &&
		ss.Sequence.Valid() &&
		ss.ValidOwn()
}

func (ss *Model) Equal(m Model) bool {
	return ss.Model.Equal(*m.Model) &&
		ss.Sequence.Equal(*m.Sequence) &&
		ss.EqualOwn(m)
}