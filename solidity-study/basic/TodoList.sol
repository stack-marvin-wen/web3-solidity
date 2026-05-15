// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
/*
每个用户有自己的待办列表
可以添加、完成、删除待办
可以查看所有待办和已完成的待办
限制每个用户最多100个待办事项
*/
contract TodoList{
    enum TodoStatus{
        Pending,
        Completed,
        Deleted
    }
    mapping (address=>string[]) public todoList;
    mapping (address=>TodoStatus[]) public todoStatus;
    uint public constant MAX_TODOCNT=100;
    uint public pendingCnt=0;
    uint public completedCnt=0;
    function addTask(string calldata task_name) public {
        string[] memory tasks=todoList[msg.sender];
        require(tasks.length<MAX_TODOCNT,"More than max tasks");
        todoList[msg.sender].push(task_name);
        todoStatus[msg.sender].push(TodoStatus.Pending);
        pendingCnt++;
    }
    function completeTask(uint task_index) public {
        require(task_index<todoList[msg.sender].length,"Task index out of range");
        require(todoStatus[msg.sender][task_index]==TodoStatus.Pending,"Task is not pending");
        todoStatus[msg.sender][task_index]=TodoStatus.Completed;
        pendingCnt--;
        completedCnt++;
    }
    function deleteTask(uint task_index) public {
        uint len=todoList[msg.sender].length;
        require(task_index<todoList[msg.sender].length,"Task index out of range");
        require(todoStatus[msg.sender][task_index]==TodoStatus.Completed,"Task is not completed");
        TodoStatus status=todoStatus[msg.sender][task_index];
        if(status==TodoStatus.Pending) pendingCnt--;
        else if(status==TodoStatus.Completed) completedCnt--;
        todoList[msg.sender][task_index]=todoList[msg.sender][len-1];
        todoStatus[msg.sender][task_index]=todoStatus[msg.sender][len-1];
        todoList[msg.sender].pop();
        todoStatus[msg.sender].pop();
    }
    function getAllCnt() view public returns(uint,uint){
        return (pendingCnt,completedCnt);
    }
}