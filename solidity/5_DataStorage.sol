// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7 <0.9;


contract DataStorage {
   
    // 1.变量的位置，
    function fCallData(uint[] calldata _x)public pure returns (uint[] calldata) {
        // _x[0]=0; // 当你修改的时候，就会报错TypeError: Calldata arrays are read-only.
        return (_x);
    }
    
    uint[] x =[1,2,3];
    // 为了更方便观察。写了一个getX方法来查询x的值，当调用了sStorage后x的第一个元素就变成了101了；
    function getX() public view returns (uint[] memory){
        return (x);
    }
    function sStorage() public {
        uint[] storage xStorage = x; // xStorage是uint[]类型，位置是storage，引用了x，当修改xStorage的时候，x也会被修改，
        xStorage[0]=101;//当我们调用了Sstorage方法的时候就xStorage以及x的第1个元素改成100，
        // 在debug的时候可以看到x的第1个元素一开是1 ，后面变成了100
        
    }

    uint[] y  =[1,2,3]; // 在函数之外定义的变量，默认就是stroage类型的，


    function yMemory() public view {
        uint[] memory ymemory = y; // y是storage类的，ymemory是引用了y，所以修改ymemory的时候不会修改到y
        ymemory[0]=101;
        
    }
    
    // 我们调用getY的时候就会发现，其实yMemory方法并不能修改y的值，这就是引用
    function getY() public view returns (uint[] memory){
        return (y);
    }
    // 2.变量的作用域
    // 2.1 状态变量： 链上的变量：写在合约里面函数外面的变量，就是状态变量，在所有合约内函数都可以访问，同样gas的消耗也高
    uint public a=1;
    uint public b;
    string public z;// 这些变量都是状态变量，同时，这些变量会自带有get方法去查看他们的值，（数组没有自带get）
    // 这些变量，可以在函数体内改变他们的值
    function foo()external  {
        a = 5;
        b=6;
        z="0xAA";
    }
    // 2.2 局部变量，这个就很好理解了，就是在函数内部定义的变量，只作用的函数内，出去了就没用啦
    function bar() external pure returns(uint){
        uint xx = 1; // 只在函数内部有用呢
        uint yy = 3;
        uint zz = xx + yy;
        return(zz);
    }
    // 2.3 全局变量，这里的全局变量和py不同，这里指的是solidity预留关键字的变量，他们可以在函数内不声明直接使用
    function global()external view returns(address,uint,bytes32,bytes memory)  {
        address sender=msg.sender; // 发送者的地址，猜测应该是调用这个合约的人的地址,address是一种类型
        uint blockNum = block.number; // 当前区块高度？
        bytes32 bHash = blockhash(blockNum - 1); // 获取区块哈希
        bytes memory data=msg.data; //这个没用过 ,byte是类型，memory是位置
        return (sender,blockNum,bHash,data);
    }
    function commonGlobal()external  payable   returns (uint,uint,uint256,bytes calldata,uint){
        uint _gaslimit=block.gaslimit;//当前区块的gaslimit
        uint _timeStamo=block.timestamp;//当前区块时间戳，为unix纪元以来的秒
        uint256 _gasleft=gasleft(); // 剩余gas
        bytes calldata _data=msg.data;// 完整的calldata
        uint _value=msg.value; //当前交易发送的wei值
        return (_gaslimit,_timeStamo,_gasleft,_data,_value);

    }
    // 以太坊没有小数点，用0代替为小数点来确保交易的精度，为了防止精度的损失，利用一台单位可以避免误算，
    // wei : 1
    // gwei : 1e9 =1000000000
    // ether : 1e18 = 100000000000000000
    function weiUit()external pure returns(uint)  {
        assert(1 wei ==1e0);
        assert(1 wei == 1);
        // return 1 wei; // 调用的时候就看到返回了1，
        return 1 gwei; // 试试看返回什么 返回 
    }
    function gweiUit()external pure returns(uint)  {
        assert(1 gwei ==1e9);
        assert(1 gwei == 1000000000); // 如果assert错误，调用的时候会报错的嘿嘿
        
        return 1 gwei; // 试试看返回什么 返回  1000000000
    }
    function etherUit()external pure returns(uint)  {
        assert(1 ether ==1e18);
        assert(1 ether == 1000000000000000000);
        // return 1 wei; // 调用的时候就看到返回了1，
        return 1 ether; // 试试看返回什么 返回  1000000000000000000
    }
    // 时间单位 可以在合约中设定一个时间必须在一个周期内完成，时间有周天时分秒

    function secondsUnit()external pure returns (uint) {
        assert(1 seconds ==1); // 也就是说 1 秒就是=1 uint
        return 1 seconds;
    }
    function minutesUnit()external pure returns (uint) {
        assert(1 minutes ==60);
        return 1 minutes;
    }
    function hoursUnit()external pure returns (uint) {
        assert(1 hours ==3600);
        return 1 hours;
    }
    function daysUnit()external pure returns (uint) {
        assert(1 days ==3600*24);
        return 1 days;
    }
    function weeksUnit()external pure returns (uint) {
        assert(1 weeks ==60*60*24*7);
        return 1 weeks;
    }


}