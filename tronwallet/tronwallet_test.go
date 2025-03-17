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
func TestSendTrc20(t *testing.T) {
	result := SendTrc20()
	assert.True(t, result)
}
func TestGetTranscation(t *testing.T) {
	getTranscation("0xa0504eac7c9a0d5dd17f48d07569b6b8c39c9d3bb3f75940cb34b1f570f937d9")
}
func TestTransactionBuilder(t *testing.T) {
	result := TransactionBuilder()
	assert.True(t, result)
}
func TestGetContractABI(t *testing.T) {
	GetContractABI("TF17BgPaZYbz8oxbjhriubPDsA7ArKoLX3")
}
func TestDelegateTrx(t *testing.T) {
	result := DelegateTrx()
	assert.True(t, result)
}
