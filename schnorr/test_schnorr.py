import schnorr as schnorr
from schnorr import SchnorrSignObj
from ecdsa import SECP256k1, SigningKey, VerifyingKey

# 初始化类
schnorr_test = schnorr.SchnorrSignObj()

# 产生密钥并对交易签名
private_key = SigningKey.generate(curve=SECP256k1).to_string()
message = b"Hello, Schnorr!2222"
signature = schnorr_test.schnorr_sign(private_key, message)
print("Signature:", signature)

# 验证交易签名
public_key = VerifyingKey.from_string(
    SigningKey.from_string(private_key, curve=SECP256k1).get_verifying_key().to_string(), curve=SECP256k1)
is_valid = schnorr_test.schnorr_verify(public_key, message, signature)

print("Signature valid:", is_valid)
