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
		t.Fatal("Unable to pass AuthChallenge\n", authKey, "\n", decrAuth)

	}

}
