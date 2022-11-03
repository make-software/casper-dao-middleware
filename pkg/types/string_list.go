package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type StringList []string

func (l StringList) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *StringList) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		err := json.Unmarshal(v, l)
		return err
	default:
		return errors.New("invalid type when trying to unmarshal StringList read from the database")
	}
}

func (l StringList) SubSlice(slice []string) bool {
	if len(slice) > len(l) {
		return false
	}
	for _, elem := range slice {
		if !l.Contains(elem) {
			return false
		}
	}
	return true
}

func (l StringList) Contains(element string) bool {
	for _, value := range l {
		if value == element {
			return true
		}
	}
	return false
}
