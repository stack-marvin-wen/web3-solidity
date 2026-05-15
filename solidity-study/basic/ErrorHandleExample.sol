// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;
contract ErrorHandle{
    error Unauthorized(address caller);
    function example() public view {
        revert Unauthorized(msg.sender);
    }
}