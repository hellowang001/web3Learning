package hdwallet

import "testing"

func Test_cosmosWalletPrivate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "cosmosWalletPrivate",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("This is a log message.")
			cosmosWalletPrivate()

		})
	}
}
