// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
interface IERC20 {
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
}
/*
使用接口调用确保类型安全
检查返回值，确保转账成功
添加事件日志，记录每次交换
实现1:1的交换比例
*/
contract Swap{
    IERC20 public  tokenA ;
    IERC20 public  tokenB ;
    event Swapped(address indexed sender, uint256 amountA, uint256 amountB);
    constructor(address _tokenA, address _tokenB){
        tokenA = IERC20(_tokenA);
        tokenB = IERC20(_tokenB);
    }
    /*用户的tokenA交换为tokenB*/
    function swap(address to,uint256 amount) public {
        require(tokenA.transfer(to, amount),"swap fail");
        require(tokenB.transferFrom(to, msg.sender, amount),"swap fail");        
        emit Swapped(msg.sender,amount,amount);
    }
       
}