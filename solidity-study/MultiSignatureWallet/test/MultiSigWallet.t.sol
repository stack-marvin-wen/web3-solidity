// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;
import "forge-std/Test.sol";
import {MultiSigWallet} from "../src/MultiSigWallet.sol";
contract MultiSigWalletTest is Test { 
    event AddUser(address indexed user);
    event RemoveUser(address indexed opUser,address indexed user);
    event Deposit(address indexed depositor, uint value); 
    event SubmitTransaction(address indexed creator, uint indexed id, address indexed to, uint value, bytes data);
    event DeleteTransaction(address indexed deletor, uint indexed id);
    event ConfirmTransaction(address indexed confirmer, uint indexed id);
    event RevokeConfirmation(address indexed confirmer, uint indexed id);
    event ExecuteTransaction(address indexed confirmer, uint indexed id);
    MultiSigWallet multiSigWallet;
    function setUp() public {
        address[] memory owners = new address[](4);
        owners[0] = address(0x1);
        owners[1] = address(0x2);
        owners[2] = address(0x3);
        owners[3] = address(this);
        multiSigWallet = new MultiSigWallet(owners, 2);
        deal(address(multiSigWallet), 100 ether);
    }
    /**
     * @dev 测试初始化
     */
    function test_init() public view {
        assertEq(multiSigWallet.getOwnerCnt(), 4);
        assertEq(multiSigWallet.getBalance(), 100 ether);
        assertEq(multiSigWallet.getRequired(), 2);
        assertEq(multiSigWallet.getOwners()[0], address(0x1));
        assertEq(multiSigWallet.getOwners()[1], address(0x2));
        assertEq(multiSigWallet.getOwners()[2], address(0x3));
        assertEq(multiSigWallet.getOwners()[3], address(this));
    }
    /**
     * @dev 测试转账
     */
    function test_transfer() public { 
        vm.expectEmit(true, false, false, true);
        emit Deposit(address(this), 1 ether);
        (bool success,)=address(multiSigWallet).call{value: 1 ether}("");
        assertEq(success,true);
        assertEq(multiSigWallet.getBalance(), 101 ether);
    }
    /**
     * @dev 测试添加用户 只有当前操作者是owner才能添加用户
     */
    function test_AddUserNotOpOwner() public {
        vm.prank(address(4));
        vm.expectRevert("Not an owner");
        multiSigWallet.addOwner(address(0x5));
    }
    /**
     * @dev 测试添加用户 已经存在的用户不能添加
     */
    function test_AddUserOwnerAddExist() public {
        vm.prank(address(3));
        vm.expectRevert("Owner is not unique");
        multiSigWallet.addOwner(address(0x1));
    }
    /**
     * @dev 测试添加用户
     */
    function test_AddUser() public {
        vm.prank(address(3));
        vm.expectEmit(true, false, false, false);
        emit AddUser(address(0x5));
        multiSigWallet.addOwner(address(0x5));
        assertEq(multiSigWallet.getOwnerCnt(), 5);
        assertEq(multiSigWallet.getOwners()[4], address(0x5));
    }
    /**
     * @dev 测试移除用户 只有当前操作者是owner才能移除用户
     */
    function test_RemoveUserNotOpOwner() public { 
        vm.prank(address(4));
        vm.expectRevert("Not an owner");
        multiSigWallet.removeOwner(address(0x1));
    }
    /**
     * @dev 测试移除用户 已经存在的用户才能移除
     */
    function test_RemoveUserRemoveNonExist() public { 
        vm.prank(address(3));
        vm.expectRevert("Not an owner");
        multiSigWallet.removeOwner(address(0x4));
    }
    /**
     * @dev 测试移除用户
     */
    function test_RemoveUser() public { 
        vm.prank(address(3));
        vm.expectEmit(true, true, false, false);
        emit RemoveUser(address(3),address(0x1));
        multiSigWallet.removeOwner(address(0x1));
        assertEq(multiSigWallet.getOwnerCnt(), 3);
        assertEq(multiSigWallet.getOwners()[0], address(this));
        assertEq(multiSigWallet.getOwners()[1], address(0x2));
    }
    /**
     * @dev 测试创建交易 只有owner才能创建交易
     */
    function test_CreateTransactionNonOwner() public{
        vm.prank(address(4));
        vm.expectRevert("Not an owner");
        multiSigWallet.submitTransaction(address(0x6), 1 ether, "");
    }
    /**
     * @dev 测试创建交易题案
     */
    function test_CreateTransaction() public{
        vm.prank(address(3));
        vm.expectEmit(true, true, true, false);
        emit SubmitTransaction(address(this), multiSigWallet.getTransactionCnt(), address(0x6), 1 ether, "");
        multiSigWallet.submitTransaction(address(0x6), 1 ether, "");
        assertEq(multiSigWallet.getTransactionCnt(), 1);
    }
    /**
     * @dev 删除交易 只有创建者能够删除
     */
    function test_DeleteTransactionNonCreatorOrExcuted() public{
        vm.prank(address(3));
        vm.expectEmit(true, true, true, false);
        emit SubmitTransaction(address(this), multiSigWallet.getTransactionCnt(), address(0x6), 1 ether, "");
        uint txIndex=multiSigWallet.submitTransaction(address(0x6), 1 ether, "");
        vm.prank(address(2));
        vm.expectRevert("Not the creator");
        multiSigWallet.deleteTransaction(txIndex);
    }
    /**
     * @dev 删除交易 
     */
    function test_DeleteTransaction() public{
        vm.expectEmit(true, true, true, false);
        emit SubmitTransaction(address(this), multiSigWallet.getTransactionCnt(), address(0x6), 1 ether, "");
        uint txIndex=multiSigWallet.submitTransaction(address(0x6), 1 ether, "");
        vm.expectEmit(true, true, false, false);
        emit DeleteTransaction(address(this), txIndex);
        multiSigWallet.deleteTransaction(txIndex);
        assertEq(multiSigWallet.getTransactionCnt(), 0);
    }
    /**
     * @dev 确认交易 只有owner才能确认交易
     */
    function test_confirmTransactionNonOwner() public{
        vm.prank(address(4));
        vm.expectRevert("Not an owner");
        multiSigWallet.confirmTransaction(0);
    }
    /**
     * @dev 确认交易 只能确认存在的交易
     */
    function test_confirmTransaction() public{
        vm.prank(address(3));
        vm.expectEmit(true, true, true, false);
        emit SubmitTransaction(address(this), multiSigWallet.getTransactionCnt(), address(0x6), 1 ether, "");
        uint txIndex=multiSigWallet.submitTransaction(address(0x6), 1 ether, "");
        vm.prank(address(2));
        vm.expectEmit(true, true, false, false);
        emit ConfirmTransaction(address(2), txIndex);
        multiSigWallet.confirmTransaction(txIndex);
    }
    /**
     * @dev 撤回确认
     */
    function test_revokeConfirmation() public{
        vm.prank(address(3));
        vm.expectEmit(true, true, true, false);
        emit SubmitTransaction(address(this), multiSigWallet.getTransactionCnt(), address(0x6), 1 ether, "");
        uint txIndex=multiSigWallet.submitTransaction(address(0x6), 1 ether, "");
        vm.prank(address(2));
        vm.expectEmit(true, true, false, false);
        emit ConfirmTransaction(address(2), txIndex);
        multiSigWallet.confirmTransaction(txIndex);
        vm.prank(address(6));
        vm.expectRevert("Not an owner");
        multiSigWallet.revokeConfirmation(txIndex);

        vm.prank(address(1));
        vm.expectRevert("Transaction has not confirmed by the yourself");
        multiSigWallet.revokeConfirmation(txIndex);

        vm.prank(address(2));
        vm.expectEmit(true, true, true, false);
        emit RevokeConfirmation(address(2), txIndex);
        multiSigWallet.revokeConfirmation(txIndex);
    }
    /**
     * @dev 测试执行交易提案
     */
    function test_executeTransaction() public{
        vm.expectEmit(true, true, true, false);
        emit SubmitTransaction(address(this), multiSigWallet.getTransactionCnt(), address(0x6), 1 ether, "");
        uint txIndex=multiSigWallet.submitTransaction(address(0x6), 1 ether, "");
        (, address creator, , , , )=multiSigWallet.txDetails(txIndex);
        assertEq(creator,address(this));
        vm.prank(address(2));
        vm.expectEmit(true, true, false, false);
        emit ConfirmTransaction(address(2), txIndex);
        multiSigWallet.confirmTransaction(txIndex);

        vm.prank(address(1));
        vm.expectEmit(true, true, false, false);
        emit ConfirmTransaction(address(1), txIndex);
        multiSigWallet.confirmTransaction(txIndex);

        vm.prank(address(this));
        vm.expectEmit(true, true, false, false);

        emit ExecuteTransaction(address(this), txIndex);
        multiSigWallet.executeTransaction(txIndex);
    }

}