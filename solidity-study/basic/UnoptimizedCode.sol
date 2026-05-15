
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract UnoptimizedCode {
    uint[] public data;
    
    function process(uint[] memory values) public {
        for(uint i = 0; i < values.length; i++) {
            if(values[i] > 10) {
                data.push(values[i]);
            }
        }
    }
}

contract OptimizedCode {
    uint[] public data;
    
    function process(uint[] calldata values) public {
        uint len=values.length;
        uint[] memory temp=new uint[](len);
        uint cnt=0;
        for(uint i = 0; i < values.length; i++) {
            if(values[i] > 10) {
                temp[cnt++]=values[i];
            }
        }
        for(uint i=0;i<cnt;i++){
            data.push(temp[i]);
        }
    }
}