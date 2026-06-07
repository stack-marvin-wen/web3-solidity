// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {AggregatorV3Interface} from "../src/IAggregatorV3Interface.sol";
import {NFTAuction} from "../src/NFTAuction.sol";

contract MockV3Aggregator is AggregatorV3Interface {
    int256 private answer;
    uint8 private _decimals;

    constructor(uint8 decimals_, int256 answer_) {
        _decimals = decimals_;
        answer = answer_;
    }

    function decimals() external view returns (uint8) {
        return _decimals;
    }

    function description() external pure returns (string memory) {
        return "mock";
    }

    function version() external pure returns (uint256) {
        return 1;
    }

    function getRoundData(
        uint80
    )
        external
        view
        returns (
            uint80,
            int256,
            uint256,
            uint256,
            uint80
        )
    {
        return (0, answer, 0, 0, 0);
    }

    function latestRoundData()
        external
        view
        returns (
            uint80,
            int256,
            uint256,
            uint256,
            uint80
        )
    {
        return (0, answer, 0, 0, 0);
    }
}

contract MockERC20 is ERC20 {
    constructor() ERC20("MockToken", "MTK") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract MockNFT is ERC721 {
    uint256 public nextId = 1;

    constructor() ERC721("MockNFT", "MNFT") {}

    function mint(address to) external returns (uint256) {
        uint256 id = nextId++;
        _mint(to, id);
        return id;
    }
}

contract NFTAuctionTest is Test {
    NFTAuction auction;
    MockV3Aggregator ethFeed;
    MockV3Aggregator tokenFeed;
    MockERC20 token;
    MockNFT nft;

    event AuctionCreated(
        address indexed seller,
        address indexed nftContract,
        uint256 indexed tokenId,
        address paymentToken,
        uint256 startPriceUsd,
        uint256 endTime
    );

    receive() external payable {}
    fallback() external payable {}

    function setUp() public {
        ethFeed = new MockV3Aggregator(8, 2000e8); // 1 ETH = 2000 USD
        tokenFeed = new MockV3Aggregator(8, 2e8);   // 1 MTK = 2 USD
        token = new MockERC20();
        nft = new MockNFT();

        auction = new NFTAuction(address(this), address(ethFeed));
        auction.setTokenUsdFeed(address(token), address(tokenFeed));

        nft.mint(address(this));
        nft.approve(address(auction), 1);

        token.mint(address(1), 1000e18);
        vm.prank(address(1));
        token.approve(address(auction), type(uint256).max);
    }

    function test_createAuction_eth() public {
        uint256 auctionId = auction.createAuction(address(nft), 1, address(0), 1000e18, 1);
        (
            address seller,
            address nftContract,
            uint256 tokenId,
            address paymentToken,
            uint256 startPriceUsd,
            ,
            ,
            ,
            uint256 endTime,
            
        ) = auction.getAuction(auctionId);

        assertEq(seller, address(this));
        assertEq(nftContract, address(nft));
        assertEq(tokenId, 1);
        assertEq(paymentToken, address(0));
        assertEq(startPriceUsd, 1000e18);
        assertGt(endTime, block.timestamp);
    }

    function test_bidEth() public {
        uint256 auctionId = auction.createAuction(address(nft), 1, address(0), 1000e18, 1);
        deal(address(1), 500 ether);
        vm.prank(address(1));
        
        auction.bidEth{value: 1 ether}(auctionId);
        
        (,,,,, uint256 highestBidAmount, uint256 highestBidUsd, address highestBidder,,) =
            auction.getAuction(auctionId);

        assertEq(highestBidAmount, 1 ether);
        assertEq(highestBidUsd, 2000e18);
        assertEq(highestBidder, address(1));
    }

    function test_bidErc20() public {
        uint256 auctionId = auction.createAuction(address(nft), 1, address(token), 10e18, 1);
        deal(address(1), 10 ether);
        vm.prank(address(1));
        auction.bidErc20(auctionId, 10e18);

        (,,,,, uint256 highestBidAmount, uint256 highestBidUsd, address highestBidder,,) =
            auction.getAuction(auctionId);

        assertEq(highestBidAmount, 10e18);
        assertEq(highestBidUsd, 20e18);
        assertEq(highestBidder, address(1));
    }

    function test_endAuction() public {
        uint256 auctionId = auction.createAuction(address(nft), 1, address(0), 1000e18, 1);

        deal(address(1), 500 ether);
        vm.prank(address(1));
        auction.bidEth{value: 1 ether}(auctionId);

        vm.warp(block.timestamp + 2 hours);
        auction.endAuction(auctionId);

        assertEq(nft.ownerOf(1), address(1));
    }
}