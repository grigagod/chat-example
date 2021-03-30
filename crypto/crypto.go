package crypto

import (
	"crypto/rand"
	"math/big"
)

type GP struct {
	p *big.Int
	g *big.Int
}

type Keys struct {
	PublicKey  *big.Int
	PrivateKey *big.Int
	group      *GP
}

func (gp *GP) SetGP() {
	gp.p, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A63A3620FFFFFFFFFFFFFFFF", 16)
	gp.g = new(big.Int).SetInt64(2)
}

func (keys *Keys) SetPublicKey() {
	if keys.PrivateKey != nil {
		// g ^ pk mod p
		keys.PublicKey = new(big.Int).Exp(keys.group.g, keys.PrivateKey, keys.group.p)
	}
}

func (keys *Keys) GenerateKeys() {
	var gp GP
	gp.SetGP()

	randReader := rand.Reader

	private, err := rand.Int(randReader, gp.p)
	if err != nil {
		return
	}
	zero := big.NewInt(0)

	for private.Cmp(zero) == 0 {
		private, err = rand.Int(randReader, gp.p)
		if err != nil {
			return
		}
	}

	keys.PrivateKey = private
	keys.group = &gp
	keys.SetPublicKey()
}

func (self *Keys) KeyMixing(other *Keys) *big.Int {
	return new(big.Int).Exp(other.PublicKey, self.PrivateKey, self.group.p)
}
