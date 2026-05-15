// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract HelloWorld{
     string public message; // 状态变量，存储在区块链上（storage中） 永久存储在区块链上
     constructor(){ // 构造函数在合约部署时自动执行，且只执行一次。 部署时自动调用，只执行一次，之后无法再调用 用于初始化合约状态
        message="Hello World";
     }
     function updateMSG(string memory newMsg) public  { 
        message=newMsg;
     }
     function getMsg() public view returns (string memory){
        return message;
     }

}