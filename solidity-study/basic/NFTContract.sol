// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
/*
定义NFT结构体（id、owner、price、forSale）
铸造NFT功能
上架/下架功能
购买功能
查询所有在售NFT
*/
contract PicNFTContract{
    struct NFT{
        uint256 id;
        string addr;
        address owner;
        uint price;
        bool forSale;
    }
    NFT[] public nfts;
    uint id=1;
    mapping (address=>NFT[]) userNfts;
    function createNFT(string calldata addr,uint price) public {
        NFT memory newNft=nfts.push();
        newNft.id=++id;
        newNft.addr=addr;
        newNft.price=price;
        newNft.forSale=false;
        newNft.owner=msg.sender;
    }
    function setSale(uint _id,bool sale) public{
        nfts[_id-1].forSale=sale;
    }
    function buy(uint _id) public payable{
        NFT memory nft=nfts[_id-1];
        require(nft.forSale==true,"not for sale");
        require(msg.value>=nft.price,"not enough money");
        nft.owner=msg.sender;
        nft.forSale=false;
    }
    function getAllNfts() public view returns(NFT[] memory,uint){
        NFT[] memory result=new NFT[](nfts.length);
        uint len=nfts.length;
        uint l=0;
        for(uint i=0;i<len;i++){
            NFT memory nft=nfts[i];
            if(nft.forSale==true){
                result[l++]=nfts[i];
            }
        }
        return (result,l);
    }
}