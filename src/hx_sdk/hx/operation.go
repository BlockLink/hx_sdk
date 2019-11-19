/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: operation.go.go, Date: 2018-10-31
 *
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package hx

type Asset struct {
	Hx_amount   int64  `json:"amount"`
	Hx_asset_id string `json:"asset_id"`
}

//
// hx  --- "1.3.0"
// btc --- "1.3.1"
// ltc --- "1.3.2"
// hc  --- "1.3.3"
func DefaultAsset() Asset {
	return Asset{
		0,
		"1.3.0",
	}
}

type Extension struct {
	extension []string
}

type Memo struct {
	Hx_from    string `json:"from"` //public_key_type  33
	Hx_to      string `json:"to"`   //public_key_type  33
	Hx_nonce   uint64 `json:"nonce"`
	Hx_message string `json:"message"`

	IsEmpty bool   `json:"-"`
	Message string `json:"-"`
}

func DefaultMemo() Memo {

	return Memo{
		"HX1111111111111111111111111111111114T1Anm",
		"HX1111111111111111111111111111111114T1Anm",
		0,
		"",
		true,
		"",
	}

}

type Authority struct {
	Hx_weight_threshold uint32          `json:"weight_threshold"`
	Hx_account_auths    []interface{}   `json:"account_auths"`
	Hx_key_auths        [][]interface{} `json:"key_auths"`
	Hx_address_auths    []interface{}   `json:"address_auths"`

	Key_auths string `json:"-"`
}

func DefaultAuthority() Authority {

	return Authority{
		1,
		[]interface{}{},
		[][]interface{}{{"", 1}},
		[]interface{}{},
		"",
	}
}

type AccountOptions struct {
	Hx_memo_key              string        `json:"memo_key"`
	Hx_voting_account        string        `json:"voting_account"`
	Hx_num_witness           uint16        `json:"num_witness"`
	Hx_num_committee         uint16        `json:"num_committee"`
	Hx_votes                 []interface{} `json:"votes"`
	Hx_miner_pledge_pay_back byte          `json:"miner_pledge_pay_back"`
	Hx_extensions            []interface{} `json:"extensions"`
}

func DefaultAccountOptions() AccountOptions {

	return AccountOptions{
		"",
		"1.2.5",
		0,
		0,
		[]interface{}{},
		10,
		[]interface{}{},
	}

}

// transfer operation tag is  0
type TransferOperation struct {
	Hx_fee          Asset  `json:"fee"`
	Hx_guarantee_id string `json:"guarantee_id,omitempty"`
	Hx_from         string `json:"from"`
	Hx_to           string `json:"to"`

	Hx_from_addr string `json:"from_addr"`
	Hx_to_addr   string `json:"to_addr"`

	Hx_amount Asset `json:"amount"`
	Hx_memo   *Memo `json:"memo,omitempty"`

	Hx_extensions []interface{} `json:"extensions"`
}

func DefaultTransferOperation() *TransferOperation {

	return &TransferOperation{
		DefaultAsset(),
		"",
		"1.2.0",
		"1.2.0",
		"",
		"",
		DefaultAsset(),
		nil,
		make([]interface{}, 0),
	}
}

// account bind operation tag is 10
type AccountBindOperation struct {
	Hx_fee               Asset  `json:"fee"`
	Hx_crosschain_type   string `json:"crosschain_type"`
	Hx_addr              string `json:"addr"`
	Hx_account_signature string `json:"account_signature"`
	Hx_tunnel_address    string `json:"tunnel_address"`
	Hx_tunnel_signature  string `json:"tunnel_signature"`
	Hx_guarantee_id      string `json:"guarantee_id,omitempty"`
}

func DefaultAccountBindOperation() *AccountBindOperation {

	return &AccountBindOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
		"",
	}
}

// account unbind operation tag is 11
type AccountUnBindOperation struct {
	Hx_fee               Asset  `json:"fee"`
	Hx_crosschain_type   string `json:"crosschain_type"`
	Hx_addr              string `json:"addr"`
	Hx_account_signature string `json:"account_signature"`
	Hx_tunnel_address    string `json:"tunnel_address"`
	Hx_tunnel_signature  string `json:"tunnel_signature"`
}

func DefaultAccountUnBindOperation() *AccountUnBindOperation {

	return &AccountUnBindOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
	}
}

// withdraw cross chain operation tag is 61
type WithdrawCrosschainOperation struct {
	Hx_fee              Asset  `json:"fee"`
	Hx_withdraw_account string `json:"withdraw_account"`
	Hx_amount           string `json:"amount"`
	Hx_asset_symbol     string `json:"asset_symbol"`

	Hx_asset_id           string `json:"asset_id"`
	Hx_crosschain_account string `json:"crosschain_account"`
	Hx_memo               string `json:"memo"`
}

func DefaultWithdrawCrosschainOperation() *WithdrawCrosschainOperation {

	return &WithdrawCrosschainOperation{
		DefaultAsset(),
		"",
		"",
		"",
		"",
		"",
		"",
	}
}

//register account operation tag is 5
type RegisterAccountOperation struct {
	Hx_fee              Asset     `json:"fee"`
	Hx_registrar        string    `json:"registrar"`
	Hx_referrer         string    `json:"referrer"`
	Hx_referrer_percent uint16    `json:"referrer_percent"`
	Hx_name             string    `json:"name"`
	Hx_owner            Authority `json:"owner"`
	Hx_active           Authority `json:"active"`
	Hx_payer            string    `json:"payer"`

	Hx_options      AccountOptions `json:"options"`
	Hx_extensions   interface{}    `json:"extensions"`
	Hx_guarantee_id string         `json:"guarantee_id,omitempty"`
}

func DefaultRegisterAccountOperation() *RegisterAccountOperation {

	return &RegisterAccountOperation{
		DefaultAsset(),
		"1.2.0",
		"1.2.0",
		0,
		"",
		DefaultAuthority(),
		DefaultAuthority(),
		"",

		DefaultAccountOptions(),
		make(map[string]interface{}, 0),
		"",
	}

}

//lock balance operation tag is 55
type LockBalanceOperation struct {
	Hx_lock_asset_id     string `json:"lock_asset_id"`
	Hx_lock_asset_amount int64  `json:"lock_asset_amount"`
	Hx_contract_addr     string `json:"contract_addr"`

	Hx_lock_balance_account string `json:"lock_balance_account"`
	Hx_lockto_miner_account string `json:"lockto_miner_account"`
	Hx_lock_balance_addr    string `json:"lock_balance_addr"`

	Hx_fee Asset `json:"fee"`
}

func DefaultLockBalanceOperation() *LockBalanceOperation {

	return &LockBalanceOperation{
		"1.3.0",
		0,
		"",
		"",
		"",
		"",
		DefaultAsset(),
	}
}

//foreclose balance operation tag is 56
type ForecloseBalanceOperation struct {
	Hx_fee Asset `json:"fee"`

	Hx_foreclose_asset_id     string `json:"foreclose_asset_id"`
	Hx_foreclose_asset_amount int64  `json:"foreclose_asset_amount"`

	Hx_foreclose_miner_account string `json:"foreclose_miner_account"`
	Hx_foreclose_contract_addr string `json:"foreclose_contract_addr"`

	Hx_foreclose_account string `json:"foreclose_account"`
	Hx_foreclose_addr    string `json:"foreclose_addr"`
}

func DefaultForecloseBalanceOperation() *ForecloseBalanceOperation {

	return &ForecloseBalanceOperation{
		DefaultAsset(),
		"1.3.0",
		0,
		"",
		"",
		"",
		"",
	}
}

//obtain pay back operation tag is 73
type ObtainPaybackOperation struct {
	Hx_pay_back_owner   string          `json:"pay_back_owner"`
	Hx_pay_back_balance [][]interface{} `json:"pay_back_balance"`
	Hx_guarantee_id     string          `json:"guarantee_id,omitempty"`
	Hx_fee              Asset           `json:"fee"`

	citizen_name []string
	obtain_asset []Asset
}

func DefaultObtainPaybackOperation() *ObtainPaybackOperation {

	return &ObtainPaybackOperation{
		"",
		[][]interface{}{{"", DefaultAsset()}},
		"",
		DefaultAsset(),
		nil,
		nil,
	}
}

// contract invoke operation tag is 79
type ContractInvokeOperation struct {
	Hx_fee           Asset  `json:"fee"`
	Hx_invoke_cost   uint64 `json:"invoke_cost"`
	Hx_gas_price     uint64 `json:"gas_price"`
	Hx_caller_addr   string `json:"caller_addr"`
	Hx_caller_pubkey string `json:"caller_pubkey"`
	Hx_contract_id   string `json:"contract_id"`
	Hx_contract_api  string `json:"contract_api"`
	Hx_contract_arg  string `json:"contract_arg"`
	//Hx_extension     []interface{} `json:"extensions"`
	Hx_guarantee_id string `json:"guarantee_id,omitempty"`
}

func DefaultContractInvokeOperation() *ContractInvokeOperation {

	return &ContractInvokeOperation{
		DefaultAsset(),
		0,
		0,
		"",
		"",
		"",
		"",
		"",
		//make([]interface{}, 0),
		"",
	}
}

// transfer to contract operation tag is 81
type ContractTransferOperation struct {
	Hx_fee           Asset  `json:"fee"`
	Hx_invoke_cost   uint64 `json:"invoke_cost"`
	Hx_gas_price     uint64 `json:"gas_price"`
	Hx_caller_addr   string `json:"caller_addr"`
	Hx_caller_pubkey string `json:"caller_pubkey"`
	Hx_contract_id   string `json:"contract_id"`
	Hx_amount        Asset  `json:"amount"`
	Hx_param         string `json:"param"`
	//Hx_extension     []interface{} `json:"extensions"`
	Hx_guarantee_id string `json:"guarantee_id,omitempty"`
}

func DefaultContractTransferOperation() *ContractTransferOperation {

	return &ContractTransferOperation{
		DefaultAsset(),
		0,
		0,
		"",
		"",
		"",
		DefaultAsset(),
		"",
		//make([]interface{}, 0),
		"",
	}
}
