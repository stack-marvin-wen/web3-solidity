// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract BoolContract{
    bool isActive=true;
    bool isPause=false;
    function switchActive() public {
        isActive=!isActive;
    }
    function checkStatus() view public returns (bool) {
        return !isPause && isActive;
    }
    function ifStatementCheck() view public returns( string memory){
        if(isActive && !isPause){
            return "Active";
        }else{
            return "Not Active";
        }
    }
}