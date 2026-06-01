// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/TodoToken.sol";
import "../src/TodoList.sol";

contract DeployAllScript is Script {
    function run() external {
        uint256 pk = vm.envUint("PRIVATE_KEY");
        
        vm.startBroadcast(pk);

        TodoToken token = new TodoToken();
        console.log("TodoToken deployed to:", address(token));

        TodoList todoList = new TodoList(address(token));
        console.log("TodoList deployed to:", address(todoList));

        vm.stopBroadcast();
    }
}