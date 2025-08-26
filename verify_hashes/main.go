package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

func main() {
	testCases := []string{"", "hello", "Hello, World!"}

	for _, testCase := range testCases {
		data := []byte(testCase)

		fmt.Printf("Input: %q\n", testCase)

		// SHA1
		sha1Hash := sha1.Sum(data)
		fmt.Printf("SHA1: %s\n", hex.EncodeToString(sha1Hash[:]))

		// SHA224
		sha224Hash := sha256.Sum224(data)
		fmt.Printf("SHA224: %s\n", hex.EncodeToString(sha224Hash[:]))

		// SHA256
		sha256Hash := sha256.Sum256(data)
		fmt.Printf("SHA256: %s\n", hex.EncodeToString(sha256Hash[:]))

		// SHA384
		sha384Hash := sha512.Sum384(data)
		fmt.Printf("SHA384: %s\n", hex.EncodeToString(sha384Hash[:]))

		// SHA512
		sha512Hash := sha512.Sum512(data)
		fmt.Printf("SHA512: %s\n", hex.EncodeToString(sha512Hash[:]))

		// SHA512_224
		sha512_224Hash := sha512.Sum512_224(data)
		fmt.Printf("SHA512_224: %s\n", hex.EncodeToString(sha512_224Hash[:]))

		// SHA512_256
		sha512_256Hash := sha512.Sum512_256(data)
		fmt.Printf("SHA512_256: %s\n", hex.EncodeToString(sha512_256Hash[:]))

		fmt.Println("---")
	}
}
