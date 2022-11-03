package types

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/crypto/blake2b"
)

type (
	KeyAlgorithm        byte
	KeyAlgorithmSetting struct {
		length int
		name   string
	}
)

const (
	ED25519   KeyAlgorithm = 1
	SECP256K1 KeyAlgorithm = 2
)

var PublicKeySettings = map[KeyAlgorithm]KeyAlgorithmSetting{
	ED25519: {
		length: 32,
		name:   "ED25519",
	},
	SECP256K1: {
		length: 33,
		name:   "SECP256K1",
	},
}

func (a KeyAlgorithm) String() string {
	return PublicKeySettings[a].name
}

var (
	ErrEmptyPublicKey            = errors.New("empty public key")
	ErrPublicKeyLength           = errors.New("invalid public key length")
	ErrNoPublicKeyAlgorithmFound = errors.New("no public key algorithm found")
	ErrScanPublicKey             = errors.New("failed to scan PublicKey")
)

// PublicKey represent account public key type
type PublicKey struct {
	algo KeyAlgorithm
	pbk  []byte
}

// NewPublicKeyFromHexString creates a PublicKey object from a hexadecimal string (containing the Key algorithm identifier).
func NewPublicKeyFromHexString(hexData string) (PublicKey, error) {
	bytesData, err := hex.DecodeString(hexData)
	if err != nil {
		return PublicKey{}, err
	}

	if len(bytesData) == 0 {
		return PublicKey{}, ErrEmptyPublicKey
	}

	algo := KeyAlgorithm(bytesData[0])
	if _, ok := PublicKeySettings[algo]; !ok {
		return PublicKey{}, ErrNoPublicKeyAlgorithmFound
	}

	return newPublicKey(bytesData[1:], algo)
}

// NewPublicKeyFromBytes creates a PublicKey object from a bytes (containing the Key algorithm identifier).
func NewPublicKeyFromBytes(data []byte) (PublicKey, error) {
	if len(data) == 0 {
		return PublicKey{}, ErrEmptyPublicKey
	}

	algo := KeyAlgorithm(data[0])
	if _, ok := PublicKeySettings[algo]; !ok {
		return PublicKey{}, ErrNoPublicKeyAlgorithmFound
	}

	return newPublicKey(data[1:], algo)
}

// NewPublicKeyFromRawBytes creates a PublicKey object from a bytes without algorithm identifier.
func NewPublicKeyFromRawBytes(data []byte) (PublicKey, error) {
	if len(data) == 0 {
		return PublicKey{}, ErrEmptyPublicKey
	}

	var algo KeyAlgorithm
	for k, setting := range PublicKeySettings {
		if len(data) != setting.length {
			continue
		}
		algo = k
	}

	if algo == 0 {
		return PublicKey{}, ErrNoPublicKeyAlgorithmFound
	}

	return newPublicKey(data, algo)
}

func (p *PublicKey) AccountHash() *Hash {
	bytesToHash := make([]byte, 0, len(p.algo.String())+1+len(p.pbk))

	bytesToHash = append(bytesToHash, []byte(strings.ToLower(p.algo.String()))...)
	bytesToHash = append(bytesToHash, byte(0))
	bytesToHash = append(bytesToHash, p.pbk...)

	blakeHash := blake2b.Sum256(bytesToHash)
	hash, _ := NewHashFromRawBytes(blakeHash[:])

	return &hash
}

// ToHex returns an on-chain account key in hex.
func (p *PublicKey) ToHex() string {
	return hex.EncodeToString(p.Bytes())
}

// Equals compare with provided PublicKey
func (p *PublicKey) Equals(pubKey PublicKey) bool {
	return p.ToHex() == pubKey.ToHex()
}

// Bytes return byte representation of PublicKey object.
func (p PublicKey) Bytes() []byte {
	res := make([]byte, 0, len(p.pbk)+1)

	res = append(res, byte(p.algo))
	res = append(res, p.pbk...)

	return res
}

func newPublicKey(rawData []byte, algo KeyAlgorithm) (PublicKey, error) {
	pbk := make([]byte, len(rawData))
	copy(pbk, rawData)

	if len(pbk) != PublicKeySettings[algo].length {
		return PublicKey{}, ErrPublicKeyLength
	}

	return PublicKey{
		algo,
		pbk,
	}, nil
}

func (p *PublicKey) UnmarshalJSON(data []byte) error {
	hashBuf := make([]byte, len(data))
	copy(hashBuf, data)

	var value string
	if err := json.Unmarshal(hashBuf, &value); err != nil {
		return err
	}

	publicKey, err := NewPublicKeyFromHexString(value)
	if err != nil {
		return err
	}

	p.algo = publicKey.algo
	p.pbk = publicKey.pbk
	return nil
}

// MarshalJSON convert PublicKey to hex during marshaling
func (p *PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToHex())
}

// Value rewrite behaviour for inserting PublicKey to db
func (p PublicKey) Value() (driver.Value, error) {
	// we save raw bytes into db
	return p.Bytes(), nil
}

// Scan rewrite behaviour for selecting PublicKey from db
func (p *PublicKey) Scan(value interface{}) error {
	if value == nil {
		return ErrEmptyPublicKey
	}

	bv, err := driver.String.ConvertValue(value)
	if err != nil {
		return ErrScanPublicKey
	}

	v, ok := bv.([]byte)
	if !ok {
		return ErrScanPublicKey
	}

	pubKeyBuf := make([]byte, len(v))
	copy(pubKeyBuf, v)

	pubKey, err := NewPublicKeyFromBytes(pubKeyBuf)
	if err != nil {
		return err
	}
	*p = pubKey
	return nil
}

func PublicKeysToHexList(publicKeys []PublicKey) []string {
	result := make([]string, 0, len(publicKeys))
	for _, key := range publicKeys {
		result = append(result, key.ToHex())
	}
	return result
}
