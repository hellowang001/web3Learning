// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

/**
 * @title 计数器合约
 * @dev 一个简单的合约，维护一个可递增的计数器值。
 */
contract Counter {
    // 公有状态变量'counter'允许外部合约和用户查看其值。
    uint public counter;

    /**
     * @dev 构造函数，将计数器初始化为0。
     */
    constructor() {
        counter = 0;
    }

    /**
     * @dev 使计数器值增加1。
     * 注意：此函数没有返回值。
     */
    function count() public {
        counter = counter + 1;
    }
}