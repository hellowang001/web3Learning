import hashlib
from ecdsa import SECP256k1, SigningKey, VerifyingKey
from ecdsa.curves import Curve


class SchnorrSignObj:
    curve: Curve
    def __init__(self):
        # 椭圆曲线参数
        self.curve = SECP256k1 # curve就是椭圆曲线对象，SECP256k1就是常见的椭圆曲线

    def schnorr_sign(self, private_key, message): # 签名函数，传私钥和msg
        # 生成私钥和公钥
        signing_key = SigningKey.from_string(private_key, curve=self.curve) # 使用私钥和椭圆曲线对象生成一个签名对象
        verifying_key = signing_key.get_verifying_key() # 从签名秘钥对象获得公钥，verifying_key就是公钥，返回一个公钥key可以用来验证签名

        # 生成随机数 r
        r = SigningKey.generate(curve=self.curve).privkey.secret_multiplier # 生成随机数作为签名的一部分
        R = VerifyingKey.from_public_point(r * self.curve.generator, curve=self.curve) # 使用随机数r乘以椭圆曲线的生成元，得到签名中的挑战点R。

        # 计算哈希 e
        e = hashlib.sha256(R.to_string() + verifying_key.to_string() + message).digest() # 这个之前学过，sha256加密返回哈希计算结果的二进制表示，传入挑战点R+公钥+msg
        e = int.from_bytes(e, 'big') # 将字节对象转成大整数，big表示大端字节数

        # 计算签名 s
        s = (r + e * signing_key.privkey.secret_multiplier) % self.curve.order # 计算签名中的s，根据Schnorr签名算法的计算公式，将随机数r和哈希e与私钥的倍数相加，并对椭圆曲线的阶数取模。
        return R.to_string(), s.to_bytes(32, 'big') # 返回挑战点 R 的字符串，和签名s的32字节大端字节串表示。

    def schnorr_verify(self, public_key, message, signature): # 签名验证，传入公钥，msg， 签名signature
        # 解析签名
        R = VerifyingKey.from_string(signature[0], curve=self.curve) # 从签名中解析挑战点 R ，传入对应的椭圆曲线对象
        s = int.from_bytes(signature[1], 'big') # 从签名消息中解析签名 s大端字节串表示

        # 计算哈希 e
        e = hashlib.sha256(signature[0] + public_key.to_string() + message).digest() # sha256加密返回哈希计算结果的二进制表示，传入签名[0] + 公钥 + msg
        e = int.from_bytes(e, 'big') # 将字节对象转成大整数，big表示大端字节数

        # 验证签名
        sG = VerifyingKey.from_public_point(s * self.curve.generator, curve=self.curve) # 从公钥获取指针（初始化验证公钥对象），传入将生成元乘以 s 值，得到一个新的点。self.curve：表示使用相同的椭圆曲线对象进行计算。
        # R：是挑战点，它是一个椭圆曲线上的点，R.pubkey.point：表示挑战点 R 的公钥点。
        # e：是哈希值，它是一个整数public_key.pubkey.point：表示公钥的点，
        # public_key.pubkey.point：表示公钥的点。
        ReP = VerifyingKey.from_public_point(R.pubkey.point + e * public_key.pubkey.point, curve=self.curve) 
        return sG.to_string() == ReP.to_string() # 验证是否通过