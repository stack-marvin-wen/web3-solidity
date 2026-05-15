// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
/*
实现重入锁保护
应用Gas限制，防止恶意调用
支持多个所有者
需要达到指定确认数才能执行
*/

contract MultiSigWallet {

    address[] public Owners;
    mapping(address => bool) public isOwner;
    struct Transactionitem{
        address to;
        uint256 value;
        bytes data;
        bool executed;
        uint256 confirmations;
    }
    Transactionitem[] public transactions;
    mapping(uint256 => mapping(address => bool)) public isConfirmed;
    uint public Requirement;
    bool private locked;
    modifier onlyOwner(){
        require(isOwner[msg.sender],"not owner");
        _;
    }
    modifier nonLocked{
        require(locked,"is locked");
        locked=true;
        _;
        locked=false;
    }

    constructor(address[] memory _owners, uint256 _required) {
        require(_owners.length > 0, "Owners required");
        require(_required > 0 && _required <= _owners.length, "Invalid required");
        
        for (uint256 i = 0; i < _owners.length; i++) {
            isOwner[_owners[i]] = true;
            Owners.push(_owners[i]);
        }
        Requirement = _required;
    }

    function submit(address _to, uint256 _value, bytes memory _data) 
        external 
        onlyOwner 
        returns (uint256) {
            uint256 txId = transactions.length;
            transactions.push(Transactionitem({
                to: _to,
                value: _value,
                data: _data,
                executed: false,
                confirmations: 0
            }));
            return txId;
        }
    function confirm(uint256 _txId) external onlyOwner {
        require(!isConfirmed[_txId][msg.sender], "Already confirmed");
        isConfirmed[_txId][msg.sender] = true;
        transactions[_txId].confirmations += 1;
    }

    function execute(uint256 _txId) external onlyOwner nonLocked {
        Transactionitem storage tx = transactions[_txId];
        require(!tx.executed , "Already executed");
        require(tx.confirmations >= Requirement, "Insufficient confirmations");
        
        tx.executed = true;
        
        // 使用Gas限制
        (bool success, ) = tx.to.call{gas: 50000, value: tx.value}(tx.data);
        require(success, "Execution failed");
    }
}