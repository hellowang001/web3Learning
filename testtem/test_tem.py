import logging
import time

import requests
import logger

s = requests.session()

URL = 'http://www.baidu.com'
URL_CSND = 'https://www.csdn.net/'


# resp = s.get(url=URL)

# resp.text

# print(resp.text)


def print_name(fun):
    """
    这个函数的作用是打印当前运行的函数的名字
    :param fun:是一个函数
    :return:fun()运行的结果
    """

    def w_fun():
        logger.logger.info(f"当前运行的函数的名字是={fun.__name__}")

        # print(fun.__name__)

        a = fun()
        return a

    return w_fun


def print_time(fun):
    def w_time():
        start = time.time()
        logger.logger.info(f"调用get_baidu前的时间是{start}")
        r = fun()
        end = time.time()
        logger.logger.info(f"调用get_baidu后的时间是{end}")
        logger.logger.info(f"调用get_baidu 所花费的总时间是{end - start}")
        return r

    return w_time


@print_time
def get_baidu():
    start = time.time()
    logger.logger.info(f"调用get_baidu前的时间是{start}")
    time.sleep(1)
    resp = s.get(url=URL)
    r = resp.text
    end = time.time()
    logger.logger.info(f"调用get_baidu后的时间是{end}")
    logger.logger.info(f"调用get_baidu 所花费的总时间是{end - start}")

    return r


@print_time
def get_test():
    time.sleep(1)
    resp = s.get(url=URL_CSND)

    return resp.text


@print_time
@print_name
def any_test():
    time.sleep(1)
    return 0


@print_name
def test_ziyong():
    a = 1
    return a


class Test_a:
    def __init__(self, value):
        self.value = value

    def test(self):
        return 1

    def __str__(self):
        return 'test'

    def __sub__(self, other):
        print('运行了减法这个方法,sub')

        return self.value - other.value


if __name__ == '__main__':
    # r = get_baidu()
    # print(f" 打印get_baidu函数响应回来的结果 {r},----时间戳是：-{time.time()}-- ")
    # r = get_test()
    # r = any_test()
    # r = test_ziyong()
    # r = test_ziyong()
    #
    # print(r)
    # list()
    list_a = list((1, 2, 3, 4, 6))  # [1,2,3,4]
    # for i in list_a:
    #     print(i)
    # print(list_a.__class__)
    # list_b = list((1, 2, 3, 4, 6))
    # list_c = list_a + list_b
    # list_a = list_b
    # # b = 1
    # b = int(1)
    c = int(2)
    # d = b + c
    # x = dict(a =1)
    # y = set()
    # z = tuple()

    # for i in b:
    #     print(i)
    # print(b)
    # c = Test_a()
    # print(c)
    # a = Test_a(1000)
    # b = Test_a(1)
    # print(a-b)
    a = "test"
    b = "b"
    # c = str("c")
    print(a+b)
