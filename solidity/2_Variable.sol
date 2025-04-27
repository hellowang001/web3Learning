// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract Variable {
    // 1、布尔值
    bool public _bool = true; 
    //类型 公有 值名称=值
    // 1.1、布尔运算
    bool public _bool1 = !_bool; // 取非
    bool public _bool2 = _bool && _bool1; // 与
    bool public _bool3 = _bool || _bool1; // 或
    bool public _bool4 = _bool == _bool1; // 相等
    bool public _bool5 = _bool != _bool1; // 不相等

    // 2、整型
    int public _int = -1; // 整数，包括负数
    uint public _uint = 1; // 无符号整数
    uint256 public _number = 20220330; // 256位无符号整数

    // 2.1\整数运算
    uint256 public _number1 = _number + 1; // +，-，*，/
    uint256 public _number2 = 2**2; // 指数
    uint256 public _number3 = 7 % 2; // 取余数
    bool public _numberbool = _number2 > _number3; // 比大小

    //3、地址
    address public _address = 0x7A58c0Be72BE218B41C608b7Fe7C5bB630736C71;
    address payable public _address1 = payable(_address); // payable address，可以转账、查余额
    // 地址类型的成员
    uint256 public balance = _address1.balance; // balance of address

    // 4、固定长度的字节数组
    bytes32 public _byte32 = "MiniSolidity"; // 字节形式
    bytes1 public _byte = _byte32[0]; // 取其第一个字节

    // 5、枚举
    // 用enum将uint 0， 1， 2表示为Buy, Hold, Sell
    enum ActionSet { Buy, Hold, Sell }
    // 创建enum变量 action
    ActionSet action = ActionSet.Buy;

        // enum可以和uint显式的转换
    function enumToUint() external view returns(uint){
        return uint(action);
    }
}