package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"practical-task/task-2/counter"
)

func main() {
	// 连接到以太坊Sepolia测试网络
	fmt.Println("正在连接到以太坊Sepolia测试网络...")
	url := "https://sepolia.infura.io/v3/4e00451dd920412090191a4315760504"
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal("❌ 连接网络失败:", err)
	}
	fmt.Println("✅ 网络连接成功")

	// 从私钥获取ECDSA私钥对象
	fmt.Println("正在解析私钥...")
	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		log.Fatal("❌ 解析私钥失败:", err)
	}
	fmt.Println("✅ 私钥解析成功")

	// 从私钥获取公钥
	fmt.Println("正在从私钥提取公钥...")
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("❌ 公钥类型断言失败: publicKey不是*ecdsa.PublicKey类型")
	}
	fmt.Println("✅ 公钥提取成功")

	// 从公钥获取发送方地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("�� 部署者地址: %s\n", fromAddress.Hex())

	// 检查账户余额
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal("❌ 获取账户余额失败:", err)
	}
	fmt.Printf("💰 账户余额: %s ETH\n", weiToEther(balance).String())

	// 获取发送方地址的nonce值（交易序号）
	fmt.Printf("正在获取地址 %s 的nonce值...\n", fromAddress.Hex())
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("❌ 获取nonce失败:", err)
	}
	fmt.Printf("✅ nonce值获取成功: %d\n", nonce)

	// 获取建议的Gas价格
	fmt.Println("正在获取建议的Gas价格...")
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("❌ 获取Gas价格失败:", err)
	}
	fmt.Printf("✅ Gas价格获取成功: %s wei\n", gasPrice.String())

	// 获取网络链ID
	fmt.Println("正在获取网络链ID...")
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("❌ 获取链ID失败:", err)
	}
	fmt.Printf("✅ 链ID获取成功: %s\n", chainId.String())

	// 创建交易授权对象
	fmt.Println("正在创建交易授权对象...")
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal("❌ 创建交易授权对象失败:", err)
	}

	// 设置交易参数
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(500000) // 增加Gas限制
	auth.GasPrice = gasPrice
	fmt.Println("✅ 交易授权对象创建成功")
	fmt.Printf("🔧 交易参数设置:\n")
	fmt.Printf("   Nonce: %d\n", auth.Nonce)
	fmt.Printf("   Value: %s wei\n", auth.Value.String())
	fmt.Printf("   GasLimit: %d\n", auth.GasLimit)
	fmt.Printf("   GasPrice: %s wei\n", auth.GasPrice.String())

	// 部署Counter合约
	fmt.Println("正在部署Counter智能合约...")
	address, tx, instance, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatal("❌ 合约部署失败:", err)
	}
	fmt.Println("✅ 合约部署交易已提交")

	// 打印合约部署信息
	fmt.Println("========== 合约部署信息 ==========")
	fmt.Printf("📄 合约地址: %s\n", address.Hex())
	fmt.Printf("🔗 交易哈希: %s\n", tx.Hash().Hex())
	fmt.Println("=================================")

	// 等待交易确认
	fmt.Println("⏳ 等待交易确认...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatal("❌ 等待交易确认失败:", err)
	}

	if receipt.Status == 1 {
		fmt.Println("✅ 合约部署成功!")
		fmt.Printf("📦 区块号: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("⛽ Gas使用量: %d\n", receipt.GasUsed)
	} else {
		log.Fatal("❌ 合约部署失败: 交易被回滚")
	}

	// 测试合约功能
	fmt.Println("\n========== 测试合约功能 ==========")

	// 获取初始计数
	count, err := instance.GetCount(nil)
	if err != nil {
		log.Fatal("❌ 获取计数失败:", err)
	}
	fmt.Printf("�� 初始计数: %s\n", count.String())

	// 增加计数
	fmt.Println("正在增加计数...")
	auth.Nonce = big.NewInt(int64(nonce + 1))
	auth.GasLimit = uint64(100000)

	tx, err = instance.Increment(auth)
	if err != nil {
		log.Fatal("❌ 增加计数失败:", err)
	}
	fmt.Printf("🔗 增加计数交易哈希: %s\n", tx.Hash().Hex())

	// 等待交易确认
	receipt, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatal("❌ 等待增加计数交易确认失败:", err)
	}

	if receipt.Status == 1 {
		fmt.Println("✅ 计数增加成功!")

		// 获取新的计数
		count, err = instance.GetCount(nil)
		if err != nil {
			log.Fatal("❌ 获取新计数失败:", err)
		}
		fmt.Printf("📊 新计数: %s\n", count.String())
	} else {
		log.Fatal("❌ 增加计数失败: 交易被回滚")
	}

	fmt.Println("🎉 所有操作完成!")
}

// 将Wei转换为Ether的辅助函数
func weiToEther(wei *big.Int) *big.Float {
	ether := new(big.Float)
	ether.SetString(wei.String())
	etherFloat := new(big.Float).Quo(ether, big.NewFloat(1e18))
	return etherFloat
}
