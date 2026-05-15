// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

contract ShopPayment{
    address immutable public owner;
    uint public constant ITEM_PRICE = 0.1 ether;
    mapping(address => uint) public purchases;
    constructor() {
        owner = msg.sender;
    }
    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }
    // 购买商品
    function buyItem(uint quantity) public payable{
        require(msg.value == quantity * ITEM_PRICE, "Incorrect amount");
        purchases[msg.sender] += quantity;
    }
    // 查询购买数量
    function getPurchases(address buyer) public view returns (uint) {
        return purchases[buyer];
    }
    // 提现（仅owner）
    function withdraw(uint amount) public payable  onlyOwner{
        require(amount<=getContractBalance(),"No balance to withdraw");
        (bool success,)=payable(owner).call{value:amount}("");
        require(success,"Transfer fail");
    }
    // 查询合约余额
    function getContractBalance() public view returns (uint) {
        return address(this).balance;
    }
}