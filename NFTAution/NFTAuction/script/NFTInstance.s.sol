// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;
import {Script, console} from "forge-std/Script.sol";

import {NFTAInstance} from "../src/NFTInstance.sol";

contract DeployNFTInstanceScript is Script{
    function run() external returns (address) {
        vm.startBroadcast();
        string memory _name="MyNFT";
        string memory _symol="MNFT";
        // 创建一个 MockV3Aggregator 合约实例
        NFTAInstance nft = new NFTAInstance(_name, _symol);
        console.log("NFTAInstance deployed to:", address(nft));
        vm.stopBroadcast();
        return address(nft);
    }
}