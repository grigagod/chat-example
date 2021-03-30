package util

import (
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

// MarshalKey marshals an DH public key to a byte slice
func MarshalKey(key *big.Int) []byte {
	pemBlock := &pem.Block{
		Type:  "DH KEY",
		Bytes: key.Bytes()}

	return pem.EncodeToMemory(pemBlock)
}

// UnmarshalKey unmarshals an DH key from byte format
func UnmarshalKey(pemBlock []byte) (key *big.Int, err error) {
	data, _ := pem.Decode(pemBlock)
	if data == nil {
		err = errors.New("Public key was not in the correct PEM format")
		return
	}
	key = new(big.Int).SetBytes(data.Bytes)
	return
}

// NowMillis returns the current unix millisecond timestamp
func NowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
