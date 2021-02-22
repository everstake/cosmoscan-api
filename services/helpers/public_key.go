package helpers

import (
	"encoding/base64"
	"fmt"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func GetHexAddressFromBase64PK(key string) (address string, err error) {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return address, fmt.Errorf("base64.DecodeString: %s", err.Error())
	}
	if len(decodedKey) != 32 {
		return address, fmt.Errorf("wrong key format")
	}
	pub := ed25519.PubKeyEd25519{}
	copy(pub[:], decodedKey)
	return pub.Address().String(), nil
}
