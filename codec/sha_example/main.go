package main

import (
	"fmt"
	"log"

	"github.com/monaco-io/lib/codec"
)

func main() {
	// 示例数据
	testData := "Hello, World! This is a test message for hash functions."

	fmt.Println("=== 哈希函数示例 ===")
	fmt.Printf("原始数据: %s\n\n", testData)

	// MD5 (不推荐用于安全目的)
	fmt.Printf("MD5:        %s\n", codec.MD5String(testData))

	// SHA-1 (不推荐用于安全目的)
	fmt.Printf("SHA-1:      %s\n", codec.SHA1String(testData))

	// SHA-2 系列 (推荐)
	fmt.Printf("SHA-224:    %s\n", codec.SHA224String(testData))
	fmt.Printf("SHA-256:    %s\n", codec.SHA256String(testData))
	fmt.Printf("SHA-384:    %s\n", codec.SHA384String(testData))
	fmt.Printf("SHA-512:    %s\n", codec.SHA512String(testData))

	// SHA-512 变种
	fmt.Printf("SHA-512/224: %s\n", codec.SHA512_224String(testData))
	fmt.Printf("SHA-512/256: %s\n", codec.SHA512_256String(testData))

	fmt.Println("\n=== 不同数据大小测试 ===")

	// 测试空字符串
	fmt.Printf("空字符串 SHA-256: %s\n", codec.SHA256String(""))

	// 测试单个字符
	fmt.Printf("单字符 'a' SHA-256: %s\n", codec.SHA256String("a"))

	// 测试重复字符
	fmt.Printf("重复字符 'aaaa' SHA-256: %s\n", codec.SHA256String("aaaa"))

	fmt.Println("\n=== 字节数组 vs 字符串 ===")

	// 测试字节数组和字符串版本
	dataBytes := []byte(testData)

	sha256FromString := codec.SHA256String(testData)
	sha256FromBytes := codec.SHA256(dataBytes)

	fmt.Printf("字符串版本: %s\n", sha256FromString)
	fmt.Printf("字节版本:   %s\n", sha256FromBytes)
	fmt.Printf("结果一致:   %t\n", sha256FromString == sha256FromBytes)

	fmt.Println("\n=== 常见用例示例 ===")

	// 密码哈希 (实际应用中应该使用专门的密码哈希函数如 bcrypt)
	password := "mySecretPassword123"
	fmt.Printf("密码哈希 (SHA-256): %s\n", codec.SHA256String(password))

	// 文件完整性校验
	fileContent := "This is the content of a file that needs integrity checking."
	checksum := codec.SHA256String(fileContent)
	fmt.Printf("文件校验和: %s\n", checksum)

	// 唯一标识符生成
	identifier := fmt.Sprintf("user_%d_session_%d", 12345, 67890)
	uniqueID := codec.SHA1String(identifier)[:16] // 取前16位作为短标识符
	fmt.Printf("唯一标识符: %s\n", uniqueID)

	fmt.Println("\n=== 哈希长度对比 ===")
	fmt.Printf("MD5:         %d 字符 (%d 位)\n", len(codec.MD5String(testData)), len(codec.MD5String(testData))*4)
	fmt.Printf("SHA-1:       %d 字符 (%d 位)\n", len(codec.SHA1String(testData)), len(codec.SHA1String(testData))*4)
	fmt.Printf("SHA-224:     %d 字符 (%d 位)\n", len(codec.SHA224String(testData)), len(codec.SHA224String(testData))*4)
	fmt.Printf("SHA-256:     %d 字符 (%d 位)\n", len(codec.SHA256String(testData)), len(codec.SHA256String(testData))*4)
	fmt.Printf("SHA-384:     %d 字符 (%d 位)\n", len(codec.SHA384String(testData)), len(codec.SHA384String(testData))*4)
	fmt.Printf("SHA-512:     %d 字符 (%d 位)\n", len(codec.SHA512String(testData)), len(codec.SHA512String(testData))*4)
	fmt.Printf("SHA-512/224: %d 字符 (%d 位)\n", len(codec.SHA512_224String(testData)), len(codec.SHA512_224String(testData))*4)
	fmt.Printf("SHA-512/256: %d 字符 (%d 位)\n", len(codec.SHA512_256String(testData)), len(codec.SHA512_256String(testData))*4)

	fmt.Println("\n=== 安全建议 ===")
	fmt.Println("• MD5 和 SHA-1 已被认为不安全，不应用于安全敏感的应用")
	fmt.Println("• SHA-256 是目前广泛推荐的哈希算法")
	fmt.Println("• SHA-512 适用于需要更高安全性的场景")
	fmt.Println("• 对于密码存储，请使用专门的密码哈希函数 (如 bcrypt, scrypt, argon2)")
	fmt.Println("• 对于数字签名，建议使用 SHA-256 或更强的算法")

	// 演示哈希碰撞检测（理论上）
	fmt.Println("\n=== 哈希特性演示 ===")

	// 微小变化导致完全不同的哈希
	original := "Hello, World!"
	modified := "Hello, world!" // 只改变了一个字母的大小写

	fmt.Printf("原始:   '%s' -> %s\n", original, codec.SHA256String(original))
	fmt.Printf("修改后: '%s' -> %s\n", modified, codec.SHA256String(modified))
	fmt.Printf("哈希相同: %t\n", codec.SHA256String(original) == codec.SHA256String(modified))

	// 错误处理示例
	defer func() {
		if r := recover(); r != nil {
			log.Printf("发生错误: %v", r)
		}
	}()

	// 所有函数都是安全的，不会出现panic
	fmt.Println("\n所有哈希函数执行完成，未发生错误！")
}
