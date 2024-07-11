# import requests
#
# url = "https://eth-sepolia.g.alchemy.com/v2/fzgfj4QuLlNEyn2LrLZsseBAClGdMnyP"
#
# payload = {
#     "id": 1,
#     "jsonrpc": "2.0",
#     "method": "eth_blockNumber"
# }
# headers = {
#     "accept": "application/json",
#     "content-type": "application/json"
# }
# 1
# response = requests.post(url, json=payload, headers=headers)
#
# print(response.text)

import requests


url = "https://eth-sepolia.g.alchemy.com/v2/fzgfj4QuLlNEyn2LrLZsseBAClGdMnyP"

payload = {
    "id": 1,
    "jsonrpc": "2.0",
    "params": ["0x9fccf5325E52747e6d4E8EeE7A0473926D47228c", "latest"],
    "method": "eth_getTransactionCount"
}
headers = {
    "accept": "application/json",
    "content-type": "application/json"
}

response = requests.post(url, json=payload, headers=headers)

print("addresNonecs=",response.text)