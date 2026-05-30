// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Script.sol";
import "../src/TodoList.sol";

contract DeployScript is Script {
    uint256 devKey = vm.envUint("PRIVATE_KEY");

    function run() external {
        
        vm.startBroadcast(devKey);

        TodoList TDL = new TodoList();

        vm.stopBroadcast();

        console.log("TodoList deployed to:", address(TDL));
    }
}
