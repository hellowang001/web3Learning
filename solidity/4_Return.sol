// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7 <0.9;

contract Return {
    // 这里最后一个返回值是一个数组，这里为什么要标明memory呢？因为数组要标明他的内存地址，相当于go的指针，因为数组不是值，是指针的合集
    // memory在Solidity中是一个关键字，它用来定义变量的数据位置。Solidity有三种数据位置：storage，memory和calldata。

    // storage：这是合约的状态变量被持久性存储的地方。存储在storage中的变量会在区块链上持久化，即使在外部函数调用结束后也不会消失。

    // memory：这是一个临时存储区域，用于存储函数调用期间的变量。一旦函数调用结束，存储在memory中的变量就会被清除。

    // calldata：这是一个只读存储区域，只能用于外部函数调用的参数。

    // 在你的代码中，uint256[3] memory _array表示_array是一个存储在内存中的数组，它的长度为3，每个元素的类型为uint256。这意味着_array只会在returnNamed函数调用期间存在，一旦函数调用结束，_array就会被清除
    function returnMultiple() public pure returns(uint256, bool, uint256[3] memory){
        return(1, true, [uint256(1),2,5]); 
        // 没有名字的返回，我们调用的时候显示：
        // 0:uint256: 1
        // 1:bool: true
        // 2:uint256[3]: 1,2,5

    }

    // 2、命名式返回,可以在returns 中标明返回变量的名称，Solidity会初始化这些变量，无需在代码中特别写明 return
        function returnNamed()public pure returns (uint256 _number,bool _bool,uint256[3] memory _array  ) {
        _number =2;
        _bool=false;
        _array=[uint256(3),2,1];
        // 有名字的返回，调用的时候显示：
        // 0: uint256: _number 2
        // 1: bool: _bool false
        // 2: uint256[3]: _array 3,2,1
    }
    // 2.2 命名式返回，当然，你也可以在命名式返回中使用return
    function returnNamed2()public pure returns (uint256 _number,bool _bool,uint256[3] memory _array  ){
        return(1, true, [uint256(1),2,5]);
        // 有名字的返回，调用的时候显示：
        // 0: uint256: _number 2
        // 1: bool: _bool true
        // 2: uint256[3]: _array 1,2,5
    }
    // 3 读取所有返回值，解构式赋值
    function readReturn()public pure{
        uint256 _number;
        bool _bool;
        bool _bool2;
        uint256[3] memory _array;
        (_number,_bool,_array)=returnNamed(); // 用着三个值去接收方法的返回值
        (,_bool2,)=returnNamed();// 当然也可以不使用
        // 当然这个方法没有返回，所以调用的时候啥都不显示，只是演示函数的调用而已
    }

}