// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
/*
定义Proposal结构体（包含voters的mapping）
支持创建提案
支持投票（每人只能投一次）
查询提案信息
获取获胜提案
*/
contract VoteManContract{
    struct Proposal{
        // 地址投票题案
        string title;
        address proposer;
        mapping(address => bool) voters;
        uint voteCount;
        uint agreeCnt;
        uint disagreeCnt;
    }
    Proposal[] public proposals;
    mapping (string=>uint) public proposalTitleToIndex;

    uint maxIndex=0;
    function createProposal(string calldata title) public {
        Proposal storage p=proposals.push();
        p.title=title;
        p.proposer=msg.sender;
        proposalTitleToIndex[title]=proposals.length-1;
    }
    function vote(string calldata title,bool agree) public {
        uint index=proposalTitleToIndex[title];
        proposals[index].voters[msg.sender]=agree;
        proposals[index].voteCount++;
        if(agree) {
            proposals[index].agreeCnt++;
            if(proposals[index].agreeCnt>proposals[maxIndex].agreeCnt) maxIndex=index;
        }
        else proposals[index].disagreeCnt++;
    }
    function getSuccessProposal() public view returns(string memory,address,uint,uint,uint){

        return (
            proposals[maxIndex].title,
            proposals[maxIndex].proposer,
            proposals[maxIndex].voteCount,
            proposals[maxIndex].agreeCnt,
            proposals[maxIndex].disagreeCnt
        );
    }
    function getProposalByTitle(string calldata title) public view returns(string memory,address,uint,uint,uint){
        uint index=proposalTitleToIndex[title];
        return (
            proposals[index].title,
            proposals[index].proposer,
            proposals[index].voteCount,
            proposals[index].agreeCnt,
            proposals[index].disagreeCnt
        );
    }
}