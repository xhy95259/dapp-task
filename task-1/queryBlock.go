package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 定义以太坊节点的 URL
	url := ""

	// 连接到以太坊客户端
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal("连接以太坊客户端失败:", err)
	}

	// 指定要查询的区块号
	blockNumber := big.NewInt(9135366)

	// 获取区块头信息
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal("获取区块头失败:", err)
	}

	// 打印区块头信息
	fmt.Println("区块头编号:", header.Number.Uint64())     // 区块号
	fmt.Println("区块头时间戳:", header.Time)               // 区块生成时间
	fmt.Println("区块头难度:", header.Difficulty.Uint64()) // 挖矿难度（通常为 0，因为使用 PoS）
	fmt.Println("区块头哈希:", header.Hash().Hex())        // 区块头的哈希值

	// 获取完整区块信息
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal("获取区块失败:", err)
	}

	// 打印区块信息
	fmt.Println("区块编号:", block.Number().Uint64())     // 区块号
	fmt.Println("区块时间戳:", block.Time())               // 区块生成时间
	fmt.Println("区块难度:", block.Difficulty().Uint64()) // 挖矿难度
	fmt.Println("区块哈希:", block.Hash().Hex())          // 区块哈希值
	fmt.Println("交易数量:", len(block.Transactions()))   // 区块中包含的交易数

	// 获取区块中的交易数量
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal("获取交易数量失败:", err)
	}

	// 打印交易数量
	fmt.Println("交易总数:", count) // 交易数量
}
