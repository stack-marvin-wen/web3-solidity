// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract VoteContract{
    enum VoteStatus{
        Yes,
        No,
        Abstain
    }
    mapping(address=>VoteStatus) userVote;
    mapping(address => bool) private hasVoted;
    uint public YesCnt=0;
    uint public NoCnt=0;
    uint public AbstainCnt=0;
    function userVoteAction(uint voteOption) public{
        require(voteOption <= uint(VoteStatus.Abstain), "Invalid vote option");
        address user=msg.sender;
        if (hasVoted[user]){
            revert("Already voted");
        }
        hasVoted[user]=true;
        userVote[user]=VoteStatus(voteOption);
        if (voteOption==uint(VoteStatus.Yes)){
            YesCnt++;
        }else if (voteOption==uint(VoteStatus.No)){
            NoCnt++;
        }else{
            AbstainCnt++;
        }
    }
    function queryVoteRes() view public returns(uint,uint,uint){
        return (YesCnt,NoCnt,AbstainCnt);
    }
}