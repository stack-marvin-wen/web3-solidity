// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {NFTAuction} from "../src/NFTAuction.sol";
import {MockV3Aggregator} from "../src/MockV3Aggregator.sol";
/**
 * @title DeployERC20
 * @dev 部署 NFT拍卖系统智能合约
 */
contract DeployNFTAuctionScript is Script {
    
    function setUp() public {}

    function run() external returns (address auctionAddress, address feedAddress) {
        vm.startBroadcast();

        MockV3Aggregator feed = new MockV3Aggregator(8, 2000e8);
        NFTAuction auction = new NFTAuction(msg.sender, address(feed));

        vm.stopBroadcast();

        console.log("Mock ETH/USD feed deployed at:", address(feed));
        console.log("NFTAuction deployed at:", address(auction));

        return (address(auction), address(feed));
    }
}