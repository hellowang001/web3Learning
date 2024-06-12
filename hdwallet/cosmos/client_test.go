package cosmos

import (
	"testing"
)

func Test_client(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test-client",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewClient()
			t.Log("This is a log message.")
			//transaction()
		})
	}
}
