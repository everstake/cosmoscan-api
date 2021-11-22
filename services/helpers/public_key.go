package helpers

import (
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types"
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
	pub := ed25519.PubKey(decodedKey)
	return pub.Address().String(), nil
}

func GetBech32FromBase64PK(pkB64 string, pkType string) (address string, err error) {
	decodedKey, err := base64.StdEncoding.DecodeString(pkB64)
	if err != nil {
		return address, fmt.Errorf("base64.DecodeString: %s", err.Error())
	}
	var hexAddress string
	switch pkType {
	case "/cosmos.crypto.secp256k1.PubKey":
		pk := secp256k1.PubKey{Key: decodedKey}
		hexAddress = pk.Address().String()
	default:
		return address, fmt.Errorf("%s - unknown PK type", pkType)
	}
	addr, err := types.AccAddressFromHex(hexAddress)
	if err != nil {
		return address, fmt.Errorf("types.AccAddressFromHex: %s", err.Error())
	}
	return addr.String(), nil
}
