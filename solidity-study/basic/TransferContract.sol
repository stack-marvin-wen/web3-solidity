// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract TransferContract{
    // 接收ETH的函数需要payable修饰符
    receive() external payable { }
    // transfer方法（推荐，失败会回退）
    function transferDemo(address payable _acceptor, uint amount) public {
        (bool success, ) = _acceptor.call{value: amount}("");
        require(success, "Transfer failed");
    }
    // send方法（不推荐，需要检查返回值）
    function sendDemo(address payable _acceptor, uint amount) public returns(bool) {
    (bool suss, ) = _acceptor.call{value: amount}("");
    require(suss, "call fail");
    return suss;
    }
    // call方法（最灵活，推荐用于转账）
    function callDemo(address payable _acceptor,uint amount) public {
        (bool suss,)=_acceptor.call{value:amount}("");
        require(suss,"call fail");
    }
}