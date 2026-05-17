
contract  MultiSigWallet{
    event AddUser(address indexed user);
    event RemoveUser(address indexed opUser,address indexed user);
    event Deposit(address indexed depositor, uint value); 
    event SubmitTransaction(address indexed creator, uint indexed id, address indexed to, uint value, bytes data);
    event DeleteTransaction(address indexed deletor, uint indexed id);
    event ConfirmTransaction(address indexed confirmer, uint indexed id);
    event RevokeConfirmation(address indexed confirmer, uint indexed id);
    event ExecuteTransaction(address indexed confirmer, uint indexed id);
    struct Transaction {
        address to;
        address creator;
        uint256 value;
        bytes data;
        bool executed;
        uint256 confirmCnt;
    }


    address[] public owners;
    uint256 public required;
    bool public isLocked= false;
    mapping(address => bool) public isOwner;
    Transaction[] public transactions;
    mapping(uint256 => mapping(address => bool)) public confirmations;
    /**
     * @dev modifier 仅限所有者
     */
    modifier onlyOwner() {
        require(isOwner[msg.sender], "Not an owner");
        _;
    }
    /**
     * @dev modifier 仅限创建者
     */
    modifier onlyCreator(uint txIndex){
        require(transactions[txIndex].creator == msg.sender, "Not the creator");
        _;
    }
    /**
     * @dev modifier 仅限未执行
     */
    modifier nonExected(uint txIndex) {
        require(!transactions[txIndex].executed, "Transaction has already been executed");
        _;
    }
    /**
     * @dev modifier 仅限未确认
     */
    modifier nonConfirmed(uint txIndex) {
        require(!confirmations[txIndex][msg.sender], "Transaction has already been confirmed");
        _;
    }
    /**
     * @dev modifier 仅限已确认
     */
    modifier confirmed(uint txIndex) {
        require(confirmations[txIndex][msg.sender], "Transaction has not confirmed by the yourself");
        _;
    }
    /**
     * @dev modifier 达到所需确认数
     */
    modifier greaterRequiredCnt(uint txIndex) {
        require(transactions[txIndex].confirmCnt >= required, "Required number of confirmations not reached");
        _;
    }
    /**
     * @dev 锁定合约(防止重入攻击)
     */
    modifier nonLock(){
        require(!isLocked, "Contract is locked");
        isLocked = true;
        _;
        isLocked = false;
    }
    /**
     * @dev 构造函数
     * @param _owners 多签钱包的所有者地址数组
     * @param _required 执行交易所需的确认数 
     */
    constructor(address[] memory _owners, uint256 _required) {
        require(_owners.length > 0, "Owners required");
        require(_required > 0 && _required <= _owners.length, "Invalid number of required confirmations");
        for (uint256 i = 0; i < _owners.length; i++) {
            address owner = _owners[i];
            require(owner != address(0), "Invalid owner");
            require(!isOwner[owner], "Owner is not unique");
        }
        owners = _owners;
        required = _required;
        for (uint256 i = 0; i < _owners.length; i++) {
            isOwner[_owners[i]] = true;
        }
    }
    /**
     * @dev 添加所有者
     * @param _owner 新所有者地址
     */
    function addOwner(address _owner) public onlyOwner {
        require(_owner != address(0), "Invalid owner");
        require(!isOwner[_owner], "Owner is not unique");
        owners.push(_owner);
        isOwner[_owner] = true;
        emit AddUser(_owner);
    }
    /**
     * @dev 获取所有者数量
     * @return 所有者数量
     */
    function getOwnerCnt() public view returns (uint256) {
        return owners.length;
    }
    /**
     * @dev 检查地址是否是所有者
     * @param _owner 检查的地址
     */
    function checkOwner(address _owner) public view returns (bool) {
        return isOwner[_owner];
    }
    /**
     * @dev 获取多签钱包余额
     * @return 多签钱包余额
     */
    function getBalance() public view returns (uint256) {
        return address(this).balance;
    }
    /**
     * @dev 移除所有者
     * @param _owner 要移除的所有者地址
     */
    function removeOwner(address _owner) public onlyOwner {
        require(isOwner[_owner], "Not an owner");
        isOwner[_owner] = false;
        for (uint256 i = 0; i < owners.length; i++) {
            if (owners[i] == _owner) {
                owners[i] = owners[owners.length - 1];
                owners.pop();
                break;
            }
        }
        emit RemoveUser(msg.sender, _owner);
    }
    /**
     * @dev 接收以太币
     * @dev 触发 Deposit 事件
     */
    receive() external payable {
        if (msg.value > 0) {
            emit Deposit(msg.sender, msg.value);
        }
    }
    /**
     * @dev 处理未匹配的函数调用
     * @dev 触发 Deposit 事件
     */
    fallback() external payable {
        if (msg.value > 0) {
            emit Deposit(msg.sender, msg.value);
        }
    }
    /**
     * @dev 创建交易提案
     * @param _to 交易接收者地址
     * @param _value 交易金额
     * @param _data 交易数据
     * @dev 触发 SubmitTransaction 事件
     */
    function submitTransaction(address _to, uint256 _value, bytes memory _data) public onlyOwner returns (uint256) {
        uint256 txIndex = transactions.length;
        transactions.push(Transaction({
            to: _to,
            creator: msg.sender,
            value: _value,
            data: _data,
            executed: false,
            confirmCnt: 0
        }));
        emit SubmitTransaction(msg.sender, txIndex, _to, _value, _data);
        return txIndex;
    }
    /** 
     * @dev 删除交易提案
     * @param _txIndex 要删除的提案ID
     * @dev 触发 DeleteTransaction 事件
    */
    function deleteTransaction(uint256 _txIndex) public onlyOwner onlyCreator(_txIndex) nonExected(_txIndex) {
        require(_txIndex < transactions.length, "Invalid transaction index");
        transactions[_txIndex] = transactions[transactions.length - 1];
        transactions.pop();
        emit DeleteTransaction(msg.sender, _txIndex);
    }
    /**
     * @dev 获取交易详情
     * @param _txIndex 要获取的提案ID
     * @return to 提案接收者地址
     * @return creator 提案创建者地址
     * @return value 提案金额
     * @return data 提案数据
     * @return executed 提案是否已执行
     * @return confirmCnt 提案已确认次数
    */
    function txDetails(uint256 _txIndex) public view returns (address to, address creator, uint256 value, bytes memory data, bool executed, uint256 confirmCnt) {
        require(_txIndex < transactions.length, "Invalid transaction index");
        Transaction storage transaction = transactions[_txIndex];
        return (transaction.to, transaction.creator, transaction.value, transaction.data, transaction.executed, transaction.confirmCnt);
    }
    /**
     * @dev 确认交易提案
     * @param _txIndex 要确认的提案ID
     * @dev 触发 ConfirmTransaction 事件
    */
    function confirmTransaction(uint256 _txIndex) public onlyOwner nonExected(_txIndex) nonConfirmed(_txIndex) {
        require(_txIndex < transactions.length, "Invalid transaction index");
        Transaction storage transaction = transactions[_txIndex];
        confirmations[_txIndex][msg.sender] = true;
        transaction.confirmCnt += 1;
        emit ConfirmTransaction(msg.sender, _txIndex);
    }
    /**
     * @dev 取消确认交易提案
     * @param _txIndex 要取消确认的提案ID
     * @dev 触发 RevokeConfirmation 事件
    */
    function revokeConfirmation(uint256 _txIndex) public onlyOwner nonExected(_txIndex) confirmed(_txIndex) {
        require(_txIndex < transactions.length, "Invalid transaction index");
        Transaction storage transaction = transactions[_txIndex];
        confirmations[_txIndex][msg.sender] = false;
        transaction.confirmCnt -= 1;
        emit RevokeConfirmation(msg.sender, _txIndex);
    }
    /**
     * @dev 执行交易提案
     * @param _txIndex 要执行的提案ID
     * @dev 触发 ExecuteTransaction 事件
    */
    function executeTransaction(uint256 _txIndex) public onlyCreator(_txIndex) nonExected(_txIndex) greaterRequiredCnt(_txIndex) nonLock { 
        require(_txIndex < transactions.length, "Invalid transaction index");
        Transaction storage transaction = transactions[_txIndex];
        transaction.executed = true;
        (bool success, ) = transaction.to.call{value: transaction.value}(transaction.data);
        require(success, "Transaction execution failed");
        emit ExecuteTransaction(msg.sender, _txIndex);
    }
    /**
     * @dev 获取交易提案数量
     * @return 交易提案数量
    */    
    function getTransactionCnt() public view returns (uint256) {
        return transactions.length;
     }
    /**
     * @dev 获取所需确认数
     * @return 所需确认数
    */
    function getRequired() public view returns (uint256) {
        return required;
    }
    /**
     * @dev 获取所有者地址列表
     * @return 所有者地址列表
    */
    function getOwners() public view returns (address[] memory) {
        return owners;
    }
}