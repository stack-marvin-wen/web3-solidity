// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract AddFeature{
    // 查询地址余额
    function getBalance(address add) view public returns(uint){
        return add.balance;
    }
    // 获取当前合约地址
    function getContractAdd() view public returns (address){
        return address(this);
    }
    // 获取合约余额
    function getContractBalance() view public returns (uint) {
        return address(this).balance;
    }
    // 检查是否为零地址
    function isZeroAdd(address add) pure public returns (bool){
        return add == address(0);
    }
}