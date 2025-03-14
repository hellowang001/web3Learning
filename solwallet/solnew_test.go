package solwallet

import (
	"fmt"
	"testing"
)

func TestCreateAccountSimpleTx(t *testing.T) {
	CreateAccountSimpleTx()
}
func TestCreateAccountSystemTx(t *testing.T) {
	CreateAccountSystemTx()
}
func TestCreateComplexTx(t *testing.T) {
	CreateComplexTx()
}
func TestCreateTokenWithMetadata(t *testing.T) {
	fmt.Println("test --------------")
	CreateTokenWithMetadata()
}
func TestMintToken(t *testing.T) {
	fmt.Println("test --------------")
	MintToken()
}
