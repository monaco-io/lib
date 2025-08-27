package codec

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// GZipEncode 压缩数据到 gzip 格式
// 支持可选的压缩级别配置
func GZipEncode(input []byte, level ...int) ([]byte, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}

	var buf bytes.Buffer

	// 设置压缩级别，默认为 DefaultCompression
	compressionLevel := gzip.DefaultCompression
	if len(level) > 0 {
		compressionLevel = level[0]
		// 验证压缩级别的有效性
		if compressionLevel < gzip.HuffmanOnly || compressionLevel > gzip.BestCompression {
			return nil, fmt.Errorf("invalid compression level: %d, must be between %d and %d",
				compressionLevel, gzip.HuffmanOnly, gzip.BestCompression)
		}
	}

	w, err := gzip.NewWriterLevel(&buf, compressionLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer: %w", err)
	}

	// 使用 defer 确保 writer 被正确关闭
	defer func() {
		if closeErr := w.Close(); closeErr != nil {
			// 如果关闭时出错，记录错误但不覆盖原始错误
			_ = closeErr
		}
	}()

	if _, err := w.Write(input); err != nil {
		return nil, fmt.Errorf("failed to write data: %w", err)
	}

	// 显式关闭 writer 以确保所有数据被刷新
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}

// GzipDecode 解压 gzip 格式的数据
func GzipDecode(input []byte) ([]byte, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}

	if len(input) == 0 {
		return []byte{}, nil
	}

	r := bytes.NewReader(input)
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}

	defer func() {
		if closeErr := gr.Close(); closeErr != nil {
			// 忽略关闭错误，因为我们已经读取了数据
			_ = closeErr
		}
	}()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, gr); err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}

	return buf.Bytes(), nil
}

// GzipEncodeString 压缩字符串到 gzip 格式
func GzipEncodeString(input string, level ...int) ([]byte, error) {
	return GZipEncode([]byte(input), level...)
}

// GzipDecodeString 解压 gzip 数据到字符串
func GzipDecodeString(input []byte) (string, error) {
	data, err := GzipDecode(input)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// IsGzipData 检查数据是否为有效的 gzip 格式
func IsGzipData(data []byte) bool {
	if len(data) < 3 {
		return false
	}
	// gzip 魔数检查：0x1f, 0x8b
	return data[0] == 0x1f && data[1] == 0x8b
}

// CompressRatio 计算压缩比率（原始大小 / 压缩后大小）
func CompressRatio(original, compressed []byte) float64 {
	if len(compressed) == 0 {
		return 0
	}
	return float64(len(original)) / float64(len(compressed))
}
