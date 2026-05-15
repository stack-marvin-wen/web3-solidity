// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract SafeContract{
    /*
    限制最大长度为100
    实现安全的添加功能（safePush）
    实现两种删除方法（保序和快速）
    实现分批求和功能（sumRange）
    实现查找功能（返回元素索引）
    实现获取所有元素功能
    */
    uint public constant MAX_ARRAY_LEN=100;
    uint[] public array;
    constructor(){
        for(uint i=0;i<20;i++){
            array.push(i);
        }
    }
    function safePush(uint a) public {
        uint len=array.length;
        require(len<MAX_ARRAY_LEN,"More than max len");
        array.push(a);
    }
    function deleteNoOrder(uint _index) public {
        require(_index<array.length,"Out of range");
        array[_index]=array[array.length-1];
        array.pop();
    }
    function deleteOrder(uint _index) public {
        uint len=array.length;
        require(_index<len,"Out of range");
        for(uint i=_index;i+1<len;i++){
            array[i]=array[i+1];
        }
        array.pop();
    }
    function sumByRange(uint begin,uint end) public view  returns(uint) {
        uint len=array.length;
        require(begin<end,"begin would be less than end");
        require(end<len,"out of range");
        uint sum=0;
        for(uint i=begin;i<=end;i++){
            sum+=array[i];
        }
        return sum;
    }
    function sumArray(uint batchSize) public view returns(uint) {
        uint sum=0;
        uint len=array.length;
        if(batchSize>len) batchSize=len;
        for(uint i=0;i*batchSize<len;i++){
            uint end=(i+1)*batchSize-1;
            if(end>=len) end=len-1;
            sum+=sumByRange(i*batchSize, (i+1)*batchSize-1);
        }
        return sum;
    }
    function searchByIndex(uint num) public view returns(uint){
        uint len=array.length;
        for(uint i=0;i<len;i++){
            if(array[i]==num) return i;
        }
        revert("Not found"); 
    }
    function getAllArray() public view returns(uint[] memory){
        return array;
    }
}