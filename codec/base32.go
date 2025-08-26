package codec

import "encoding/base32"

func Base32Encode(data []byte) string {
	return base32.StdEncoding.EncodeToString(data)
}

func Base32Decode(encoded string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(encoded)
}
