package types

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyUref  = errors.New("empty uref")
	ErrUrefLength = errors.New("invalid bytes uref representation")
	ErrScanUref   = errors.New("failed to scan uref")
)

type Uref string

func NewUrefFromBytes(rawUref []byte) (Uref, error) {
	if len(rawUref) != 32 {
		return "", ErrUrefLength
	}

	return Uref(fmt.Sprintf("uref-%s-007", hex.EncodeToString(rawUref))), nil
}

func (u Uref) Bytes() ([]byte, error) {
	splits := strings.Split(string(u), "-")
	if len(splits) == 1 {
		return nil, errors.New("invalid uref format")
	}

	return hex.DecodeString(splits[1])
}

func (u Uref) String() string {
	return string(u)
}

func (u Uref) Value() (driver.Value, error) {
	// we save raw bytes into db
	return u.Bytes()
}

func (u *Uref) Scan(value interface{}) error {
	if value == nil {
		return ErrEmptyUref
	}

	bv, err := driver.String.ConvertValue(value)
	if err != nil {
		return ErrScanUref
	}

	v, ok := bv.([]byte)
	if !ok {
		return ErrScanUref
	}

	urefBuf := make([]byte, len(v))
	copy(urefBuf, v)

	uref, err := NewUrefFromBytes(urefBuf)
	if err != nil {
		return err
	}
	*u = uref
	return nil
}
