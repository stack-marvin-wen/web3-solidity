// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;
import {Test} from "forge-std/Test.sol";
import {IERC721Receiver} from "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import {NFTAInstance} from "../src/NFTInstance.sol";
contract NFTInstanceTest is Test, IERC721Receiver { 
    /**
     * @dev NFT铸造事件
     * @param minter 铸造者地址
     * @param tokenId 新创建的Token ID
     * @param uri 元数据URI
     */
    event NFTMinted(
        address indexed minter, 
        uint256 indexed tokenId, 
        string uri
    );
    /**
     * @dev NFT拥有者变更事件
     * @param tokenId 变更的Token ID
     * @param oldOwner 旧拥有者地址
     * @param newOwner 新拥有者地址
     */
    event OwnerChanged(
        uint256 indexed tokenId,
        address indexed oldOwner,
        address indexed newOwner
    );
    /**
     * @dev 提现事件
     * @param owner 提现者地址
     * @param amount 提现金额
     */
    event Withdraw(address indexed owner, uint256 amount);
    // NFT实例合约
    NFTAInstance nft;
    /** 
     * @dev 测试环境设置
     */
    function setUp() public {
        nft = new NFTAInstance("TestNFT", "TNFT");
        deal(address(1), 100 ether);
    }
    /**
     * @dev 测试铸造NFT
     */
    function onERC721Received(
        address,
        address,
        uint256,
        bytes calldata
    ) external pure override returns (bytes4) {
        return IERC721Receiver.onERC721Received.selector;
    }
    /**
     * @dev 测试接收ETH
     */
    receive() external payable {}
    fallback() external payable {}
    /**
     * @dev 测试初始化
     */
    function test_init() public view{
        assertEq(nft.getName(), "TestNFT");
        assertEq(nft.getSybol(), "TNFT");
        assertEq(nft.totalNFTCnt(),0);
    }
    /**
     * @dev 测试铸造NFT
     */
    function test_mint() public {
        vm.expectEmit(true, true, false, false);
        emit NFTMinted(address(this) ,1,"https://example.com/1.json"); 
        uint tokenId=nft.mint{value : 0.05 ether}("https://example.com/1");
        assertEq(nft.totalNFTCnt(),1);
        assertEq(nft.ownerOf(tokenId),address(this));
    }
    /**
     * @dev 测试修改NFT拥有者
     */
    function test_changeOwner() public {
        vm.prank(address(1));
        vm.expectEmit(true, true, false, false);
        emit NFTMinted(address(1) ,1,"https://example.com/1");
        uint tokenId=nft.mint{value : 0.05 ether}("https://example.com/1");
        assertEq(nft.ownerOf(tokenId),address(1));

        vm.expectEmit(true, true, false, false);
        emit OwnerChanged(tokenId ,address(1),address(this));
        nft.changeOwner(address(this),tokenId);
        assertEq(nft.ownerOf(tokenId),address(this));
    }
    /**
     * @dev 测试提取铸造费用
     */
    function test_WithDraw() public {
        nft.mint{value : 0.05 ether}("https://example.com/1");
        uint256 balanceBefore = address(this).balance;
        vm.expectEmit(true, true, false, false);
        emit Withdraw(address(this), 0.05 ether);
        nft.withdraw();
        uint256 balanceAfter = address(this).balance;
        assertEq(balanceAfter, balanceBefore + 0.05 ether);
    }
}