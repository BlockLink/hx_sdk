/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: serialize.go
 * Date: 2018-09-07
 *
 */

package hx

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// inferface for serialize hx transaction
type HxSearilze interface {
	Serialize() []byte
}

/**
 *  some basic type serialization function
 */
//func PackUint32(writer *bytes.Buffer, val uint32) ([]byte, error) {
//
//	uint64_val := uint64(val)
//
//	for {
//		uint8_val := uint8(uint64_val) & 0x7F
//
//		uint64_val >>= 7
//
//		if uint64_val > 0 {
//			uint8_val |= 0x1 << 7
//		} else {
//			uint8_val |= 0x0 << 7
//		}
//
//		err := writer.WriteByte(uint8_val)
//		if err != nil {
//			return nil, fmt.Errorf("in PackUint32 function, write byte failed: %v", err)
//		}
//
//		if uint64_val == 0 {
//			break
//		}
//
//	}
//
//	return writer.Bytes(), nil
//
//}
//
//
//func UnPackUint32(reader *bytes.Reader) (uint32, error) {
//
//	var uint32_val uint32 = 0
//	var by uint8 = 0
//	for {
//		uint8_val, err := reader.ReadByte()
//		if err != nil {
//			return 0, fmt.Errorf("in UnPackUint32 function, read byte failed: %v", err)
//		}
//
//		uint32_val |= uint32(uint8_val & 0x7F) << by
//
//		by += 7
//
//		if (uint8_val & 0x80) == 0 {
//			break
//		}
//
//	}
//
//	return uint32_val, nil
//}

func PackUint16(val uint16, isLittleEndian bool) []byte {

	res := make([]byte, 2)

	if isLittleEndian {
		binary.LittleEndian.PutUint16(res, val)
	} else {
		binary.BigEndian.PutUint16(res, val)
	}

	return res

}

func UnPackUint16(bytes []byte, isLittleEndian bool) uint16 {

	var res uint16

	if isLittleEndian {
		res = binary.LittleEndian.Uint16(bytes)
	} else {
		res = binary.BigEndian.Uint16(bytes)
	}

	return res
}

func PackUint32(val uint32, isLittleEndian bool) []byte {

	res := make([]byte, 4)

	if isLittleEndian {
		binary.LittleEndian.PutUint32(res, val)
	} else {
		binary.BigEndian.PutUint32(res, val)
	}

	return res

}

func UnPackUint32(bytes []byte, isLittleEndian bool) uint32 {

	var res uint32

	if isLittleEndian {
		res = binary.LittleEndian.Uint32(bytes)
	} else {
		res = binary.BigEndian.Uint32(bytes)
	}

	return res
}

func PackInt64(val int64, isLittleEndian bool) []byte {

	res := make([]byte, 8)

	if isLittleEndian {
		binary.LittleEndian.PutUint64(res, uint64(val))
	} else {
		binary.BigEndian.PutUint64(res, uint64(val))
	}

	return res
}

func UnPackInt64(bytes []byte, isLittleEndian bool) int64 {

	var res int64

	if isLittleEndian {
		res = int64(binary.LittleEndian.Uint64(bytes))
	} else {
		res = int64(binary.BigEndian.Uint64(bytes))
	}

	return res
}

func PackVarUint32(val uint32) []byte {

	res := make([]byte, 0)

	//one byte
	if val < 0x80 {

		res = append(res, byte(val))

		return res
	} else if val < 0x4000 { //two byte

		byte1 := val / 0x80
		byte2 := val%0x80 + 0x80

		res = append(res, byte(byte2))
		res = append(res, byte(byte1))

	} else if val < 0x200000 { //three byte

		byte1 := val / 0x4000
		byte2 := val%0x4000/0x80 + 0x80
		byte3 := val%0x80 + 0x80

		res = append(res, byte(byte3))
		res = append(res, byte(byte2))
		res = append(res, byte(byte1))

	} else if val < 0x10000000 { //four byte

		byte1 := val / 0x200000
		byte2 := val%0x200000/0x4000 + 0x80
		byte3 := val%0x4000/0x80 + 0x80
		byte4 := val%0x80 + 0x80

		res = append(res, byte(byte4))
		res = append(res, byte(byte3))
		res = append(res, byte(byte2))
		res = append(res, byte(byte1))
	} else {

		byte1 := val / 0x10000000
		byte2 := val%0x10000000/0x200000 + 0x80
		byte3 := val%0x200000/0x4000 + 0x80
		byte4 := val%0x4000/0x80 + 0x80
		byte5 := val%0x80 + 0x80

		res = append(res, byte(byte5))
		res = append(res, byte(byte4))
		res = append(res, byte(byte3))
		res = append(res, byte(byte2))
		res = append(res, byte(byte1))

	}

	return res
}

func (asset *Asset) Serialize() []byte {

	byte_int64 := PackInt64(asset.Hx_amount, true)

	//byte for asset_id_type, default to zero
	tmp_id, err := GetId(asset.Hx_asset_id)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 := PackVarUint32(tmp_id)
	byte_int64 = append(byte_int64, byte_uint32...)

	return byte_int64
}

func (memo *Memo) Serialize() []byte {

	if memo == nil {
		return []byte{0}
	} else {

		//byte for optional, have element default to one
		var res []byte
		res = append(res, byte(1))
		byte_pub := make([]byte, 74)
		res = append(res, byte_pub...)
		// memo message
		res = append(res, byte(len(memo.Message)+4))
		byte_pub = make([]byte, 4)
		res = append(res, byte_pub...)
		res = append(res, []byte(memo.Message)...)
		return res

	}

}

func (authority *Authority) Serialize() []byte {

	var res []byte
	res = append(res, PackUint32(authority.Hx_weight_threshold, true)...)
	res = append(res, byte(0))
	res = append(res, byte(len(authority.Hx_key_auths)))
	tmpByte, _ := GetPubkeyBytes(authority.Key_auths)
	res = append(res, tmpByte...)
	res = append(res, PackUint16(1, true)...)
	res = append(res, byte(0))

	return res
}

func (acc *AccountOptions) Serialize() []byte {

	var res []byte
	tmpByte, _ := GetPubkeyBytes(acc.Hx_memo_key)
	res = append(res, tmpByte...)
	res = append(res, byte(5))
	res = append(res, PackUint16(0, true)...)
	res = append(res, PackUint16(0, true)...)
	res = append(res, byte(0))
	res = append(res, byte(10))
	res = append(res, byte(0))

	return res
}

func (tranferOp *TransferOperation) Serialize() []byte {

	res := tranferOp.Hx_fee.Serialize()

	if tranferOp.Hx_guarantee_id != "" {
		res = append(res, byte(1))
		tmp_id, err := GetId(tranferOp.Hx_guarantee_id)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)

		byteTmp := make([]byte, 2)
		res = append(res, byteTmp...)

	} else {
		byteTmp := make([]byte, 3)
		res = append(res, byteTmp...)
	}

	byteTmp, _ := GetAddressBytes(tranferOp.Hx_from_addr)
	res = append(res, byteTmp...)
	byteTmp, _ = GetAddressBytes(tranferOp.Hx_to_addr)
	res = append(res, byteTmp...)

	byteTmp = tranferOp.Hx_amount.Serialize()
	res = append(res, byteTmp...)

	byteTmp = tranferOp.Hx_memo.Serialize()
	res = append(res, byteTmp...)
	res = append(res, byte(0))

	return res

}

func (bindOp *AccountBindOperation) Serialize() []byte {

	res := bindOp.Hx_fee.Serialize()
	res = append(res, byte(len(bindOp.Hx_crosschain_type)))
	res = append(res, []byte(bindOp.Hx_crosschain_type)...)
	tmpByte, _ := GetAddressBytes(bindOp.Hx_addr)
	res = append(res, tmpByte...)
	tmpByte, _ = hex.DecodeString(bindOp.Hx_account_signature)
	res = append(res, tmpByte...)
	res = append(res, byte(len(bindOp.Hx_tunnel_address)))
	res = append(res, []byte(bindOp.Hx_tunnel_address)...)

	tmpByte = PackVarUint32(uint32(len(bindOp.Hx_tunnel_signature)))
	res = append(res, tmpByte...)
	res = append(res, []byte(bindOp.Hx_tunnel_signature)...)

	if bindOp.Hx_guarantee_id != "" {
		res = append(res, byte(1))
		tmp_id, err := GetId(bindOp.Hx_guarantee_id)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)
	} else {
		res = append(res, byte(0))
	}

	return res

}

func (unbindOp *AccountUnBindOperation) Serialize() []byte {

	res := unbindOp.Hx_fee.Serialize()
	res = append(res, byte(len(unbindOp.Hx_crosschain_type)))
	res = append(res, []byte(unbindOp.Hx_crosschain_type)...)
	tmpByte, _ := GetAddressBytes(unbindOp.Hx_addr)
	res = append(res, tmpByte...)
	tmpByte, _ = hex.DecodeString(unbindOp.Hx_account_signature)
	res = append(res, tmpByte...)
	res = append(res, byte(len(unbindOp.Hx_tunnel_address)))
	res = append(res, []byte(unbindOp.Hx_tunnel_address)...)
	tmpByte = PackVarUint32(uint32(len(unbindOp.Hx_tunnel_signature)))
	res = append(res, tmpByte...)
	res = append(res, []byte(unbindOp.Hx_tunnel_signature)...)

	res = append(res, byte(0))

	return res

}

func (withdraw *WithdrawCrosschainOperation) Serialize() []byte {

	var res []byte
	res = append(res, withdraw.Hx_fee.Serialize()...)
	tmpByte, _ := GetAddressBytes(withdraw.Hx_withdraw_account)
	res = append(res, tmpByte...)
	res = append(res, byte(len(withdraw.Hx_amount)))
	res = append(res, []byte(withdraw.Hx_amount)...)
	res = append(res, byte(len(withdraw.Hx_asset_symbol)))
	res = append(res, []byte(withdraw.Hx_asset_symbol)...)

	//byte for asset_id_type, default to zero
	tmp_id, err := GetId(withdraw.Hx_asset_id)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 := PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)

	res = append(res, byte(len(withdraw.Hx_crosschain_account)))
	res = append(res, []byte(withdraw.Hx_crosschain_account)...)
	res = append(res, byte(len(withdraw.Hx_memo)))
	res = append(res, []byte(withdraw.Hx_memo)...)

	return res
}

func (register *RegisterAccountOperation) Serialize() []byte {

	var res []byte
	res = append(res, register.Hx_fee.Serialize()...)

	tmpByte := make([]byte, 2)
	res = append(res, tmpByte...)
	tmpByte = PackUint16(0, true)
	res = append(res, tmpByte...)
	res = append(res, byte(len(register.Hx_name)))
	res = append(res, []byte(register.Hx_name)...)

	res = append(res, register.Hx_owner.Serialize()...)
	res = append(res, register.Hx_active.Serialize()...)

	tmpByte, _ = GetAddressBytes(register.Hx_payer)
	res = append(res, tmpByte...)

	res = append(res, register.Hx_options.Serialize()...)
	res = append(res, byte(0))

	if register.Hx_guarantee_id != "" {
		res = append(res, byte(1))
		tmp_id, err := GetId(register.Hx_guarantee_id)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)
	} else {
		res = append(res, byte(0))
	}

	return res
}

func (lockOp *LockBalanceOperation) Serialize() []byte {

	var res []byte
	tmp_id, err := GetId(lockOp.Hx_lock_asset_id)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 := PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)
	res = append(res, PackInt64(lockOp.Hx_lock_asset_amount, true)...)

	//tmpByte, _ := GetAddressBytes(lockOp.Hx_contract_addr)
	//res = append(res, tmpByte...)
	var invalid_address_byte []byte
	invalid_address_byte = append(invalid_address_byte, byte(0x35))
	tmpByte := make([]byte, 20)
	invalid_address_byte = append(invalid_address_byte, tmpByte...)
	res = append(res, invalid_address_byte...)

	tmp_id, err = GetId(lockOp.Hx_lock_balance_account)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 = PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)

	tmp_id, err = GetId(lockOp.Hx_lockto_miner_account)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 = PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)

	tmpByte, _ = GetAddressBytes(lockOp.Hx_lock_balance_addr)
	res = append(res, tmpByte...)

	res = append(res, lockOp.Hx_fee.Serialize()...)

	return res
}

func (obtainOp *ObtainPaybackOperation) Serialize() []byte {

	var res []byte

	tmpByte, _ := GetAddressBytes(obtainOp.Hx_pay_back_owner)
	res = append(res, tmpByte...)

	res = append(res, byte(len(obtainOp.Hx_pay_back_balance)))
	for i := 0; i < len(obtainOp.Hx_pay_back_balance); i++ {
		res = append(res, byte(len(obtainOp.citizen_name[i])))
		res = append(res, []byte(obtainOp.citizen_name[i])...)
		res = append(res, obtainOp.obtain_asset[i].Serialize()...)
	}

	if obtainOp.Hx_guarantee_id != "" {
		res = append(res, byte(1))
		tmp_id, err := GetId(obtainOp.Hx_guarantee_id)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)
	} else {
		res = append(res, byte(0))
	}

	res = append(res, obtainOp.Hx_fee.Serialize()...)

	return res
}

func (forecloseOp *ForecloseBalanceOperation) Serialize() []byte {

	var res []byte
	res = append(res, forecloseOp.Hx_fee.Serialize()...)

	tmp_id, err := GetId(forecloseOp.Hx_foreclose_asset_id)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 := PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)
	res = append(res, PackInt64(forecloseOp.Hx_foreclose_asset_amount, true)...)

	tmp_id, err = GetId(forecloseOp.Hx_foreclose_miner_account)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 = PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)

	//tmpByte, _ := GetAddressBytes(forecloseOp.Hx_foreclose_contract_addr)
	//res = append(res, tmpByte...)
	var invalid_address_byte []byte
	invalid_address_byte = append(invalid_address_byte, byte(0x35))
	tmpByte := make([]byte, 20)
	invalid_address_byte = append(invalid_address_byte, tmpByte...)
	res = append(res, invalid_address_byte...)

	tmp_id, err = GetId(forecloseOp.Hx_foreclose_account)
	if err != nil {
		fmt.Println(err)
		panic(tmp_id)
	}
	byte_uint32 = PackVarUint32(tmp_id)
	res = append(res, byte_uint32...)
	tmpByte, _ = GetAddressBytes(forecloseOp.Hx_foreclose_addr)
	res = append(res, tmpByte...)

	return res
}

func (contractOp *ContractInvokeOperation) Serialize() []byte {

	var res []byte
	res = append(res, contractOp.Hx_fee.Serialize()...)

	byte_int64 := PackInt64(int64(contractOp.Hx_invoke_cost), true)
	res = append(res, byte_int64...)
	byte_int64 = PackInt64(int64(contractOp.Hx_gas_price), true)
	res = append(res, byte_int64...)

	tmpByte, _ := GetAddressBytes(contractOp.Hx_caller_addr)
	res = append(res, tmpByte...)
	tmpByte, _ = hex.DecodeString(contractOp.Hx_caller_pubkey)
	res = append(res, tmpByte...)
	tmpByte, _ = GetAddressBytes(contractOp.Hx_contract_id)
	res = append(res, tmpByte...)
	res = append(res, byte(len(contractOp.Hx_contract_api)))
	res = append(res, []byte(contractOp.Hx_contract_api)...)
	res = append(res, byte(len(contractOp.Hx_contract_arg)))
	res = append(res, []byte(contractOp.Hx_contract_arg)...)

	if contractOp.Hx_guarantee_id != "" {
		res = append(res, byte(1))
		tmp_id, err := GetId(contractOp.Hx_guarantee_id)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)
	} else {
		res = append(res, byte(0))
	}

	return res
}

func (contractOp *ContractTransferOperation) Serialize() []byte {
	var res []byte
	res = append(res, contractOp.Hx_fee.Serialize()...)

	byte_int64 := PackInt64(int64(contractOp.Hx_invoke_cost), true)
	res = append(res, byte_int64...)
	byte_int64 = PackInt64(int64(contractOp.Hx_gas_price), true)
	res = append(res, byte_int64...)

	tmpByte, _ := GetAddressBytes(contractOp.Hx_caller_addr)
	res = append(res, tmpByte...)
	tmpByte, _ = hex.DecodeString(contractOp.Hx_caller_pubkey)
	res = append(res, tmpByte...)
	tmpByte, _ = GetAddressBytes(contractOp.Hx_contract_id)
	res = append(res, tmpByte...)

	res = append(res, contractOp.Hx_amount.Serialize()...)
	res = append(res, byte(len(contractOp.Hx_param)))
	res = append(res, []byte(contractOp.Hx_param)...)

	if contractOp.Hx_guarantee_id != "" {
		res = append(res, byte(1))
		tmp_id, err := GetId(contractOp.Hx_guarantee_id)
		if err != nil {
			fmt.Println(err)
			panic(tmp_id)
		}
		byte_uint32 := PackVarUint32(tmp_id)
		res = append(res, byte_uint32...)
	} else {
		res = append(res, byte(0))
	}

	return res

}

func (trx *Transaction) Serialize() []byte {

	var res []byte
	res = append(res, PackUint16(trx.Hx_ref_block_num, true)...)
	res = append(res, PackUint32(trx.Hx_ref_block_prefix, true)...)
	res = append(res, PackUint32(trx.Expiration, true)...)

	//operations
	res = append(res, byte(len(trx.Operations)))
	for _, v := range trx.Operations {

		if transferOp, ok := v.(TransferOperation); ok {
			res = append(res, byte(0))
			res = append(res, transferOp.Serialize()...)
		} else if bindOp, ok := v.(AccountBindOperation); ok {
			res = append(res, byte(10))
			res = append(res, bindOp.Serialize()...)
		} else if unbindOp, ok := v.(AccountUnBindOperation); ok {
			res = append(res, byte(11))
			res = append(res, unbindOp.Serialize()...)
		} else if withdrawOp, ok := v.(WithdrawCrosschainOperation); ok {
			res = append(res, byte(61))
			res = append(res, withdrawOp.Serialize()...)
		} else if registerOp, ok := v.(RegisterAccountOperation); ok {
			res = append(res, byte(5))
			res = append(res, registerOp.Serialize()...)
		} else if lockOp, ok := v.(LockBalanceOperation); ok {
			res = append(res, byte(55))
			res = append(res, lockOp.Serialize()...)
		} else if forecloseOp, ok := v.(ForecloseBalanceOperation); ok {
			res = append(res, byte(56))
			res = append(res, forecloseOp.Serialize()...)
		} else if obtainOp, ok := v.(ObtainPaybackOperation); ok {
			res = append(res, byte(73))
			res = append(res, obtainOp.Serialize()...)
		} else if contractOp, ok := v.(ContractInvokeOperation); ok {
			res = append(res, byte(79))
			res = append(res, contractOp.Serialize()...)
		} else if contractOp, ok := v.(ContractTransferOperation); ok {
			res = append(res, byte(81))
			res = append(res, contractOp.Serialize()...)
		}

	}

	//extension
	res = append(res, byte(0))

	//signature
	if len(trx.Hx_signatures) > 0 {
		res = append(res, byte(len(trx.Hx_signatures)))
	}

	return res
}
