import os
import hashlib
import bip39


# 1.生成随机熵
entropy = os.urandom(16)

# 2.计算校验和
# 2.1 先加密转化
hash_bytes = hashlib.sha256(entropy).digest()

# 2.2 转化二进制且填充0
checksum_bits = bin(hash_bytes[0])[2:].zfill(8)[:4]

# 3 组合熵和校验和
# 3.1 先遍历后拼接字符串
entropy_bits = ''.join([bin(byte)[2:].zfill(8) for byte in entropy])

# 3.2 然后组合进去
combined_bits = entropy_bits + checksum_bits

# 4.分割为助记词索引
indices = [int(combined_bits[i:i+11], 2) for i in range(0, len(combined_bits), 11)]


# 5.去助记词库里面映射出助记词
wordlist = bip39.INDEX_TO_WORD_TABLE
mnemonic = ' '.join([wordlist[index] for index in indices])

print(mnemonic)
