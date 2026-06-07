// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {AggregatorV3Interface} from "./IAggregatorV3Interface.sol";
/**
 * @title NFTAuction
 * @dev NFT 拍卖合约，支持 ETH / ERC20 出价，并按 USD 价格统一比较
 */
contract NFTAuction is ReentrancyGuard {
    using SafeERC20 for IERC20;

    /**
     * @dev 拍卖创建事件
     */
    event AuctionCreated(
        address indexed seller,
        address indexed nftContract,
        uint256 indexed tokenId,
        address paymentToken,
        uint256 startPriceUsd,
        uint256 endTime
    );

    /**
     * @dev 最高出价增加事件
     */
    event HighestBidIncreased(
        uint256 indexed auctionId,
        address indexed bidder,
        uint256 amount,
        uint256 amountUsd
    );

    /**
     * @dev 拍卖结束事件
     */
    event AuctionEnded(
        uint256 indexed auctionId,
        address winner,
        uint256 amount,
        uint256 amountUsd
    );

    /**
     * @dev 退款事件
     */
    event Refund(
        uint256 indexed auctionId,
        address indexed bidder,
        uint256 amount
    );

    /**
     * @dev ETH/USD 预言机更新事件
     */
    event EthUsdFeedUpdated(address indexed feed);

    /**
     * @dev ERC20/USD 预言机更新事件
     */
    event TokenUsdFeedUpdated(address indexed token, address indexed feed);

    /**
     * @dev 拍卖结构体
     */
    struct Auction {
        address seller;
        address nftContract;
        uint256 tokenId;
        address paymentToken;     // address(0) 表示 ETH
        uint256 startPriceUsd;    // USD，18 位
        uint256 highestBidAmount;  // 支付币种原始数量
        uint256 highestBidUsd;     // USD，18 位
        address highestBidder;
        uint256 endTime;
        bool active;
    }

    /// @dev 所有拍卖
    Auction[] public auctions;

    /// @dev ETH/USD 价格预言机
    AggregatorV3Interface public ethUsdFeed;

    /// @dev ERC20/USD 价格预言机
    mapping(address => AggregatorV3Interface) public tokenUsdFeed;

    /// @dev 每个拍卖的待退款金额
    mapping(uint256 => mapping(address => uint256)) public pendingWithdrawals;

    /// @dev 每个拍卖的竞拍者列表
    mapping(uint256 => address[]) public bidders;

    /// @dev 已上架 NFT 防重复
    mapping(bytes32 => bool) public listed;

    /// @dev 平台手续费，基点，250 = 2.5%
    uint256 public platformFee = 250;

    /// @dev 手续费接收地址
    address public feeRecipient;

    /**
     * @dev 构造函数
     * @param _feeRecipient 手续费接收地址
     * @param _ethUsdFeed ETH/USD 价格预言机地址
     */
    constructor(address _feeRecipient, address _ethUsdFeed) {
        require(_feeRecipient != address(0), "Invalid recipient");
        require(_ethUsdFeed != address(0), "Invalid feed");

        feeRecipient = _feeRecipient;
        ethUsdFeed = AggregatorV3Interface(_ethUsdFeed);
    }

    /**
     * @dev 生成 NFT 唯一键
     * @param nftContract NFT 合约地址
     * @param tokenId NFT TokenId
     */
    function _nftKey(address nftContract, uint256 tokenId) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(nftContract, tokenId));
    }

    /**
     * @dev 获取支付币种对应的 USD 价格预言机
     * @param paymentToken 支付币种，address(0) 表示 ETH
     */
    function _getFeed(address paymentToken) internal view returns (AggregatorV3Interface) {
        if (paymentToken == address(0)) {
            return ethUsdFeed;
        }

        AggregatorV3Interface feed = tokenUsdFeed[paymentToken];
        require(address(feed) != address(0), "Feed not set");
        return feed;
    }

    /**
     * @dev 将支付金额转换为 USD（18 位精度）
     * @param paymentToken 支付币种，address(0) 表示 ETH
     * @param amount 支付数量
     */
    function amountToUsd(address paymentToken, uint256 amount) public view returns (uint256) {
        AggregatorV3Interface feed = _getFeed(paymentToken);
        (, int256 price, , , ) = feed.latestRoundData();
        require(price > 0, "Invalid price");

        uint8 feedDecimals = feed.decimals();
        uint8 tokenDecimals = paymentToken == address(0)
            ? 18
            : IERC20Metadata(paymentToken).decimals();

        return Math.mulDiv(
            amount,
            uint256(price) * 1e18,
            10 ** (uint256(tokenDecimals) + uint256(feedDecimals))
        );
    }

    /**
     * @dev 创建拍卖
     * @param nftContract NFT 合约地址
     * @param tokenId NFT TokenId
     * @param paymentToken 支付币种，address(0) 表示 ETH
     * @param startPriceUsd 起拍价（USD，18 位）
     * @param durationHours 拍卖持续小时数
     * @return auctionId 拍卖 ID
     */
    function createAuction(
        address nftContract,
        uint256 tokenId,
        address paymentToken,
        uint256 startPriceUsd,
        uint256 durationHours
    ) external returns (uint256 auctionId) {
        require(nftContract != address(0), "Invalid nft address");
        require(startPriceUsd > 0, "Invalid start price");

        bytes32 key = _nftKey(nftContract, tokenId);
        require(!listed[key], "NFT has existed");

        if (paymentToken != address(0)) {
            require(address(tokenUsdFeed[paymentToken]) != address(0), "Feed not set");
        }

        IERC721 nft = IERC721(nftContract);
        require(nft.ownerOf(tokenId) == msg.sender, "You are not the owner of this NFT");
        require(
            nft.getApproved(tokenId) == address(this) ||
                nft.isApprovedForAll(msg.sender, address(this)),
            "NFT is not approved for this contract"
        );

        uint256 endTime = block.timestamp + durationHours * 3600;

        auctions.push(
            Auction({
                seller: msg.sender,
                nftContract: nftContract,
                tokenId: tokenId,
                paymentToken: paymentToken,
                startPriceUsd: startPriceUsd,
                highestBidAmount: 0,
                highestBidUsd: 0,
                highestBidder: address(0),
                endTime: endTime,
                active: true
            })
        );

        listed[key] = true;
        auctionId = auctions.length - 1;

        emit AuctionCreated(msg.sender, nftContract, tokenId, paymentToken, startPriceUsd, endTime);
    }

    /**
     * @dev ETH 出价
     * @param auctionId 拍卖 ID
     */
    function bidEth(uint256 auctionId) external payable nonReentrant {
        Auction storage auction = auctions[auctionId];
        require(auction.paymentToken == address(0), "Use ERC20 bid");
        require(msg.value > 0, "No ETH sent");
        _placeBid(auctionId, msg.value);
    }

    /**
     * @dev ERC20 出价
     * @param auctionId 拍卖 ID
     * @param amount ERC20 数量
     */
    function bidErc20(uint256 auctionId, uint256 amount) external nonReentrant {
        Auction storage auction = auctions[auctionId];
        require(auction.paymentToken != address(0), "Use ETH bid");
        require(amount > 0, "Zero amount");

        IERC20(auction.paymentToken).safeTransferFrom(msg.sender, address(this), amount);
        _placeBid(auctionId, amount);
    }

    /**
     * @dev 内部出价逻辑
     * @param auctionId 拍卖 ID
     * @param amount 支付数量
     */
    function _placeBid(uint256 auctionId, uint256 amount) internal {
        Auction storage auction = auctions[auctionId];

        require(auction.active, "Auction is not active");
        require(block.timestamp < auction.endTime, "Auction has ended");
        require(msg.sender != auction.seller, "Shouldn't auction owner");

        uint256 bidUsd = amountToUsd(auction.paymentToken, amount);
        require(bidUsd >= auction.startPriceUsd, "Bid is too low");
        require(bidUsd > auction.highestBidUsd, "Bid is too low");

        address previousBidder = auction.highestBidder;
        uint256 previousAmount = auction.highestBidAmount;

        if (previousBidder != address(0)) {
            pendingWithdrawals[auctionId][previousBidder] += previousAmount;
        }

        if (previousBidder != msg.sender) {
            bidders[auctionId].push(msg.sender);
        }

        auction.highestBidAmount = amount;
        auction.highestBidUsd = bidUsd;
        auction.highestBidder = msg.sender;

        emit HighestBidIncreased(auctionId, msg.sender, amount, bidUsd);
    }

    /**
     * @dev 结束拍卖
     * @param auctionId 拍卖 ID
     */
    function endAuction(uint256 auctionId)
        public
        inActive(auctionId)
        nonReentrant
        notInValid(auctionId)
    {
        Auction storage auction = auctions[auctionId];
        auction.active = false;
        listed[_nftKey(auction.nftContract, auction.tokenId)] = false;

        if (auction.highestBidder != address(0)) {
            uint256 fee = Math.mulDiv(auction.highestBidAmount, platformFee, 10000);
            uint256 sellerAmount = auction.highestBidAmount - fee;

            IERC721(auction.nftContract).safeTransferFrom(
                auction.seller,
                auction.highestBidder,
                auction.tokenId
            );

            if (auction.paymentToken == address(0)) {
                (bool successSeller, ) = payable(auction.seller).call{value: sellerAmount}("");
                require(successSeller, "Failed to transfer Ether");

                (bool successFee, ) = payable(feeRecipient).call{value: fee}("");
                require(successFee, "Failed to transfer fee");
            } else {
                IERC20 token = IERC20(auction.paymentToken);
                token.safeTransfer(auction.seller, sellerAmount);
                token.safeTransfer(feeRecipient, fee);
            }

            emit AuctionEnded(auctionId, auction.highestBidder, auction.highestBidAmount, auction.highestBidUsd);
        } else {
            emit AuctionEnded(auctionId, address(0), 0, 0);
        }
    }

    /**
     * @dev 领取退款
     * @param auctionId 拍卖 ID
     */
    function refund(uint256 auctionId) external nonReentrant {
        Auction storage auction = auctions[auctionId];
        uint256 amount = pendingWithdrawals[auctionId][msg.sender];
        require(amount > 0, "No funds to withdraw");

        pendingWithdrawals[auctionId][msg.sender] = 0;

        if (auction.paymentToken == address(0)) {
            (bool success, ) = payable(msg.sender).call{value: amount}("");
            require(success, "Failed to transfer Ether");
        } else {
            IERC20(auction.paymentToken).safeTransfer(msg.sender, amount);
        }

        emit Refund(auctionId, msg.sender, amount);
    }

    /**
     * @dev 查询拍卖信息
     * @param auctionId 拍卖 ID
     * @return seller 卖家
     * @return nftContract NFT 合约
     * @return tokenId TokenId
     * @return paymentToken 支付币种
     * @return startPriceUsd 起拍价（USD）
     * @return highestBidAmount 当前最高出价
     * @return highestBidUsd 当前最高出价（USD）
     * @return highestBidder 当前最高出价者
     * @return endTime 结束时间
     * @return active 是否激活
     */
    function getAuction(uint256 auctionId)
        external
        view
        returns (
            address seller,
            address nftContract,
            uint256 tokenId,
            address paymentToken,
            uint256 startPriceUsd,
            uint256 highestBidAmount,
            uint256 highestBidUsd,
            address highestBidder,
            uint256 endTime,
            bool active
        )
    {
        Auction memory auction = auctions[auctionId];
        return (
            auction.seller,
            auction.nftContract,
            auction.tokenId,
            auction.paymentToken,
            auction.startPriceUsd,
            auction.highestBidAmount,
            auction.highestBidUsd,
            auction.highestBidder,
            auction.endTime,
            auction.active
        );
    }

    /**
     * @dev 设置平台手续费
     * @param newFee 新手续费，基点
     */
    function setPlatformFee(uint256 newFee) external {
        require(msg.sender == feeRecipient, "Not fee recipient");
        require(newFee <= 1000, "Fee too high");
        platformFee = newFee;
    }

    /**
     * @dev 更新手续费接收地址
     * @param newRecipient 新接收地址
     */
    function updateFeeRecipient(address newRecipient) external {
        require(msg.sender == feeRecipient, "Not fee recipient");
        require(newRecipient != address(0), "Invalid address");
        feeRecipient = newRecipient;
    }

    /**
     * @dev 更新 ETH/USD 预言机
     * @param newFeed 新预言机地址
     */
    function setEthUsdFeed(address newFeed) external {
        require(msg.sender == feeRecipient, "Not fee recipient");
        require(newFeed != address(0), "Invalid feed");
        ethUsdFeed = AggregatorV3Interface(newFeed);
        emit EthUsdFeedUpdated(newFeed);
    }

    /**
     * @dev 设置某个 ERC20 的 USD 预言机
     * @param token ERC20 地址
     * @param feed USD 预言机地址
     */
    function setTokenUsdFeed(address token, address feed) external {
        require(msg.sender == feeRecipient, "Not fee recipient");
        require(token != address(0), "Invalid token");
        require(feed != address(0), "Invalid feed");
        tokenUsdFeed[token] = AggregatorV3Interface(feed);
        emit TokenUsdFeedUpdated(token, feed);
    }

    /**
     * @dev 获取 ETH/USD 价格
     */
    function getEthPriceUsd() public view returns (uint256) {
        (, int256 price, , , ) = ethUsdFeed.latestRoundData();
        require(price > 0, "Invalid price");
        return uint256(price);
    }

    /**
     * @dev 获取某个拍卖当前最高出价的 USD 价值
     * @param auctionId 拍卖 ID
     */
    function bidHighUsdValue(uint256 auctionId) public view returns (uint256) {
        return auctions[auctionId].highestBidUsd;
    }

    /**
     * @dev 将 ETH 数量转换为 USD
     * @param amountWei ETH 数量（wei）
     */
    function bidValueUsd(uint256 amountWei) public view returns (uint256) {
        return ethToUsd(amountWei);
    }
    /**
     * @dev 获取所有的拍卖
     */
    function getAuctions() public view returns (Auction[] memory) {
        return auctions;
    }
    /**
     * @dev 将 ETH 数量转换为 USD
     * @param amountWei ETH 数量（wei）
     */
    function ethToUsd(uint256 amountWei) public view returns (uint256) {
        uint256 ethPrice = getEthPriceUsd();
        return (amountWei * ethPrice * 1e18) / 1e18;
    }

    /**
     * @dev 当前拍卖是否激活
     * @param _auctionId 拍卖 ID
     */
    modifier inActive(uint256 _auctionId) {
        require(auctions[_auctionId].active, "Auction is not active");
        _;
    }

    /**
     * @dev 当前拍卖是否已结束
     * @param _auctionId 拍卖 ID
     */
    modifier notInValid(uint256 _auctionId) {
        require(block.timestamp >= auctions[_auctionId].endTime, "Auction is not ended");
        _;
    }

    /**
     * @dev 当前拍卖是否未结束
     * @param _auctionId 拍卖 ID
     */
    modifier inValid(uint256 _auctionId) {
        require(block.timestamp < auctions[_auctionId].endTime, "Auction has ended");
        _;
    }
     
}