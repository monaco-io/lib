package codec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

type CodecAES interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

// aesCipher AES加密器
type aesCipher struct {
	key        []byte
	block      cipher.Block
	cipherType AESCipherType
}

type AESCipherType int

const (
	GCM AESCipherType = iota
	CBC
	CTR
)

/**
AES 中的 128、192、256 主要是根据密钥长度来区分的，这是三者最核心的区别。具体对应关系如下：

AES-128：使用 128 位（16 字节）的密钥
AES-192：使用 192 位（24 字节）的密钥
AES-256：使用 256 位（32 字节）的密钥

除了密钥长度不同，三者的加密轮数也不同（轮数越多，加密过程越复杂）：

AES-128：10 轮加密
AES-192：12 轮加密
AES-256：14 轮加密
**/
// NewAESCipher 创建AES加密器
func NewAESCipher(key string, cipherType AESCipherType) (*aesCipher, error) {
	keyBytes := []byte(key)

	// 检查密钥长度，必须是16、24或32字节
	keyLen := len(keyBytes)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, fmt.Errorf("invalid key length: %d, must be 16, 24, or 32 bytes", keyLen)
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	return &aesCipher{
		block:      block,
		key:        keyBytes,
		cipherType: cipherType,
	}, nil
}

// Encrypt 根据加密器类型进行加密
func (a *aesCipher) Encrypt(plaintext []byte) ([]byte, error) {
	switch a.cipherType {
	case GCM:
		return a.encryptGCM(plaintext)
	case CBC:
		return a.encryptCBC(plaintext)
	case CTR:
		return a.encryptCTR(plaintext)
	default:
		return nil, errors.New("unsupported cipher type")
	}
}

// Decrypt 根据加密器类型进行解密
func (a *aesCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	switch a.cipherType {
	case GCM:
		return a.decryptGCM(ciphertext)
	case CBC:
		return a.decryptCBC(ciphertext)
	case CTR:
		return a.decryptCTR(ciphertext)
	default:
		return nil, errors.New("unsupported cipher type")
	}
}

// encryptGCM 使用AES-GCM加密
func (a *aesCipher) encryptGCM(plaintext []byte) ([]byte, error) {
	aesGCM, err := cipher.NewGCM(a.block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decryptGCM 使用AES-GCM解密
func (a *aesCipher) decryptGCM(ciphertext []byte) ([]byte, error) {
	aesGCM, err := cipher.NewGCM(a.block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// encryptCBC 使用AES-CBC加密
func (a *aesCipher) encryptCBC(plaintext []byte) ([]byte, error) {
	// PKCS7填充
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(a.block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// decryptCBC 使用AES-CBC解密
func (a *aesCipher) decryptCBC(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(a.block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 移除PKCS7填充
	plaintext, err := pkcs7Unpad(ciphertext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// encryptCTR 使用AES-CTR加密
func (a *aesCipher) encryptCTR(plaintext []byte) ([]byte, error) {
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(a.block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// decryptCTR 使用AES-CTR解密
func (a *aesCipher) decryptCTR(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(a.block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

// pkcs7Pad PKCS7填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs7Unpad 移除PKCS7填充
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding")
	}

	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, errors.New("invalid padding")
	}

	for i := length - unpadding; i < length; i++ {
		if data[i] != byte(unpadding) {
			return nil, errors.New("invalid padding")
		}
	}

	return data[:(length - unpadding)], nil
}
