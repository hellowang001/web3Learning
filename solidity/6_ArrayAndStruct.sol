// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7 <0.9;

contract ArrayAndStruct {
    // 数组 Array 分为定长数组和可变长数组两种，类似go的数组和切片
    // 定长数组：声明的实际必须标明长度,T[k] T是type类型，k是长度
    uint[8] array1;// 长度为8的数组，类型是uint
    bytes1[5] array2;
    address[10] array3;
    // 可变长度数组，也是动态数组，在声明的时候不指定数组长度，T[]
    uint[] array4;
    bytes1[] array5;
    address[] array6;
    bytes array7; // 注意bytes 已经标明了，不需要加[]
    // 2 创建数组
    // 对于menory修饰的动态数组，可以用new 操作符来创建，但是必须声明长度，且声明长度后不可以发送改变
    // memory动态数组

    // uint[] memory array8 = new uint[](5); // 这里报错，因为memory类型的变量只能在函数内部声明，
    // bytes memory array9 = new bytes(9); // 这里报错，因为memory类型的变量只能在函数内部声明，
    function c()public pure  {
        uint[] memory array8 = new uint[](5); // 在函数里面就必须要写memory或者calldata
        array8[0]=1; // 动态数组需要一个一个赋值
        array8[1]=2; // 动态数组需要一个一个赋值
        array8[2]=3; // 动态数组需要一个一个赋值

        bytes memory array9 = new bytes(9);
        array9[0]=0x00;
    }
    // 数组字面常熟 市协作表达式形式的数据，你只需要写明数组的第一个元素的类型，后面的元素会根据上下文推断
    function f()public pure  {
        // uint[] memory a;
        // a= [uint(1),2,3];
        g([uint(1),2,3]);
    }

    function g(uint[3] memory _data) public pure {
        // ...
    }
    // 3 数组成员，指的是数组类型的数据包含哪些属性和方法
    // 3.1 length 长度，
    // 3.2 push() 动态数组拥有push方法，可以在数组最后添加一个0元素，并返回该元素的引用。
    // 3,3 push(x) 可以在数据最后一个添加 x 元素 ； 
    // 3.4 pop() 移除数组最后一个元素
    function arrayPush()public  returns (uint[]memory) {
        uint[2] memory a =[uint(1),2];
        array4=a;
        array4.push(3);
        return array4;
        //{"0": "uint256[]: 1,2,3"}
    }
    function arrayLen()public  returns (uint) {
        uint[2] memory a =[uint(1),2];
        array4=a;
        return array4.length;  // {"0": "uint256: 2"}
    }
    // 4 结构体，顾名思义，各种数据类型组合的结构体，结构体中的元素可以是原始类型，也可以是引用类型；结构体可以作为数组或映射的元素。创建结构体的方法：
    struct Student{
        uint256 id;
        uint256 score;
    }
    Student student; // 舒适化这个结构体，待会下面会用到的
    // 4.1 方法1: 给结构体赋值方式一：
    function initStudent1()external  {
        Student storage _student=student ;// storage状态变量，链上存储
        _student.id=11;
        _student.score=100; // debug的时候会发现，当程序走到这里的时候 student的属性发生了改变
    }
    // 4.2 方法2:直接引用状态变量的struct
    function initStudent2() external{   
        student.id = 1;
        student.score = 80;
    }
    // 4.3 方法3:构造函数式，类似与类的初始化
    function initStudent3() external {
        student = Student(3, 90);
    }
    // 方法4:key value
    function initStudent4() external {
        student =Student(
            {
                id:4,
                score:50
            }
        );
    }
}