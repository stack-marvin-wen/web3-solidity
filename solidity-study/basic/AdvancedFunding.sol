// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;
contract AdvancedCrowdfunding {
    enum State { 
        Fundraising, 
        Successful, 
        Failed, 
        PaidOut 
    }
    State public state;
    address immutable public OWNER;
    uint immutable public GOAL;
    uint immutable public MINIMUM_CONTRIBUTION;
    uint immutable public DURATIONDATE;
    mapping (address=>uint) public contributions;
    uint public total;
    address[] public conAddr;
    modifier inState(State s){
        require(s==state,"Not in right state");
        _;
    }
    modifier Owner(){
        require(msg.sender==OWNER,"not the owner, and not allow this operation");
        _;
    }
    constructor(uint goalAmount, uint durationDays){
        OWNER=msg.sender;
        MINIMUM_CONTRIBUTION=0.1 ether;
        GOAL=goalAmount;
        DURATIONDATE=block.timestamp+durationDays*24*60*60;
    }
    
    // 贡献资金
    function contributor() public payable inState(State.Fundraising){
        uint amount=msg.value;
        require(amount>=MINIMUM_CONTRIBUTION,"Amount is not enough");
        require(block.timestamp<=DURATIONDATE,"Out of date");
        address addr=msg.sender;
        if(contributions[addr]>0){
            conAddr.push(addr);
        }
        total+=amount;
        contributions[addr]+=amount;
        
        checkState();
    }
    // 检查并更新状态
    function checkState() public {
        if(address(this).balance>=GOAL){
                state=State.Successful;
            }
        if(block.timestamp>DURATIONDATE){
            if(address(this).balance<GOAL){
                state=State.Failed;
            }
        }
    }
    // 创建者提取资金
    function withdraw()public payable  Owner inState(State.Successful) {
        (bool success,)=OWNER.call{value:GOAL}("");
        require(success,"Withdraw fail");
    }
    // 退款
    function refund() public payable {
        uint len=conAddr.length;
        for(uint i=0;i<len;i++){
            address addr=conAddr[i];
            total-=contributions[addr];
            (bool success,)=addr.call{value:contributions[addr]}("");
            require(success,"Refund fail");
        }
        state=State.Failed;
    }
    // 查询函数
    function getInfo() public view returns (
        State _state,
        uint goal,
        uint funded,
        uint deadline,
        uint timeRemaining,
        uint contributors
    ) {
        uint remaining = 0;
        if (block.timestamp < DURATIONDATE) {
            remaining = DURATIONDATE - block.timestamp;
        }
        
        return (
            state,
            GOAL,
            total,
            DURATIONDATE,
            remaining,
            conAddr.length
        );
    }

    // getProgress
    function getProgress() view public returns(uint){
        return (total * 100) / GOAL;
    }
    // isActive
    function isActive() view public returns(bool){
        return state == State.Fundraising &&  block.timestamp <= DURATIONDATE;
    }
}