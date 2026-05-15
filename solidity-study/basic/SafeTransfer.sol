// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;
contract SafeTransfer{
    mapping(address => uint) public balances;
    uint public constant MAX_BATCH_SIZE = 50;
    function deposit() public payable{
        balances[msg.sender]+=msg.value;
    }
    function batchTransfer(
        address[] memory recipients,
        uint[] memory amounts
    ) public {
        require(recipients.length==amounts.length,"Length mismatch");
        require(recipients.length<=MAX_BATCH_SIZE,"Batch too large");
        uint totalAmount=0;
        for(uint i=0;i<recipients.length;i++){
            totalAmount+=amounts[i];
        }
        require(balances[msg.sender]>=totalAmount,"Insufficient balance");
        for(uint i=0;i<recipients.length;i++){
            balances[msg.sender]-=amounts[i];
            balances[recipients[i]]+=amounts[i];
        }
    }
    function getBalance(address user) public view returns (uint) {
        return balances[user];
    }
}