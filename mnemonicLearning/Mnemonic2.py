# 生成助记词

# 1.1.随机熵生成
# 首先生成一段随机熵（Entropy）。熵的长度可以是 128 到 256 位，并且是 32 的倍数。常见的熵长度有 128 位（12 个助记词）和 256 位（24 个助记词）。


import hashlib
import os
import bip39

# 1.1.随机熵生成
# entropy = os.urandom(16) # os模块的生成16位的随机字节
entropy = b'\xfb\xe50\xc5\x84\xb6\x06{.n\x12\x06SN\xb5\x82'
print(entropy)

"""
1.2.计算校验和
对熵进行 SHA-256 哈希计算，并取哈希值的前几位作为校验和。sub scrip table
校验和的长度取决于熵的长度。例如，128 位熵需要 4 位校验和（因为 128 / 32 = 4），256 位熵需要 8 位校验和
最后截取前4位出来即可
"""
entropy_sha256 = hashlib.sha256(entropy).digest()
print(entropy_sha256)
# print(bin(entropy_sha256[0])[2:].zfill(8))
entropy_bin = bin(entropy_sha256[0])[2:].zfill(8)[:4]
print(entropy_bin)
"""
1.3.组合熵和校验
将校验和附加到熵的末尾，形成一个新的二进制序列。这个序列的总长度为 (熵的长度 + 校验和的长度)。
"""
entropy_bits = ''.join([bin(byte)[2:].zfill(8)
                       for byte in entropy])  # 这一步是生成校验和
# 把校验和拼接到熵的末尾
combined_bits = entropy_bits+entropy_bin


"""
1.4 分割为助记词索引，等下去助记词库里面那助记词
"""
indices = [int(combined_bits[i:i + 11], 2)
           for i in range(0, len(combined_bits), 11)]
print(indices)  # 这是一个列表，里面有12个数字索引，

"""
1.5 映射为助记词，就是通过上面的索引列表去助记词库里面拿到对应的助记词
"""
wordlist = bip39.INDEX_TO_WORD_TABLE  # bip39.INDEX_TO_WORD_TABLE 返回出来的是2048个助记词单词组成的元组

mnemonic = ' '.join(wordlist[index] for index in indices)
print(mnemonic)
