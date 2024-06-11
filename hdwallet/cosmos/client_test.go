package cosmos

import (
	"fmt"
	"testing"
)

func Test_client(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"test-client",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//NewClient()
			account, _ := queryAccountInfo()
			fmt.Println(account.Address)
		})
	}
}
