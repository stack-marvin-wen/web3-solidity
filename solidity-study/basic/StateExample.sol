// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract StateExample{
    enum State{
        Preparing,   // 准备中
        Active,      // 进行中
        Checking,    // 检查中
        Success,     // 成功
        Failed,      // 失败
        Cancelled    // 已取消
    }

    State public currentState;
    // 状态检查modifier
    modifier inState(State expected) {
        require(currentState == expected, "Invalid state for this operation");
        _;
    }
    
    constructor() {
        currentState = State.Preparing;
    }
    // 只能在Preparing状态执行
    function start() public inState(State.Preparing) {
        currentState = State.Active;
    }
    
    // 只能在Active状态执行
    function contribute() public payable inState(State.Active) {
        // 贡献资金
    }
    
    // 只能在Active状态执行
    function check() public inState(State.Active) {
        currentState = State.Checking;
    }
    
    // 状态转换
    function finalize() public inState(State.Checking) {
        if (address(this).balance >= 100 ether) {
            currentState = State.Success;
        } else {
            currentState = State.Failed;
        }
    }
    
    // 紧急取消
    function cancel() public {
        require(
            currentState == State.Preparing || currentState == State.Active,
            "Cannot cancel at this stage"
        );
        currentState = State.Cancelled;
    }
}