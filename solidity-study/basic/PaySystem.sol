// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
/*
创建一个完整的支付合约：

支持存款（deposit）
支持提款（withdraw）
支持紧急停止（pause）
Owner可以暂停/恢复合约
查询余额
限制最小存款金额
*/
contract PaySystem{
    address public Owner;
    bool public is_paused;
    mapping(address=>uint) public balances;
    uint public minDepoAmount;
    constructor(){
        Owner=msg.sender;
        is_paused=false;
        minDepoAmount=1 ether;
    }
    modifier isOwner(){
        require(Owner==msg.sender,"Not the Owner");
        _;
    }
    modifier isPaused(){
        require(is_paused,"Contract was paused by owner");
        _;
    }
    modifier isNotPaused(){
        require(!is_paused,"Contract is running");
        _;
    }
    modifier isMinDepoAmount(){
        require(msg.value>=minDepoAmount,"The amount is less than minDepoAmount");
        _;
    }
    modifier hasAmount(uint value){
        require(balances[msg.sender]>=value,"Not enough money");
        _;
    }
    function pauseContract() public isOwner isNotPaused{
        is_paused=true;
    }
    function resumeContract() public isOwner isPaused{
        is_paused=false;
    }
    function setMinDepoAmount(uint value) public isOwner{
        minDepoAmount=value;
    }

    function queryAmount(address addr) public view returns(uint value){
        value=balances[addr];
    }

    function deposit() public payable isNotPaused isMinDepoAmount{
        balances[msg.sender]+=msg.value;
    }
    function withdraw(uint value) public payable isNotPaused hasAmount(value) {
        balances[msg.sender]-=value;
        (bool success, ) = payable(msg.sender).call{value: value}("");
        require(success, "Transfer failed");
    }
}