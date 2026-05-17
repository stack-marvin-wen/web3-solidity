// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import {CrowdfundingCampaign} from "./CrowdfundingCampaign.sol";

contract CrowdfundingFactory {
    CrowdfundingCampaign[] public campaigns;
    mapping(address => uint256[]) public userIndex;
    event CampaignCreated(
        address indexed user, address indexed campaignAddress, uint256 goal, uint256 duration, uint256 minContribution
    );

    /**
     * @dev 创建众筹
     * @param _goal 目标金额
     * @param _duration 众筹持续时间
     * @param _minContribution 最小贡献金额
     */
    function createCampaign(uint256 _goal, uint256 _duration, uint256 _minContribution) public returns (address) {
        CrowdfundingCampaign campaign = new CrowdfundingCampaign(_goal, _duration, _minContribution);
        campaigns.push(campaign);
        userIndex[msg.sender].push(campaigns.length - 1);
        emit CampaignCreated(msg.sender, address(campaign), _goal, block.timestamp + _duration, _minContribution);
        return address(campaign);
    }

    /**
     * @dev 获取用户众筹列表
     */
    function getUserCrowdfunding() public view returns (address[] memory) {
        uint256[] memory userIndexs = userIndex[msg.sender];
        uint256 length = userIndexs.length;
        if (length == 0) {
            return new address[](0);
        }
        address[] memory userCampaigns = new address[](length);
        for (uint256 i = 0; i < length; i++) {
            userCampaigns[i] = address(campaigns[userIndexs[i]]);
        }
        return userCampaigns;
    }

    /**
     * @dev 获取所有众筹列表
     */
    function getAllCrowdfunding() public view returns (address[] memory) {
        // 创建与活动数组长度相同的地址数组
        address[] memory campaignAddresses = new address[](campaigns.length);
        // 遍历所有活动，将地址存入数组
        for (uint256 i = 0; i < campaigns.length; i++) {
            campaignAddresses[i] = address(campaigns[i]);
        }
        // 返回地址数组
        return campaignAddresses;
    }

    /**
     * @dev 获取众筹数量
     */
    function getCampaignCount() public view returns (uint256) {
        return campaigns.length;
    }
}
