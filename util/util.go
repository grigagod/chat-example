package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"io"

	//"encoding/pem"
	//"errors"
	"math/big"
	"time"

	"github.com/grigagod/chat-example/crypto"
)

func GenKeyPair() (keys crypto.Keys) {
	keys.GenerateKeys()
	return
}

// MarshalKey marshals an DH public key to a byte slice
func MarshalKey(key *big.Int) []byte {
	return key.Bytes()
}

// UnmarshalKey unmarshals an DH key from byte format
func UnmarshalKey(slice []byte) (key *big.Int, err error) {
	key = new(big.Int).SetBytes(slice)
	return
}

// GenAuthChallenge generates Challenge using first 32 bytes of public keys as chiper
func GenAuthChallenge(pubKey *big.Int) ([]byte, []byte) {
	h := sha1.New()
	encKey := h.Sum(MarshalKey(pubKey))[:16]
	authKey := make([]byte, len(encKey))
	rand.Read(authKey)

	c, _ := aes.NewCipher(encKey)

	encAuthKey := make([]byte, len(authKey))

	c.Encrypt(encAuthKey, authKey)

	return encAuthKey, authKey
}

// DecryptChallenge decryptes Challenge using first 32 bytes of pubKey
func DecryptChallenge(pubKey *big.Int, msg []byte) []byte {
	h := sha1.New()
	decrKey := h.Sum(MarshalKey(pubKey))[:16]
	authKey := make([]byte, len(msg))

	c, _ := aes.NewCipher(decrKey)

	c.Decrypt(authKey, msg)

	return authKey

}

func EncryptDirectMessage(shared_key *big.Int, msg string) []byte {
	h := sha1.New()
	encrKey := h.Sum(MarshalKey(shared_key))[:16]

	block, _ := aes.NewCipher(encrKey)

	encrMsg := make([]byte, block.BlockSize()+len(msg))

	iv := encrMsg[:aes.BlockSize]
	io.ReadFull(rand.Reader, iv)
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrMsg[block.BlockSize():], []byte(msg))

	return encrMsg
}

func DecryptDirectMessage(shared_key *big.Int, msg []byte) string {
	h := sha1.New()
	decrKey := h.Sum(MarshalKey(shared_key))[:16]

	block, _ := aes.NewCipher(decrKey)

	decrMsg := make([]byte, len(msg)-block.BlockSize())

	iv := msg[:block.BlockSize()]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decrMsg, msg[block.BlockSize():])

	return string(decrMsg)
}

// NowMillis returns the current unix millisecond timestamp
func NowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
