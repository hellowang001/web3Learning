import requests
from passlib.context import CryptContext

pwd_context = CryptContext(schemes=['bcrypt'], deprecated='auto')


def get_password_hash(password: str) -> str:
    """
    生成哈希密码
    :param password: 原始密码
    :return: 哈希密码
    """
    return pwd_context.hash(password)


pas = get_password_hash("a123456")

s = requests.session()

# URL = "http://127.0.0.1:9000/auth/api/login"
URL = "http://127.0.0.1:9000/auth/login"
header={
    'Content-Type': 'application/json',
    'Accept': 'application/json',
}
req_data = {
    "telephone": "15812345678",
    "password": pas,
    "method": "0",
    "platform": "0"
}
# req_data = {
#     "telephone": "15020221010",
#     "password": "kinit2022",
#     "method": "0",
#     "platform": "0"
# }
users_URL = "http://127.0.0.1:9000/vadmin/auth/users"

if __name__ == '__main__':
    print(req_data)

    res = s.post(url=URL, json=req_data, headers=header)
    print(res.json())
