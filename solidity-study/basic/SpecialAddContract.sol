// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract SpecialAddContract{
    function getSpecialAdd() view public returns(address, address,address){
        return (
            msg.sender, //当前调用者的地址
            tx.origin,// 交易发起者的地址
            address(this) //当前合约的地址
        );
    }

    // 用户 -> 合约A -> 合约B
    // 在合约B中：
    // msg.sender = 合约A的地址
    // tx.origin = 用户的地址
}