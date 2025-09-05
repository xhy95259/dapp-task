package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接到以太坊Sepolia测试网络
	url := ""
	fmt.Println("正在连接以太坊Sepolia测试网络...")
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal("❌ 连接以太坊网络失败:", err)
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
	fmt.Printf("📬 发送方地址: %s\n", fromAddress.Hex())

	// 获取发送方地址的nonce值（交易序号）
	fmt.Printf("正在获取地址 %s 的nonce值...\n", fromAddress.Hex())
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("❌ 获取nonce失败:", err)
	}
	fmt.Printf("✅ nonce值获取成功: %d\n", nonce)

	// 设置转账金额（单位：wei，这里是0.001 ETH）
	value := big.NewInt(1e15) // in wei (0.001 eth)
	fmt.Printf("💸 转账金额: %s wei (0.001 ETH)\n", value.String())

	// 设置Gas限制
	gasLimit := uint64(21000) // in units
	fmt.Printf("⛽ Gas限制: %d\n", gasLimit)

	// 获取建议的Gas价格
	fmt.Println("正在获取建议的Gas价格...")
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("❌ 获取Gas价格失败:", err)
	}
	fmt.Printf("✅ Gas价格获取成功: %s wei\n", gasPrice.String())

	// 设置接收方地址
	toAddress := common.HexToAddress("0x56161e6389eD71C3D4a3C60a3a0a1C17D77Ef031")
	fmt.Printf("📧 接收方地址: %s\n", toAddress.Hex())

	// 设置交易数据（这里为空）
	var data []byte
	fmt.Printf("📄 交易数据: %x (空数据)\n", data)

	// 创建交易对象
	fmt.Println("正在创建交易对象...")
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	fmt.Println("✅ 交易对象创建成功")

	// 获取网络链ID
	fmt.Println("正在获取网络链ID...")
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("❌ 获取链ID失败:", err)
	}
	fmt.Printf("✅ 链ID获取成功: %s\n", chainID.String())

	// 使用EIP155签名规则对交易进行签名
	fmt.Println("正在对交易进行签名...")
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("❌ 交易签名失败:", err)
	}
	fmt.Println("✅ 交易签名成功")

	// 发送交易到网络
	fmt.Printf("正在发送交易 %s ...\n", signedTx.Hash().Hex())
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("❌ 发送交易失败:", err)
	}

	// 打印交易哈希
	fmt.Printf("🎉 交易已成功发送!\n")
	fmt.Printf("🔗 交易哈希: %s\n", signedTx.Hash().Hex())
	fmt.Printf("📋 交易详情:\n")
	fmt.Printf("   发送方: %s\n", fromAddress.Hex())
	fmt.Printf("   接收方: %s\n", toAddress.Hex())
	fmt.Printf("   金额: %s wei\n", value.String())
	fmt.Printf("   Gas限制: %d\n", gasLimit)
	fmt.Printf("   Gas价格: %s wei\n", gasPrice.String())
	fmt.Printf("   Nonce: %d\n", nonce)
}
