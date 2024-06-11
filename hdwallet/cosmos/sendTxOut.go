package main




import (
	"github.com/cosmos/cosmos-sdk/simapp"
)

func main() error {
	// Choose your codec: Amino or Protobuf. Here, we use Protobuf, given by the following function.
	app := simapp.NewSimApp(...)

	// Create a new TxBuilder.
	txBuilder := app.TxConfig().NewTxBuilder()

	// --snip--
}
//func sendTx() error {
	//	// --snip--
	//

	//priv1, _, addr1 := testdata.KeyTestPubAddr()
	//priv2, _, addr2 := testdata.KeyTestPubAddr()
	//privs := []cryptotypes.PrivKey{priv1, priv2}
	//account, _ := queryAccountInfo() // The accounts' account numbers
	//fromAddress := account.Address
	//accountNumber := account.AccountNumber
	//accountSequence := account.Sequence
	//
	////
	//// First round: we gather all the signer infos. We use the "set empty
	////	// signature" hack to do that.
	//var sigsV2 []signing.SignatureV2
	//for i, priv := range privs {
	//	sigV2 := signing.SignatureV2{
	//		PubKey: priv.PubKey(),
	//		Data: &signing.SingleSignatureData{
	//			SignMode:  &typetx.TxConfig.SignModeHandler().DefaultMode(),
	//			Signature: nil,
	//		},
	//		Sequence: accSeqs[i],
	//	}
	//	//
	//	sigsV2 = append(sigsV2, sigV2)



		//	}
		//	err := txBuilder.SetSignatures(sigsV2...)
		//	if err != nil {
		//		return err
		//	}
		//
		//	// Second round: all signer infos are set, so each signer can sign.
		//	sigsV2 = []signing.SignatureV2{}
		//	for i, priv := range privs {
		//		signerData := xauthsigning.SignerData{
		//			ChainID:       chainID,
		//			AccountNumber: accNums[i],
		//			Sequence:      accSeqs[i],
		//		}
		//		sigV2, err := tx.SignWithPrivKey(
		//			encCfg.TxConfig.SignModeHandler().DefaultMode(), signerData,
		//			txBuilder, priv, encCfg.TxConfig, accSeqs[i])
		//		if err != nil {
		//			return nil, err
		//		}
		//
		//		sigsV2 = append(sigsV2, sigV2)
		//	}
		//	err = txBuilder.SetSignatures(sigsV2...)
		//	if err != nil {
		//		return err
		//	}
	//}
