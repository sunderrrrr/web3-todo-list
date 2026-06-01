// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/TodoList.sol";

contract TodoListTest is Test {
    TodoList public todoList;
    
    address alice = address(0x1);
    address bob = address(0x2);

    function setUp() public {
        todoList = new TodoList();
    }

    function testCreateTodo() public {
        vm.startPrank(alice);
        
        todoList.create("Buy ETH");

        TodoList.Todo memory todo = todoList.getTodo(1);
        
        assertEq(todo.text, "Buy ETH");
        assertEq(todo.isCompleted, false);
        assertEq(todo.owner, alice);
        
        vm.stopPrank();
    }

    function testToggleComplete() public {
        vm.startPrank(alice);
        todoList.create("Buy ETH");
        
        todoList.toggleComplete(1);
        TodoList.Todo memory todo = todoList.getTodo(1);
        assertTrue(todo.isCompleted);
        
        todoList.toggleComplete(1);
        todo = todoList.getTodo(1);
        assertFalse(todo.isCompleted);
        
        vm.stopPrank();
    }

    function testCannotToggleOthersTodo() public {
        vm.prank(alice);
        todoList.create("Alice's todo");

        vm.prank(bob);
        vm.expectRevert("Not ur todo");
        todoList.toggleComplete(1);
    }

    function testDeleteIsSoftDelete() public {
        vm.prank(alice);
        todoList.create("Secret note");
        
        vm.prank(alice);
        todoList.deleteTodo(1);

        TodoList.Todo memory todo = todoList.getTodo(1);
        assertTrue(todo.isDeleted);
        
        assertEq(todo.text, "Secret note");
        assertEq(todo.owner, alice);
    }

    function testgetAllTodos() public {
        vm.prank(alice);
        todoList.create("Alice todo 1");
        
        vm.prank(alice);
        todoList.create("Alice todo 2");
        
        vm.prank(bob);
        todoList.create("Bob todo");

        vm.prank(alice);
        TodoList.Todo[] memory aliceTodos = todoList.getAllTodos();
        
        assertEq(aliceTodos.length, 2);
        assertEq(aliceTodos[0].text, "Alice todo 1");
        assertEq(aliceTodos[1].text, "Alice todo 2");
    }
    function testUpdateText() public {
        vm.prank(alice);
        todoList.create("Old text");

        vm.prank(alice);
        todoList.updateText(1, "New text");

        TodoList.Todo memory todo = todoList.getTodo(1);
        assertEq(todo.text, "New text");
    }

    function testCannotUpdateOthersTodo() public {
        vm.prank(alice);
        todoList.create("Alice's todo");

        vm.prank(bob);
        vm.expectRevert("Not ur todo");
        todoList.updateText(1, "Hacked text");
    }

    function testCannotUpdateDeletedTodo() public {
        vm.prank(alice);
        todoList.create("Will be deleted");

        vm.prank(alice);
        todoList.deleteTodo(1);

        vm.prank(alice);
        vm.expectRevert("Todo already deleted");
        todoList.updateText(1, "Ghost edit");
    }
}