package dogewallet

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/wire"
	"testing"
)

var (
	// 新建两个输入地址
	fromPrikey3, _     = hex.DecodeString("ac31d9781c013a02a14374a3aaf395629ac8df090fa473c463c7808ac466ae0f")
	fromPubkey3, _     = hex.DecodeString("03f00657c01decffd467f1e5a02c506420b0c589c7dcc7d3525da557eadd95dea2")
	fromLegacyAddress3 = "DGEPAvRQWYSTj29DrnQctCHt6bYcGi5Sx7"

	fromPrikey4, _     = hex.DecodeString("226039030d0db31dbba70135d05d7468dff0d25179c613723dfab6e7ef4963d0")
	fromPubkey4, _     = hex.DecodeString("03bbd68b6948c72adb96abba3a6ac0ffcc99dac2f99f5b56f93cc82f4bddc13f70")
	fromLegacyAddress4 = "DUUXfaJxGPTeFmKDpFn7FsE7kvSi3UxGSM"
	ac2                = TestAccount2{

		Prikey:        fromPrikey3,
		PubKey:        fromPubkey3,
		LegacyAddress: fromLegacyAddress3,
	}
	// 两个输出地址
	//toLegacyAddress3 = "D8Tve2b8AQrJyJG1GMuzNu6vRvAQNBxkgY"
	//toLegacyAddress3 = "DC4a3yz5Bm7Y96gHz6oRTro83dgSBm6G6g"
	toLegacyAddress3 = "DGEPAvRQWYSTj29DrnQctCHt6bYcGi5Sx7"
	toLegacyAddress4 = "D5r1YqMaimKtGT9wWJmNXRXJqdYdk7adSU"
	toLegacyAddress5 = "D9VEpes5Teq9oszWW9S7BtXSCjUpCmTbm2"
)

func Test_GetBalance(t *testing.T) {
	// 构建交易
	outputAmounts := []uint64{2456974, 951234} // 前面是发出金额，50000000 后面是找零金额1=100000000

	// 构建输入的交易结构，
	tx2 := &wire.MsgTx{
		Version:  1,
		LockTime: 0,
	}
	txIn2, err := txbuilder.NewInput("85ed74e9306fe56cc908f8650ba469dd552a419ad8329fa47f674649736c53db", 0)
	require.NoError(t, err)
	//tx2.AddTxIn(txIn1)
	tx2.AddTxIn(txIn2)

	// fee 假设 0.0001 // 添加交易输出。
	//output1, err := txbuilder.NewOutput(net, toLegacyAddress3, outputAmounts[0]) // 指定接收方，
	//output1, err := txbuilder.NewOutput(net, toLegacyAddress3, outputAmounts[0]) // 指定接收方，
	output2, err := txbuilder.NewOutput(net, toLegacyAddress4, outputAmounts[1]) // 指定接收方，
	//output2, err := txbuilder.NewOutput(net, fromLegacyAddress3, outputAmounts[1]) // 这是一笔找零
	//tx2.AddTxOut(output1)
	tx2.AddTxOut(output2)
	//tx2.AddTxOut(output2)
	// 准备签名脚本和公钥脚本的映射。
	var i2SignScript = make(map[uint32][]byte)
	var i2PkScript = make(map[uint32][]byte)
	// 为输入生成签名脚本。
	encodedAddress, err := btcutil.DecodeAddress(ac2.LegacyAddress, net) // 输入放地址解析出来 pkScript
	require.NoError(t, err)
	pkScript, err := txscript.PayToAddrScript(encodedAddress) // 把解析出来的地址，转成pkScript
	require.NoError(t, err)
	// 签名要签整个交易，要经过一次哈希，把未签名的交易哈希，
	hash, err := txscript.CalcSignatureHash(pkScript, txscript.SigHashAll, tx2, 0) // 生成一个待签名哈希
	require.NoError(t, err)

	// 哈希映射，把哈希存起来，后面签名的时候要用
	signHashMap[uint32(0)] = hash
	i2PkScript[uint32(0)] = pkScript

	// 签名
	sig, err := ac2.Key.Sign(hash) // 签名已经是der编码
	require.NoError(t, err)
	verify := sig.Verify(hash, ac2.Pk) // 验证签名，拿公钥验签，是否正确。
	require.True(t, verify)
	signature := append(sig.Serialize(), byte(txscript.SigHashAll))  //签名转成byte, 加上签名类型(doge 特性)
	class, _, _, err := txscript.ExtractPkScriptAddrs(pkScript, net) // 提取地址类型，拿到class
	require.NoError(t, err)
	switch class {
	case txscript.PubKeyHashTy:
		// 拿到脚本。signScript
		signScript, err := txscript.NewScriptBuilder().AddData(signature).AddData(ac2.PubKey).Script()
		require.NoError(t, err)
		i2SignScript[uint32(0)] = signScript
	default:
		t.Fatalf("un support class: %v error", class)
	}
}
