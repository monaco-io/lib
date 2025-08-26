package codec

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// MD5 计算数据的MD5哈希值
func MD5(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

// MD5String 计算字符串的MD5哈希值
func MD5String(s string) string {
	return MD5([]byte(s))
}

// SHA1 计算数据的SHA1哈希值
func SHA1(data []byte) string {
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}

// SHA1String 计算字符串的SHA1哈希值
func SHA1String(s string) string {
	return SHA1([]byte(s))
}

// SHA256 计算数据的SHA256哈希值
func SHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// SHA256String 计算字符串的SHA256哈希值
func SHA256String(s string) string {
	return SHA256([]byte(s))
}

// SHA384 计算数据的SHA384哈希值
func SHA384(data []byte) string {
	hash := sha512.Sum384(data)
	return hex.EncodeToString(hash[:])
}

// SHA384String 计算字符串的SHA384哈希值
func SHA384String(s string) string {
	return SHA384([]byte(s))
}

// SHA512 计算数据的SHA512哈希值
func SHA512(data []byte) string {
	hash := sha512.Sum512(data)
	return hex.EncodeToString(hash[:])
}

// SHA512String 计算字符串的SHA512哈希值
func SHA512String(s string) string {
	return SHA512([]byte(s))
}

// SHA224 计算数据的SHA224哈希值
func SHA224(data []byte) string {
	hash := sha256.Sum224(data)
	return hex.EncodeToString(hash[:])
}

// SHA224String 计算字符串的SHA224哈希值
func SHA224String(s string) string {
	return SHA224([]byte(s))
}

// SHA512_224 计算数据的SHA512/224哈希值
func SHA512_224(data []byte) string {
	hash := sha512.Sum512_224(data)
	return hex.EncodeToString(hash[:])
}

// SHA512_224String 计算字符串的SHA512/224哈希值
func SHA512_224String(s string) string {
	return SHA512_224([]byte(s))
}

// SHA512_256 计算数据的SHA512/256哈希值
func SHA512_256(data []byte) string {
	hash := sha512.Sum512_256(data)
	return hex.EncodeToString(hash[:])
}

// SHA512_256String 计算字符串的SHA512/256哈希值
func SHA512_256String(s string) string {
	return SHA512_256([]byte(s))
}
