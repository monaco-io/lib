package main

import (
	"compress/gzip"
	"fmt"
	"log"
	"strings"

	"github.com/monaco-io/lib/xfile"
)

func main() {
	// 示例数据
	originalText := strings.Repeat("Hello, World! This is a test for compression. ", 100)
	originalData := []byte(originalText)

	fmt.Printf("Original data size: %d bytes\n", len(originalData))
	fmt.Printf("Original text preview: %.50s...\n", originalText)
	fmt.Println()

	// 测试不同压缩级别
	levels := []struct {
		name  string
		level int
	}{
		{"No Compression", gzip.NoCompression},
		{"Best Speed", gzip.BestSpeed},
		{"Default", gzip.DefaultCompression},
		{"Best Compression", gzip.BestCompression},
	}

	for _, lvl := range levels {
		fmt.Printf("=== %s (Level %d) ===\n", lvl.name, lvl.level)

		// 压缩
		compressed, err := xfile.Encode(originalData, lvl.level)
		if err != nil {
			log.Printf("Compression failed: %v", err)
			continue
		}

		// 计算压缩比
		ratio := xfile.CompressRatio(originalData, compressed)

		fmt.Printf("Compressed size: %d bytes\n", len(compressed))
		fmt.Printf("Compression ratio: %.2f:1\n", ratio)
		fmt.Printf("Space saved: %.1f%%\n", (1-float64(len(compressed))/float64(len(originalData)))*100)

		// 验证是否为有效的 gzip 数据
		if xfile.IsGzipData(compressed) {
			fmt.Printf("✓ Valid gzip format\n")
		} else {
			fmt.Printf("✗ Invalid gzip format\n")
		}

		// 解压验证
		decompressed, err := xfile.Decode(compressed)
		if err != nil {
			log.Printf("Decompression failed: %v", err)
			continue
		}

		if string(decompressed) == originalText {
			fmt.Printf("✓ Round-trip successful\n")
		} else {
			fmt.Printf("✗ Round-trip failed\n")
		}

		fmt.Println()
	}

	// 测试字符串 API
	fmt.Println("=== String API Test ===")
	testString := "你好，世界！🌍 Hello, World!"
	fmt.Printf("Original string: %s\n", testString)

	compressedBytes, err := xfile.GzipEncodeString(testString)
	if err != nil {
		log.Printf("String compression failed: %v", err)
		return
	}

	decompressedString, err := xfile.GzipDecodeString(compressedBytes)
	if err != nil {
		log.Printf("String decompression failed: %v", err)
		return
	}

	fmt.Printf("Decompressed string: %s\n", decompressedString)
	fmt.Printf("String round-trip: %t\n", testString == decompressedString)
	fmt.Printf("Compressed size: %d bytes\n", len(compressedBytes))
}
