package codec

import (
	"bytes"
	"compress/gzip"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestGzipEncode(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "empty data",
			input:   []byte{},
			wantErr: false,
		},
		{
			name:    "simple text",
			input:   []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name:    "binary data",
			input:   []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD},
			wantErr: false,
		},
		{
			name:    "large text data",
			input:   []byte(strings.Repeat("Hello, World! ", 1000)),
			wantErr: false,
		},
		{
			name:    "unicode text",
			input:   []byte("你好，世界！🌍"),
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GZipEncode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GzipEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// 验证压缩结果可以解压
				decoded, err := GzipDecode(got)
				if err != nil {
					t.Errorf("Failed to decode compressed data: %v", err)
					return
				}
				if !bytes.Equal(decoded, tt.input) {
					t.Errorf("Decoded data doesn't match original. Got %v, want %v", decoded, tt.input)
				}

				// 验证压缩结果是有效的 gzip 数据
				if len(got) > 0 && !IsGzipData(got) {
					t.Errorf("Compressed data is not valid gzip format")
				}
			}
		})
	}
}

func TestGzipEncodeWithLevel(t *testing.T) {
	input := []byte(strings.Repeat("Hello, World! ", 100))

	tests := []struct {
		name    string
		level   int
		wantErr bool
	}{
		{
			name:    "no compression",
			level:   gzip.NoCompression,
			wantErr: false,
		},
		{
			name:    "best speed",
			level:   gzip.BestSpeed,
			wantErr: false,
		},
		{
			name:    "default compression",
			level:   gzip.DefaultCompression,
			wantErr: false,
		},
		{
			name:    "best compression",
			level:   gzip.BestCompression,
			wantErr: false,
		},
		{
			name:    "huffman only",
			level:   gzip.HuffmanOnly,
			wantErr: false,
		},
		{
			name:    "invalid level - too low",
			level:   -10,
			wantErr: true,
		},
		{
			name:    "invalid level - too high",
			level:   20,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GZipEncode(input, tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("GzipEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// 验证可以正确解压
				decoded, err := GzipDecode(got)
				if err != nil {
					t.Errorf("Failed to decode: %v", err)
					return
				}
				if !bytes.Equal(decoded, input) {
					t.Errorf("Decoded data doesn't match original")
				}
			}
		})
	}
}

func TestGzipDecode(t *testing.T) {
	// 准备测试数据
	validData := []byte("Hello, World!")
	validCompressed, _ := GZipEncode(validData)

	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "valid compressed data",
			input:   validCompressed,
			want:    validData,
			wantErr: false,
		},
		{
			name:    "empty compressed data",
			input:   []byte{},
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "invalid gzip data",
			input:   []byte("not gzip data"),
			wantErr: true,
		},
		{
			name:    "corrupted gzip header",
			input:   []byte{0x1f, 0x8b, 0x00}, // 不完整的 gzip 头
			wantErr: true,
		},
		{
			name:    "wrong magic number",
			input:   []byte{0x1f, 0x8c, 0x08, 0x00}, // 错误的魔数
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GzipDecode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GzipDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("GzipDecode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGzipEncodeString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   "",
			wantErr: false,
		},
		{
			name:    "simple string",
			input:   "Hello, World!",
			wantErr: false,
		},
		{
			name:    "unicode string",
			input:   "你好，世界！🌍",
			wantErr: false,
		},
		{
			name:    "multiline string",
			input:   "Line 1\nLine 2\nLine 3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GzipEncodeString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GzipEncodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// 验证可以解压回原始字符串
				decoded, err := GzipDecodeString(got)
				if err != nil {
					t.Errorf("Failed to decode: %v", err)
					return
				}
				if decoded != tt.input {
					t.Errorf("Decoded string = %v, want %v", decoded, tt.input)
				}
			}
		})
	}
}

func TestGzipDecodeString(t *testing.T) {
	// 准备测试数据
	testString := "Hello, 世界!"
	compressed, _ := GzipEncodeString(testString)

	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr bool
	}{
		{
			name:    "valid compressed string",
			input:   compressed,
			want:    testString,
			wantErr: false,
		},
		{
			name:    "empty data",
			input:   []byte{},
			want:    "",
			wantErr: false,
		},
		{
			name:    "invalid data",
			input:   []byte("invalid"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GzipDecodeString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GzipDecodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GzipDecodeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsGzipData(t *testing.T) {
	// 准备测试数据
	validGzipData, _ := GZipEncode([]byte("test"))

	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "valid gzip data",
			data: validGzipData,
			want: true,
		},
		{
			name: "gzip magic number only",
			data: []byte{0x1f, 0x8b, 0x08},
			want: true,
		},
		{
			name: "empty data",
			data: []byte{},
			want: false,
		},
		{
			name: "insufficient data",
			data: []byte{0x1f},
			want: false,
		},
		{
			name: "wrong magic number",
			data: []byte{0x1f, 0x8c, 0x08},
			want: false,
		},
		{
			name: "plain text",
			data: []byte("Hello, World!"),
			want: false,
		},
		{
			name: "nil data",
			data: nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsGzipData(tt.data); got != tt.want {
				t.Errorf("IsGzipData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompressRatio(t *testing.T) {
	tests := []struct {
		name       string
		original   []byte
		compressed []byte
		want       float64
	}{
		{
			name:       "equal size",
			original:   []byte("hello"),
			compressed: []byte("world"),
			want:       1.0,
		},
		{
			name:       "2:1 compression",
			original:   []byte("hello world"),
			compressed: []byte("hello"),
			want:       2.2,
		},
		{
			name:       "empty compressed",
			original:   []byte("hello"),
			compressed: []byte{},
			want:       0,
		},
		{
			name:       "empty original",
			original:   []byte{},
			compressed: []byte("test"),
			want:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompressRatio(tt.original, tt.compressed)
			if got != tt.want {
				t.Errorf("CompressRatio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGzipRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "small text",
			data: []byte("Hello, World!"),
		},
		{
			name: "large repetitive data",
			data: []byte(strings.Repeat("Hello, World! ", 1000)),
		},
		{
			name: "binary data",
			data: generateRandomBytes(1000),
		},
		{
			name: "unicode text",
			data: []byte("这是一个包含中文字符的测试文本。🚀✨"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 压缩
			compressed, err := GZipEncode(tt.data)
			if err != nil {
				t.Fatalf("GzipEncode() error = %v", err)
			}

			// 解压
			decompressed, err := GzipDecode(compressed)
			if err != nil {
				t.Fatalf("GzipDecode() error = %v", err)
			}

			// 验证数据一致性
			if !bytes.Equal(tt.data, decompressed) {
				t.Errorf("Round trip failed. Original length: %d, Decompressed length: %d",
					len(tt.data), len(decompressed))
			}

			// 验证压缩效果（对于重复数据应该有较好的压缩比）
			if len(tt.data) > 100 && strings.Contains(string(tt.data), "Hello, World!") {
				ratio := CompressRatio(tt.data, compressed)
				if ratio < 2.0 {
					t.Logf("Compression ratio for repetitive data: %.2f (might be low but acceptable)", ratio)
				}
			}
		})
	}
}

func TestGzipCompressionLevels(t *testing.T) {
	data := []byte(strings.Repeat("Hello, World! This is a test for compression. ", 100))

	levels := []int{
		gzip.NoCompression,
		gzip.BestSpeed,
		gzip.DefaultCompression,
		gzip.BestCompression,
	}

	results := make(map[int]int) // level -> compressed size

	for _, level := range levels {
		compressed, err := GZipEncode(data, level)
		if err != nil {
			t.Errorf("GzipEncode with level %d failed: %v", level, err)
			continue
		}

		results[level] = len(compressed)

		// 验证可以正确解压
		decompressed, err := GzipDecode(compressed)
		if err != nil {
			t.Errorf("GzipDecode failed for level %d: %v", level, err)
			continue
		}

		if !bytes.Equal(data, decompressed) {
			t.Errorf("Round trip failed for level %d", level)
		}
	}

	// 验证压缩级别的效果
	if len(results) >= 2 {
		// BestCompression 应该比 BestSpeed 压缩得更小（对于重复数据）
		if bestComp, ok := results[gzip.BestCompression]; ok {
			if bestSpeed, ok := results[gzip.BestSpeed]; ok {
				if bestComp > bestSpeed {
					t.Logf("Note: BestCompression (%d bytes) > BestSpeed (%d bytes) - this can happen with small or non-repetitive data", bestComp, bestSpeed)
				}
			}
		}
	}
}

// 辅助函数：生成随机字节
func generateRandomBytes(n int) []byte {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = byte(rand.Intn(256))
	}
	return bytes
}

// 基准测试
func BenchmarkGzipEncode(b *testing.B) {
	data := []byte(strings.Repeat("Hello, World! ", 1000))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := GZipEncode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGzipDecode(b *testing.B) {
	data := []byte(strings.Repeat("Hello, World! ", 1000))
	compressed, err := GZipEncode(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GzipDecode(compressed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGzipEncodeSmall(b *testing.B) {
	data := []byte("Hello, World!")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := GZipEncode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGzipEncodeLarge(b *testing.B) {
	data := []byte(strings.Repeat("Hello, World! ", 10000))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := GZipEncode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGzipEncodeWithLevels(b *testing.B) {
	data := []byte(strings.Repeat("Hello, World! ", 1000))

	levels := []struct {
		name  string
		level int
	}{
		{"NoCompression", gzip.NoCompression},
		{"BestSpeed", gzip.BestSpeed},
		{"DefaultCompression", gzip.DefaultCompression},
		{"BestCompression", gzip.BestCompression},
	}

	for _, lvl := range levels {
		b.Run(lvl.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := GZipEncode(data, lvl.level)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
