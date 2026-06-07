package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("===========连接节点===========")
	ctx, cannel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cannel()
	rpcURL := "https://mainnet.infura.io/v3/9edc71dd25e5412b8f973b8651981df3"
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Println("连接节点失败:", err)
		return
	}
	defer client.Close()
	fmt.Println("===========连接成功===========")
	// 连接以太坊节点，打印链 ID 和最新区块高度。
	chainID, err := client.ChainID(ctx)
	if err != nil {
		fmt.Println("获取链 ID 失败:", err)
		return
	}
	fmt.Println("链 ID:", chainID)
	fmt.Println("===========区块详细信息===========")
	block, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		fmt.Println("获取最新区块失败:", err)
		return
	}
	fmt.Println("区块高度: ", block.Number.Uint64())
	fmt.Println("区块哈希: ", block.Hash())
	fmt.Println("=== Ethereum Node Info ===")
	fmt.Printf("RPC URL       : %s\n", rpcURL)
	fmt.Printf("Chain ID      : %s\n", chainID.String())
	fmt.Println("\n⚠️  注意: 'Latest' 区块是节点当前认为的最新区块，可能尚未被所有节点确认")
	fmt.Println("   不同RPC节点可能返回不同的 'latest' 区块，导致与浏览器不匹配")
	fmt.Println("   建议对比 'Safe' 或 'Finalized' 区块（已确认的区块）")
	fmt.Println()
	fmt.Printf("Latest Block  : %d\n", block.Number.Uint64())
	fmt.Printf("Block Hash    : %s\n", block.Hash().Hex())
	fmt.Printf("Block Time    : %s\n", time.Unix(int64(block.Time), 0).Format(time.RFC3339))
	fmt.Println("==========================")

	// 示例：也可以获取任意指定高度的区块头
	if block.Number.Uint64() > 0 {
		num := new(big.Int).Sub(block.Number, big.NewInt(1))
		prevHeader, err := client.HeaderByNumber(ctx, num)
		if err == nil {
			fmt.Printf("Prev Block    : %d (%s)\n", prevHeader.Number.Uint64(), prevHeader.Hash().Hex())
		}
	}
	// 查询 safe 区块（浏览器通常显示这个）
	safeHeader, safeHash, err := getBlockByTag(ctx, client, "safe")
	if err != nil {
		log.Printf("failed to get safe block: %v (this may not be supported by all nodes)", err)
	} else {
		fmt.Println("\n=== Safe Block (推荐对比) ===")
		fmt.Printf("Block Number  : %d\n", safeHeader.Number.Uint64())
		fmt.Printf("Block Hash    : %s (RPC提供的hash，与浏览器一致)\n", safeHash.Hex())
		fmt.Printf("Calculated    : %s (计算出的hash，可能不匹配)\n", safeHeader.Hash().Hex())
		fmt.Printf("Block Time    : %s\n", time.Unix(int64(safeHeader.Time), 0).Format(time.RFC3339))
		fmt.Printf("Confirmations : %d\n", block.Number.Uint64()-safeHeader.Number.Uint64())
		fmt.Println("=============================")
	}

	// 查询 finalized 区块
	finalizedHeader, finalizedHash, err := getBlockByTag(ctx, client, "finalized")
	if err != nil {
		log.Printf("failed to get finalized block: %v (this may not be supported by all nodes)", err)
	} else {
		fmt.Println("\n=== Finalized Block ===")
		fmt.Printf("Block Number  : %d\n", finalizedHeader.Number.Uint64())
		fmt.Printf("Block Hash    : %s (RPC提供的hash，与浏览器一致)\n", finalizedHash.Hex())
		fmt.Printf("Calculated    : %s (计算出的hash，可能不匹配)\n", finalizedHeader.Hash().Hex())
		fmt.Printf("Block Time    : %s\n", time.Unix(int64(finalizedHeader.Time), 0).Format(time.RFC3339))
		fmt.Printf("Confirmations : %d\n", block.Number.Uint64()-finalizedHeader.Number.Uint64())
		fmt.Println("========================")
	}
}

// getBlockByTag 查询指定标签的区块头（safe, finalized, latest 等）
// 返回 Header、RPC 提供的 Hash 和错误
// 注意：需要使用底层 RPC 调用，因为 ethclient 的高级 API 不直接支持这些标签
func getBlockByTag(ctx context.Context, client *ethclient.Client, tag string) (*types.Header, common.Hash, error) {
	// 获取底层 RPC 客户端
	rpcClient := client.Client()

	// 获取区块头数据（使用 false 只获取 header，不包含交易）
	var raw json.RawMessage
	err := rpcClient.CallContext(ctx, &raw, "eth_getBlockByNumber", tag, false)
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("RPC call failed: %w", err)
	}

	if len(raw) == 0 || string(raw) == "null" {
		return nil, common.Hash{}, fmt.Errorf("%s block not found", tag)
	}

	// 解析完整的区块头字段
	var blockData struct {
		Number      string         `json:"number"`
		Hash        common.Hash    `json:"hash"`
		ParentHash  common.Hash    `json:"parentHash"`
		UncleHash   common.Hash    `json:"sha3Uncles"`
		Coinbase    common.Address `json:"miner"`
		Root        common.Hash    `json:"stateRoot"`
		TxHash      common.Hash    `json:"transactionsRoot"`
		ReceiptHash common.Hash    `json:"receiptsRoot"`
		Bloom       hexutil.Bytes  `json:"logsBloom"`
		Difficulty  *hexutil.Big   `json:"difficulty"`
		GasLimit    hexutil.Uint64 `json:"gasLimit"`
		GasUsed     hexutil.Uint64 `json:"gasUsed"`
		Time        hexutil.Uint64 `json:"timestamp"`
		Extra       hexutil.Bytes  `json:"extraData"`
		MixDigest   common.Hash    `json:"mixHash"`
		Nonce       hexutil.Bytes  `json:"nonce"`
		BaseFee     *hexutil.Big   `json:"baseFeePerGas"`
	}
	if err := json.Unmarshal(raw, &blockData); err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to unmarshal block header: %w", err)
	}

	// 解析区块号
	num, ok := new(big.Int).SetString(blockData.Number[2:], 16)
	if !ok {
		return nil, common.Hash{}, fmt.Errorf("invalid block number: %s", blockData.Number)
	}

	// 构造完整的 Header
	header := &types.Header{
		ParentHash:  blockData.ParentHash,
		UncleHash:   blockData.UncleHash,
		Coinbase:    blockData.Coinbase,
		Root:        blockData.Root,
		TxHash:      blockData.TxHash,
		ReceiptHash: blockData.ReceiptHash,
		Bloom:       types.BytesToBloom(blockData.Bloom),
		Difficulty:  big.NewInt(0),
		Number:      num,
		GasLimit:    uint64(blockData.GasLimit),
		GasUsed:     uint64(blockData.GasUsed),
		Time:        uint64(blockData.Time),
		Extra:       blockData.Extra,
		MixDigest:   blockData.MixDigest,
		BaseFee:     nil,
	}
	// 设置 Difficulty
	if blockData.Difficulty != nil {
		header.Difficulty = blockData.Difficulty.ToInt()
	}

	// 设置 BaseFee（EIP-1559）
	if blockData.BaseFee != nil {
		header.BaseFee = blockData.BaseFee.ToInt()
	}

	// 设置 Nonce
	if len(blockData.Nonce) >= 8 {
		var nonceBytes [8]byte
		copy(nonceBytes[:], blockData.Nonce[:8])
		header.Nonce = types.BlockNonce(nonceBytes)
	}

	// 返回 Header 和 RPC 提供的 hash
	// 注意：手动构造的 Header 计算出的 hash 可能不准确，因为：
	// 1. RPC 返回的某些字段可能格式不完全匹配 go-ethereum 的内部格式
	// 2. Header 的内部缓存字段可能未正确初始化
	// 因此，我们应该直接使用 RPC 返回的 hash，它与浏览器显示的 hash 一致
	return header, blockData.Hash, nil
}

/*
safe, finalized, latest 区块的区别
latest区块是网络刚刚产出的“最新鲜”的区块。它位于链的顶端，但还没有经过充分的网络验证。有可能会被拒绝，从而导致该区块在链上消失（发生链重组）。仅用于实时监控或数据展示（如区块浏览器显示最新交易），绝对不可用于确认交易结果。
safe是已经被下一个区块所引用（即确认了一次）。这意味着它不再是孤立的，被网络抛弃的概率大幅降低。 对于绝大多数日常交易（如转账、NFT买卖）来说，1-2个确认（即safe状态）已经足够安全。
finalized是在权益证明（PoS）机制下，经过特定轮次投票后被永久锁定的区块。要回滚它，需要攻击者销毁巨额的质押资产（理论上可行，但经济上不可行）。 跨链桥资产转移、大额金融结算、法律层面的确权。追求最高安全级别的场景。
*/
