package util

import (
	"bytes"
	"testing"
)

func TestAuthChallenge(t *testing.T) {
	keys := GenKeyPair()

	encAuth, authKey := GenAuthChallenge(keys.PublicKey)
	decrAuth := DecryptChallenge(keys.PublicKey, encAuth)
	if !bytes.Equal(authKey, decrAuth) {
		t.Log(decrAuth)
		t.Log(authKey)
		t.Fatal("Unable to pass AuthChallenge")

	}

}

func TestDirectMessage(t *testing.T) {
	keys1 := GenKeyPair()
	keys2 := GenKeyPair()

	shared_key := keys1.KeyMixing(keys2.PublicKey)

	msg := "Hi dude its me gregory "

	encrMsg := EncryptDirectMessage(shared_key, msg)

	decrMsg := DecryptDirectMessage(shared_key, encrMsg)

	if decrMsg != msg {
		t.Log(decrMsg)
		t.Fatal("Unable to decrypt direct message ")
	}
}
