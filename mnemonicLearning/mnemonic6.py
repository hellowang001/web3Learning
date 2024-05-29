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
#首先，hashlib.sha256(entropy)创建一个SHA-256哈希算法对象，并将entropy作为输入进行哈希计算。
# 然后，.digest()方法被调用，返回计算得到的哈希结果的二进制表示。这个二进制表示是一个字节数组（bytes），长度为32字节（256位），表示SHA-256算法对输入数据的哈希结果。
hash_bytes = hashlib.sha256(entropy).digest()
print(hash_bytes) # b's\xf7\xc9z\x88\xf0\x7f\x0f\xf0`(\xe4\xd6\x92@\xa9?\xa2os~\xbe*\x86\xe7w\x97D\xf8Z"\xb4'
# entropy_sha256[0]表示entropy_sha256的第一个字节，它是一个整数值。
# bin(entropy_sha256[0])将该整数值转换为对应的二进制字符串，返回的字符串以0b开头，例如0b11001010。
# [2:]是对二进制字符串进行切片操作，去掉开头的0b，例如11001010。
# .zfill(8)将二进制字符串填充到8位，不足的位数在开头补零，例如11001010变为11001010。
# [:4]是对二进制字符串进行切片操作，截取前4位，例如11001010变为1100。
# 结果被赋值给变量entropy_bin，表示对entropy_sha256的第一个字节进行处理后得到的二进制字符串，长度为4位。
checksum_bits = bin(hash_bytes[0])[2:].zfill(8)[:4]
print(checksum_bits) # 0111
"""
1.3.组合熵和校验
将校验和附加到熵的末尾，形成一个新的二进制序列。这个序列的总长度为 (熵的长度 + 校验和的长度)。
"""
# 这里bin(byte)后的[2:]是切掉0b开头，zfill(8)是填充，如果长度没有8位就填充0
entropy_bits = ''.join([bin(byte)[2:].zfill(8)
                       for byte in entropy])  # 这一步是生成校验和
print(entropy_bits)
# 把校验和拼接到熵的末尾
combined_bits = entropy_bits + checksum_bits
print(combined_bits)

"""
1.4 分割为助记词索引，等下去助记词库里面那助记词
"""
# range(0, len(combined_bits), 11)创建一个迭代范围，从0开始，每次增加11，直到达到combined_bits的长度之前的最大值。
# combined_bits[i:i + 11]是对combined_bits进行切片，从索引i开始，切片长度为11，得到一个11位的二进制字符串。
# int(combined_bits[i:i + 11], 2)将切片得到的二进制字符串作为参数传递给int()函数，并指定进制为2，将二进制字符串转换为对应的整数值。
indices = [int(combined_bits[i:i + 11], 2)
           for i in range(0, len(combined_bits), 11)] # 11是步长，[i:i+11]表示往后面取11位
print(indices)  # 这是一个列表，里面有12个数字索引，

"""
1.5 映射为助记词，就是通过上面的索引列表去助记词库里面拿到对应的助记词
"""
wordlist = bip39.INDEX_TO_WORD_TABLE  # bip39.INDEX_TO_WORD_TABLE 返回出来的是2048个助记词单词组成的元组

mnemonic = ' '.join(wordlist[index] for index in indices)
print(mnemonic)
print(mnemonic)
print(mnemonic)
print(mnemonic)
print(mnemonic)


