package codec

import (
	"testing"
)

func TestMD5(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"hello", "5d41402abc4b2a76b9719d911017c592"},
		{"Hello, World!", "65a8e27d8879283831b664bd8b7f0ad4"},
		{"The quick brown fox jumps over the lazy dog", "9e107d9d372bb6826bd81d3542a419d6"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := MD5String(tc.input)
			if result != tc.expected {
				t.Errorf("MD5(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSHA1(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"hello", "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{"Hello, World!", "0a0a9f2a6772942557ab5355d76af442f8f65e01"},
		{"The quick brown fox jumps over the lazy dog", "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := SHA1String(tc.input)
			if result != tc.expected {
				t.Errorf("SHA1(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSHA224(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "d14a028c2a3a2bc9476102bb288234c415a2b01f828ea62ac5b3e42f"},
		{"hello", "ea09ae9cc6768c50fcee903ed054556e5bfc8347907f12598aa24193"},
		{"Hello, World!", "72a23dfa411ba6fde01dbfabf3b00a709c93ebf273dc29e2d8b261ff"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := SHA224String(tc.input)
			if result != tc.expected {
				t.Errorf("SHA224(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSHA256(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"Hello, World!", "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"},
		{"The quick brown fox jumps over the lazy dog", "d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := SHA256String(tc.input)
			if result != tc.expected {
				t.Errorf("SHA256(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSHA384(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da274edebfe76f65fbd51ad2f14898b95b"},
		{"hello", "59e1748777448c69de6b800d7a33bbfb9ff1b463e44354c3553bcdb9c666fa90125a3c79f90397bdf5f6a13de828684f"},
		{"Hello, World!", "5485cc9b3365b4305dfb4e8337e0a598a574f8242bf17289e0dd6c20a3cd44a089de16ab4ab308f63e44b1170eb5f515"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := SHA384String(tc.input)
			if result != tc.expected {
				t.Errorf("SHA384(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSHA512(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
		{"hello", "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"},
		{"Hello, World!", "374d794a95cdcfd8b35993185fef9ba368f160d8daf432d08ba9f1ed1e5abe6cc69291e0fa2fe0006a52570ef18c19def4e617c33ce52ef0a6e5fbe318cb0387"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := SHA512String(tc.input)
			if result != tc.expected {
				t.Errorf("SHA512(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSHA512_224(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "6ed0dd02806fa89e25de060c19d3ac86cabb87d6a0ddd05c333b84f4"},
		{"hello", "fe8509ed1fb7dcefc27e6ac1a80eddbec4cb3d2c6fe565244374061c"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := SHA512_224String(tc.input)
			if result != tc.expected {
				t.Errorf("SHA512_224(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSHA512_256(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "c672b8d1ef56ed28ab87c3622c5114069bdd3ad7b8f9737498d0c01ecef0967a"},
		{"hello", "e30d87cfa2a75db545eac4d61baf970366a8357c7f72fa95b52d0accb698f13a"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := SHA512_256String(tc.input)
			if result != tc.expected {
				t.Errorf("SHA512_256(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

// 测试字节数组和字符串版本的一致性
func TestHashConsistency(t *testing.T) {
	testData := "Test data for consistency check"
	testBytes := []byte(testData)

	// 测试MD5
	if MD5(testBytes) != MD5String(testData) {
		t.Error("MD5 byte and string versions are inconsistent")
	}

	// 测试SHA1
	if SHA1(testBytes) != SHA1String(testData) {
		t.Error("SHA1 byte and string versions are inconsistent")
	}

	// 测试SHA224
	if SHA224(testBytes) != SHA224String(testData) {
		t.Error("SHA224 byte and string versions are inconsistent")
	}

	// 测试SHA256
	if SHA256(testBytes) != SHA256String(testData) {
		t.Error("SHA256 byte and string versions are inconsistent")
	}

	// 测试SHA384
	if SHA384(testBytes) != SHA384String(testData) {
		t.Error("SHA384 byte and string versions are inconsistent")
	}

	// 测试SHA512
	if SHA512(testBytes) != SHA512String(testData) {
		t.Error("SHA512 byte and string versions are inconsistent")
	}

	// 测试SHA512_224
	if SHA512_224(testBytes) != SHA512_224String(testData) {
		t.Error("SHA512_224 byte and string versions are inconsistent")
	}

	// 测试SHA512_256
	if SHA512_256(testBytes) != SHA512_256String(testData) {
		t.Error("SHA512_256 byte and string versions are inconsistent")
	}
}

// 性能基准测试

func BenchmarkMD5(b *testing.B) {
	data := []byte("Hello, World! This is a test message for benchmarking hash functions.")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MD5(data)
	}
}

func BenchmarkSHA1(b *testing.B) {
	data := []byte("Hello, World! This is a test message for benchmarking hash functions.")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SHA1(data)
	}
}

func BenchmarkSHA224(b *testing.B) {
	data := []byte("Hello, World! This is a test message for benchmarking hash functions.")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SHA224(data)
	}
}

func BenchmarkSHA256(b *testing.B) {
	data := []byte("Hello, World! This is a test message for benchmarking hash functions.")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SHA256(data)
	}
}

func BenchmarkSHA384(b *testing.B) {
	data := []byte("Hello, World! This is a test message for benchmarking hash functions.")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SHA384(data)
	}
}

func BenchmarkSHA512(b *testing.B) {
	data := []byte("Hello, World! This is a test message for benchmarking hash functions.")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SHA512(data)
	}
}

func BenchmarkSHA512_224(b *testing.B) {
	data := []byte("Hello, World! This is a test message for benchmarking hash functions.")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SHA512_224(data)
	}
}

func BenchmarkSHA512_256(b *testing.B) {
	data := []byte("Hello, World! This is a test message for benchmarking hash functions.")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SHA512_256(data)
	}
}

// 测试不同数据大小的性能
func BenchmarkHashFunctions_1KB(b *testing.B) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	b.Run("MD5", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			MD5(data)
		}
	})

	b.Run("SHA1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SHA1(data)
		}
	})

	b.Run("SHA256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SHA256(data)
		}
	})

	b.Run("SHA512", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SHA512(data)
		}
	})
}
