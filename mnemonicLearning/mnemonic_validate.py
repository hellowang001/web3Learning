import bip39
def validate_mnemonic(mnemonic, wordlist):
    words = mnemonic.split()
    
    # 检查单词数量
    if len(words) not in [12, 15, 18, 21, 24]:
        return False
    
    # 检查单词是否在词汇表中
    for word in words:
        if word not in wordlist:
            return False
    
    # 将助记词转化成位串
    binary_string = ''
    for word in words:
        index = wordlist.index(word)
        binary_string += format(index, '011b') # format(index, '011b') 的含义是将 index 格式化为 11 位的二进制字符串，不足 11 位的部分在前面用零填充。
    print(binary_string)
    # 提取种子和校验和
    # 计算种子部分的长度。根据给定的单词列表 words 的长度，每个单词占用 11 位的二进制字符串，除以 3 取整数部分表示校验和的长度。通过减去校验和的长度，得到种子部分的长度。
    seed_bits_length = (len(words) * 11) - (len(words) // 3) # 
    # 从 binary_string 中提取种子部分。使用切片操作 [:seed_bits_length]，获取从开头到种子部分长度的子字符串。
    seed_bits = binary_string[:seed_bits_length]
    # 从 binary_string 中提取校验和部分。使用切片操作 [seed_bits_length:]，获取从种子部分长度开始到末尾的子字符串。
    checksum_bits = binary_string[seed_bits_length:]
    
    
    # 计算校验和
    import hashlib
    # 将种子部分的二进制字符串 seed_bits 转换为字节串。
    # 首先，使用 int(seed_bits, 2) 将二进制字符串转换为整数。
    # 然后，使用 to_bytes(len(seed_bits) // 8, byteorder='big') 将整数转换为字节串，
    # 其中 len(seed_bits) // 8 表示字节串的长度，byteorder='big' 表示使用大端字节序。
    seed_bytes = int(seed_bits, 2).to_bytes(len(seed_bits) // 8, byteorder='big')
    # 对种子字节串 seed_bytes 进行 SHA-256 哈希计算，并获取哈希结果的十六进制表示。
    # 使用 hashlib.sha256() 创建 SHA-256 哈希算法对象，然后调用 .hexdigest() 方法获取哈希结果的十六进制字符串表示。
    hash_value = hashlib.sha256(seed_bytes).hexdigest()
    # 将哈希结果的十六进制字符串 hash_value 转换为 256 位的二进制字符串。
    # 首先，使用 int(hash_value, 16) 将十六进制字符串转换为整数。
    # 然后，使用 bin() 将整数转换为二进制字符串，去掉开头的 '0b'，并使用 zfill(256) 在前面填充零，使其长度达到 256 位。
    hash_bits = bin(int(hash_value, 16))[2:].zfill(256)
    # 提取校验和部分。使用切片操作 [:len(words) // 3]，获取从二进制字符串 hash_bits 的开头到校验和长度的子字符串。校验和长度为单词列表 words 的长度除以 3。
    calculated_checksum = hash_bits[:len(words) // 3]
    
    # 验证校验和，判断助记词转换的位串与计算得来的校验和一致
    return checksum_bits == calculated_checksum

# Example usage:
mnemonic = "legal winner thank year wave sausage worth useful legal winner thank yellow"
wordlist = bip39.INDEX_TO_WORD_TABLE  # BIP-39 wordlist
is_valid = validate_mnemonic(mnemonic, wordlist)
print("Is valid mnemonic:", is_valid)