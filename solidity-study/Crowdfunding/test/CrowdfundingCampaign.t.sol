// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import {Test} from "forge-std/Test.sol";
import {CrowdfundingCampaign} from "../src/CrowdfundingCampaign.sol";

contract CrowdfundingCampaignTest is Test {
    CrowdfundingCampaign public crowdfundingCampaign;
    event Created(address indexed owner, uint256 goal, uint256 deadline);
    event Start();

    receive() external payable {}
    fallback() external payable {}

    /**
     * @dev 测试初始化
     */
    function setUp() public {
        crowdfundingCampaign = new CrowdfundingCampaign(1 ether, 10 days, 0.01 ether);
        crowdfundingCampaign.start();
    }

    /**
     * @dev 测试构造函数
     */
    function test_Constructor()  public view {
        assertEq(crowdfundingCampaign.goal(), 1 ether);
        assertEq(crowdfundingCampaign.MIN_CONTRIBUTION(), 0.01 ether);
        assertEq(uint256(crowdfundingCampaign.state()), uint256(CrowdfundingCampaign.State.Funding));
        assertEq(crowdfundingCampaign.owner(), address(this));
        assertEq(crowdfundingCampaign.total(), 0);
    }

    /**
     * @dev 测试过期时间是否正确
     */
    function test_Deadline() public {
        vm.warp(1_700_000_000);
        CrowdfundingCampaign cc = new CrowdfundingCampaign(1 ether, 10 days, 0.01 ether);
        assertEq(cc.deadline(), 1_700_000_000 + 10 days);
    }

    /**
     * @dev 测试启动众筹
     */
    function test_Start() public {
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Preparing))
        );
        vm.expectEmit(false, false, false, false);
        emit Start();
        crowdfundingCampaign.start();
        assertEq(uint256(crowdfundingCampaign.state()), uint256(CrowdfundingCampaign.State.Funding));
    }

    /**
     * @dev 测试贡献 状态不是Funding时 抛出异常
     */
    function test_ContributionWithNoRightState() public {
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Preparing))
        );
        vm.expectRevert(bytes("Not in current state"));
        crowdfundingCampaign.contribution{value: 0.5 ether}();
    }

    /**
     * @dev 测试贡献 小于最小贡献金额时 抛出异常
     */
    function test_ContributionWithLessThanMinContribution() public {
        assertEq(uint256(crowdfundingCampaign.state()), uint256(CrowdfundingCampaign.State.Funding));
        deal(address(1), 1 ether);
        vm.expectRevert(bytes("Less than min contribute"));
        vm.prank(address(1));
        crowdfundingCampaign.contribution{value: 0.0095 ether}();
    }

    /**
     * @dev 测试贡献 当大于截止日期的时候 抛出异常
     */
    function test_ContributionWithDeadline() public {
        vm.warp(1_700_000_000);
        deal(address(1), 1 ether);
        vm.expectRevert(bytes("Due date"));
        vm.prank(address(1));
        crowdfundingCampaign.contribution{value: 0.5 ether}();
    }

    /**
     * @dev 测试贡献 正常情况
     */
    function test_Contribution() public {
        deal(address(1), 1 ether);
        vm.prank(address(1));
        crowdfundingCampaign.contribution{value: 0.5 ether}();
        assertEq(crowdfundingCampaign.total(), 0.5 ether);
        assertEq(crowdfundingCampaign.contributions(address(1)), 0.5 ether);
    }

    /**
     * @dev 测试提现 状态不是Success时 抛出异常
     */
    function test_WithdrawWithNoRightState() public {
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Funding))
        );
        vm.expectRevert(bytes("Not in current state"));
        crowdfundingCampaign.withdraw();
    }

    /**
     * @dev 测试提现 不是众筹发起人时 抛出异常
     */
    function test_WithdrawNotOwner() public {
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Success))
        );
        vm.expectRevert(bytes("Not owner"));
        vm.prank(address(1));
        crowdfundingCampaign.withdraw();
    }

    /**
     * @dev 测试提现 状态不是Success时 抛出异常
     */
    function test_WithdrawNotInSuccess() public {
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Funding))
        );
        vm.expectRevert(bytes("Not in current state"));
        crowdfundingCampaign.withdraw();
    }

    /**
     * @dev 测试提现 正常情况
     */
    function test_Withdraw() public {
        // 先贡献
        deal(address(1), 1 ether);
        vm.prank(address(1));
        crowdfundingCampaign.contribution{value: 1 ether}();
        // 直接修改状态为成功，然后提现
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Success))
        );

        uint256 balanceBefore = address(this).balance;
        crowdfundingCampaign.withdraw();
        uint256 balanceAfter = address(this).balance;
        assertEq(balanceAfter - balanceBefore, 1 ether);
        assertEq(crowdfundingCampaign.total(), 0);
    }

    /**
     * @dev 测试退款 状态不是Failed时 抛出异常
     */
    function test_RefoundWithNoRightState() public {
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Funding))
        );
        vm.expectRevert(bytes("Not in current state"));
        crowdfundingCampaign.refund();
    }

    /**
     * @dev 测试退款 没有贡献过时 抛出异常
     */
    function test_RefoundWithNoContribution() public {
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Failed))
        );
        vm.expectRevert(bytes("No contributions"));
        vm.prank(address(1));
        crowdfundingCampaign.refund();
    }

    /**
     * @dev 测试退款 正常情况
     */
    function test_Refound() public {
        deal(address(1), 1 ether);
        vm.prank(address(1));
        crowdfundingCampaign.contribution{value: 0.5 ether}();
        assertEq(crowdfundingCampaign.getContribution(address(1)), 0.5 ether);
        // 直接修改状态为失败，然后退款
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Failed))
        );
        vm.prank(address(1));
        uint256 balanceBefore = address(1).balance;
        crowdfundingCampaign.refund();
        uint256 balanceAfter = address(1).balance;
        assertEq(balanceAfter - balanceBefore, 0.5 ether);
        assertEq(crowdfundingCampaign.getContribution(address(1)), 0);
    }

    function test_Status() public {
        assertEq(uint256(crowdfundingCampaign.state()), uint256(CrowdfundingCampaign.State.Funding));
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Preparing))
        );
        assertEq(crowdfundingCampaign.status(), "Preparing");
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Success))
        );
        assertEq(crowdfundingCampaign.status(), "Success");
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Failed))
        );
        assertEq(crowdfundingCampaign.status(), "Failed");
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Closed))
        );
        assertEq(crowdfundingCampaign.status(), "Closed");
    }

    /**
     * @dev 测试进度
     */
    function test_Progress() public {
        assertEq(crowdfundingCampaign.progress(), 0);
        deal(address(1), 1 ether);
        vm.store(
            address(crowdfundingCampaign), bytes32(uint256(6)), bytes32(uint256(CrowdfundingCampaign.State.Funding))
        );
        vm.prank(address(1));
        crowdfundingCampaign.contribution{value: 0.5 ether}();
        assertEq(crowdfundingCampaign.total(), 0.5 ether);
        assertEq(crowdfundingCampaign.progress(), 50);
    }
}
