package hx

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"errors"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/hcec/secp256k1"
	"github.com/HcashOrg/hcd/hcutil"
	"github.com/HcashOrg/hcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ebfe/keccak"
	"hx_sdk/common"
)

const (
	CoinHC     = "HC"
	CoinBTC    = "BTC"
	CoinBCH    = "BCH"
	CoinLTC    = "LTC"
	CoinETH    = "ETH"
	CoinUSDT   = "USDT"
	CoinPAX    = "PAX"
	CoinERCPAX = "ERCPAX"
	CoinERCELF = "ERCELF"
	CoinELF    = "ELF"
)

var (
	//ethSigSuffix0     = "25"
	//ethSigSuffix1     = "26"
	//ethSigSuffixByte0 = byte(0x25)
	//ethSigSuffixByte1 = byte(0x26)
	ethSigSuffix0     = "1b"
	ethSigSuffix1     = "1c"
	ethSigSuffixByte0 = byte(0x1b)
	ethSigSuffixByte1 = byte(0x1c)
)

func SetTestnetEthSig() {
	ethSigSuffix0 = "1b"
	ethSigSuffix1 = "1c"

	ethSigSuffixByte0 = byte(0x1b)
	ethSigSuffixByte1 = byte(0x1c)
}

// SignAddress sign address to bind  to hx chain
func SignAddress(wif, address, coin string) (string, error) {
	switch coin {
	case CoinHC:
		return hcSignAddress(wif, address)

	case CoinBTC:
		fallthrough
	case CoinUSDT:
		return btcSignAddress(wif, address)

	case CoinBCH:
		return "", errors.New("not support!")
		//return bchSignAddress(wif, address)

	case CoinLTC:
		return ltcSignAddress(wif, address)

	case CoinETH:
		fallthrough
	case CoinPAX:
		fallthrough
	case CoinERCPAX:
		fallthrough
	case CoinELF:
		return ethSignAddress2(wif, address)
	default:
		return ethSignAddress2(wif, address)
	}

	// return "", fmt.Errorf("SignAddress: invalid coin: %s", coin)
}

func hcSignAddress(wif, addr string) (sig string, err error) {
	w, err := hcutil.DecodeWIF(wif)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, "Hc Signed Message:\n")
	wire.WriteVarString(&buf, 0, addr)
	messageHash := chainhash.HashB(buf.Bytes())

	pkCast, ok := w.PrivKey.(*secp256k1.PrivateKey)
	if !ok {
		fmt.Printf("Unable to create secp256k1.PrivateKey" +
			"from chainec.PrivateKey")
		return
	}
	res, err := secp256k1.SignCompact(secp256k1.S256(), pkCast, messageHash, true)

	return base64.StdEncoding.EncodeToString(res), nil
}

// fast hash
func Keccak256(data ...[]byte) []byte {
	h := keccak.New256()
	for _, b := range data {
		h.Write(b)
	}
	r := h.Sum(nil)

	return r
}

// eth 签名
func Sign2(wif string, msg []byte) (sig []byte, err error) {
	buf, err := hex.DecodeString(wif)
	if err != nil {
		fmt.Println("decode wif failed: ", err)
		return
	}

	key2, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), buf)

	s2, err := secp256k1.SignCompact(secp256k1.S256(), key2, []byte(msg), false)
	return s2, err
}

// use bts sign
func ethSignAddress2(wif, addr string) (data string, err error) {
	addr = strings.ToLower(addr)
	if strings.HasPrefix(addr, "0x") {
		addr = addr[2:]
	}

	baddr, _ := hex.DecodeString(addr)
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(baddr))

	h := Keccak256(append([]byte(msg), baddr...))

	sig, err := Sign2(wif, h)
	// Convert to Ethereum signature format with 'recovery id' v at the end.
	v := sig[0] - 27
	copy(sig, sig[1:])

	// 测试链 00 -> 1b   01 -> 1c
	// 正式链 00 -> 25   01 -> 26
	if v == byte(0) {
		sig[64] = ethSigSuffixByte0

		// res = res[0:len(res)-2] + "1b"
	} else {
		sig[64] = ethSigSuffixByte1
		// res = res[0:len(res)-2] + "1c"
	}

	return "0x" + hex.EncodeToString(sig), nil
}

/*
func ethSignAddress(wif, addr string) (sig string, err error) {
	baddr, _ := hex.DecodeString(addr)
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(baddr))

	h := crypto.Keccak256(append([]byte(msg), baddr...))
	// fmt.Println("msg addr:", hex.EncodeToString(h))

	buf, err := hex.DecodeString(wif)
	if err != nil {
		fmt.Println("decode wif failed: ", err)
		return
	}

	key, err := crypto.ToECDSA(buf)

	data, err := crypto.Sign(h, key)

	if err != nil {
		fmt.Printf("sign eth failed: %v", err)
		return
	}
	// fmt.Println("signed:", hex.EncodeToString(data))

	// TODO: 测试链和正式链的结尾不同
	// 测试链 00 -> 1b   01 -> 1c
	// 正式链 00 -> 25   01 -> 26
	res := hex.EncodeToString(data)
	suffix := res[len(res)-2 : len(res)]
	if suffix == "00" {
		res = res[0:len(res)-2] + "1b"
	} else if suffix == "01" {
		res = res[0:len(res)-2] + "1c"
	} else {
		return "", fmt.Errorf("invalid signature suffix: %v", suffix)
	}

	return "0x" + res, nil
}
*/

func DoubleHashB(b []byte) []byte {
	first := sha256.Sum256(b)
	second := sha256.Sum256(first[:])
	return second[:]
}

//func bchSignAddress(wif, addr string) (sig string, err error) {
//	w, err := btcutil.DecodeWIF(wif)
//	if err != nil {
//		return
//	}
//
//	var buf bytes.Buffer
//	wire.WriteVarString(&buf, 0, "Bitcoin Signed Message:\n")
//	wire.WriteVarString(&buf, 0, addr)
//
//	messageHash := DoubleHashB(buf.Bytes())
//
//	pkCast := secp256k1.PrivateKey(*w.PrivKey)
//
//	res, err := secp256k1.SignCompactBCH(secp256k1.S256(), &pkCast, messageHash, true)
//
//	return base64.StdEncoding.EncodeToString(res), nil
//}

func btcSignAddress(wif, addr string) (sig string, err error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, "Bitcoin Signed Message:\n")
	wire.WriteVarString(&buf, 0, addr)

	messageHash := DoubleHashB(buf.Bytes())

	pkCast := secp256k1.PrivateKey(*w.PrivKey)

	res, err := secp256k1.SignCompact(secp256k1.S256(), &pkCast, messageHash, true)

	return base64.StdEncoding.EncodeToString(res), nil
}

func ltcSignAddress(wif, addr string) (sig string, err error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, "Litecoin Signed Message:\n")
	wire.WriteVarString(&buf, 0, addr)

	messageHash := DoubleHashB(buf.Bytes())

	pkCast := secp256k1.PrivateKey(*w.PrivKey)
	res, err := secp256k1.SignCompact(secp256k1.S256(), &pkCast, messageHash, true)

	return base64.StdEncoding.EncodeToString(res), nil
}

func btsSign(wif string, data []byte) (res []byte, err error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return
	}

	pkCast := secp256k1.PrivateKey(*w.PrivKey)
	fmt.Println("ecPrivkey", *pkCast.D)
	res, err = secp256k1.SignCompact(secp256k1.S256(), &pkCast, data, true)
	return
}

// convert der sign to hx signature
func DerSignToHxSign(sig string, recid int) ([]byte, error) {
	bsig, err := common.ConvertDerSig(sig)
	if err != nil {
		return nil, err
	}

	var hxsig []byte = make([]byte, 65)
	// bsig = r + s
	hxsig[0] = byte(27 + 4 + int(recid))
	copy(hxsig[1:], bsig[:])

	return hxsig, nil
}
