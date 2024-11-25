from email import header

import requests
from passlib.context import CryptContext

pwd_context = CryptContext(schemes=['bcrypt'], deprecated='auto')

s = requests.session()


def get_password_hash(password: str) -> str:
    """
    生成哈希密码
    :param password: 原始密码
    :return: 哈希密码
    """
    return pwd_context.hash(password)


pas = get_password_hash("a123456")

URL = "http://127.0.0.1:9000/vadmin/auth/users"
header_test = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
    "Authorization": "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxNTAyMDIyMTAxMCIsImlzX3JlZnJlc2giOmZhbHNlLCJwYXNzd29yZCI6IiQyYiQxMiRDZTdlU1VLSUlsOERNS2VEeU5IeXIuRHA0YWVzUUNNNzBSZVBpZ1JWRW55MUVxbDMxUjBDcSIsImV4cCI6MTcyODM4MDI4NH0.38apKtbH5W_XeYp2Xs0iXYqbyLrYZBkVUpgIKjgEAwI"

}
req_data = {
    "name": "test123",
    "telephone": "15812345678",
    "email": "123456@qq.com",
    "nickname": "test12345",
    "avatar": None,
    "is_active": True,
    "is_staff": True,
    "gender": "0",
    "is_wx_server_openid": False,
    "role_ids": [1],
    "dept_ids": [1],
    "password": pas
}

if __name__ == '__main__':
    print(req_data)

    res = s.post(url=URL, json=req_data, headers=header_test)
    print(res.json())
