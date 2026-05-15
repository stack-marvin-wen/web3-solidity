// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
/*
用户注册（包含name、email）
更新个人资料
存款功能（payable）
查询用户信息
获取所有用户列表
分批查询用户
限制最多1000个用户
*/
contract UserManContract{
    struct User{
        string name;
        string email;
        uint banance;
    }
    
    uint public constant MAX_USER_CNT=1000;
    User[] public users;
    mapping (string=>bool) isExist;
    mapping (string=>uint) mapIndex;
    mapping (address=>string) addressToEmail;
    function register(string calldata name,string calldata email) public {
        uint len=users.length;
        require(!isExist[email], "Email exist");
        require(len<MAX_USER_CNT,"Out of max storage");
        User memory user=User({name:name,email:email,banance:0});
        users.push(user);
        isExist[email]=true;
        mapIndex[email]=len-1;
        addressToEmail[msg.sender]=email;
    }

    function updateUserInfo(string calldata newName,string calldata newEmail)public {
        string memory originalEmail=addressToEmail[msg.sender];
        require(isExist[originalEmail],"User not exist");
        uint _index=mapIndex[originalEmail];
        users[_index].name=newName;
        if (keccak256(bytes(originalEmail))!=keccak256(bytes(newEmail))){
            isExist[originalEmail]=false;
            isExist[newEmail]=true;
            mapIndex[newEmail]=mapIndex[originalEmail];
            users[_index].email=newEmail;
        }
    }
    function deposit(uint amount) public payable{
        string memory originalEmail=addressToEmail[msg.sender];
        uint index=mapIndex[originalEmail];
        users[index].banance+=amount;
    }
    function getUserInfo(string calldata email) public view returns(User memory){
        string memory _email=email;
        if(keccak256(bytes(email))==keccak256("")) {
            _email=addressToEmail[msg.sender];
        }
        require(isExist[email],"User not exist");
        uint _index=mapIndex[email];
        return users[_index];

    }
    function getAllUser() public view returns(User[] memory){
        return users;
    }
    function batchQuery(uint begin,uint end) public view returns(User[] memory){
        require(end >= begin, "Invalid range");
        require(end < users.length, "Index out of bounds");
        User[] memory res=new User[](end-begin+1);
        uint j=0;
        for(uint i=begin;i<end;i++){
            res[j]=users[i];
            j++;
        }
        return res;
    } 
}