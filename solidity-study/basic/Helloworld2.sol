// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract HelloWorld{
    address public  add;
    constructor(){
        add=msg.sender;
    } 
    function getOwner() public view returns(address){
        return add;
    }
    function isOwner() public view returns(bool){
        return add==msg.sender;
    }
}