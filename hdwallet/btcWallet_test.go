package hdwallet

import "testing"

func Test_btcWallet(t *testing.T) {
	tests := []struct {
		name string
	}{{
		name: "btc",
	},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := btcWallet()
			P2SHAddr(key)
		})
	}
}
