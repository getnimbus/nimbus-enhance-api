package encoder

import (
	"github.com/mr-tron/base58"
)

func IsBase58(str string) bool {
	decodedBytes, _ := base58.Decode(str)
	if decodedBytes == nil {
		return false
	}

	// re-encode the decoded bytes and compare with the original string
	reencodedStr := base58.Encode(decodedBytes)
	return reencodedStr == str
}
