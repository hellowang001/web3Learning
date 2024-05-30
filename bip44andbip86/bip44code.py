# -*- coding: utf-8 -*-

import bip39
import bech32
from bip32 import BIP32, HARDENED_INDEX
from hashlib import sha256
from ecdsa import SECP256k1, SigningKey
def create_address(mnemonic):
    seed = bip39.phrase_to_seed(mnemonic)
    bip32_root_key = BIP32Key.fromEntropy(seed)
    path = "m/44'/0'/0'/0/0"
    key = bip32_root_key
    for level in path.split('/')[1:]:
        if level.endswith("'"):
            index = BIP32_HARDEN + int(level[:-1])
        else:
            index = int(level)
        key = key.ChildKey(index)
    address = key.Address()
    return address

# ERROR: Could not build wheels for coincurve, which is required to install pyproject.toml-based projects