// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;
contract FundSystem{
    enum State{
        Fundraising, 
        Success, 
        Failed, 
        PaidOut
    }
    State public state = State.Fundraising;
    uint public goal=100 ether;
    modifier inState(State s){
        require(s==state,"Mismatch state");
        _;
    }
    function contribute() public payable inState(State.Fundraising){
        require(msg.value>0,"You must send some ether");
        goal-=msg.value;
    }
    function checkGoalReached() public inState(State.Fundraising){
        if(address(this).balance >=goal){
            state=State.Success;
        }else{
            state=State.Failed;
        }
    }
}