package common

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"
)

const (
	bip44prefix          uint32 = 44
	HardenedKeyZeroIndex uint32 = 0x80000000
)

//
func DerivePubkey(mnemonic, password string, coinIdx, actIdx, addrIdx int) (*btcec.PublicKey, error) {
	seed := MnemonicToSeed(mnemonic, password)
	netp := &chaincfg.MainNetParams
	masterKey, err := hdkeychain.NewMaster(seed, netp)
	if err != nil {
		return nil, fmt.Errorf("create master key failed: %v", err)
	}

	addrKey, _, err := BIP44AccountKey(masterKey, netp, uint32(coinIdx), uint32(actIdx), uint32(addrIdx))

	pubKey, err := addrKey.ECPubKey()
	if err != nil {
		return nil, err
	}
	return pubKey, nil

}

func DerivePubkeyBytes(mnemonic, password string, coinIdx, actIdx, addrIdx int, compress bool) ([]byte, error) {
	pubKey, err := DerivePubkey(mnemonic, password, coinIdx, actIdx, addrIdx)
	if err != nil {
		return nil, err
	}

	var key []byte
	if compress {
		key = pubKey.SerializeCompressed()
	} else {
		key = pubKey.SerializeUncompressed()
	}
	return key, nil
}

func DerivePrivateKey(mnemonic, password string, coinIdx, actIdx, addrIdx int) (*btcec.PrivateKey, error) {
	seed := MnemonicToSeed(mnemonic, password)
	netp := &chaincfg.MainNetParams
	masterKey, err := hdkeychain.NewMaster(seed, netp)
	if err != nil {
		return nil, fmt.Errorf("create master key failed: %v", err)
	}

	addrKey, _, err := BIP44AccountKey(masterKey, netp, uint32(coinIdx), uint32(actIdx), uint32(addrIdx))

	privKey, err := addrKey.ECPrivKey()
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func DerivePrivateKeyBytes(mnemonic, password string, coinIdx, actIdx, addrIdx int) ([]byte, error) {
	key, err := DerivePrivateKey(mnemonic, password, coinIdx, actIdx, addrIdx)
	if err != nil {
		return nil, err
	}
	return key.Serialize(), nil
}

func MnemonicToSeed(mnemonic, password string) []byte {
	return bip39.NewSeed(mnemonic, password)
}

// 由 master key 生成一个 HD extendedKey
// m / 44' / coinIdx' / accIdx' / 0 / addrIdx
// return: address, path, error
func BIP44AccountKey(key *hdkeychain.ExtendedKey, netp *chaincfg.Params,
	coinIdx, accIdx, addrIdx uint32) (*hdkeychain.ExtendedKey, string, error) {
	// 先导出account key

	purposeIndex := HardenedKeyZeroIndex + bip44prefix // bip44
	coinTypeIndex := HardenedKeyZeroIndex + coinIdx    // bip32 coin type
	accIdx = HardenedKeyZeroIndex + accIdx

	// m / 44'
	purposeK, err := key.Child(purposeIndex)
	if err != nil {
		return nil, "", fmt.Errorf("derive m/44' key failed: %v", err)
	}

	// m / 44' / coinIdx'
	cTypeK, err := purposeK.Child(coinTypeIndex)
	if err != nil {
		return nil, "", fmt.Errorf("derive m/44'/%d' key failed: %v", coinIdx, err)
	}

	// m / 44' / coinIdx' / accIdx
	accK, err := cTypeK.Child(accIdx)
	if err != nil {
		return nil, "", fmt.Errorf("derive m/44'/%d'/%d key failed: %v", coinIdx, accIdx, err)
	}

	// external chain: m / 44' / coinIdx' / accIdx / 0
	accExtK, err := accK.Child(0)
	if err != nil {
		return nil, "", fmt.Errorf("derive m/44'/%d'/%d/0 key failed: %v", coinIdx, accIdx, err)
	}

	addrKey, err := accExtK.Child(addrIdx)
	if err != nil {
		return nil, "", fmt.Errorf("derive m/44'/%d'/%d/0/%d key failed: %v", coinIdx, accIdx, addrIdx, err)
	}
	path := fmt.Sprintf("m/44'/%d'/%d/0/%d", coinIdx, accIdx-HardenedKeyZeroIndex, addrIdx)
	return addrKey, path, nil
}
