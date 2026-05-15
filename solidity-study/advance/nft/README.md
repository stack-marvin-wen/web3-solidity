# NFT市场需求分析
NFT遵循ERC721标准：
- 每个代币都有一个唯一的Token ID
- 这使得每个NFT都是独一无二的，不可互换的
- 就像艺术品一样，CryptoPunk 1号和CryptoPunk 2号虽然都是CryptoPunk系列，但它们是完全不同的资产
- 价值可能相差巨大

这种唯一性特征使得ERC721非常适合表示：

- 数字艺术品
- 游戏道具
- 虚拟土地
- 域名
- 身份证明
- 其他具有独特性的资产

## ERC721核心接口
ERC721标准定义了一系列核心接口函数，理解这些函数对于开发NFT市场至关重要，因为我们的市场合约会频繁调用这些接口。
- balanceOf函数：
```solidity
function balanceOf(address owner) external view returns (uint256 balance);
```
1. 用于查询某个地址拥有的NFT数量
2. 注意，虽然它返回的是数量，但每个NFT仍然是独特的
3. 比如一个地址拥有5个NFT，这5个NFT可能完全不同
- ownerOf函数
```solidity
function ownerOf(uint256 tokenId) external view returns (address owner);
```
1. 用于查询某个特定tokenId的所有者
2. 在交易前，我们需要用这个函数验证卖家确实拥有这个NFT
3. 这是交易安全的基础

- safeTransferFrom函数
```solidity
function safeTransferFrom(
    address from,
    address to,
    uint256 tokenId,
    bytes calldata data
) external;
```
1. 安全转移NFT的标准方式
2. 为什么叫"安全"呢？因为它会检查接收方是否能够处理NFT
3. 避免NFT被转移到无法操作的合约地址中
4. 在实现买卖功能时，我们会使用这个函数来转移NFT

- approve函数：
```solidity
function setApprovalForAll(address operator, bool approved) external;
```

1. 批量授权，允许某个地址（通常是市场合约）操作你所有的NFT
2. 用户只需要授权一次，就可以在市场上交易任意数量的NFT
3. 这大大提升了用户体验

## ERC721元数据扩展接口
- name函数
```solidity
function name() external public returns(string memory);
```
1. 返回NFT集合的名称
2. 比如"CryptoPunks"或"Bored Ape Yacht Club"
- symbol函数：
```solidity
function symbol() external view returns (string memory);
```
1. 返回符号
2. 比如"PUNKS"或"BAYC"
- tokenURI函数：
```solidity
function tokenURI(uint256 tokenId) external view returns (string memory);
```
1. 这是最重要的元数据函数
2. 返回一个URI，通常指向一个JSON文件
3. 这个JSON文件包含了NFT的所有元数据信息
4. 当你调用tokenURI并传入一个tokenId时，它会返回一个URL
5. 可能是https链接，也可能是IPFS链接
6. 这个链接指向一个JSON文件，JSON文件里包含了NFT的名称、描述、图片链接以及各种属性

元数据格式
```json
{
  "name": "My Awesome NFT #1",
  "description": "This is a description of my awesome NFT",
  "image": "ipfs://QmZ4tDuvesekSs4qM5ZBKpXiZGun7S2CYtEZRB3DYXkjGx",
  "attributes": [
    {
      "trait_type": "Color",
      "value": "Blue"
    },
    {
      "trait_type": "Rarity",
      "value": "Legendary"
    }
  ]
}
```


## NFT操作
### 铸造（使用OpenZeppelin库）

`OpenZeppelin`提供了专业经过审计的ERC721的非同质化代币接口
1. 定义NFT需要继承ERC721，ERC721URIStorage，Ownable三个标准合约
```solidity
contract MyNFT is ERC721, ERC721URIStorage, Ownable{

}
```
- ERC721：提供了标准的NFT功能，包括所有权管理、转移等
- ERC721URIStorage：让我们可以为每个tokenId设置独立的URI
- Ownable：提供了所有权管理功能，只有owner可以执行某些操作

2. 合约状态变量和构造函数
```solidity
uint256 private _tokenIdCounter;  // 用于生成唯一的tokenId
uint256 public constant MAX_SUPPLY = 10000;  // 最大供应量
uint256 public mintPrice = 0.01 ether;  // 铸造价格

event NewMint(address addr,uint tokenId,string uri);
constructor() ERC721("MyNFT", "MNFT") Ownable(msg.sender) {}

```
构造函数中，我们设置了NFT集合的名称为"MyNFT"，符号为"MNFT"。当然，在实际项目中，你会使用更有意义的名称。
3. 基础铸造功能
```solidity
function mint(string memory uri) returns(uint){
    _tokenIdCounter++;
    newTokenId=_tokenIdCounter;
    require(newTokenId<=MAX_SUPPLY,"超过了最大供应量");
    _safeMint(msg.sender, newTokenId);
    _safeTokenURI(newTokenId,uri);
    return newTokenId;
}
```
`_safeMint`是`OpenZeppelin`提供的内部函数，它会创建NFT并转移到`to`地址，同时进行必要的安全检查。`_setTokenURI`设置元数据URI，这样外部应用就可以通过`tokenURI`函数获取NFT的展示信息。

4. 增强铸造
```solidity
function mint(string memory uri) returns(uint){
    require(msg.value>=mintPrice,"余额不足");
    _tokenIdCounter++;
    newTokenId=_tokenIdCounter;
    require(newTokenId<=MAX_SUPPLY,"超过了最大供应量");
    _safeMint(msg.sender, newTokenId);
    _safeTokenURI(newTokenId,uri);
    emit NewMint(msg.sender,newTokenId,uri);
    return newTokenId;
}
```
5. 其他重要函数
- tokenURI函数重写：
```solidity
function tokenURI(uint256 tokenId) public view override(ERC721,ERC721URIStorage) returns(string memory){
    return super.tokenURI(tokenId);
}
```
重写了tokenURI函数。这是因为我们同时继承了ERC721和ERC721URIStorage，两个合约都实现了这个函数，所以需要明确指定使用哪个实现。这里我们使用super.tokenURI，会按照继承顺序调用正确的实现。
- supportsInterface函数：
```solidity
function supportsInterface(bytes4 interfaceId)
    public
    view
    override(ERC721, ERC721URIStorage)
    returns (bool)
{
    return super.supportsInterface(interfaceId);
}
```
这个函数用于检查合约是否支持某个接口。这对于ERC165标准兼容性很重要，也是ERC2981版税标准的基础。

- withdraw函数：
```solidity
function withdraw() public onlyOwner {
    uint256 balance = address(this).balance;
    payable(owner()).transfer(balance);
}
```
这个函数允许合约所有者提取铸造费用。使用onlyOwner修饰符确保只有所有者可以调用。

### 市场合约
市场合约是一个独立的合约，它会与各种NFT合约进行交互。这种设计有几个优势：
- 通用性：可以支持多种NFT合约，不局限于某个特定合约
- 可升级性：市场逻辑可以独立升级，不影响NFT合约
- 安全性：市场合约和NFT合约分离，降低风险
1. 核心数据结构
```solidity
struct Listing {
    address seller;        // 卖家地址
    address nftContract;   // NFT合约地址
    uint256 tokenId;       // Token ID
    uint256 price;         // 售价（wei）
    bool active;           // 是否激活
}
mapping(uint256 => Listing) public listings;  // 挂单映射
uint256 public listingCounter;                // 挂单计数器

uint256 public platformFee = 250;  // 2.5% 手续费
address public feeRecipient;        // 手续费接收地址
```
- seller：存储卖家地址，只有卖家本人可以下架或修改价格
- nftContract：NFT合约的地址。因为市场要支持各种NFT，不能写死某个特定合约，所以需要记录是哪个NFT合约
- tokenId：指定是哪个具体的NFT。结合nftContract和tokenId，就能唯一确定一个NFT
- price：售价，以wei为单位
- active：表示挂单是否有效。当NFT被购买或卖家主动下架时，这个字段会被设为false，防止重复购买

2. 合约事件
- 新NFT上架
- 删除NFT
- 价格更新
- NFT卖出

3. 上架
```solidity
function listNFT(
    address nftContract,
    uint tokenId,
    uint price
) external returns(uint){
    require(price>0,"Price would be greater than 0");
    IERC721 nft= IERC721(nftContract);
    require(nft.ownerOf(tokenId)==msg.sender,"You aren't the owner");
    require(
        nft.getApproved(tokenId) == address(this) ||
        nft.isApprovedForAll(msg.sender, address(this)),
        "Marketplace not approved"
    );
    listingCounter++;
    listings[listingCounter]=Listing({
        seller:msg.sender,
        nftContract:nftContract,
        tokenId:tokenId,
        price:price,
        active:true
    })
    emit NFTListed(
        listingCounter,
        msg.sender,
        nftContract,
        tokenId,
        price
    );
    return listingCounter;
}
```
安全检查：
- 价格检查：确保价格大于0。虽然理论上可以设置为0实现赠送，但在实际市场中，0价格可能导致一些问题，所以我们要求价格必须大于0。
- 所有权检查：这非常重要。我们调用NFT合约的ownerOf函数，确认调用者确实拥有这个NFT。如果有人试图上架别人的NFT，这里会失败。
- 授权检查：市场合约需要有权限来转移用户的NFT，否则在购买时无法完成转移。我们检查两种授权方式：
    - 通过approve对单个NFT的授权
    - 通过setApprovalForAll对所有NFT的授权
    - 只要有其中一种授权，就可以上架
4. 下架功能
```solidity
function delistNFT(uint listingId){
    require(listingId>0 && listingId<=listingCounterm,"No listId");
    Listing storage temp=listings[listingId];
    require(temp.active,"NFT is not active, not allow de list");
    require(temp.seller==msg.sender,"Not the owner")
    temp.active=false;
    emit NFTDeListed(
        msg.sender,
        listingId
    );
}
```
5. 价格更新
```solidity
function updatePrice(uint256 listingId, uint256 newPrice) external {
    require(newPrice > 0, "Price must be greater than 0");
    
    Listing storage listing = listings[listingId];
    require(listing.active, "Listing not active");
    require(listing.seller == msg.sender, "Not the seller");
    
    listing.price = newPrice;
    
    emit PriceUpdated(listingId, newPrice);
}
```
6. 买卖NFT
```
function buyNFT(uint256 listingId) external payable nonReentrant{
    Listing storage temp=listings[listId];
    
    require(temp.active, "Listing not active");
    require(msg.sender != temp.seller, "Cannot buy your own NFT");
    require(msg.sender>=temp.price,"Unavailiable fee");
    temp.active=false;
    uint fee=temp.price*platformFee/1000;
    uint amount=temp.price-fee;
    IERC721(temp.nftContract).safeTransferFrom(temp.seller,msg.sender,temp.tokenId);
    (bool success,)=temp.seller.call{value: sellerAmount}("");
    require(successSeller, "Transfer to seller failed");
    
    (bool successFee, ) = feeRecipient.call{value: fee}("");
    require(successFee, "Transfer fee failed");
    if (msg.value > temp.price) {
        (bool successRefund, ) = msg.sender.call{
            value: msg.value - listing.price
        }("");
        require(successRefund, "Refund failed");
    }
    emit NFTSold(listingId, msg.sender, listing.seller, listing.price);
}
```
### 版税系统

版税系统是NFT市场的一个重要特性，它能让NFT创作者在每次二次交易中都获得一定比例的收益。

版税系统的作用：

- 保障创作者权益：让创作者能够从NFT的增值中持续受益
- 激励创作：创作者知道即使首次销售后，还能从后续交易中获得收益
- 行业标准：主流NFT市场都支持版税系统

1. ERC2981标准

版税系统基于ERC2981标准实现。ERC2981定义了royaltyInfo函数，接收tokenId和售价作为参数，返回版税接收地址和版税金额。
```solidity
interface IERC2981 is IERC165 {
    function royaltyInfo(
        uint256 tokenId,
        uint256 salePrice
    ) external view returns (
        address receiver,
        uint256 royaltyAmount
    );
}
```
2. 版税检查函数
- _getRoyaltyInfo函数：
```solidity
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
```

- 集成版税的购买函数
```
function buyNFT(uint256 listingId) external payable nonReentrant{
    Listing storage temp=listings[listId];
    
    require(temp.active, "Listing not active");
    require(msg.sender != temp.seller, "Cannot buy your own NFT");
    require(msg.sender>=temp.price,"Unavailiable fee");
    temp.active=false;
    uint fee=temp.price*platformFee/1000;
    // 获取版税信息
    (address royaltyReceiver, uint256 royaltyAmount) = _getRoyaltyInfo(
        listing.nftContract,
        listing.tokenId,
        listing.price
    );




    uint amount=temp.price-fee-royaltyAmount;


    IERC721(temp.nftContract).safeTransferFrom(temp.seller,msg.sender,temp.tokenId);

    // 资金分配顺序：版税 -> 平台手续费 -> 卖家收益
    if (royaltyAmount > 0 && royaltyReceiver != address(0)) {
        (bool successRoyalty, ) = royaltyReceiver.call{value: royaltyAmount}("");
        require(successRoyalty, "Royalty transfer failed");
    }
    (bool successFee, ) = feeRecipient.call{value: fee}("");
    require(successFee, "Transfer fee failed");

    (bool success,)=temp.seller.call{value: sellerAmount}("");
    require(successSeller, "Transfer to seller failed");
    
    
    if (msg.value > temp.price) {
        (bool successRefund, ) = msg.sender.call{
            value: msg.value - listing.price
        }("");
        require(successRefund, "Refund failed");
    }
    emit NFTSold(listingId, msg.sender, listing.seller, listing.price);
}
```

### 拍卖NFT

1. 英式拍卖（English Auction）：

- 起拍价由卖家设定
- 买家可以出价
- 每次出价必须高于当前最高出价
- 拍卖时间结束后，出价最高者获得NFT
- 这是最常见的拍卖方式

2. 荷兰式拍卖（Dutch Auction）：

- 价格从高到低逐渐降低
- 第一个接受当前价格的买家获得NFT
- 适用于快速销售场景

3. 拍卖数据结构
```solidity
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
// 待退款映射：

mapping(uint256 => mapping(address => uint256)) public pendingReturns;
uint public platformFee = 250;  // 2.5% 手续费 
address public feeRecipient;        // 手续费接收地址
```

4. 创建拍卖(createAuction)
```solidity

```
5. 出价功能(placeBid)
6. 退款机制(withdrawBid)
7. 结束拍卖(endAuction)


