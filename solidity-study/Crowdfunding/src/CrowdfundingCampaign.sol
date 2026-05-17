// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

contract CrowdfundingCampaign {
    enum State {
        Preparing,
        Funding,
        Success,
        Failed,
        Closed
    }
    uint256 public immutable MIN_CONTRIBUTION;
    address public owner;
    uint256 public goal;
    uint256 public deadline;
    uint256 public total;
    bool public unlocked = true;
    mapping(address => uint256) public contributions;
    State public state;

    event Contribute(address indexed contributor, uint256 amount);
    event WithDraw(address indexed owner, uint256 amount);
    event Start();
    event Refund(address indexed addr, uint256 amount);
    event Finalize(uint256 time, uint256 total, bool success);
    /**
     * @dev 状态机
     * @param _state 状态
     */
    modifier inState(State _state) {
        require(state == _state, "Not in current state");
        _;
    }
    /**
     * @dev 权限控制,仅Owner可调用
     */
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }
    /**
     * @dev 锁定
     */
    modifier unLocked() {
        require(unlocked, "Locked");
        unlocked = false;
        _;
        unlocked = true;
    }
    /**
     * @dev 验证时间
     */
    modifier isValid() {
        require(deadline >= block.timestamp, "Due date");
        _;
    }
    /**
     * @dev 验证金额
     */
    modifier MoreThanMinContribute() {
        require(msg.value >= MIN_CONTRIBUTION, "Less than min contribute");
        _;
    }

    /**
     * @dev 构造函数
     * @param _goal 目标金额
     * @param _duration 众筹持续时间
     */
    constructor(uint256 _goal, uint256 _duration, uint256 _minContribution) {
        owner = msg.sender;
        goal = _goal;
        deadline = block.timestamp + _duration;
        state = State.Preparing;
        MIN_CONTRIBUTION = _minContribution;
    }

    /**
     * @dev 贡献函数
     */
    function contribution() public payable inState(State.Funding) isValid MoreThanMinContribute {
        contributions[msg.sender] += msg.value;
        total += msg.value;
        emit Contribute(msg.sender, msg.value);
    }

    /**
     * @dev 提现函数
     */
    function withdraw() public inState(State.Success) onlyOwner unLocked {
        total = 0;
        state = State.Closed;
        (bool success,) = msg.sender.call{value: address(this).balance}("");
        require(success, "Withdraw failed");
        emit WithDraw(msg.sender, address(this).balance);
    }

    /**
     * @dev 退款函数
     */
    function refund() public inState(State.Failed) unLocked {
        uint256 contributed = contributions[msg.sender];
        require(contributed > 0, "No contributions");
        contributions[msg.sender] = 0;
        total -= contributed;
        state = State.Closed;
        (bool success,) = msg.sender.call{value: contributed}("");
        require(success, "Refund failed");
        emit Refund(msg.sender, contributed);
    }

    /**
     * @dev 启动众筹
     */
    function start() public onlyOwner inState(State.Preparing) {
        state = State.Funding;
        emit Start();
    }

    /**
     * @dev 结束众筹
     */
    function finalize() external inState(State.Funding) {
        require(block.timestamp > deadline, "Not due");
        if (total >= goal) {
            state = State.Success;
        } else {
            state = State.Failed;
        }
        emit Finalize(block.timestamp, total, state == State.Success);
    }

    /**
     * @dev 获取某个人的贡献
     */
    function getContribution(address addr) public view returns (uint256) {
        return contributions[addr];
    }

    /**
     * @dev 获取众筹状态
     */
    function status() public view returns (string memory) {
        if (state == State.Preparing) {
            return "Preparing";
        } else if (state == State.Funding) {
            return "Funding";
        } else if (state == State.Success) {
            return "Success";
        } else if (state == State.Failed) {
            return "Failed";
        } else {
            return "Closed";
        }
    }

    /**
     * @dev 获取进度
     */
    function progress() public view returns (uint256) {
        return total * 100 / goal;
    }
}
