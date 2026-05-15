// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
contract EnumContract{
    enum OrderStatus{
        Created,
        Paid,
        Shipped,
        Devlivered,
        Canceled
    }
    OrderStatus status;
    function createOrder() public {
        status=OrderStatus.Created;
    }
    function payOrder() public {
        require(status==OrderStatus.Created,"Order is not in created state");
        status=OrderStatus.Paid;
    }
    function shipOrder() public {
        require(status==OrderStatus.Paid,"No paid");
        status=OrderStatus.Shipped;
    }
    function isPaid() public view returns(bool){
        return status==OrderStatus.Paid;
    }
    function getStatusAsUint() public view returns(uint){
        return uint(status);
    }

    function setStatusFromUint(uint _status) public {
        require(_status <= uint(OrderStatus.Canceled), "Invalid status");
        status=OrderStatus(_status);
    }
}
contract Crowdfunding{
    enum ProjectStatus {
        Fundraising,  // 募资中
        Successful,   // 成功
        Failed        // 失败
    }
    ProjectStatus public status = ProjectStatus.Fundraising;
    uint public goal = 100 ether;
    uint public raised;
    function contribute() public payable {
        require(status == ProjectStatus.Fundraising, "Not fundraising");
        raised += msg.value;
    }
    
    function finalize() public {
        require(status == ProjectStatus.Fundraising, "Already finalized");
        
        if (raised >= goal) {
            status = ProjectStatus.Successful;
        } else {
            status = ProjectStatus.Failed;
        }
    }
}