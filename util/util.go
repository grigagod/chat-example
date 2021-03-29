package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"
)

// MarshalPublic marshals an RSA public key to a byte slice
func MarshalPublic(key *rsa.PublicKey) []byte {
	pemBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(key)}

	return pem.EncodeToMemory(pemBlock)
}

// UnmarshalPublic unmarshals an RSA public key from byte format
func UnmarshalPublic(pemBlock []byte) (key *rsa.PublicKey, err error) {
	data, _ := pem.Decode(pemBlock)
	if data == nil {
		err = errors.New("Public key was not in the correct PEM format")
		return
	}
	key, err = x509.ParsePKCS1PublicKey(data.Bytes)
	return
}

// NowMillis returns the current unix millisecond timestamp
func NowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
