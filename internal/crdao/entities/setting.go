package entities

import (
	"time"
)

type Setting struct {
	Name           string     `json:"name" db:"name"`
	Value          string     `json:"value" db:"value"`
	NextValue      *string    `json:"-" db:"next_value"`
	ActivationTime *time.Time `json:"-" db:"activation_time"`
}

func NewSetting(name, value string, nextValue *string, activationTime *time.Time) Setting {
	return Setting{
		Name:           name,
		Value:          value,
		NextValue:      nextValue,
		ActivationTime: activationTime,
	}
}
