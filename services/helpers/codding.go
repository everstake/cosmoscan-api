package helpers

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func B64ToHex(b64Str string) (hexStr string, err error) {
	bts, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		return hexStr, fmt.Errorf("base64.StdEncoding.DecodeString: %s", err.Error())
	}
	return hex.EncodeToString(bts), nil
}
