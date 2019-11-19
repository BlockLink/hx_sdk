package common

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
)

// 由私钥计算公钥
func Privkey2Pubkey(key []byte, compressed bool) []byte {
	_, pub := btcec.PrivKeyFromBytes(btcec.S256(), key)
	if compressed {
		return pub.SerializeCompressed()
	}

	return pub.SerializeUncompressed()
}

func PrivkeyFromBytes(key []byte) *btcec.PrivateKey {
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), key)
	return priv
}

func PrivkeyFromString(skey string) (*btcec.PrivateKey, error) {
	key, err := hex.DecodeString(skey)
	if err != nil {
		return nil, err
	}

	return PrivkeyFromBytes(key), nil
}
