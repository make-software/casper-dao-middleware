package entities

type TotalReputation struct {
	AvailableAmount *uint64 `json:"available_amount" db:"available_amount"`
	StakedAmount    *uint64 `json:"staked_amount" db:"staked_amount"`
}
