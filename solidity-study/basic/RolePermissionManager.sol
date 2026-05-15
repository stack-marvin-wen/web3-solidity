// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
/*
定义三种角色：Owner、Admin、User
实现角色分配和检查
不同角色有不同权限
Owner可以添加Admin
Admin可以添加User
所有人可以查询角色
*/
contract RolePermission{
    address public Owner;
    mapping (address=>string) AddrRole;
    constructor(){
        Owner=msg.sender;
    }
    modifier isOwner(){
        require(msg.sender==Owner,"Not Owner");
        _;
    }
    modifier isAdmin(){
        require(keccak256(bytes(AddrRole[msg.sender]))==keccak256(bytes("Admin")),"Not Admin");
        _;
    }
    modifier isUser(){
        require(keccak256(bytes(AddrRole[msg.sender]))==keccak256(bytes("User")),"Not User");
        _;
    }

    function addAdmin(address _addr) public isOwner{
        AddrRole[_addr]="Admin";
    }
    function addUser(address _addr) public isAdmin{
        AddrRole[_addr]="User";
    }
    function getRole(address _addr) public view isUser returns(string memory){
        return AddrRole[_addr];
    }
}