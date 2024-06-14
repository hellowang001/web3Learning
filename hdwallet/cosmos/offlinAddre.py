import hashlib
import bech32
from secp256k1 import PrivateKey
from ripemd.ripemd160 import ripemd160,RIPEMD160Hash
def generate_wallet():
    privkey = PrivateKey().serialize()
    return {
        "private_key": privkey,
        "public_key": privkey_to_pubkey(privkey),
        "address": privkey_to_address(privkey),
    }

def privkey_to_pubkey(privkey: str) -> str:
    privkey_obj = PrivateKey(bytes.fromhex(privkey))
    return privkey_obj.pubkey.serialize().hex()

def pubkey_to_address(pubkey: str) -> str:
    pubkey_bytes = bytes.fromhex(pubkey)
    s = hashlib.new("sha256", pubkey_bytes).digest()
    # print("s",type(s))
    # r = hashlib.new("ripemd160", s).digest()
    # ripemd160()
    r2=RIPEMD160Hash().new(s).digest()
    return bech32.bech32_encode("cosmos", bech32.convertbits(r2, 8, 5))

def privkey_to_address(privkey: str) -> str:
    pubkey = privkey_to_pubkey(privkey)
    return pubkey_to_address(pubkey)
# if __name__ == "__main__":
    # wallet = generate_wallet()

    # print(wallet)
    # print(privkey_to_pubkey(wallet[""]))