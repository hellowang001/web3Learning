
import bip39,os,hashlib
class Bip39Mnemonic:
    def __init__(self):
        pass

    def createMnemonic(self, number):
        mnemonic = bip39.get_entropy_bits(number)
        return mnemonic

    def mnemonicToEntropy(self,  mnemonic):
        decode_words = bip39.decode_phrase( mnemonic)
        return decode_words
    def entropyToMnemonic(self, entropy):
        nemonic_entropy = bip39.encode_bytes(entropy)
        return nemonic_entropy

    def mnemonicToSeed(self, mnemonic):
        nemonic_to_seed = bip39.phrase_to_seed(mnemonic)
        return nemonic_to_seed

    def validateMnemonic(self, mnemonic):
        return bip39.check_phrase(mnemonic)

    def generateMnemonic(self):
        # 1. 生成 128 位随机熵 (16 字节)
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
    
if __name__=='__main__':
    print("hhh")
    mne = Bip39Mnemonic()
    mnemonicmnemonic_phrase=mne.generateMnemonic()
    print(f"Generated mnemonic phrase: {mnemonic_phrase}")
    # print(mnemonic)
    mnemonic_12mnemonic_phrase = mne.createMnemonic(12)
    print(f"create mnemonic phrase: {mnemonic_phrase}")
    print(mnemonic_12)