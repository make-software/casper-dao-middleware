package types

import (
	"errors"
	"strings"
)

var (
	ErrInvalidKeyString = errors.New("invalid key string")
)

const (
	KeyAccount KeyType = iota
	KeyHash
	KeyURef
	KeyTransfer
	KeyDeployInfo
	KeyEraInfo
	KeyBalance
	KeyBid
	KeyWithdraw
	KeyDictionary
	KeySystemContractRegistry
)

type (
	KeyType       byte
	KeyStringType string
	Key           struct {
		Type        KeyType
		AccountHash *Hash
		Hash        *Hash
	}
)

func (k Key) GetHashValue() (*Hash, error) {
	switch k.Type {
	case KeyAccount:
		return k.AccountHash, nil
	case KeyHash:
		return k.Hash, nil
	default:
		return nil, errors.New("no hash value found")
	}
}

func ParseKeyFromBytes(bytes []byte) (Key, []byte, error) {
	if len(bytes) == 0 {
		return Key{}, nil, errors.New("empty bytes provided")
	}
	// take first byte to see KeyType
	keyType := KeyType(bytes[0])
	remainder := make([]byte, len(bytes)-1)
	copy(remainder, bytes[1:])

	switch keyType {
	case KeyAccount:
		hash, err := NewHashFromRawBytes(remainder[:32])
		if err != nil {
			return Key{}, nil, err
		}
		return Key{
			Type:        KeyAccount,
			AccountHash: &hash,
		}, remainder[32:], nil
	case KeyHash:
		hash, err := NewHashFromRawBytes(remainder[:32])
		if err != nil {
			return Key{}, nil, err
		}
		return Key{
			Hash: &hash,
		}, remainder[32:], nil
	}

	return Key{}, nil, nil
}

// Example: "Key::Account(<hash>)", "Key::Hash(<hash>)"
func ParseKeyFromString(key string) (Key, error) {
	idx := strings.Index(key, "::")
	if idx == -1 {
		return Key{}, ErrInvalidKeyString
	}

	typeIdxStart := idx + 2

	openBracketIndex := strings.Index(key, "(")
	if openBracketIndex == -1 {
		return Key{}, ErrInvalidKeyString
	}

	keyTypeStr := key[typeIdxStart:openBracketIndex]
	keyValue := key[openBracketIndex+1 : len(key)-1]

	switch keyTypeStr {
	case "Account":
		hash, err := NewHashFromHexString(keyValue)
		if err != nil {
			return Key{}, err
		}
		return Key{
			Type:        KeyAccount,
			AccountHash: &hash,
		}, nil
	case "Hash":
		hash, err := NewHashFromHexString(keyValue)
		if err != nil {
			return Key{}, err
		}
		return Key{
			Type: KeyHash,
			Hash: &hash,
		}, nil
	}

	return Key{}, nil
}
