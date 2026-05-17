// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
/**
 * @title NFTAInstance
 * @dev 这是一个ERC721代币合约，用于实现NFTA的实例化功能。
 */
contract NFTAInstance is ERC721,ERC721URIStorage,Ownable{ 
    // tokenId计数器
    uint256 public tokenIdCounter;
    // 最大供应量
    uint256 public constant MAX_SUPPLY = 10000;
    // 铸造价格
    uint256 public MINT_PRICE = 0.05 ether;

    /**
     * @dev NFT铸造事件
     * @param minter 铸造者地址
     * @param tokenId 新创建的Token ID
     * @param uri 元数据URI
     */
    event NFTMinted(
        address indexed minter, 
        uint256 indexed tokenId, 
        string uri
    );
    /**
     * @dev NFT拥有者变更事件
     * @param tokenId 变更的Token ID
     * @param oldOwner 旧拥有者地址
     * @param newOwner 新拥有者地址
     */
    event OwnerChanged(
        uint256 indexed tokenId,
        address indexed oldOwner,
        address indexed newOwner
    );
    /**
     * @dev 提现事件
     * @param owner 提现者地址
     * @param amount 提现金额
     */
    event Withdraw(address indexed owner, uint256 amount);
    /**
     * @dev 构造函数
     * @param _name NFT名称
     * @param _symol NFT符号
     */
    constructor(string memory _name,string memory _symol) ERC721(_name,_symol) Ownable(msg.sender) {
        tokenIdCounter = 0;
    }
    /**
     * @dev 铸造NFT
     * @param uri NFT的元数据URI（通常是IPFS链接）
     * @return 新创建的Token ID
     * @notice 需要支付mintPrice的ETH才能铸造
     */
    function mint(string memory uri) public payable returns (uint256) {
        // 检查供应量限制
        require(tokenIdCounter < MAX_SUPPLY, "Max supply reached");
        
        // 检查支付金额
        require(msg.value >= MINT_PRICE, "Insufficient payment");
        
        // 递增计数器
        tokenIdCounter++;
        uint256 newTokenId = tokenIdCounter;
        
        // 安全铸造NFT
        _safeMint(msg.sender, newTokenId);
        
        // 设置元数据URI
        _setTokenURI(newTokenId, uri);
        
        // 触发事件
        emit NFTMinted(msg.sender, newTokenId, uri);
        
        return newTokenId;
    }
    /**
     * @dev 修改NFT拥有者
     */
    function changeOwner(address newOwner,uint _tokenId) public onlyOwner {
        require(newOwner != address(0), "Zero address");
        address old = ownerOf(_tokenId);
        _transfer(ownerOf(_tokenId), newOwner, _tokenId);
        emit OwnerChanged(_tokenId, old, newOwner);
    }

    /**
     * @dev 重写tokenURI函数
     * @param tokenId Token ID
     * @return 元数据URI
     * @notice 需要重写以解决多重继承的冲突
     */
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    /**
     * @dev 检查接口支持
     * @param interfaceId 接口ID
     * @return 是否支持该接口
     * @notice 实现ERC165标准，支持接口查询
     */
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }

    /**
     * @dev 查询总供应量
     * @return 已铸造的NFT数量
     */
    function totalNFTCnt() public view returns (uint256) {
        return tokenIdCounter;
    }
    /**
     * @dev 提取铸造费用
     * @notice 只有合约所有者可以调用
     */
    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No balance to withdraw");
        (bool success,)=owner().call{value:balance}("");
        require(success,"With draw fail");
        emit Withdraw(owner(), balance);
    }
    /**
     * @dev 设置铸造价格
     * @param newPrice 新的铸造价格（wei）
     * @notice 只有合约所有者可以调用
     */
    function setMintPrice(uint256 newPrice) public onlyOwner {
        MINT_PRICE = newPrice;
    }

    function getName() public view returns (string memory) {
        return super.name();
    }

    function getSybol() public view returns (string memory) {
        return super.symbol();
        
    }

    function getTotalSupply() public pure returns (uint256) {
        return MAX_SUPPLY;
    }

    function getAllNFTs() public view returns (uint256[] memory) {
        uint256[] memory tokenIds = new uint256[](tokenIdCounter);
        for (uint256 i = 0; i < tokenIdCounter; i++) {
            tokenIds[i] = i + 1; // Token IDs从1开始
        }
        return tokenIds;
    }
}