# main.py
from web3 import Web3

# Setup
alchemy_url = "https://eth-sepolia.g.alchemy.com/v2/fzgfj4QuLlNEyn2LrLZsseBAClGdMnyP"
w3 = Web3(Web3.HTTPProvider(alchemy_url))

# Print if web3 is successfully connected11222
# print(w3.isConnected())

# Get the latest block number
latest_block = w3.eth.block_number
print(latest_block)

# Get the balance of an account
balance = w3.eth.get_balance('0x9fccf5325E52747e6d4E8EeE7A0473926D47228c')
print("balances = ",balance)

# Get the information of a transaction
# tx = w3.eth.get_transaction('0x5c504ed432cb51138bcf09aa5e8a410dd4a1e204ef84bfed1be16dfba1b22060')
# print(tx)