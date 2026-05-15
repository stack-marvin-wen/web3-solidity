// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;
/*
支持创建多个提案
每个提案有截止时间
只有owner可以创建提案
每个地址只能投一次票
可以查询投票结果
可以获取获胜提案
*/

contract VoteSystem{
    struct Proposal {
        string description;
        uint voteCount;
        uint deadline;
        bool exists;
    }
    
    address public owner;
    uint public proposalCount;
    
    mapping(uint => Proposal) public proposals;
    mapping(uint => mapping(address => bool)) public hasVoted;
    
    constructor() {
        owner = msg.sender;
    }
    modifier Owner(){
        require(msg.sender==owner,"Not the Owner");
        _;
    }
    modifier ProposalExist(uint proposalId){
        require(proposals[proposalId].exists,"Proposal not exist");
        _;
    }
    modifier NotVote(uint proposalId,address addr){
        require(!hasVoted[proposalId][addr],"User has voted");
        _;
    }
    modifier NotDue(uint proposalId){
        require(block.timestamp<=proposals[proposalId].deadline,"Proposal has due");
        _;
    }
    // TODO: 实现创建提案
    function createProposal(string memory description, uint durationDays) 
        public Owner
    {
        require(keccak256(bytes(description))!=keccak256(bytes("")),"Description not allow empty");
        require(durationDays>0,"Duration Days should more than 0");

        proposals[proposalCount]=Proposal({
            description:description,
            voteCount: 0,  // Add this line
            exists:true,
            deadline:block.timestamp+durationDays*24*60*60
        });
        proposalCount++;
    }
    
    // TODO: 实现投票
    function vote(uint proposalId) public ProposalExist(proposalId) NotVote(proposalId,msg.sender) NotDue(proposalId) {
        // 检查提案存在
        // 检查是否已投票
        // 检查是否已截止
        // 执行投票
        proposals[proposalId].voteCount++;
        hasVoted[proposalId][msg.sender]=true;
    }
    
    // TODO: 获取获胜提案
    function getWinner() public view returns (uint) {
        // 遍历所有提案
        // 找出票数最多的
        uint maxVoteCnt=0;
        uint winnerId=0;
        for(uint i=0;i<proposalCount;i++){
            if(proposals[i].voteCount>maxVoteCnt){
                maxVoteCnt=proposals[i].voteCount;
                winnerId=i;
            }
        }
        return winnerId;
    }
}