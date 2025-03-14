package tronwallet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTronWallet(t *testing.T) {
	piv, addr := TronWallet()
	fmt.Println("piv:", piv)
	fmt.Println("addr:", addr)
}
func TestSendTrx(t *testing.T) {
	result := SendTrx()
	assert.True(t, result)

}
func TestWalletGen(t *testing.T) {
	WalletGen()
}
func TestFindAccount(t *testing.T) {
	FindAccount()
}
func TestVerifyAccount(t *testing.T) {
	VerifyAccount()
}
