// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract FixbyteDemo{
    bytes1 public b1 = 0x12;
    bytes4 public b4 = 0x12345678;
    bytes32 public b32 = 0x1234567890123456789012345678901234567890123456789012345678901234;
    // 获取长度
    function getLength() public pure returns(uint,uint,uint){
        bytes1 a;
        bytes2 b;
        bytes32 c;
        return (a.length,b.length,c.length);
    }
}
contract FixbyteUseCase{
    // 1. 存储哈希值
    bytes32 public fileHash;
    function storeHash(string memory a) public {
        fileHash=keccak256(bytes(a));
    } 
    // 2. 存储签名
    bytes32 public r;
    bytes32 public s;
    uint8 public v;
    // 3. 紧凑数据存储
    bytes4 public functionSelector = 0x70a08231;  // balanceOf(address)的函数选择器
    // 4. 存储合约地址
    bytes20 public contractAddress;
}