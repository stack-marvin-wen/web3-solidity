// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract DynamicBytes{
    bytes public data;
    function pushByte() public {
        data.push(0x12);
    }
    function getBytesLength() public view returns(uint){
        return data.length;
    }
    function getByte(uint index) public view returns(bytes1){
        require(index<data.length,"index out of bounds");
        return data[index];
    }
    function popByte() public {
        data.pop();
    }
}