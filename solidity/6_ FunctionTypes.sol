// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7 <0.9;

// function <function name>([parameter types[, ...]]) {internal|external|public|private} [pure|view|payable] [virtual|override] [<modifiers>]
// [returns (<return types>)]{ <function body> }
// 1、`<function name>`：函数名。
// 2、`{internal|external|public|private}`：函数可见性说明符，共有4种。
// 2.1 `public`：内部和外部均可见。
// 2.2 `private`：只能从本合约内部访问，继承的合约也不能使用。
// 2.3 `external`：只能从合约外部访问（但内部可以通过 `this.f()` 来调用，`f`是函数名）。
// 2.4 `internal`: 只能从合约内部访问，继承的合约可以用。
// 注意：合约中定义的函数需要明确指定可见性，它们没有默认值。
// 3、`[pure|view|payable]`：决定函数权限/功能的关键字。`payable`（可支付的）很好理解，带着它的函数，运行的时候可以给合约转入 ETH。`pure` 和 `view` 的介绍见下一节。
// 4、`[virtual|override]`: 方法是否可以被重写，或者是否是重写方法。`virtual`用在父合约上，标识的方法可以被子合约重写。`override`用在子合约上，表名方法重写了父合约的方法。
// 5、`<modifiers>`: 自定义的修饰器，可以有0个或多个修饰器。
// 6、`[returns ()]`：函数返回的变量类型和名称。
// 7、`<function body>`: 函数体。
contract FunctionTypes{
    uint256 public number = 1; // 先定义一个初始化的变量，待会定义方法去改变他
     // 默认function，这里用到了2的 external，只能从合约外部访问（但内部可以通过 `this.f()` 来调用，`f`是函数名）。
    function add() external{
        // 普通函数，使用了number变量，调用的时候会让number加1 这个方法不需要声明权限即可改变变量的值
        number=number+1;
    }
    // addPure 这里用到了2的external 表示外部访问，pure表示权限，returns表示有返回值
    // 3.1 pure: 就是一个纯纯牛马，不能读也不能写
    function addPure(uint256 _number) external pure returns(uint256 new_number ){

        // 这个Pure的add没有使用函数外的任何变量和方法，只在内部玩，试试调用外部的变量,
        // 报错提示我们Function declared as pure, but this expression (potentially) reads from the environment or state and thus requires "view". 
        // 方法是pure类型，是一个纯函数，但该表达式（可能）从环境或状态读取了其他变量或者方法，因此需要将方法调整为 “view”类型。
        // new_number=number+1; 

        // 有一个入餐，传入一个uint256 返回一个新的变量 是传入的+1 试试
        new_number=_number +1 ;
    }
    // 3.2 view: 是一只舔狗，只能看，只能拿过来引用，不可以改变其他变量或者方法，
    function addView() external view returns (uint256 new_number){
        // 试试改变内部变量
        // 这里报错提示： Function declared as view, but this expression (potentially) modifies the state and thus requires non-payable (the default) or payable.
        // 报错提示我们这是一个view类型方法，但是检测到你在改变其他变量，这样不行，需要变成payable 类型的方法才能去改变变量，
        // number=2; // 报错提示view类型方法不可改变变量，只有payble类型方法才可
        new_number=number +1 ; // 始终等于number+1
    }

    // 3.3 payable: 支付了，付钱了就是大哥啊，大哥当然可以为所欲为了，想怎么改变变量就怎么改变，嘿嘿
    function addPayable()external payable returns (uint256 new_number){
        number=100;
        new_number=number+1;
    }
    // 3.4 payable 递钱，能给合约支付eth的函数
    function minusPayable()external payable returns (uint256 balance){
        miuns();
        balance=address(this).balance; //address(this) 就是引用合约地址，balance就是合约的余额
    }
    // 2.3 internal： 内部只可以在合约内部调用,这个函数在部署后调试的时候直接不会显示出来
    function miuns()internal {
        number=number - 1;  
    }
    // 合约内的函数调用内部函数
    function minutesCall() external {
        miuns();
    }
}


