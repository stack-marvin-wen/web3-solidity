// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/utils/introspection/IERC165.sol";
interface IERC2981 is IERC165 {
    function royaltyInfo(
        uint256 tokenId,
        uint256 salePrice
    ) external view returns (
        address receiver,
        uint256 royaltyAmount
    );
}
contract NFTAutionEn{
    struct Auction {
        address seller;           // 卖家地址
        address nftContract;      // NFT合约地址
        uint256 tokenId;          // Token ID
        uint256 startPrice;       // 起拍价
        uint256 highestBid;       // 当前最高出价
        address highestBidder;    // 当前最高出价者
        uint256 endTime;          // 拍卖结束时间
        bool active;              // 是否激活
    }
    
    Auction[] public auctions ;
    // 待退款映射：
    mapping(uint256 => mapping(address => uint256)) public pendingReturns;
    uint public platformFee = 250;  // 2.5% 手续费 
    address public feeRecipient;        // 手续费接收地址

    event NFTAuction(address indexed seller,address indexed nftContract,uint indexed tokenId,uint startPrice,uint endTime);
    event NewBid(uint indexed  accid,address indexed  accor,uint price,uint time);
    event AuctionEnded(uint indexed id,address winner,uint amount);
    // 创建拍卖(createAuction)
    function createAuction(address nftContract, uint tokenId, uint startPrice, uint duration) public {
        require(nftContract != address(0), "Invalid NFT contract");
        require(duration>0,"Duration would more than 0");
        IERC721 nft=IERC721(nftContract);
        
        require(nft.ownerOf(tokenId)==msg.sender,"Not the owner");
        require(nft.getApproved(tokenId) == address(this) || nft.isApprovedForAll(msg.sender, address(this)), "Marketplace not approved");
        uint auctionsCnt=auctions.length;
        auctionsCnt++;
        auctions[auctionsCnt] = Auction({
            seller:msg.sender,
            nftContract:nftContract,
            tokenId:tokenId,
            startPrice:startPrice,
            highestBid:startPrice,
            highestBidder:address(0),
            endTime:block.timestamp+duration,
            active:true
        });
        emit NFTAuction(msg.sender,nftContract,tokenId,startPrice,duration);
    }
    // 出价功能(placeBid)
    function placeBid(uint auctionId,uint price) public  returns(bool){
        Auction storage au=auctions[auctionId];
        require(au.active,"Invalid");
        require(au.endTime>=block.timestamp,"Due date");
        require(price>au.highestBid,"Price would more than highestBid");
        require(price>au.startPrice,"Price would more than startPrice");
        if(au.highestBidder!=address(0)){
            pendingReturns[auctionId][au.highestBidder]+=au.highestBid;
        }
        au.highestBidder=msg.sender;
        au.highestBid=price;
        emit NewBid(auctionId,msg.sender,price,block.timestamp);
        return true;
    }
    // 退款机制(withdrawBid)
    function withdrawBid(uint auctionId) public {
        uint256 amount = pendingReturns[auctionId][msg.sender];
        require(amount > 0, "No pending return");
        
        pendingReturns[auctionId][msg.sender] = 0;
        
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
    }
    // 结束拍卖(endAuction)
    function endAuction(uint auctionId) public {
        Auction storage auction = auctions[auctionId];
    
        require(auction.active, "Auction not active");
        require(block.timestamp >= auction.endTime, "Auction not ended");
        
        auction.active = false;
        if(auction.highestBidder!=address(0)){
            uint256 fee = (auction.highestBid * platformFee) / 1000;

            (address royaltyReceiver, uint256 royaltyAmount) = _getRoyaltyInfo(
                auction.nftContract,
                auction.tokenId,
                auction.highestBid
            );
            uint amount=auction.highestBid-fee-royaltyAmount;
            (bool royasucc,)=royaltyReceiver.call{value:royaltyAmount}("");
            require(royasucc,"Royalty pay fail");
            (bool sellsucc,)=auction.seller.call{value:amount}("");
            require(sellsucc,"Seller pay fail");

            (bool feesuccess,)=feeRecipient.call{value:fee}("");
            require(feesuccess,"Platform fee pay fail");
            emit AuctionEnded(
                auctionId,
                auction.highestBidder,
                auction.highestBid
            );
        }
        else{
            emit AuctionEnded(
                auctionId,
                address(0),
                0
            );
        }
    }

    function _getRoyaltyInfo(
        address nftContract,
        uint256 tokenId,
        uint256 salePrice
    ) internal view returns (address receiver, uint256 royaltyAmount) {
        // 检查NFT合约是否支持ERC2981
        if (IERC165(nftContract).supportsInterface(type(IERC2981).interfaceId)) {
            (receiver, royaltyAmount) = IERC2981(nftContract).royaltyInfo(
                tokenId,
                salePrice
            );
        } else {
            // 不支持版税，返回零地址和零金额
            receiver = address(0);
            royaltyAmount = 0;
        }
    }

}
