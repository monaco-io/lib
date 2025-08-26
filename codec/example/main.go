package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/monaco-io/lib/codec"
)

func main() {
	// 示例1: 使用AESCipher对象进行不同模式的加密
	fmt.Println("=== AESCipher 对象示例 ===")
	key := "mysecretkey12345" // 16字节密钥
	plaintext := "Hello, World! This is a secret message."

	// GCM模式
	fmt.Println("--- GCM模式 ---")
	gcmCipher, err := codec.NewAESCipher(key, codec.GCM)
	if err != nil {
		log.Fatal(err)
	}

	gcmEncrypted, err := gcmCipher.Encrypt([]byte(plaintext))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("原文: %s\n", plaintext)
	fmt.Printf("GCM加密长度: %d bytes\n", len(gcmEncrypted))
	fmt.Printf("GCM加密(hex): %s\n", hex.EncodeToString(gcmEncrypted))

	gcmDecrypted, err := gcmCipher.Decrypt(gcmEncrypted)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GCM解密: %s\n\n", string(gcmDecrypted))

	// CBC模式
	fmt.Println("--- CBC模式 ---")
	cbcCipher, err := codec.NewAESCipher(key, codec.CBC)
	if err != nil {
		log.Fatal(err)
	}

	cbcEncrypted, err := cbcCipher.Encrypt([]byte(plaintext))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("CBC加密长度: %d bytes\n", len(cbcEncrypted))
	fmt.Printf("CBC加密(hex): %s\n", hex.EncodeToString(cbcEncrypted))

	cbcDecrypted, err := cbcCipher.Decrypt(cbcEncrypted)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("CBC解密: %s\n\n", string(cbcDecrypted))

	// CTR模式
	fmt.Println("--- CTR模式 ---")
	ctrCipher, err := codec.NewAESCipher(key, codec.CTR)
	if err != nil {
		log.Fatal(err)
	}

	ctrEncrypted, err := ctrCipher.Encrypt([]byte(plaintext))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("CTR加密长度: %d bytes\n", len(ctrEncrypted))
	fmt.Printf("CTR加密(hex): %s\n", hex.EncodeToString(ctrEncrypted))

	ctrDecrypted, err := ctrCipher.Decrypt(ctrEncrypted)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("CTR解密: %s\n\n", string(ctrDecrypted))

	// 示例2: 验证不同密钥长度
	fmt.Println("=== 不同密钥长度示例 ===")
	testPlaintext := "Test message"

	// 16字节密钥 (AES-128)
	key128 := "1234567890123456"
	cipher128, err := codec.NewAESCipher(key128, codec.GCM)
	if err != nil {
		log.Fatal(err)
	}
	encrypted128, _ := cipher128.Encrypt([]byte(testPlaintext))
	decrypted128, _ := cipher128.Decrypt(encrypted128)
	fmt.Printf("AES-128: %s -> %s\n", testPlaintext, string(decrypted128))

	// 24字节密钥 (AES-192)
	key192 := "123456789012345678901234"
	cipher192, err := codec.NewAESCipher(key192, codec.GCM)
	if err != nil {
		log.Fatal(err)
	}
	encrypted192, _ := cipher192.Encrypt([]byte(testPlaintext))
	decrypted192, _ := cipher192.Decrypt(encrypted192)
	fmt.Printf("AES-192: %s -> %s\n", testPlaintext, string(decrypted192))

	// 32字节密钥 (AES-256)
	key256 := "12345678901234567890123456789012"
	cipher256, err := codec.NewAESCipher(key256, codec.GCM)
	if err != nil {
		log.Fatal(err)
	}
	encrypted256, _ := cipher256.Encrypt([]byte(testPlaintext))
	decrypted256, _ := cipher256.Decrypt(encrypted256)
	fmt.Printf("AES-256: %s -> %s\n", testPlaintext, string(decrypted256))
}
