// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;
contract MyERCToken{
    enum State{
        PAUSE,
        UNPAUSE
    }
    uint immutable public MAX_BATCH_SIZE;

    string public  name;
    string public symbol;
    uint public decimals;
    uint public totalSupply; 


    mapping (address=>uint) public balanceOf;
    mapping (address=>mapping (address=>uint)) public allowance;


    address public OWNER;


    event Transfer(address indexed from, address indexed to, uint value);
    event Approval(address indexed owner, address indexed spender, uint value);

    State public state;

    modifier onlyOwner() {
        require(msg.sender == OWNER, "Only owner can call this");
        _;
    }

    modifier inState(State s){
        require(state==s,"Not in right state");
        _;
    }

    constructor(
        string memory _name,
        string memory _symbol,
        uint8 _decimals,
        uint256 _initialSupply,
        uint max_batch_size
    ) {
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        totalSupply = _initialSupply * 10**_decimals;
        OWNER = msg.sender;
        balanceOf[msg.sender] = totalSupply;
        MAX_BATCH_SIZE=max_batch_size;
        emit Transfer(address(0), msg.sender, totalSupply);
    }

    function transfer(address to, uint256 amount) public inState(State.UNPAUSE) returns (bool){
        require(amount>=0, "Amount should be greater than zero");
        require(balanceOf[msg.sender]>=amount, "Not enough tokens");
        require(to!=address(0),"Not allow to address as zero");
        balanceOf[msg.sender]-=amount;
        balanceOf[to]+=amount;
        emit Transfer(msg.sender, to, amount);
        return true;
    }

    function approve(address spender, uint256 amount) public inState(State.UNPAUSE) returns (bool){
        require(spender!=address(0),"Address is empty");
        require(amount>0,"Amount is zero");
        allowance[msg.sender][spender]=amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }

    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) public inState(State.UNPAUSE) returns (bool){
        require(from!=address(0),"Address is empty");
        require(to!=address(0),"Address is empty");
        require(amount>0,"Amount is zero");
        require(balanceOf[from]>=amount,"Not enough tokens");
        require(allowance[from][msg.sender]>=amount,"Not enough allowance");
        balanceOf[from]-=amount;
        balanceOf[to]+=amount;
        allowance[from][msg.sender]-=amount;
        emit Transfer(from,to,amount);
        return true;
    }

    function mint(address to, uint256 amount) public onlyOwner{
        require(amount>0,"Amount is zero");
        require(to!=address(0),"Address is empty");
        totalSupply+=amount;
        balanceOf[to]+=amount;
        emit Transfer(address(0),to,amount);
    }

    function burn(uint256 amount) public{
        require(balanceOf[msg.sender]>=amount,"Amount is zero");
        require(amount>0,"Amount is zero");
        totalSupply-=amount;
        balanceOf[msg.sender]-=amount;
        emit Transfer(msg.sender,address(0),amount);
    }

    function batchTransfer(address[] memory recipients, uint256[] memory amounts) public returns(bool) {
        uint len=recipients.length;
        require(len==amounts.length,"Mis match the length");
        require(len<=MAX_BATCH_SIZE,"Greater than Max Batch Size");
        uint total=0;
        for(uint i=0;i<len;i++){
            total+=amounts[i];
        }
        require(total<=balanceOf[msg.sender],"Not enough tokens");
        for (uint256 i = 0; i < recipients.length; i++) {
            require(recipients[i] != address(0), "Invalid address");
            require(amounts[i] > 0, "Invalid amount");
        }

        for(uint i=0;i<len;i++){
            balanceOf[msg.sender]-=amounts[i];
            balanceOf[recipients[i]]+=amounts[i];
            emit Transfer(msg.sender,recipients[i],amounts[i]);
        }
        return true;
    }

    function pause() public onlyOwner{
        state=State.PAUSE;
    }
    function unpause() public onlyOwner{
        state=State.UNPAUSE;
    }
}