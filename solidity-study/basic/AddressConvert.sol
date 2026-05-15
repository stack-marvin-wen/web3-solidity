// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract AddressConvertContract{
    receive() external payable { }
    // address转为address payable
    function addTopayable(address add) pure public returns(address payable ){
        return payable(add);
    }
    // uint160转为address
    function uint160Toaddress(uint160 add) pure public returns(address ){
        return address(add);
    }
    // address转为uint160
    function addTouint160(address add) pure public returns (uint160){
        return uint160(add);
    }
    // 示例：使用转换
    function payDemo(address add) public {
        address payable accpetor=payable(add);
        (bool success,)=accpetor.call{value:1000}("");
        require(success,"pay fail");
    }
}