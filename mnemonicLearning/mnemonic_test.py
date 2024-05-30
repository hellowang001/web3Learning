# -*- coding: utf-8 -*-
import hashlib
import os

# 生成12个单词的助记词，
# 12个对应的是128位熵，然后助记词每加3个就加32为

# 1、随机熵生成，
# 用os库，先生成对应的位数
entropy = os.urandom(12)
print(entropy)
# 2、计算校验和，
# 先对熵进行SHA-256的哈希计算，并取前4位做校验和，
hash_byte=hashlib.sha256(entropy).digest()
print(hash_byte)
# 这个前4位是怎么来的，就是128÷32=4

# 3、组合熵和校验和

# 4、分割为助记词索引

# 5、映射为助记词
