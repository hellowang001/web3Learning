from pip import BIP32Key
from bip32utils import BIP32_HARDEN
from bip39 import Mnemonic
from ethereum.utils import privtoaddr

def eth_wallet():
    # 生成助记词和种子
    entropy = Mnemonic().generate_entropy(128)
    mnemonic = Mnemonic().to_mnemonic(entropy)
    seed = Mnemonic().to_seed(mnemonic)

    # 派生主私钥
    bip32_key = BIP32Key.fromEntropy(seed)

    # 派生以太坊子私钥
    path = "m/44'/60'/0'/0/0"
    key = bip32_key.ChildKey(
        BIP32_HARDEN + 44
    ).ChildKey(
        BIP32_HARDEN + 60
    ).ChildKey(
        BIP32_HARDEN + 0
    ).ChildKey(
        0
    ).ChildKey(
        0
    )

    # 获取以太坊私钥
    eth_private_key = key.WalletImportFormat()

    # 获取以太坊地址
    eth_address = privtoaddr(key.PrivateKey().to_string()).hex()

    print("助记词:", mnemonic)
    print("私钥:", eth_private_key)
    print("地址:", eth_address)

eth_wallet()