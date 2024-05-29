import bip

import bip

bip39_mnemonic = bip.Bip39Mnemonic()
mnemonic_phrase = bip39_mnemonic.generateMnemonic()
print(f"Generated mnemonic phrase: {mnemonic_phrase}")

mnemonic_12_phrase = bip39_mnemonic.createMnemonic(12)
print(f"create mnemonic phrase: {mnemonic_12_phrase}")