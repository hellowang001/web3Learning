from web3 import Web3
import json

# 设置连接和账户信息
private_key = "c535facd6873ca2b3718e3ede4f626ae126c99b4a26b4354da704d8dc78b43c1"  # 发送方私钥
from_address = "0xa3856a939A623EdBde8f908037d3F33FceBC5408"                       # 发送方地址
to_address = "0x38B59D6D4ef6A4991926Cf04c7c2092a0E86140F"                         # 接收方地址
eth_url = "https://eth-sepolia.g.alchemy.com/v2/fzgfj4QuLlNEyn2LrLZsseBAClGdMnyP" # API URL

# 连接到以太坊节点
web3 = Web3(Web3.HTTPProvider(eth_url))
if not web3.is_connected():
    print("Failed to connect to Ethereum node.")
    exit()

# 检查账户余额
balance = web3.eth.get_balance(from_address)
print(f"Account balance: {web3.from_wei(balance, 'ether')} ETH")

# 获取最新区块号
block = web3.eth.get_block('latest')
print(f"Latest block: {block.number}")

# 获取交易计数 (nonce)
nonce = web3.eth.get_transaction_count(from_address)
print(f"Nonce: {nonce}")

# 建议的 gas 价格
gas_price = web3.eth.gas_price
print(f"Gas price: {web3.from_wei(gas_price, 'gwei')} Gwei")

# 创建交易
transaction = {
    'to': to_address,
    # 'value': web3.toWei(0.0011, 'ether'),  # 发送的以太币数量
    'value': 1100000000000000,  # 发送的以太币数量
    'gas': 21320,  # Gas limit
    'gasPrice': 2000000000,  # Gas price
    'nonce': nonce,
    'data': web3.to_hex(b"Hello, World! python"),  # 附加的数据
    'chainId': web3.eth.chain_id  # 链 ID
}

# 对交易进行签名
signed_tx = web3.eth.account.sign_transaction(transaction, private_key)
print(f"Transaction hash: {signed_tx.hash.hex()}")
、
# 发送交易
tx_hash = web3.eth.send_raw_transaction(signed_tx.rawTransaction)
print(f"Transaction sent, hash: {tx_hash.hex()}")
# 第一笔交易：0x22af000c45f51c03e88ea21281ccd571359e653341aef2a7fadbdc96b41eee49
# 第二笔交易：0xa502a7a7f589555eb0ce9bbd2d48220c860bbd0b7bf33b97819932ac0b092d6e
# 第三笔交易：0x775b97666b65444f2bb1b2bad852dbb9b56020fdbab513336cdec46bad1c2a36
