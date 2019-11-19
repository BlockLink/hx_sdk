/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: trx.go
 * Date: 2018-09-04
 *
 */

package hx

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"hx_sdk/btssign"
)

const expireTimeout = 86000

// define trx structure
type Transaction struct {
	Hx_ref_block_num    uint16 `json:"ref_block_num"`
	Hx_ref_block_prefix uint32 `json:"ref_block_prefix"`
	Hx_expiration       string `json:"expiration"`

	Hx_operations [][]interface{} `json:"operations"`
	Hx_extensions []interface{}   `json:"extensions"`
	Hx_signatures []string        `json:"signatures"`

	Expiration uint32        `json:"-"`
	Operations []interface{} `json:"-"`
}

func DefaultTransaction() *Transaction {

	return &Transaction{
		0,
		0,
		"",
		nil,
		nil,
		nil,
		0,
		nil,
	}
}

func GetId(id string) (uint32, error) {

	idSlice := strings.Split(id, ".")

	if len(idSlice) != 3 {
		return 0, fmt.Errorf("in GetId function, get account id failed")
	}

	res, err := strconv.ParseUint(idSlice[2], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("in GetId function, Parse id error %v", err)
	}

	return uint32(res), nil

}

func Str2Time(str string) int64 {

	str += "Z"
	t, err := time.Parse(time.RFC3339, str)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	return t.Unix()

}

func Time2Str(t int64) string {

	l_time := time.Unix(t, 0).UTC()
	timestr := l_time.Format(time.RFC3339)

	timestr = timestr[:len(timestr)-1]

	return timestr
}

// in multiple precision mode
func CalculateFee(basic_op_fee int64, len_memo int64) int64 {

	var basic_memo_fee int64 = 1
	return basic_op_fee + len_memo*basic_memo_fee
}

func (asset *Asset) SetAssetBySymbol(symbol string) {
	symbol = strings.ToUpper(symbol)

	if symbol == "HX" {
		asset.Hx_asset_id = "1.3.0"
	} else if symbol == "BTC" {
		asset.Hx_asset_id = "1.3.1"
	} else if symbol == "LTC" {
		asset.Hx_asset_id = "1.3.2"
	} else if symbol == "HC" {
		asset.Hx_asset_id = "1.3.3"
	}

}

func GetRefblockInfo(info string) (uint16, uint32, error) {

	refinfo := strings.Split(info, ",")
	// refinfo := []string{"21771", "761216631"}

	if len(refinfo) != 2 {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, get refblockinfo failed")
	}
	ref_block_num_str, ref_block_prefix_str := refinfo[0], refinfo[1]
	ref_block_num, err := strconv.ParseUint(ref_block_num_str, 10, 16)
	if err != nil {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, convert ref_block_num failed: %v", err)
	}

	ref_block_prefix, err := strconv.ParseUint(ref_block_prefix_str, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("in GetRefblockInfo function, convert ref_block_prefix failed: %v", err)
	}

	return uint16(ref_block_num), uint32(ref_block_prefix), nil
}

func GetSignature(wif string, hash []byte) ([]byte, error) {

	ecPrivkey, err := ImportWif(wif)
	if err != nil {
		return nil, fmt.Errorf("in GetSignature function, get ecprivkey failed: %v", err)
	}

	ecPrivkeyByte := ecPrivkey.Serialize()
	return btssign.SignCompact(hash, ecPrivkeyByte, true)
	//fmt.Println("the uncompressed pubkey is: ", hex.EncodeToString(ecPrivkey.PubKey().SerializeUncompressed()))
	//fmt.Println("the compressed pubkey is: ", hex.EncodeToString(ecPrivkey.PubKey().SerializeCompressed()))
	/*
		for {
			sig, err := bts.SignCompact(hash, ecPrivkeyByte, true)
			if err != nil {
				return nil, fmt.Errorf("in GetSignature function, sign compact failed: %v", err)
			}

			pubkey_byte, err := bts.RecoverPubkey(hash, sig, true)
			if err != nil {
				return nil, fmt.Errorf("in GetSignature function, sign compact failed: %v", err)
			}
			fmt.Println("recoverd pubkey is: ", hex.EncodeToString(pubkey_byte))

			if bytes.Compare(ecPrivkey.PubKey().SerializeCompressed(), pubkey_byte) == 0 {
				return sig, nil
			}

		}
	*/
}

func BuildUnsignedTx(refinfo, from, to, memo, assetId string, amount, fee int64, guarantee_id string) (*Transaction, error) {
	// build unsigned tx hash
	asset_amount := DefaultAsset()
	asset_amount.Hx_amount = amount
	asset_amount.Hx_asset_id = assetId // SetAssetBySymbol(symbol)

	asset_fee := DefaultAsset()
	asset_fee.Hx_amount = fee
	asset_fee.SetAssetBySymbol("HX")

	transferOp := DefaultTransferOperation()
	transferOp.Hx_fee = asset_fee
	transferOp.Hx_from_addr = from
	transferOp.Hx_to_addr = to
	transferOp.Hx_amount = asset_amount

	if memo == "" {
		transferOp.Hx_memo = nil
	} else {
		memo_trx := DefaultMemo()
		memo_trx.Message = memo
		memo_trx.IsEmpty = false
		memo_trx.Hx_message = hex.EncodeToString(append(make([]byte, 4), []byte(memo_trx.Message)...))
		transferOp.Hx_memo = &memo_trx
	}

	if guarantee_id != "" {
		transferOp.Hx_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		// panic("get refinfo failed!")
		return nil, err
	}

	return &Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{0, transferOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*transferOp},
	}, nil

}

func BuildUnsignedTxHash(refinfo, from, to, memo, assetId string, amount, fee int64,
	guarantee_id, chain_id string) ([]byte, error) {
	tx, err := BuildUnsignedTx(refinfo, from, to, memo, assetId, amount, fee, guarantee_id)
	if err != nil {
		return nil, err
	}
	res := tx.Serialize()
	fmt.Printf("hex before sign: %v\n", hex.EncodeToString(res))
	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	return toSign[:], nil
}

func RebuildTxWithSign(refinfo, from, to, memo, assetId string, amount, fee int64,
	guarantee_id, sig string) ([]byte, error) {
	tx, err := BuildUnsignedTx(refinfo, from, to, memo, assetId, amount, fee, guarantee_id)
	if err != nil {
		return nil, err
	}

	tx.Hx_signatures = append(tx.Hx_signatures, sig)
	fmt.Printf("RebuildTxWithSign: signature=%v\n", sig)

	b, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	return b, nil
}

func BuildTransferTransaction(refinfo, wif string, from, to, memo, assetId string, amount, fee int64,
	symbol string, guarantee_id, chain_id string) (b []byte, err error) {

	asset_amount := DefaultAsset()
	asset_amount.Hx_amount = amount
	asset_amount.Hx_asset_id = assetId // SetAssetBySymbol(symbol)

	asset_fee := DefaultAsset()
	//asset_fee.Hx_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Hx_amount = fee
	asset_fee.SetAssetBySymbol("HX")

	transferOp := DefaultTransferOperation()
	transferOp.Hx_fee = asset_fee
	transferOp.Hx_from_addr = from
	transferOp.Hx_to_addr = to
	transferOp.Hx_amount = asset_amount

	if memo == "" {
		transferOp.Hx_memo = nil
	} else {
		memo_trx := DefaultMemo()
		memo_trx.Message = memo
		memo_trx.IsEmpty = false
		memo_trx.Hx_message = hex.EncodeToString(append(make([]byte, 4), []byte(memo_trx.Message)...))
		transferOp.Hx_memo = &memo_trx
	}

	if guarantee_id != "" {
		transferOp.Hx_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-09-26T09:14:40"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{0, transferOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*transferOp},
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Hx_signatures = append(transferTrx.Hx_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return
}

// bind tunnel address fee is not needed, always 0
func BuildBindAccountTransaction(refinfo, wif, addr string, fee int64,
	crosschain_addr, crosschain_symbol, crosschain_wif string, guarantee_id, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Hx_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Hx_amount = fee
	asset_fee.SetAssetBySymbol("HX")

	bindOp := DefaultAccountBindOperation()
	bindOp.Hx_fee = asset_fee
	bindOp.Hx_crosschain_type = crosschain_symbol
	bindOp.Hx_addr = addr
	if guarantee_id != "" {
		bindOp.Hx_guarantee_id = guarantee_id
	}

	// sign the addr
	addrByte, err := GetAddressBytes(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	toSign := sha256.Sum256(addrByte)
	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}
	bindOp.Hx_account_signature = hex.EncodeToString(sig)
	bindOp.Hx_tunnel_address = crosschain_addr
	crosschain_sig, err := SignAddress(crosschain_wif, crosschain_addr, crosschain_symbol)
	if err != nil {
		fmt.Println(err)
		return
	}
	bindOp.Hx_tunnel_signature = crosschain_sig

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-11-01T02:59:30"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		// panic("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{10, bindOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*bindOp},
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign = sha256.Sum256(append(chainid_byte, res...))

	sig, err = GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Hx_signatures = append(transferTrx.Hx_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return
}

// BuildUnBindAccountTransaction bind tunnel address
// wif: hx wif
// addr: hx address
// fee:
// crosschain_addr: btc/eth/ltc/hc address
// crosschain_symbol: btc/eth/ltc/hc
// crosschain_wif: btc/eth/ltc/hc wif
// chain_id
func BuildUnBindAccountTransaction(refinfo, wif, addr string, fee int64,
	crosschain_addr, crosschain_symbol, crosschain_wif, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Hx_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Hx_amount = fee
	asset_fee.SetAssetBySymbol("HX")

	unbindOp := DefaultAccountUnBindOperation()
	unbindOp.Hx_fee = asset_fee
	unbindOp.Hx_crosschain_type = crosschain_symbol
	unbindOp.Hx_addr = addr
	//sign the addr
	addrByte, err := GetAddressBytes(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	toSign := sha256.Sum256(addrByte)
	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}
	unbindOp.Hx_account_signature = hex.EncodeToString(sig)
	unbindOp.Hx_tunnel_address = crosschain_addr
	crosschain_sig, err := SignAddress(crosschain_wif, crosschain_addr, crosschain_symbol)
	if err != nil {
		fmt.Println(err)
		return
	}
	unbindOp.Hx_tunnel_signature = crosschain_sig

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-11-01T02:59:30"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		// panic("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{11, unbindOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*unbindOp},
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign = sha256.Sum256(append(chainid_byte, res...))

	sig, err = GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Hx_signatures = append(transferTrx.Hx_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return
}

func BuildWithdrawCrosschainTransaction(refinfo, wif, addr string, fee int64,
	crosschain_addr, crosschain_symbol, assetId, crosschain_amount, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Hx_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Hx_amount = fee
	asset_fee.SetAssetBySymbol("HX")

	withdrawOp := DefaultWithdrawCrosschainOperation()
	withdrawOp.Hx_withdraw_account = addr
	withdrawOp.Hx_amount = crosschain_amount
	withdrawOp.Hx_asset_symbol = crosschain_symbol
	withdrawOp.Hx_asset_id = assetId
	/*
		if crosschain_symbol == "BTC" {
			withdrawOp.Hx_asset_symbol = "BTC"
			withdrawOp.Hx_asset_id = "1.3.1"
		} else if crosschain_symbol == "LTC" {
			withdrawOp.Hx_asset_symbol = "LTC"
			withdrawOp.Hx_asset_id = "1.3.2"
		} else if crosschain_symbol == "HC" {
			withdrawOp.Hx_asset_symbol = "HC"
			withdrawOp.Hx_asset_id = "1.3.3"
		}
	*/
	withdrawOp.Hx_crosschain_account = crosschain_addr

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-11-01T02:59:30"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{61, withdrawOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*withdrawOp},
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Hx_signatures = append(transferTrx.Hx_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	return
}

func BuildRegisterAccountTransaction(refinfo, wif, addr, public_key string, fee int64,
	guarantee_id, register_name, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Hx_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Hx_amount = fee
	asset_fee.SetAssetBySymbol("HX")

	registerOp := DefaultRegisterAccountOperation()
	registerOp.Hx_fee = asset_fee
	registerOp.Hx_payer = addr
	registerOp.Hx_name = register_name
	registerOp.Hx_owner.Hx_key_auths = [][]interface{}{{public_key, 1}}
	registerOp.Hx_owner.Key_auths = public_key
	registerOp.Hx_active.Hx_key_auths = registerOp.Hx_owner.Hx_key_auths
	registerOp.Hx_active.Key_auths = public_key
	registerOp.Hx_owner.Key_auths = public_key
	registerOp.Hx_options.Hx_memo_key = public_key

	if guarantee_id != "" {
		registerOp.Hx_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	// expir_str := "2018-11-06T06:21:33"
	// expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{5, registerOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*registerOp},
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Hx_signatures = append(transferTrx.Hx_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	return
}

func BuildLockBalanceTransaction(refinfo, wif, addr, account_id, lock_asset_id string,
	lock_asset_amount, fee int64, miner_id, miner_address, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Hx_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Hx_amount = fee
	asset_fee.SetAssetBySymbol("HX")

	lockOp := DefaultLockBalanceOperation()
	lockOp.Hx_fee = asset_fee
	lockOp.Hx_lock_asset_id = lock_asset_id
	lockOp.Hx_lock_asset_amount = lock_asset_amount

	if account_id == "" {
		lockOp.Hx_lock_balance_account = "1.2.0"
	} else {
		lockOp.Hx_lock_balance_account = account_id
	}
	lockOp.Hx_lock_balance_addr = addr
	lockOp.Hx_lockto_miner_account = miner_id
	lockOp.Hx_contract_addr = miner_address

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-11-07T02:18:30"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{55, lockOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*lockOp},
	}

	res := transferTrx.Serialize()

	//seed := MnemonicToSeed("venture lazy digital aware plug hire acquire abuse chunk know gloom snow much employ glow rich exclude allow", "123")
	//addrkey, _:= GetAddressKey(seed, 0, 0)
	//addr, _ := GetAddress(seed,0,0, 0x35)
	//fmt.Println("addr is: ", addr)
	//wif, _ := ExportWif(seed, 0, 0)
	//fmt.Println("wif is: ", wif)

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Hx_signatures = append(transferTrx.Hx_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	return
}

func BuildRedeemBalanceTransaction(refinfo, wif, addr, account_id, foreclose_asset_id string,
	foreclose_asset_amount, fee int64, miner_id, miner_address, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Hx_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Hx_amount = fee
	asset_fee.SetAssetBySymbol("HX")

	forecloseOp := DefaultForecloseBalanceOperation()
	forecloseOp.Hx_fee = asset_fee
	forecloseOp.Hx_foreclose_asset_id = foreclose_asset_id
	forecloseOp.Hx_foreclose_asset_amount = foreclose_asset_amount

	forecloseOp.Hx_foreclose_miner_account = miner_id
	forecloseOp.Hx_foreclose_contract_addr = miner_address

	if account_id == "" {
		forecloseOp.Hx_foreclose_account = "1.2.0"
	} else {
		forecloseOp.Hx_foreclose_account = account_id
	}

	forecloseOp.Hx_foreclose_addr = addr

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	//expir_str := "2018-11-07T02:18:30"
	//expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{56, forecloseOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*forecloseOp},
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Hx_signatures = append(transferTrx.Hx_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

// obtain_asset_arr format: []string{"citizen10,100,1.3.0", "citizen11,101,1.3.0"}
func BuildObtainPaybackTransaction(refinfo, wif, addr string, fee int64,
	obtain_asset_arr []string, guarantee_id, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Hx_amount = CalculateFee(2000, int64(len(memo) + 3))
	asset_fee.Hx_amount = fee
	asset_fee.SetAssetBySymbol("HX")

	obtainOp := DefaultObtainPaybackOperation()
	obtainOp.Hx_pay_back_owner = addr
	obtainOp.Hx_fee = asset_fee

	obtainOp.Hx_pay_back_balance = [][]interface{}{}
	if len(obtain_asset_arr) == 0 {
		return nil, fmt.Errorf("obtain asset arr forma error")
	}
	for i := 0; i < len(obtain_asset_arr); i++ {
		obtain_assets := strings.Split(obtain_asset_arr[i], ",")
		if len(obtain_assets) != 3 {
			return nil, fmt.Errorf("obtain asset arr forma error")
		}

		obtain_asset := DefaultAsset()
		amount, err := strconv.ParseInt(obtain_assets[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse obtain asset amount error")
		}
		obtain_asset.Hx_amount = amount
		obtain_asset.Hx_asset_id = obtain_assets[2]
		tmp_pay_back := [][]interface{}{{obtain_assets[0], obtain_asset}}
		obtainOp.Hx_pay_back_balance = append(obtainOp.Hx_pay_back_balance, tmp_pay_back...)
		obtainOp.citizen_name = append(obtainOp.citizen_name, obtain_assets[0])
		obtainOp.obtain_asset = append(obtainOp.obtain_asset, obtain_asset)
	}

	if guarantee_id != "" {
		obtainOp.Hx_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	// expir_str := "2018-11-07T06:20:30"
	// expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		fmt.Println("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{73, obtainOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*obtainOp},
	}

	res := transferTrx.Serialize()

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Hx_signatures = append(transferTrx.Hx_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return
}

// fee is basic fee of HX chain
func BuildContractInvokeTransaction(refinfo, wif, addr string, fee int64, gas_price, gas_limit int64, contract_id, contract_api, contract_arg string,
	guarantee_id, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Hx_amount = CalculateFee(2000, int64(len(memo) + 3))
	gas_count := gas_limit / 100 * gas_price
	if gas_limit%100 != 0 {
		gas_count += gas_price
	}
	asset_fee.Hx_amount = fee + gas_count
	asset_fee.SetAssetBySymbol("HX")

	contractOp := DefaultContractInvokeOperation()
	contractOp.Hx_fee = asset_fee

	contractOp.Hx_invoke_cost = uint64(gas_limit)
	contractOp.Hx_gas_price = uint64(gas_price)
	contractOp.Hx_caller_addr = addr
	priv, err := getPrivKey(wif)
	if err != nil {
		return nil, fmt.Errorf("get private key from wif error")
	}
	buf := priv.PubKey().SerializeCompressed()
	contractOp.Hx_caller_pubkey = hex.EncodeToString(buf)
	contractOp.Hx_contract_id = contract_id
	contractOp.Hx_contract_api = contract_api
	contractOp.Hx_contract_arg = contract_arg

	if guarantee_id != "" {
		contractOp.Hx_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	// expir_str := "2018-11-07T06:20:30"
	// expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		// panic("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{79, contractOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*contractOp},
	}

	res := transferTrx.Serialize()

	//seed := MnemonicToSeed("venture lazy digital aware plug hire acquire abuse chunk know gloom snow much employ glow rich exclude allow", "123")
	//addrkey, _:= GetAddressKey(seed, 0, 0)
	//addr, _ := GetAddress(seed,0,0, 0x35)
	//fmt.Println("addr is: ", addr)
	//wif, _ := ExportWif(seed, 0, 0)
	//fmt.Println("wif is: ", wif)

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Hx_signatures = append(transferTrx.Hx_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	return
}

// transfer to contract
func BuildContractTransferTransaction(refinfo, wif, addr string, fee int64, amount int64, assetId string, gas_price, gas_limit int64, contract_id, param string,
	guarantee_id, chain_id string) (b []byte, err error) {

	asset_fee := DefaultAsset()
	//asset_fee.Hx_amount = CalculateFee(2000, int64(len(memo) + 3))
	gas_count := gas_limit / 100 * gas_price
	if gas_limit%100 != 0 {
		gas_count += gas_price
	}
	asset_fee.Hx_amount = fee + gas_count
	asset_fee.SetAssetBySymbol("HX")

	asset_amount := DefaultAsset()
	asset_amount.Hx_amount = amount
	asset_amount.Hx_asset_id = assetId

	contractOp := DefaultContractTransferOperation()
	contractOp.Hx_fee = asset_fee
	contractOp.Hx_amount = asset_amount

	contractOp.Hx_invoke_cost = uint64(gas_limit)
	contractOp.Hx_gas_price = uint64(gas_price)
	contractOp.Hx_caller_addr = addr
	priv, err := getPrivKey(wif)
	if err != nil {
		return nil, fmt.Errorf("get private key from wif error")
	}
	buf := priv.PubKey().SerializeCompressed()
	contractOp.Hx_caller_pubkey = hex.EncodeToString(buf)
	contractOp.Hx_contract_id = contract_id
	contractOp.Hx_param = param

	if guarantee_id != "" {
		contractOp.Hx_guarantee_id = guarantee_id
	}

	expir_sec := time.Now().Unix() + expireTimeout
	expir_str := Time2Str(expir_sec)
	// expir_str := "2018-11-07T06:20:30"
	// expir_sec := Str2Time(expir_str)

	ref_block_num, ref_block_prefix, err := GetRefblockInfo(refinfo)
	if err != nil {
		// panic("get refinfo failed!")
		return
	}

	transferTrx := Transaction{
		ref_block_num,
		ref_block_prefix,
		expir_str,
		[][]interface{}{{81, contractOp}},
		make([]interface{}, 0),
		make([]string, 0),
		uint32(expir_sec),
		[]interface{}{*contractOp},
	}

	res := transferTrx.Serialize()

	//seed := MnemonicToSeed("venture lazy digital aware plug hire acquire abuse chunk know gloom snow much employ glow rich exclude allow", "123")
	//addrkey, _:= GetAddressKey(seed, 0, 0)
	//addr, _ := GetAddress(seed,0,0, 0x35)
	//fmt.Println("addr is: ", addr)
	//wif, _ := ExportWif(seed, 0, 0)
	//fmt.Println("wif is: ", wif)

	fmt.Println("chain_id:", chain_id)
	fmt.Println("res:", hex.EncodeToString(res))

	chainid_byte, _ := hex.DecodeString(chain_id)
	toSign := sha256.Sum256(append(chainid_byte, res...))
	fmt.Println("wif:", wif, "toSign:", hex.EncodeToString(toSign[:]))

	sig, err := GetSignature(wif, toSign[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	transferTrx.Hx_signatures = append(transferTrx.Hx_signatures, hex.EncodeToString(sig))

	b, err = json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("tx:", string(b))
	return
}
