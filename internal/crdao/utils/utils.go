package utils

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/blake2b"
)

func ToDictionaryItemKey(key string) (string, error) {
	res := make([]byte, 0)
	blake, err := blake2b.New256(res)
	if err != nil {
		return "", err
	}

	keyBytes := []byte(key)
	blake.Write(binary.LittleEndian.AppendUint32(nil, uint32(len(keyBytes))))
	blake.Write(keyBytes)

	return hex.EncodeToString(blake.Sum(nil)), nil
}

func ToDictionaryKey(eventsUref string, index uint32) (string, error) {
	urefsParts := strings.Split(eventsUref, "-")
	// uref format uref-d1a68e4ae2c8ffe65cafcfc172caf1179bc5fa820d25eb4574a48f89225820a0-007
	if len(urefsParts) != 3 {
		return "", errors.New("invalid uref format provided")
	}
	urefHashBytes, err := hex.DecodeString(urefsParts[1])
	if err != nil {
		return "", err
	}

	res := make([]byte, 0)
	key, err := blake2b.New256(res)
	if err != nil {
		return "", err
	}

	key.Write(urefHashBytes)
	key.Write(calculateDictionaryIndexHash(index))
	dictionaryKey := fmt.Sprintf("dictionary-%s", hex.EncodeToString(key.Sum(nil)))
	return dictionaryKey, nil
}

func calculateDictionaryIndexHash(index uint32) []byte {
	indexBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexBytes, index)

	indexHash := blake2b.Sum256(indexBytes)
	return []byte(hex.EncodeToString(indexHash[:]))
}
