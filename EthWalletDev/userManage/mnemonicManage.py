# -*- coding: utf-8 -*-

import bip39, os, hashlib
# import bip32
# from bip32 import bip32
from bip32.bip32 import BIP32
# from bip32 import BIP32


# 生成助记词
class Bip39Mnemonic(object):

    # 种子短语生成公钥
    def seedtoPub(self,seed):
        # bip32.from_seed
        # bip=BIP32.from_seed(seed=seed)

        b = BIP32()
        return None

    # 助记词转种子短语
    def mnemonicToSeed(self, phrase):
        seed = bip39.phrase_to_seed(phrase=phrase)
        return seed

        pass


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
    print(test_mne)
    test_seed = bp.mnemonicToSeed(phrase=test_mne)
    print(test_seed)
    ttt=bp.seedtoPub(seed=test_seed)
    print(ttt)
