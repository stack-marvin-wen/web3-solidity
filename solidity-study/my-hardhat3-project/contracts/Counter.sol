// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

/**
 * @title Counter
 * @dev A simple counter contract that allows incrementing a counter by a specified value.
 * The contract emits an event whenever the counter is incremented.
 */
contract Counter {
  /**
    * @dev 公开的状态变量x 表示当前计数器的值.
    * @notice 状态变量x 是公开可见的，其他合约可以读取它的值。
   */
  uint public x;
  
  event Increment(uint by);

  function inc() public {
    x++;
    emit Increment(1);
  }

  function incBy(uint by) public {
    require(by > 0, "incBy: increment should be positive");
    x += by;
    emit Increment(by);
  }
}
