package types

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

const BytesHashLength = 32

var (
	ErrEmptyHash         = errors.New("empty hash")
	ErrInvalidHashLength = errors.New("invalid hash length")
	ErrScanHash          = errors.New("failed to scan Hash")
)

// Hash represent bytes network hash representation
type Hash []byte

func NewHashFromHexStringWithPrefix(hexData string, prefix string) (Hash, error) {
	return NewHashFromHexString(strings.TrimPrefix(hexData, prefix))
}

// NewHashFromHexString creates a Hash object from a hexadecimal string.
func NewHashFromHexString(hexData string) (Hash, error) {
	bytesData, err := hex.DecodeString(hexData)
	if err != nil {
		return Hash{}, err
	}

	return newHash(bytesData)
}

// NewHashFromRawBytes creates a Hash object from a bytes.
func NewHashFromRawBytes(data []byte) (Hash, error) {
	return newHash(data)
}

func (h Hash) String() string {
	return h.ToHex()
}

// ToHex returns an on-chain account key in hex.
func (h *Hash) ToHex() string {
	return hex.EncodeToString(*h)
}

// Bytes return byte representation of Hash object.
func (h *Hash) Bytes() []byte {
	return *h
}

// MarshalJSON convert Hash to hex during marshaling
func (h *Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ToHex())
}

func (h *Hash) UnmarshalJSON(data []byte) error {
	hashBuf := make([]byte, len(data))
	copy(hashBuf, data)

	var value string
	if err := json.Unmarshal(hashBuf, &value); err != nil {
		return err
	}

	hash, err := NewHashFromHexString(value)
	if err != nil {
		return err
	}

	*h = hash
	return nil
}

func newHash(rawData []byte) (Hash, error) {
	if len(rawData) == 0 {
		return Hash{}, ErrEmptyHash
	}

	if len(rawData) != BytesHashLength {
		return Hash{}, ErrInvalidHashLength
	}

	return rawData, nil
}

// Value rewrite behaviour for inserting Hash to db
// Note: Value receiver for squirrel
func (h Hash) Value() (driver.Value, error) {
	if h == nil {
		return nil, nil
	}
	// we save raw bytes into db
	return h.Bytes(), nil
}

// Scan rewrite behaviour for selecting Hash from db
func (h *Hash) Scan(value interface{}) error {
	if value == nil {
		return ErrEmptyHash
	}

	bv, err := driver.String.ConvertValue(value)
	if err != nil {
		return ErrScanHash
	}

	v, ok := bv.([]byte)
	if !ok {
		return ErrScanHash
	}

	hashBuf := make([]byte, len(v))
	copy(hashBuf, v)

	hash, err := NewHashFromRawBytes(hashBuf)
	if err != nil {
		return err
	}
	*h = hash
	return nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	hash, err := NewHashFromHexString(string(data))
	if err != nil {
		return err
	}

	*h = hash
	return nil
}
