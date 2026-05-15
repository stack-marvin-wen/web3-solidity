// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";

contract MyNFT is ERC721,ERC721URIStorage,Ownable{
    uint256 private _tokenIdCounter;  // 用于生成唯一的tokenId
    uint256 public constant MAX_SUPPLY = 10000;  // 最大供应量
    uint256 public mintPrice = 0.01 ether;  // 铸造价格
    address public Owner;
    event NewMint(address addr,uint tokenId,string uri);
    event WithDraw();
    constructor() ERC721("MyNFT", "MNFT") Ownable(msg.sender) {}
    function mint(string memory uri) public payable returns(uint){
        // 检查供应量限制
        require(_tokenIdCounter < MAX_SUPPLY, "Max supply reached");
        
        // 检查支付金额
        require(msg.value >= mintPrice, "Insufficient payment");
        
        // 递增计数器
        _tokenIdCounter++;
        uint256 newTokenId = _tokenIdCounter;
        
        // 安全铸造NFT
        _safeMint(msg.sender, newTokenId);
        
        // 设置元数据URI
        _setTokenURI(newTokenId, uri);
        
        // 触发事件
        emit NewMint(msg.sender, newTokenId, uri);
        
        return newTokenId;
    }
    function tokenURI(uint256 tokenId) public view override(ERC721, ERC721URIStorage) returns (string memory){
        return super.tokenURI(tokenId);
    }
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }

    function totalSupply() public view returns (uint256) {
        return _tokenIdCounter;
    }

    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No balance to withdraw");
        (bool success,)=owner().call{value:balance}("");
        require(success,"Withdraw fail");
        emit WithDraw();
    }
    function setMintPrice(uint256 newPrice) public onlyOwner {
        mintPrice = newPrice;
    }
}