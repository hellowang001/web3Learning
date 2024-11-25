"""
Example no.1
============

Just an example.

AUTHOR
    mgpai22@GitHub

CREATED AT
    Wed. 20 Apr. 2022 23:43
"""
# Import packages
import json
import logging

# Import `ergpy`
from ergpy import helper_functions, appkit

# Logging utility
LOGGING_FORMAT = '[%(asctime)s] - [%(levelname)-8s] -  %(message)s'
logging.basicConfig(format=LOGGING_FORMAT)
logger: logging.Logger = logging.getLogger()
logger.setLevel(logging.INFO)

# Create connection to the blockchain
node_url: str = "http://213.239.193.208:9052/" # MainNet or TestNet
ergo = appkit.ErgoAppKit(node_url=node_url)

# Wallet mnemonic
# wallet_mnemonic = "decline reward asthma enter three clean borrow repeat identify wisdom horn pull entire adapt neglect"
wallet_mnemonic = "copy crowd youth comic hello atom amateur jewel climb concert rule scissors ostrich shoulder visa"

# (optional) If you have a mnemonic password
# mnemonicPassword = "my password"

"""
Note how the following parameters are inputed :
- In simple send, 3WwdX would get 0.22 ERG, 3WwuG9 would get 0.33 ERG, 3Wxk5o would get 0.11 ERG
- In token send, 3WwdX would get 0.22 ERG and tokens A, B, C, 3WwuG9 would get 0.33 ERG and token D and so on
注意以下参数是如何输入的:
-在简单发送中，3WwdX将获得0.22 ERG, 3WwuG9将获得0.33 ERG, 3wxk50将获得0.11 ERG 
-在令牌发送中，3WwdX将获得0.22 ERG，令牌A, B, C, 3WwuG9将获得0.33 ERG和令牌D，以此类推“”
"""
# receiver_addresses = [
#     "3WwdXmYP39DLmDWJ6grH9ArXbWuCt2uGAh46VTfeGPrHKJJY6cSJ",
#     "3WwuG9amNVDwkJdgT5Ce7aJCfeoafVmd9tag9AEiAZwgPi7pYX3w",
#     "3Wxk5oofZ3Laq2CpFW4Fi9YQiaep9bZr6QFg4s4xpzz4bi9tZq2U"
# ]
receiver_addresses = [
    "3WwdXmYP39DLmDWJ6grH9ArXbWuCt2uGAh46VTfeGPrHKJJY6cSJ",
]

# amount = [0.22, 0.33, 0.11]
amount = [0.22]

tokens = [
    ["tokenID_A", "tokenID_B", "tokenID_C"],
    ["tokenID_D"],
    ["tokenID_E", "tokenID_F"]
]

nft_name = "pythonWrapper NFT"
description = "created by MGpai's python wrapper"
image_link = "ipfs://bafkreihvadi66kddokso7x2i4kmqjhj5e3j44rn347i7drghekr2mfw3nu"
image_hash = appkit.sha256caster("f500d1ef286372a4efdf48e299049d3d26d3ce45bbe7d1f1c4c722a3a616db6d")


"""
base64 reduced transactions are also possible, the only current problem is the boxes are preselected so
if another tx gets processed before submitting the base64 tx, it will not go through.
Just add the `base64reduced=True` to get this

Example
Base64减少的事务也是可能的，目前唯一的问题是盒子是预先选择的，所以如果在提交Base64 tx之前处理另一个tx，它将不会通过。只需添加' base64reduced=True '即可获得此示例
-------
```
print(helperMethods.simpleSend(ergo=ergo,
    amount=amount,
    walletMnemonic=wallet_mnemonic,
    receiverAddresses=receiver_addresses,
    base64reduced=True
))
````
Note
----
    Specifying `senderAddress` is also possible. Once specified, the boxes will be selected from this address and change
    will be sent here as well. Make sure the prover index is specified as well (same index as address derivation).

"""

# Assertions / tests
# assert json.loads(helper_functions.get_wallet_address(ergo=ergo, amount=1, wallet_mnemonic=wallet_mnemonic))[0] == "3Wx2YrSVcrPvC7uXQRp6ZQfRd7VxjZr6fjhFEX5r1yiM8nHkGv93"
# logging.info("Tests passed")

# Logging to console
logging.info(f'Result: {helper_functions.get_wallet_address(ergo=ergo, amount=1, wallet_mnemonic=wallet_mnemonic)}')
# print(helper_functions.get_box_info(ergo=ergo, index=0, sender_address=receiver_addresses[0]))
print(helper_functions.simple_send(ergo=ergo, amount=amount, wallet_mnemonic=wallet_mnemonic, receiver_addresses=receiver_addresses))
# print(helper_functions.send_token(ergo=ergo, amount=amount, receiver_addresses=receiver_addresses, tokens=tokens, wallet_mnemonic=wallet_mnemonic))
# print(helper_functions.create_token(ergo=ergo, token_name=nft_name, description=description, token_amount=1, token_decimals=0, wallet_mnemonic=wallet_mnemonic))
# print(helper_functions.create_nft(ergo=ergo, nft_name=nft_name, description=description, image_link=image_link, image_hash=image_hash, wallet_mnemonic=wallet_mnemonic))

# Proper exit()
helper_functions.exit()