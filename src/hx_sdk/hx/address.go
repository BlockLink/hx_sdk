/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: address.go
 * Date: 2018-08-31
 *
 */

package hx

import (
	"crypto/sha512"
	"fmt"
	"strings"

	"bytes"

	"github.com/HcashOrg/hcd/hcec/secp256k1"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/golangcrypto/ripemd160"
	"github.com/tyler-smith/go-bip39"
)

const (
	CoinHX string = "HX"

	VersionNormalAddr   = 0x35
	VersionMultisigAddr = 0x32
	VersionContractAddr = 0x1c
)

/**
 * HX address struct
 */
type HxAddr struct {
	Addr string //hx main chain address string
}

// using BIP44 to manage hx address
// m / purpose' / coin_type' / account' / change / address_index
// https://github.com/satoshilabs/slips/blob/master/slip-0044.md

func MnemonicToSeed(mnemonic, password string) []byte {
	return bip39.NewSeed(mnemonic, password)
}

func getMasterkey(seed []byte, mainnet bool) (*hdkeychain.ExtendedKey, error) {

	// main net
	if mainnet {

		return hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)

	} else { //test net

		return hdkeychain.NewMaster(seed, &chaincfg.TestNet3Params)

	}

}

func getHxAddressByWif(wif string, version uint32) (addr string, err error) {
	addr = ""
	wif_key, err := getPrivKey(wif)
	if err != nil {
		return
	}
	addr = GetAddressByPubkey(wif_key.PubKey().SerializeCompressed(), "main", version)

	return
}
func GetNewPrivate() (privWif string, pubWif string, addr string, err error) {
	priv, err := btcec.NewPrivateKey(secp256k1.S256())
	if err != nil {
		return
	}
	tmp_Wif, err := btcutil.NewWIF(priv, &chaincfg.MainNetParams, true)
	if err != nil {
		return
	}
	privWif = tmp_Wif.String()

	pubByte := priv.PubKey().SerializeCompressed()
	myRipemd := ripemd160.New()

	myRipemd.Write(pubByte[:])
	checksum := myRipemd.Sum(nil)

	pubByte = append(pubByte, checksum[0:4]...)
	pubWif = "HX" + base58.Encode(pubByte)
	addr, err = getHxAddressByWif(privWif, VersionNormalAddr)
	if err != nil {
		return
	}
	return
}

func getAccountExtentkey(masterKey *hdkeychain.ExtendedKey, account uint32, addrIndex uint32) (*hdkeychain.ExtendedKey, string, error) {
	path := fmt.Sprintf("")
	// purpose & coin_tyep & change
	purpose := uint32(0x8000002C)
	coinType := uint32(0x80000000 + 999)
	change := uint32(0)

	// m / 44'
	purposeKey, err := masterKey.Child(purpose)
	if err != nil {
		return nil, path, fmt.Errorf("create purpose key failed: %v", err)
	}

	// m / 44' / 999'
	coinTypeKey, err := purposeKey.Child(coinType)
	if err != nil {
		return nil, path, fmt.Errorf("create coin type key failed: %v", err)
	}

	// m / 44' / 999' / account
	accountKey, err := coinTypeKey.Child(account)
	if err != nil {
		return nil, path, fmt.Errorf("create account key failed: %v", err)
	}

	// m / 44' / 999' / account / 0
	changeKey, err := accountKey.Child(change)
	if err != nil {
		return nil, path, fmt.Errorf("create change key failed: %v", err)
	}

	// m / 44' / 999' / account / 0 / addr
	addressKey, err := changeKey.Child(addrIndex)
	if err != nil {
		return nil, path, fmt.Errorf("create address key failed: %v", err)
	}

	return addressKey, path, err
}

func getAccountAddr(addressKey *hdkeychain.ExtendedKey, nettype string, version uint32) (string, error) {

	ecPubkey, err := addressKey.ECPubKey()
	if err != nil {
		return "", fmt.Errorf("get ecPubkey failed: %v", err)
	}

	// sha512 & ripemd160
	pubkeyByte := ecPubkey.SerializeCompressed()
	return GetAddressByPubkey(pubkeyByte, nettype, version), nil
}

func GetAddressByPubkey(pubkeyByte []byte, nettype string, version uint32) string {
	var res []byte
	res = append(res, byte(version))

	myRipemd := ripemd160.New()
	sha512Byte := sha512.Sum512(pubkeyByte)
	myRipemd.Write(sha512Byte[:])
	addrByte := myRipemd.Sum(nil)
	//fmt.Println(len(addrByte))
	res = append(res, addrByte...)

	myRipemd.Reset()
	myRipemd.Write(res)
	addrByteChecksum := myRipemd.Sum(nil)

	res = append(res, addrByteChecksum[0:4]...)
	//fmt.Println(len(addrByte))

	prefix := CoinHX
	nettype = strings.ToLower(nettype)
	if nettype == "main" || nettype == "mainnet" {
		prefix = CoinHX
	} else {
		// todo: test net with differect prefix
	}

	return prefix + string(base58.Encode(res))
}

func ValidateAddress(address, net string) bool {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	if len(address) <= 2 {
		return false
	}
	// res := address[2:]
	buf := base58.Decode(address[2:])
	if len(buf) < 8 {
		return false
	}

	tocheck := buf[0 : len(buf)-4]
	check := buf[len(buf)-4 : len(buf)]

	myRipemd := ripemd160.New()
	myRipemd.Write(tocheck)
	addrByteChecksum := myRipemd.Sum(nil)

	return bytes.Compare(check, addrByteChecksum[0:4]) == 0
}

func getAddressBytes(addressKey *hdkeychain.ExtendedKey) ([]byte, error) {

	ecPubkey, err := addressKey.ECPubKey()
	if err != nil {
		return nil, fmt.Errorf("get ecPubkey failed: %v", err)
	}

	// sha512 & ripemd160
	pubkeyByte := ecPubkey.SerializeCompressed()
	myRipemd := ripemd160.New()
	sha512Byte := sha512.Sum512(pubkeyByte)
	myRipemd.Write(sha512Byte[:])
	addrByte := myRipemd.Sum(nil)

	return addrByte, nil
}

func getWifkey(addressKey *hdkeychain.ExtendedKey) (string, error) {

	ecPrivkey, err := addressKey.ECPrivKey()
	if err != nil {
		return "", fmt.Errorf("get ecPrivkey failed: %v", err)
	}

	wif, err := btcutil.NewWIF(ecPrivkey, &chaincfg.MainNetParams, false)
	if err != nil {
		return "", fmt.Errorf("get wif failed: %v", err)
	}

	return wif.String(), nil
}

// normal addr version: 0x35
// multisig addr version: 0x32
// contract addr version: 0x1c
func GetAddress(seed []byte, nettype string, account uint32, addrIndex uint32, version uint32) (string, string, error) {

	mastkey, err := getMasterkey(seed, true) //hx using btc mainchain cfg
	if err != nil {
		return "", "", fmt.Errorf("in GetAddress function, get mastkey failed: %v", err)
	}

	accountExtendKey, path, err := getAccountExtentkey(mastkey, account, addrIndex)
	if err != nil {
		return "", "", fmt.Errorf("in GetAddress function, get accountExtensionKey failed: %v", err)
	}

	accountAddr, err := getAccountAddr(accountExtendKey, nettype, version)
	if err != nil {
		return "", "", fmt.Errorf("in GetAddress function, get accountAddr failed: %v", err)
	}

	return accountAddr, path, nil
}

func GetAddressBytes(addr string) ([]byte, error) {

	if len(addr) <= 2 {
		return nil, fmt.Errorf("in GetAddressBytes function, wrong addr format")
	}

	base58_addr := addr[2:]

	addrBytes := base58.Decode(base58_addr)

	return addrBytes[:len(addrBytes)-4], nil
}

func GetPubkeyBytes(pub string) ([]byte, error) {

	if len(pub) <= 2 {
		return nil, fmt.Errorf("in GetAddressBytes function, wrong addr format")
	}

	base58_addr := pub[2:]

	pubBytes := base58.Decode(base58_addr)

	return pubBytes[:len(pubBytes)-4], nil
}

func DerivePubkey(wif string) (pub string, err error) {
	priv, err := getPrivKey(wif)
	if err != nil {
		return
	}

	// buf := priv.PubKey().SerializeUncompressed() // SerializeCompressed()
	buf := priv.PubKey().SerializeCompressed()

	myRipemd := ripemd160.New()

	myRipemd.Write(buf[:])
	checksum := myRipemd.Sum(nil)

	buf = append(buf, checksum[0:4]...)
	return "HX" + base58.Encode(buf), nil
}

func GetAddressKey(seed []byte, account uint32, addrIndex uint32) (*hdkeychain.ExtendedKey, error) {

	mastkey, err := getMasterkey(seed, true) //hx using btc mainchain cfg
	if err != nil {
		return nil, fmt.Errorf("in GetAddress function, get mastkey failed: %v", err)
	}

	accountExtendKey, _, err := getAccountExtentkey(mastkey, account, addrIndex)
	if err != nil {
		return nil, fmt.Errorf("in GetAddress function, get accountExtensionKey failed: %v", err)
	}

	return accountExtendKey, nil
}

func ExportWif(seed []byte, account uint32, addrIndex uint32) (string, error) {

	mastkey, err := getMasterkey(seed, true) //hx using btc mainchain cfg
	if err != nil {
		return "", fmt.Errorf("in ExportWif function, get mastkey failed: %v", err)
	}

	accountExtendKey, _, err := getAccountExtentkey(mastkey, account, addrIndex)
	if err != nil {
		return "", fmt.Errorf("in ExportWif function, get accountExtensionKey failed: %v", err)
	}

	wifKey, err := getWifkey(accountExtendKey)
	if err != nil {
		return "", fmt.Errorf("in ExportWif function, get wif failed: %v", err)
	}

	return wifKey, nil

}

func getPrivKey(wif string) (*btcec.PrivateKey, error) {
	wifstruct, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, fmt.Errorf("decode wif string failed: %v", err)
	}

	return wifstruct.PrivKey, err
}

func ImportWif(wifstr string) (*btcec.PrivateKey, error) {

	return getPrivKey(wifstr)
}

func IsCanonical(sig []byte) bool {
	/*
		tmp := !(sig[1] & 0x80) &&
			!(sig[1] == 0 && !(sig[2] & 0x80)) &&
			!(sig[33] & 0x80) &&
			!(sig[33] == 0 && !(sig[34] & 0x80))

		return tmp
	*/
	return true
}
