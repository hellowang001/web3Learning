# -*- coding: utf-8 -*-
import bip32
import bip39, os, hashlib
# import bip32
# from bip32 import bip32
# from bip32.bip32 import BIP32
from bip32 import BIP32
# b = BIP32()
from ecdsa import SECP256k1, SigningKey
from web3 import Web3
from eth_account import Account


# 生成助记词
class Bip39Mnemonic(object):

    def privateKeyToAccount(self, privateKey):
        account = Account.from_key(privateKey)
        return account

    def mnemonic_to_account(self, mnemonic):
        Account.enable_unaudited_hdwallet_features()
        account = Account.from_mnemonic(mnemonic)
        return account

    # 生成随机助记词
    def generateMnemonic(self):
        """
        随机生成助记词
        :return:
        """
        # 1. 生成 128 位随机熵 (16 字节)->12位助记词
        entropy = os.urandom(16)

        # 2. 计算校验和 (SHA-256)
        hash_bytes = hashlib.sha256(entropy).digest()
        checksum_bits = bin(hash_bytes[0])[2:].zfill(8)[:4]  # 取前 4 位

        # 3. 组合熵和校验和
        entropy_bits = ''.join([bin(byte)[2:].zfill(8) for byte in entropy])
        combined_bits = entropy_bits + checksum_bits

        # 4. 分割为助记词索引
        indices = [int(combined_bits[i:i + 11], 2) for i in range(0, len(combined_bits), 11)]

        # 5. 映射为助记词
        wordlist = bip39.INDEX_TO_WORD_TABLE
        mnemonic = ' '.join([wordlist[index] for index in indices])

        return mnemonic
        pass


# 助记词拿公钥

# 拿地址

if __name__ == '__main__':
    bp = Bip39Mnemonic()
    test_mne = bp.generateMnemonic()
    account = bp.mnemonic_to_account(test_mne)
    print(account)
    a = account.address

    print(a)
    k=account.key
    print(k)
