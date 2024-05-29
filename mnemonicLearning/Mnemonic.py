import os
import hashlib
import bip39
def generate_mnemonic():
    # 1. 生成 128 位随机熵 (16 字节)
    # entropy = os.urandom(16)
    entropy = b'\x90nE\xbecy\x9bH@\x9e\x1e\x01\xff\xe5\x00?'
    print(entropy)
    
    # 2. 计算校验和 (SHA-256)
    hash_bytes = hashlib.sha256(entropy).digest() # sha256加密
    print(hash_bytes)
    # bin方法会转成0b开头的二进制，[2:]切掉前两个，
    # 使用 .zfill() 方法在字符串的左侧填充零字符（'0'），使其总长度达到 8。如果二进制字符串的长度小于 8，则在左侧填充零字符。例如，'1010' 变成了 '00001010'。
    # [:4]：使用切片操作 [:4]，截取填充后的二进制字符串的前四个字符。这样可以得到一个长度为 4 的子字符串。例如，'00001010' 变成了 '0000'
    checksum_bits = bin(hash_bytes[0])[2:].zfill(8)[:4]  # bin方法转成二进制，然后切片取前 4 位
    print(checksum_bits)

    # 3. 组合熵和校验和
    entropy_bits = ''.join([bin(byte)[2:].zfill(8) for byte in entropy])
    print(entropy_bits)
    combined_bits = entropy_bits + checksum_bits
    print(combined_bits)

    # 4. 分割为助记词索引
    indices = [int(combined_bits[i:i + 11], 2) for i in range(0, len(combined_bits), 11)]

    # 5. 映射为助记词
    wordlist = bip39.INDEX_TO_WORD_TABLE
    mnemonic = ' '.join([wordlist[index] for index in indices])

    return mnemonic


# 生成并打印助记词
mnemonic_phrase = generate_mnemonic()
print(f"Generated mnemonic phrase: {mnemonic_phrase}")