package util

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/pem"
	"errors"
	"github.com/grigagod/chat-example/crypto"
	"math/big"
	"time"
)

const authKeyLen = 64

func GenKeyPair() (keys *crypto.Keys) {
	keys.GenerateKeys()
	return
}

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

func GenAuthChallenge(pubKey *big.Int) ([]byte, []byte) {
	encKey := MarshalKey(pubKey)[:32]
	authKey := make([]byte, authKeyLen)
	rand.Read(authKey)

	c, _ := aes.NewCipher(encKey)

	encAuthKey := make([]byte, authKeyLen)

	c.Encrypt(encAuthKey, authKey)

	return encAuthKey, authKey
}

// NowMillis returns the current unix millisecond timestamp
func NowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
