package codec

import (
	"bytes"
	"crypto/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestAESCipher_GCM(t *testing.T) {
	key := "mysecretkey12345" // 16 bytes key
	plaintext := "Hello, World! This is a test message for AES encryption."

	cipher, err := NewAESCipher(key, GCM)
	if err != nil {
		t.Fatalf("Failed to create AES cipher: %v", err)
	}

	// Test Encrypt and Decrypt with GCM mode
	encrypted, err := cipher.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := cipher.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Fatalf("Decrypted text doesn't match original. Got: %s, Want: %s", string(decrypted), plaintext)
	}
}

func TestAESCipher_CBC(t *testing.T) {
	key := "mysecretkey12345" // 16 bytes key
	plaintext := "Hello, World! This is a test message for AES encryption."

	cipher, err := NewAESCipher(key, CBC)
	if err != nil {
		t.Fatalf("Failed to create AES cipher: %v", err)
	}

	// Test Encrypt and Decrypt with CBC mode
	encrypted, err := cipher.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := cipher.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Fatalf("Decrypted text doesn't match original. Got: %s, Want: %s", string(decrypted), plaintext)
	}
}

func TestAESCipher_CTR(t *testing.T) {
	key := "mysecretkey12345" // 16 bytes key
	plaintext := "Hello, World! This is a test message for AES encryption."

	cipher, err := NewAESCipher(key, CTR)
	if err != nil {
		t.Fatalf("Failed to create AES cipher: %v", err)
	}

	// Test Encrypt and Decrypt with CTR mode
	encrypted, err := cipher.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := cipher.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Fatalf("Decrypted text doesn't match original. Got: %s, Want: %s", string(decrypted), plaintext)
	}
}

func TestInvalidKeySize(t *testing.T) {
	invalidKey := "tooshort"
	_, err := NewAESCipher(invalidKey, GCM)
	if err == nil {
		t.Fatal("Expected error for invalid key size, got nil")
	}
}

func TestUnsupportedCipherType(t *testing.T) {
	key := "mysecretkey12345" // 16 bytes key

	// Create cipher with valid type first
	cipher, err := NewAESCipher(key, GCM)
	if err != nil {
		t.Fatalf("Failed to create AES cipher: %v", err)
	}

	// Manually set invalid cipher type to test error handling
	cipher.cipherType = AESCipherType(999)

	_, err = cipher.Encrypt([]byte("test"))
	if err == nil {
		t.Fatal("Expected error for unsupported cipher type, got nil")
	}
}

// TestAESCipher_DifferentKeySizes 测试不同密钥长度
func TestAESCipher_DifferentKeySizes(t *testing.T) {
	testCases := []struct {
		name    string
		keySize int
		key     string
	}{
		{"AES-128", 16, "1234567890123456"},
		{"AES-192", 24, "123456789012345678901234"},
		{"AES-256", 32, "12345678901234567890123456789012"},
	}

	plaintext := "Test message for different key sizes"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, cipherType := range []AESCipherType{GCM, CBC, CTR} {
				cipher, err := NewAESCipher(tc.key, cipherType)
				if err != nil {
					t.Fatalf("Failed to create %s cipher with %s: %v", tc.name, getCipherTypeName(cipherType), err)
				}

				encrypted, err := cipher.Encrypt([]byte(plaintext))
				if err != nil {
					t.Fatalf("Failed to encrypt with %s %s: %v", tc.name, getCipherTypeName(cipherType), err)
				}

				decrypted, err := cipher.Decrypt(encrypted)
				if err != nil {
					t.Fatalf("Failed to decrypt with %s %s: %v", tc.name, getCipherTypeName(cipherType), err)
				}

				if string(decrypted) != plaintext {
					t.Fatalf("%s %s: decrypted text doesn't match original", tc.name, getCipherTypeName(cipherType))
				}
			}
		})
	}
}

// TestAESCipher_EmptyData 测试空数据和边界条件
func TestAESCipher_EmptyData(t *testing.T) {
	key := "mysecretkey12345"

	testCases := []struct {
		name       string
		data       []byte
		shouldFail bool
	}{
		{"Empty data", []byte{}, false},
		{"Single byte", []byte{0x01}, false},
		{"Block size data", make([]byte, 16), false},
		{"Large data", make([]byte, 1024*1024), false}, // 1MB
	}

	for _, cipherType := range []AESCipherType{GCM, CBC, CTR} {
		t.Run(getCipherTypeName(cipherType), func(t *testing.T) {
			cipher, err := NewAESCipher(key, cipherType)
			if err != nil {
				t.Fatalf("Failed to create cipher: %v", err)
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					// Fill test data with random bytes
					if len(tc.data) > 0 {
						rand.Read(tc.data)
					}

					encrypted, err := cipher.Encrypt(tc.data)
					if tc.shouldFail {
						if err == nil {
							t.Fatal("Expected encryption to fail, but it succeeded")
						}
						return
					}
					if err != nil {
						t.Fatalf("Failed to encrypt: %v", err)
					}

					decrypted, err := cipher.Decrypt(encrypted)
					if err != nil {
						t.Fatalf("Failed to decrypt: %v", err)
					}

					if !bytes.Equal(decrypted, tc.data) {
						t.Fatal("Decrypted data doesn't match original")
					}
				})
			}
		})
	}
}

// TestAESCipher_InvalidCiphertext 测试无效密文
func TestAESCipher_InvalidCiphertext(t *testing.T) {
	key := "mysecretkey12345"

	testCases := []struct {
		name       string
		cipherType AESCipherType
		ciphertext []byte
	}{
		{"GCM - too short", GCM, []byte{0x01, 0x02}},
		{"CBC - too short", CBC, []byte{0x01, 0x02}},
		{"CTR - too short", CTR, []byte{0x01, 0x02}},
		{"CBC - not block aligned", CBC, make([]byte, 17)}, // 17 bytes, not multiple of 16
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cipher, err := NewAESCipher(key, tc.cipherType)
			if err != nil {
				t.Fatalf("Failed to create cipher: %v", err)
			}

			_, err = cipher.Decrypt(tc.ciphertext)
			if err == nil {
				t.Fatal("Expected decryption to fail with invalid ciphertext, but it succeeded")
			}
		})
	}
}

// TestAESCipher_ConcurrentSafety 测试多线程安全性
func TestAESCipher_ConcurrentSafety(t *testing.T) {
	key := "mysecretkey12345"
	plaintext := "Concurrent test message for thread safety validation"

	numGoroutines := runtime.NumCPU() * 2
	numOperations := 100

	for _, cipherType := range []AESCipherType{GCM, CBC, CTR} {
		t.Run(getCipherTypeName(cipherType), func(t *testing.T) {
			cipher, err := NewAESCipher(key, cipherType)
			if err != nil {
				t.Fatalf("Failed to create cipher: %v", err)
			}

			var wg sync.WaitGroup
			errChan := make(chan error, numGoroutines*numOperations)

			// 启动多个goroutine并发执行加密解密
			for i := 0; i < numGoroutines; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					for j := 0; j < numOperations; j++ {
						// 每次使用稍微不同的数据
						testData := []byte(plaintext + string(rune(id)) + string(rune(j)))

						encrypted, err := cipher.Encrypt(testData)
						if err != nil {
							errChan <- err
							return
						}

						decrypted, err := cipher.Decrypt(encrypted)
						if err != nil {
							errChan <- err
							return
						}

						if !bytes.Equal(decrypted, testData) {
							errChan <- err
							return
						}
					}
				}(i)
			}

			wg.Wait()
			close(errChan)

			// 检查是否有错误
			for err := range errChan {
				t.Fatalf("Concurrent operation failed: %v", err)
			}

			t.Logf("Successfully completed %d concurrent operations with %d goroutines",
				numGoroutines*numOperations, numGoroutines)
		})
	}
}

// TestAESCipher_LongRunning 长时间运行测试
func TestAESCipher_LongRunning(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long running test in short mode")
	}

	key := "mysecretkey12345"
	plaintext := "Long running test message"

	cipher, err := NewAESCipher(key, GCM)
	if err != nil {
		t.Fatalf("Failed to create cipher: %v", err)
	}

	start := time.Now()
	iterations := 10000

	for i := 0; i < iterations; i++ {
		encrypted, err := cipher.Encrypt([]byte(plaintext))
		if err != nil {
			t.Fatalf("Encryption failed at iteration %d: %v", i, err)
		}

		decrypted, err := cipher.Decrypt(encrypted)
		if err != nil {
			t.Fatalf("Decryption failed at iteration %d: %v", i, err)
		}

		if string(decrypted) != plaintext {
			t.Fatalf("Data corruption at iteration %d", i)
		}
	}

	duration := time.Since(start)
	t.Logf("Completed %d encrypt/decrypt cycles in %v (%.2f ops/sec)",
		iterations, duration, float64(iterations)/duration.Seconds())
}

// TestAESCipher_MemoryUsage 内存使用测试
func TestAESCipher_MemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory usage test in short mode")
	}

	key := "mysecretkey12345"
	largeData := make([]byte, 10*1024*1024) // 10MB
	rand.Read(largeData)

	cipher, err := NewAESCipher(key, GCM)
	if err != nil {
		t.Fatalf("Failed to create cipher: %v", err)
	}

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	encrypted, err := cipher.Encrypt(largeData)
	if err != nil {
		t.Fatalf("Failed to encrypt large data: %v", err)
	}

	decrypted, err := cipher.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt large data: %v", err)
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)

	if !bytes.Equal(decrypted, largeData) {
		t.Fatal("Large data corruption")
	}

	allocDiff := m2.TotalAlloc - m1.TotalAlloc
	t.Logf("Memory allocated for 10MB encrypt/decrypt: %d bytes", allocDiff)
}

// 性能基准测试

// BenchmarkAESCipher_Encrypt_GCM GCM模式加密性能测试
func BenchmarkAESCipher_Encrypt_GCM(b *testing.B) {
	benchmarkEncrypt(b, GCM)
}

// BenchmarkAESCipher_Encrypt_CBC CBC模式加密性能测试
func BenchmarkAESCipher_Encrypt_CBC(b *testing.B) {
	benchmarkEncrypt(b, CBC)
}

// BenchmarkAESCipher_Encrypt_CTR CTR模式加密性能测试
func BenchmarkAESCipher_Encrypt_CTR(b *testing.B) {
	benchmarkEncrypt(b, CTR)
}

// BenchmarkAESCipher_Decrypt_GCM GCM模式解密性能测试
func BenchmarkAESCipher_Decrypt_GCM(b *testing.B) {
	benchmarkDecrypt(b, GCM)
}

// BenchmarkAESCipher_Decrypt_CBC CBC模式解密性能测试
func BenchmarkAESCipher_Decrypt_CBC(b *testing.B) {
	benchmarkDecrypt(b, CBC)
}

// BenchmarkAESCipher_Decrypt_CTR CTR模式解密性能测试
func BenchmarkAESCipher_Decrypt_CTR(b *testing.B) {
	benchmarkDecrypt(b, CTR)
}

// BenchmarkAESCipher_EncryptDecrypt_Various 不同数据大小的性能测试
func BenchmarkAESCipher_EncryptDecrypt_Various(b *testing.B) {
	sizes := []int{16, 64, 256, 1024, 4096, 16384, 65536} // 16B to 64KB

	for _, size := range sizes {
		b.Run(formatBytes(size), func(b *testing.B) {
			benchmarkEncryptDecryptSize(b, GCM, size)
		})
	}
}

// BenchmarkAESCipher_KeySizes 不同密钥长度的性能测试
func BenchmarkAESCipher_KeySizes(b *testing.B) {
	keys := map[string]string{
		"AES128": "1234567890123456",
		"AES192": "123456789012345678901234",
		"AES256": "12345678901234567890123456789012",
	}

	for name, key := range keys {
		b.Run(name, func(b *testing.B) {
			benchmarkWithKey(b, key, GCM)
		})
	}
}

// BenchmarkAESCipher_Parallel 并行性能测试
func BenchmarkAESCipher_Parallel(b *testing.B) {
	key := "mysecretkey12345"
	plaintext := []byte("Parallel benchmark test message for performance evaluation")

	cipher, err := NewAESCipher(key, GCM)
	if err != nil {
		b.Fatalf("Failed to create cipher: %v", err)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			encrypted, err := cipher.Encrypt(plaintext)
			if err != nil {
				b.Fatalf("Encryption failed: %v", err)
			}

			_, err = cipher.Decrypt(encrypted)
			if err != nil {
				b.Fatalf("Decryption failed: %v", err)
			}
		}
	})
}

// 辅助函数

func benchmarkEncrypt(b *testing.B, cipherType AESCipherType) {
	key := "mysecretkey12345"
	plaintext := []byte("Benchmark test message for encryption performance measurement")

	cipher, err := NewAESCipher(key, cipherType)
	if err != nil {
		b.Fatalf("Failed to create cipher: %v", err)
	}

	b.ResetTimer()
	b.SetBytes(int64(len(plaintext)))

	for i := 0; i < b.N; i++ {
		_, err := cipher.Encrypt(plaintext)
		if err != nil {
			b.Fatalf("Encryption failed: %v", err)
		}
	}
}

func benchmarkDecrypt(b *testing.B, cipherType AESCipherType) {
	key := "mysecretkey12345"
	plaintext := []byte("Benchmark test message for decryption performance measurement")

	cipher, err := NewAESCipher(key, cipherType)
	if err != nil {
		b.Fatalf("Failed to create cipher: %v", err)
	}

	// Pre-encrypt data for decryption benchmark
	encrypted, err := cipher.Encrypt(plaintext)
	if err != nil {
		b.Fatalf("Failed to encrypt: %v", err)
	}

	b.ResetTimer()
	b.SetBytes(int64(len(plaintext)))

	for i := 0; i < b.N; i++ {
		// Create a copy of encrypted data for each iteration
		encryptedCopy := make([]byte, len(encrypted))
		copy(encryptedCopy, encrypted)

		_, err := cipher.Decrypt(encryptedCopy)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}

func benchmarkEncryptDecryptSize(b *testing.B, cipherType AESCipherType, size int) {
	key := "mysecretkey12345"
	plaintext := make([]byte, size)
	rand.Read(plaintext)

	cipher, err := NewAESCipher(key, cipherType)
	if err != nil {
		b.Fatalf("Failed to create cipher: %v", err)
	}

	b.ResetTimer()
	b.SetBytes(int64(size))

	for i := 0; i < b.N; i++ {
		encrypted, err := cipher.Encrypt(plaintext)
		if err != nil {
			b.Fatalf("Encryption failed: %v", err)
		}

		_, err = cipher.Decrypt(encrypted)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}

func benchmarkWithKey(b *testing.B, key string, cipherType AESCipherType) {
	plaintext := []byte("Benchmark test message for key size performance comparison")

	cipher, err := NewAESCipher(key, cipherType)
	if err != nil {
		b.Fatalf("Failed to create cipher: %v", err)
	}

	b.ResetTimer()
	b.SetBytes(int64(len(plaintext)))

	for i := 0; i < b.N; i++ {
		encrypted, err := cipher.Encrypt(plaintext)
		if err != nil {
			b.Fatalf("Encryption failed: %v", err)
		}

		_, err = cipher.Decrypt(encrypted)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}

func getCipherTypeName(cipherType AESCipherType) string {
	switch cipherType {
	case GCM:
		return "GCM"
	case CBC:
		return "CBC"
	case CTR:
		return "CTR"
	default:
		return "Unknown"
	}
}

func formatBytes(bytes int) string {
	if bytes < 1024 {
		return string(rune(bytes)) + "B"
	}
	return string(rune(bytes/1024)) + "KB"
}
