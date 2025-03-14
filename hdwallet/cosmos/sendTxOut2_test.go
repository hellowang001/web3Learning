package cosmos

import (
	"testing"
)

func Test_transaction2(t *testing.T) {
	tests := []struct {
		name string
	}{{
		name: "test_tx1",
	},

	// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//transaction()
		})
	}
}

//func TestImport2(t *testing.T) {
//	type args struct {
//		privateKeyStr string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *secp256k1.PrivKey
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := Import2(tt.args.privateKeyStr)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Import2(%v) error = %v, wantErr %v", tt.args.privateKeyStr, err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Import2(%v) = %v, want %v", tt.args.privateKeyStr, got, tt.want)
//			}
//		})
//	}
//}
