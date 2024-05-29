import os
import hashlib
import bip39


# 1. 生成随机熵
entropy = os.urandom(16)


# 2. 计算校验和（sha-256）
hash_entropy=hashlib.sha256(entropy).digest() # sha256加密

# 转换二进制
checksum = bin(hash_entropy[0])[2:].zfill(8)[:4] # 这一步有点麻烦

# 3. 拼接二进制
entropy_bits = ''.join([bin(byte)[2:].zfill(8) for byte in entropy])

# 组合熵和校验和
combined_bits = entropy_bits+ checksum

# 4.分割为助记词索引
indices = [int(combined_bits[i:i+11],2)for i in range(0,len(combined_bits),11)]

# 5. 到词库里面去映射成助记词
wordlist = bip39.INDEX_TO_WORD_TABLE
mneminic= " ".join([wordlist[index] for index in indices])
print(mneminic)
