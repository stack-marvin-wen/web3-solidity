// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract StringDemo{
    string hello;
    // 设置字符串
    function setString(string calldata a) public{
        hello=a;
    }
    // 获取字符串
    function getString() public view returns(string memory){
        return hello;
    }
    // 错误：不能直接比较
    // function errorCompare(string a,string b) public pure returns (bool){
    //     return a==b;
    // }
    // 错误：不能直接获取长度
    // function getLenght(string a) public returns(uint){
    //     return a.length;
    // }
    // 错误：不能直接拼接（0.8.12之前）
    // function concatString(string memory a,string memory b) public view  returns(string memory){
    //     return a+b;
    // }
    // 字符串比较
    function compareString(string memory a,string memory b) public pure returns ( bool){
        return keccak256(bytes(a))==keccak256(bytes(b));
    }
    // 字符串拼接
    function concat(string memory a,string memory b) public pure returns(string memory){
        return string.concat(a,b);
    }
    // 字符串与bytes转换
    function stringTobytes(string memory a) public pure returns(bytes memory){
        return bytes(a);
    }
    function bytesTostring(bytes memory a) public pure returns(string memory){
        return string(a);
    }
    function getStringLength(string memory a) public pure returns(uint){
        return bytes(a).length;
    }
}