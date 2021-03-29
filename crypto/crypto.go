package main

import (
//	"crypto/rand"
	"fmt"
	"math/big"
)


type gp struct {
	p *big.Int
	g *big.Int
}


func (self *gp) P() *big.Int {
	p := new(big.Int)
	p.Set(self.p)
	return p
}

func (self *gp) G() *big.Int {
	g := new(big.Int)
	g.Set(self.g)
	return g
}
