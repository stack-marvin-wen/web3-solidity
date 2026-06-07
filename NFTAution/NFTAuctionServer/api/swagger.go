package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const swaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "NFTAuctionServer API",
    "description": "NFT auction backend API for list, detail, and create operations.",
    "version": "1.0.0"
  },
  "basePath": "/",
  "schemes": ["http"],
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "paths": {
    "/nftauction/logs": {
      "get": {
        "summary": "获取拍卖全量日志",
        "tags": ["NFTAuction"],
        "produces": ["application/json"],
        "parameters": [
          {"name": "event_type", "in": "query", "required": false, "type": "string"},
          {"name": "auction_id", "in": "query", "required": false, "type": "string"},
          {"name": "seller", "in": "query", "required": false, "type": "string"},
          {"name": "nft_contract", "in": "query", "required": false, "type": "string"},
          {"name": "token_id", "in": "query", "required": false, "type": "string"},
          {"name": "bidder", "in": "query", "required": false, "type": "string"},
          {"name": "winner", "in": "query", "required": false, "type": "string"},
          {"name": "tx_hash", "in": "query", "required": false, "type": "string"},
          {"name": "from_block", "in": "query", "required": false, "type": "string"},
          {"name": "to_block", "in": "query", "required": false, "type": "string"}
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {"$ref": "#/definitions/AuctionLogsResponse"}
          }
        }
      }
    },
    "/nftauction/{id}/bids": {
      "get": {
        "summary": "获取出价日志",
        "tags": ["NFTAuction"],
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {"$ref": "#/definitions/BidLogsResponse"}
          }
        }
      }
    },
    "/nftauction/end": {
      "post": {
        "summary": "结束拍卖",
        "tags": ["NFTAuction"],
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {"$ref": "#/definitions/EndNFTAuctionRequest"}
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {"$ref": "#/definitions/CreateNFTAuctionResponse"}
          }
        }
      }
    },
    "/nftauction/bid": {
      "post": {
        "summary": "拍卖出价",
        "tags": ["NFTAuction"],
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {"$ref": "#/definitions/BidNFTAuctionRequest"}
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {"$ref": "#/definitions/CreateNFTAuctionResponse"}
          }
        }
      }
    },
    "/nft/mint": {
      "post": {
        "summary": "铸造 NFT",
        "tags": ["NFTInstance"],
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {"$ref": "#/definitions/MintNFTRequest"}
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {"$ref": "#/definitions/CreateNFTAuctionResponse"}
          }
        }
      }
    },
    "/nft/approve": {
      "post": {
        "summary": "授权 NFT 给拍卖合约",
        "tags": ["NFTInstance"],
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {"$ref": "#/definitions/ApproveNFTRequest"}
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {"$ref": "#/definitions/CreateNFTAuctionResponse"}
          }
        }
      }
    },
    "/nftauction": {
      "get": {
        "summary": "获取拍卖列表",
        "tags": ["NFTAuction"],
        "produces": ["application/json"],
        "responses": {
          "200": {
            "description": "success",
            "schema": {"$ref": "#/definitions/NFTAuctionListResponse"}
          }
        }
      },
      "post": {
        "summary": "创建拍卖",
        "tags": ["NFTAuction"],
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "auction",
            "in": "body",
            "required": true,
            "schema": {"$ref": "#/definitions/CreateNFTAuctionRequest"}
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {"$ref": "#/definitions/CreateNFTAuctionResponse"}
          },
          "400": {
            "description": "bad request",
            "schema": {"$ref": "#/definitions/CreateNFTAuctionResponse"}
          }
        }
      }
    },
    "/nftauction/{id}": {
      "get": {
        "summary": "获取拍卖详情",
        "tags": ["NFTAuction"],
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {"$ref": "#/definitions/NFTAuctionDetailResponse"}
          }
        }
      }
    }
  },
  "definitions": {
    "NFTAuction": {
      "type": "object",
      "properties": {
        "id": {"type": "string"},
        "name": {"type": "string"},
        "start_time": {"type": "string"},
        "end_time": {"type": "string"},
        "lowest_bid": {"type": "string"},
        "highest_bid": {"type": "string"},
        "description": {"type": "string"},
        "status": {"type": "string"},
        "image": {"type": "string"},
        "tags": {"type": "string"}
      }
    },
    "CreateNFTAuctionRequest": {
      "type": "object",
      "required": ["nft_contract", "token_id", "payment_token", "start_price_usd", "duration_hours"],
      "properties": {
        "nft_contract": {"type": "string"},
        "token_id": {"type": "string"},
        "payment_token": {"type": "string"},
        "start_price_usd": {"type": "string"},
        "duration_hours": {"type": "string"}
      }
    },
    "MintNFTRequest": {
      "type": "object",
      "required": ["uri"],
      "properties": {
        "uri": {"type": "string"}
      }
    },
    "ApproveNFTRequest": {
      "type": "object",
      "required": ["token_id"],
      "properties": {
        "token_id": {"type": "string"}
      }
    },
    "BidNFTAuctionRequest": {
      "type": "object",
      "required": ["auction_id", "bid_amount"],
      "properties": {
        "auction_id": {"type": "string"},
        "bid_amount": {"type": "string"}
      }
    },
    "EndNFTAuctionRequest": {
      "type": "object",
      "required": ["auction_id"],
      "properties": {
        "auction_id": {"type": "string"}
      }
    },
    "CreateNFTAuctionResult": {
      "type": "object",
      "properties": {
        "tx_hash": {"type": "string"}
      }
    },
    "NFTAuctionListResponse": {
      "type": "object",
      "properties": {
        "code": {"type": "integer"},
        "message": {"type": "string"},
        "data": {
          "type": "array",
          "items": {"$ref": "#/definitions/NFTAuction"}
        }
      }
    },
    "NFTAuctionDetailResponse": {
      "type": "object",
      "properties": {
        "code": {"type": "integer"},
        "message": {"type": "string"},
        "data": {"$ref": "#/definitions/NFTAuction"}
      }
    },
    "CreateNFTAuctionResponse": {
      "type": "object",
      "properties": {
        "code": {"type": "integer"},
        "message": {"type": "string"},
        "data": {"$ref": "#/definitions/CreateNFTAuctionResult"}
      }
    },
    "BidLog": {
      "type": "object",
      "properties": {
        "auction_id": {"type": "string"},
        "bidder": {"type": "string"},
        "amount": {"type": "string"},
        "amount_usd": {"type": "string"},
        "tx_hash": {"type": "string"},
        "block_number": {"type": "string"},
        "log_index": {"type": "string"}
      }
    },
    "BidLogsResponse": {
      "type": "object",
      "properties": {
        "code": {"type": "integer"},
        "message": {"type": "string"},
        "data": {
          "type": "array",
          "items": {"$ref": "#/definitions/BidLog"}
        }
      }
    },
    "AuctionLog": {
      "type": "object",
      "properties": {
        "event_type": {"type": "string"},
        "auction_id": {"type": "string"},
        "seller": {"type": "string"},
        "nft_contract": {"type": "string"},
        "token_id": {"type": "string"},
        "payment_token": {"type": "string"},
        "start_price_usd": {"type": "string"},
        "end_time": {"type": "string"},
        "bidder": {"type": "string"},
        "winner": {"type": "string"},
        "amount": {"type": "string"},
        "amount_usd": {"type": "string"},
        "tx_hash": {"type": "string"},
        "block_number": {"type": "string"},
        "log_index": {"type": "string"}
      }
    },
    "AuctionLogsResponse": {
      "type": "object",
      "properties": {
        "code": {"type": "integer"},
        "message": {"type": "string"},
        "data": {
          "type": "array",
          "items": {"$ref": "#/definitions/AuctionLog"}
        }
      }
    }
  }
}`

func GetSwaggerDoc(c *gin.Context) {
	c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(swaggerJSON))
}
