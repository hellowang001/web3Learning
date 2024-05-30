import bip39, os, hashlib


class Bip39Mnemonic:
    def __init__(self):
        pass

    def createMnemonic(self, number: int):
        """
        指定字节数生成熵
        :param number: [12,15,18,21,24]
        :return: 熵位数，12->128,15->160.....
        """
        entropy_bits = bip39.get_entropy_bits(number)
        return entropy_bits

    def mnemonicToEntropy(self, mnemonic):
        """
        将助记词短语转化成字节序列
        :param mnemonic:
        :return:
        """
        decode_words = bip39.decode_phrase(mnemonic)
        return decode_words

    def entropyToMnemonic(self, entropy):
        """
        和上面相反，将字节序列转化成对应的助记词短语？
        :param entropy:
        :return:
        """
        mnemonic_entropy = bip39.encode_bytes(entropy)
        return mnemonic_entropy

    def mnemonicToSeed(self, mnemonic):
        """
        将助记词转化成种子短语？
        :param mnemonic:
        :return:
        """
        mnemonic_to_seed = bip39.phrase_to_seed(mnemonic)
        return mnemonic_to_seed

    def validateMnemonic(self, mnemonic):
        """
        验证助记词是否符合规则
        :param mnemonic:
        :return:bool
        """
        return bip39.check_phrase(mnemonic)

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


if __name__ == '__main__':
    print("hhh")
    mne = Bip39Mnemonic()

    # mnemonic = mne.createMnemonic()
    # print(mnemonic)
    mnemonic = mne.generateMnemonic()
    print(mnemonic)
    print(mne.validateMnemonic(mnemonic))
    m2 = mne.mnemonicToSeed(mnemonic)
    print(m2)
    entropy = os.urandom(16)
    m3 = mne.entropyToMnemonic(entropy)
    print(m3)
    # m1 = mne.createMnemonic(12)
    # print(m1)
    # mnemonicmnemonic_phrase = mne.generateMnemonic()
    # print(f"Generated mnemonic phrase: {mnemonic_phrase}")
    # # print(mnemonic)
    # mnemonic_12mnemonic_phrase = mne.createMnemonic(12)
    # print(f"create mnemonic phrase: {mnemonic_phrase}")
    # print(mnemonic_12)
