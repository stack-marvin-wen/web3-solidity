// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
import "forge-std/Script.sol";
import {CrowdfundingFactory} from "../src/CrowdfundingFactory.sol";

contract CrowdfundingFactoryScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);
        new CrowdfundingFactory();
        vm.stopBroadcast();
    }
}
