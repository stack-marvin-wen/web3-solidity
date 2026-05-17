// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.28;
import {Test} from "forge-std/Test.sol";
import {CrowdfundingFactory} from "../src/CrowdfundingFactory.sol";

contract CrowdfundingFactoryTest is Test {
    CrowdfundingFactory factory;
    event CampaignCreated(
        address indexed user, address indexed campaignAddress, uint256 goal, uint256 duration, uint256 minContribution
    );

    function setUp() public {
        factory = new CrowdfundingFactory();
    }

    /**
     * @dev 测试创建众筹
     */
    function test_createCampaign() public {
        vm.warp(1_700_000_000);

        vm.expectEmit(true, false, false, true);
        emit CampaignCreated(address(this), address(0), 1 ether, block.timestamp + 1 days, 0.01 ether);
        address addr = factory.createCampaign(1 ether, 1 days, 0.01 ether);
        assertEq(factory.getCampaignCount(), 1);
        assertTrue(addr != address(0));
    }
}
