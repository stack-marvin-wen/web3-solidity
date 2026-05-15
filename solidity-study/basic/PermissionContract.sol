// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract Ownable{
    address public Owner;
    event OwnerShipChange(address indexed oldOwner, address indexed newOwner);
    constructor(){
        Owner=msg.sender;
    }
    modifier onlyOwner(){
        require(msg.sender==Owner,"Operator is not the owner");
        _;
    }
    function TransferOwnerShip(address newOwner) public onlyOwner{
        require(newOwner!=address(0),"New Owner can not be 0x0");
        address oldOwner=Owner;
        Owner=newOwner;
        emit OwnerShipChange(oldOwner,newOwner);
    }
}
contract PauseAble{
    bool public isPause;
    event Pause();
    event Unpause();
    constructor(){
        isPause=false;
    }
    modifier whenPause(){
        require(isPause,"Aleardy in un pause");
        _;
    }
    modifier whenNotPause(){
        require(!isPause,"Already in pause");
        _;
    }
    function _pause() public whenNotPause {
        isPause=true;
        emit Pause();
    }
    function _unPause() public whenPause {
        isPause=false;
        emit Unpause();
    }
}
contract PermissionContract is Ownable,PauseAble{
    constructor() Ownable() PauseAble(){

    }
    function pause() public onlyOwner{
        _pause();
    }
    function unPause() public onlyOwner{
        _unPause();
    }
}