package hx_sdk

import (
	"fmt"
	"testing"

	"./hx"
	"bytes"
	"crypto/tls"
	ghex "encoding/hex"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
)

var (
	//hxuri      = "http://192.168.1.220/server/process"
	hxuri     = "http://192.168.1.220/server/process"
	walletURI = "http://192.168.1.122:10046" // broadcast wallet rpc api
	// zhengqinpeng wid
	hxwalletid = "dcc02a08142e230814e108053eb40c6180f908da"
	//myhxaddr   = "HXNhNthqhgkEfPjzhLQ3cWnNEpEjpjtvmKzw"
	//myhcaddr   = "TsTaA9DT2z24Wg4U31QLeSvcg1agu6e4d88"
	// zhengqinpeng address
	myhxaddr = "" // "HXNWjKv1PUbZ6dgVoerkohdoCHRck9LAZh3Y"
	// myhxaddr = "HXNXPFrzwsTo7wT5QceDLEjpd2VZ8BZR3UCY"  // just for bind usdt
	// wifhx = "5KZTntQJ9AAARnuKUzgUG4MHiPPVU7FuMq4d3WNH4rhasDxRa4a" // for bind usdt

	myltcaddr = "LP2JKjy9WmSygMdoe2CzEHabXrPSPXspNF" // "LQz1gokPokkj6hK3dHt1JphZbPh5G9KpbV"
	dsthxaddr = "HXNL3kJ4prkHUsHsnwW4HGSDUfYxncWcfgDn"

	mycointype = "HX"                                  //  "HC(HX)"
	myhcaddr   = "TsZj5Cx1p94izYoxVuZyRGvskSdGaasCHir" // hc testnet
	myxwif     = "PtWUvDcPuchqQA5vfoUqoNJw6JyCSnYrbHM37XXVUbfLT4wdJj7hU"

	//mycointype = "ETH(HX)"
	//myhcaddr  = "0x1891025831596418915523e786334b2b44985272"
	//myxwif = "3f2153c638e857ae4b5ef132c1ee09c24bb48484d2dea91a5071b202be2e2a90" // eth

	//mycointype = "PAX(HX)"
	//myhcaddr  = "0x1891025831596418915523e786334b2b44985272"
	//myxwif = "3f2153c638e857ae4b5ef132c1ee09c24bb48484d2dea91a5071b202be2e2a90" // eth

	widbtc       = "a8eb57aacdb63a6ed8d485c0304260d7e627d704"
	cointypebtc  = "BTC(HX)"
	coinaddrbtc  = "1PNvmFKPGDADmrPXQcLVhyqFPvSZ1czHa"
	wifbtc       = "Kx9XqhpA2UTJm21HgVc3yuvcGMQNfyKi133MEUNGhb6NvNniLzH8" // btc
	cointypeusdt = "USDT(HX)"
	coinaddrusdt = "1PNvmFKPGDADmrPXQcLVhyqFPvSZ1czHa"
	wifusdt      = "L5QrMWSdpU8CQer7ss5zofTzmzgcB7spmBSSoXF29NKNZVYKxwkB"

	//mycointype = "LTC(HX)"
	//myhcaddr  = "LiSt5dsfy7WB5Yq1AUaxjjHPpNQU8upKES"
	//myxwif = "6uHv9Ru7dFrv1easCC1sRwvFSf9gFm2FUpBC1P2TCT9pVQyUZCo" // ltc

	//mycointype = "LTC(HX)"
	//myhcaddr  = "LP2JKjy9WmSygMdoe2CzEHabXrPSPXspNF"
	//myxwif = "T5LN3q4o3CUtKpoi3TtiL4JMPNjqaBhkoWzQJ2F4HLWuu94Y6G8y" // ltc

	hxreq = &AddrReq{
		// Mnemonic: "beyond stage sleep clip because twist token leaf atom beauty genius food business side grid unable middle armed observe pair crouch tonight away coconut",
		// zhengqinpeng mnemonic1
		// Mnemonic: "clip fox defense river cigar love sword keen omit reward keep since donkey unlock unique flip hobby apple",
		Mnemonic: "hobby actual sadness know copy achieve bulb message unhappy snack giggle core reason enroll boat magic aim sea front capital text science green joy",
		// bind usdt
		Net: "testnet",
	}
)

const (
	CoinHX string = "HX"

	VersionNormalAddr   = 0x35
	VersionMultisigAddr = 0x32
	VersionContractAddr = 0x1c
)

func init() {
	resp, err := GetAddress(hxreq, "HX")
	if err != nil {
		panic(err.Error())
	}
	// fmt.Printf("derived hx address is %v\n", resp.Address)
	myhxaddr = resp.Address
}

func testHXAddress(t *testing.T) {

	var (
		//mnemonic = "beyond stage sleep clip because twist token leaf atom beauty genius food business side grid unable middle armed observe pair crouch tonight away coconut"
		// zhengqinpeng 1
		mnemonic = "clip fox defense river cigar love sword keen omit reward keep since donkey unlock unique flip hobby apple"

		req AddrReq
	)
	req.Mnemonic = mnemonic
	req.Net = "testnet"

	for i := 0; i < 5; i++ {
		req.Index = i
		resp, err := GetAddress(&req, "HX")
		assert.Nil(t, err)
		fmt.Printf("HX index %d Address: %v\n", i, resp.Address)
	}
}

func testHXPubkey(t *testing.T) {
	var (
		mnemonic = "beyond stage sleep clip because twist token leaf atom beauty genius food business side grid unable middle armed observe pair crouch tonight away coconut"
		req      AddrReq
	)
	req.Mnemonic = mnemonic
	req.Net = "testnet"
	seed := hx.MnemonicToSeed(req.Mnemonic, "")
	wif, err := hx.ExportWif(seed, uint32(req.Account), uint32(req.Index))
	assert.Nil(t, err)

	pk, err := hx.DerivePubkey(wif)
	assert.Nil(t, err)
	fmt.Println(pk)

	// test
	wif = "5KR6ocp5eUdWWYPX7mYp4XLGBcZ2xHVHVsNaco6K2YZSWQTqES7"
	pk, err = hx.DerivePubkey(wif)
	assert.Nil(t, err)
	fmt.Println(pk)
	assert.Equal(t, "HX5cj54KW1TDK94GCcrSovYDUUYQ1FgCbECnG2CaEv5ub7Vghx5N", pk)
}

func _doReq(t *testing.T, uri string, ms map[string]interface{}) []byte {
	body, err := json.Marshal(ms)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	assert.Nil(t, err)

	req.Header.Add("content-type", "application/json")
	req.Header.Set("Connection", "close")
	req.Close = true
	// req.

	// 跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, res.StatusCode, 200)

	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	trancode := ""
	iheader := ms["header"]
	if iheader != nil {
		if header, ok := iheader.(map[string]interface{}); ok {
			itrancode, ok := header["trancode"]
			if ok {
				trancode = itrancode.(string)
			}
		}
	}
	fmt.Printf("request %s %s: req body=%v\ncode=%d response=%v\n", uri, trancode, string(body), res.StatusCode, string(buf))
	assert.Nil(t, err)

	return buf
}

func doReq(t *testing.T, ms map[string]interface{}) []byte {
	return _doReq(t, hxuri, ms)
}

func makeheader(transcode, walletid string) map[string]interface{} {
	m := map[string]interface{}{}
	m["header"] = map[string]string{
		"version":    "2.4.0",
		"language":   "zh-Hans",
		"trancode":   transcode,
		"clienttype": "Android",
		"walletid":   walletid,
		"random":     "abc",
		"handshake":  "efg",
		"imie":       "3511",
		"source":     "sdk",
	}
	m["body"] = map[string]interface{}{}

	return m
}

func setBodyData(t *testing.T, dst map[string]interface{}, buf []byte) {
	var m map[string]interface{}

	err := json.Unmarshal(buf, &m)
	assert.Nil(t, err)

	dst["data"] = m
}

func directPostWallet(t *testing.T, action string, buf []byte) {
	_directPostWallet(t, walletURI, buf)
}

func _directPostWallet(t *testing.T, uri string, buf []byte) []byte {
	// uri := walletURI // "http://192.168.1.128:10033"
	param := map[string]interface{}{
		"id":     1,
		"method": "lightwallet_broadcast",
	}
	var m map[string]interface{}
	err := json.Unmarshal(buf, &m)
	assert.Nil(t, err)

	param["params"] = []interface{}{m}
	res := _doReq(t, uri, param)
	return res
	// fmt.Printf("action %v response: %v\n", action, string(res))
}

func testUnbind(t *testing.T, wallet bool) {
	ms1 := makeheader("hxbinding_init", hxwalletid)
	bm := ms1["body"].(map[string]interface{})
	bm["origAddr"] = myhcaddr
	bm[fieldCoinAddr] = myhxaddr
	bm[fieldCoinType] = mycointype // "HC(HX)"
	buf1 := doReq(t, ms1)
	fmt.Println("bind response:", string(buf1))

	xwif, err := DumpPrivKey(hxreq, "HC")
	assert.Nil(t, err)
	_ = xwif

	content := map[string]string{
		"origAddr":    myhcaddr,
		fieldCoinAddr: myhxaddr,
		"coinType":    mycointype, //"HC(HX)",
		fieldCrossWif: myxwif,     // xwif,
	}

	cbuf, err := json.Marshal(content)
	assert.Nil(t, err)

	hex, err := HXTransaction(hxreq, HXUbind, buf1, cbuf)
	assert.Nil(t, err)

	directPostWallet(t, "unbind", hex)
	return
}

// test bind BTC
// test bind USDT, address is same with BTC
func testBindBTCUSDT(t *testing.T, wallet bool) {
	bindAddress(t, cointypebtc, myhxaddr, coinaddrbtc, wifbtc, hxwalletid, true)
	bindAddress(t, cointypeusdt, myhxaddr, coinaddrusdt, wifusdt, hxwalletid, true)
}

func bindAddress(t *testing.T, cointype string, hxaddr, coinaddr string, wif, walletid string, wallet bool) {
	ms1 := makeheader("hxbinding_init", hxwalletid)
	bm := ms1["body"].(map[string]interface{})
	bm["origAddr"] = coinaddr
	bm[fieldCoinAddr] = hxaddr
	bm[fieldCoinType] = cointype
	buf1 := doReq(t, ms1)
	fmt.Printf("bind %s %s request response: %s\n", cointype, coinaddr, string(buf1))

	content := map[string]string{
		"origAddr":    myhcaddr,
		fieldCoinAddr: myhxaddr,
		"coinType":    cointype,
		fieldCrossWif: wif,
	}

	cbuf, err := json.Marshal(content)
	assert.Nil(t, err)

	hex, err := HXTransaction(hxreq, "bind", buf1, cbuf)
	assert.Nil(t, err)

	ms2 := makeheader("hxbinding_send", hxwalletid)
	bm2 := ms2["body"].(map[string]interface{})

	delete(content, fieldCrossWif)
	bm2["tranContent"] = content
	// bm2["data"] = string(hex)
	setBodyData(t, bm2, hex)
	fmt.Println("broadcast bind data:", string(hex))
	fmt.Println("broadcast whole data:", ms2)

	if wallet {
		directPostWallet(t, "bind", hex)
		return
	}
	buf2 := doReq(t, ms2)
	fmt.Printf("bind %s %s response: %s\n", cointype, coinaddr, string(buf2))
}

func testHXBind(t *testing.T, wallet bool) {
	cointype := mycointype // "LTC(HX)"

	ms1 := makeheader("hxbinding_init", hxwalletid)
	bm := ms1["body"].(map[string]interface{})
	bm["origAddr"] = myhcaddr
	bm[fieldCoinAddr] = myhxaddr
	bm[fieldCoinType] = cointype
	buf1 := doReq(t, ms1)
	fmt.Println("bind response:", string(buf1))

	//xwif, err := DumpPrivKey(hxreq, "HC")
	//assert.Nil(t, err)
	//xwif := "PtWUvDcPuchqQA5vfoUqoNJw6JyCSnYrbHM37XXVUbfLT4wdJj7hU"
	xwif := myxwif // "T5LN3q4o3CUtKpoi3TtiL4JMPNjqaBhkoWzQJ2F4HLWuu94Y6G8y" // "6vVgBExZTVUtrYFPLj1PNfPCtXDqnLS6Ki8oxykMpuTS4hk8CPv"
	//cointype := "HC(HX)"

	content := map[string]string{
		"origAddr":    myhcaddr,
		fieldCoinAddr: myhxaddr,
		"coinType":    cointype,
		fieldCrossWif: xwif,
	}

	cbuf, err := json.Marshal(content)
	assert.Nil(t, err)

	hex, err := HXTransaction(hxreq, "bind", buf1, cbuf)
	assert.Nil(t, err)

	ms2 := makeheader("hxbinding_send", hxwalletid)
	bm2 := ms2["body"].(map[string]interface{})

	delete(content, fieldCrossWif)
	bm2["tranContent"] = content
	// bm2["data"] = string(hex)
	setBodyData(t, bm2, hex)
	fmt.Println("broadcast bind data:", string(hex))
	fmt.Println("broadcast whole data:", ms2)

	if wallet {
		directPostWallet(t, "bind", hex)
		return
	}
	buf2 := doReq(t, ms2)
	fmt.Println("bind response:", string(buf2))
}

func getFieldIn(t *testing.T, buf []byte, f string) interface{} {
	var j map[string]interface{}

	err := json.Unmarshal(buf, &j)
	assert.Nil(t, err)

	j2, ok := j["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.True(t, ok)

	i3, ok := j2[f]
	assert.True(t, ok)
	return i3
}

func getFieldString(t *testing.T, buf []byte, f string) string {
	i3 := getFieldIn(t, buf, f)

	s3, ok := i3.(string)
	assert.True(t, ok)
	return s3
}

func testTransfer(t *testing.T, ct string, wallet bool) {
	ms1 := makeheader("hx_init", hxwalletid)
	bm := ms1["body"].(map[string]interface{})
	bm[fieldCoinAddr] = myhxaddr
	bm[fieldCoinType] = ct
	buf1 := doReq(t, ms1)
	fmt.Println("hx_init response:", string(buf1))

	content := map[string]string{
		"fromAddr": myhxaddr,
		"toAddr":   dsthxaddr,
		"coinType": ct,
		"tranAmt":  "1.005",
		"tranFee":  getFieldString(t, buf1, "tranFee"),
		"bak":      "",
	}

	cbuf, err := json.Marshal(content)
	assert.Nil(t, err)

	hex, err := HXTransaction(hxreq, "transfer", buf1, cbuf)
	assert.Nil(t, err)

	ms2 := makeheader("hx_send", hxwalletid)
	bm2 := ms2["body"].(map[string]interface{})
	bm2["tranContent"] = content
	// bm2["data"] = string(hex)
	setBodyData(t, bm2, hex)
	fmt.Println("broadcast hx_send data:", string(hex))

	if wallet {
		directPostWallet(t, "bind", hex)
		return
	}
	buf2 := doReq(t, ms2)
	fmt.Println("hx_send response:", string(buf2))
}

func testWithdraw(t *testing.T, wallet bool) {
	ms1 := makeheader("withdraw_init", hxwalletid)
	bm := ms1["body"].(map[string]interface{})
	bm[fieldCoinAddr] = myhxaddr
	bm[fieldCoinType] = mycointype //"HC(HX)"
	buf1 := doReq(t, ms1)
	fmt.Println("withdraw_init response:", string(buf1))

	content := map[string]string{
		"coinType": mycointype, // CoinHX,
		"fromAddr": myhxaddr,
		"toAddr":   myhcaddr,
		"tranAmt":  "5.5",
		"tranFee":  getFieldString(t, buf1, "tranFee"),
	}

	cbuf, err := json.Marshal(content)
	assert.Nil(t, err)

	hex, err := HXTransaction(hxreq, "withdraw", buf1, cbuf)
	assert.Nil(t, err)
	ms2 := makeheader("withdraw_send", hxwalletid)
	bm2 := ms2["body"].(map[string]interface{})
	bm2["tranContent"] = content
	// bm2["data"] = string(hex)
	setBodyData(t, bm2, hex)
	fmt.Println("broadcast withdraw_send data:", string(hex))

	if wallet {
		directPostWallet(t, "bind", hex)
		return
	}

	buf2 := doReq(t, ms2)
	fmt.Println("withdraw_send response:", string(buf2))
}

func testRegister(t *testing.T, wallet bool, name string) {
	ms1 := makeheader("register_init", hxwalletid)
	bm := ms1["body"].(map[string]interface{})
	bm[fieldCoinAddr] = myhxaddr
	bm[fieldCoinType] = "HX"
	buf1 := doReq(t, ms1)
	fmt.Println("register_init response:", string(buf1))

	content := map[string]string{
		"coinType":    mycointype, // CoinHX,
		"tranFee":     getFieldString(t, buf1, "tranFee"),
		"accountName": name,
		"coinAddr":    myhxaddr,
	}

	cbuf, err := json.Marshal(content)
	assert.Nil(t, err)

	hex, err := HXTransaction(hxreq, "register", buf1, cbuf)
	assert.Nil(t, err)
	ms2 := makeheader("register_send", hxwalletid)
	bm2 := ms2["body"].(map[string]interface{})
	bm2["tranContent"] = content
	//bm2["data"] = string(hex)
	setBodyData(t, bm2, hex)
	fmt.Println("broadcast register_send data:", string(hex))

	if wallet {
		directPostWallet(t, "register", hex)
		return
	}
	buf2 := doReq(t, ms2)
	fmt.Println("register_send response:", string(buf2))
}

func testMining(t *testing.T, wallet bool) {
	ms1 := makeheader("mortgage_init", hxwalletid)
	bm := ms1["body"].(map[string]interface{})
	bm[fieldCoinAddr] = myhxaddr
	bm[fieldCoinType] = CoinHX
	// bm["citizenName"] = "a0"
	buf1 := doReq(t, ms1)
	// fmt.Println("mortgage_init response:", string(buf1))

	citizens := citizenList(t, buf1)
	iassets := getFieldIn(t, buf1, "asset").(map[string]interface{})
	num := iassets["num"].(float64)
	fmt.Println("mortgage asset to:", citizens[0])

	content := map[string]interface{}{
		"coinType": CoinHX,
		"coinAddr": myhxaddr,
		"citizen":  citizens[0],
		"tranAmt":  fmt.Sprint(num - 10),
		"tranFee":  getFieldString(t, buf1, "tranFee"),
	}

	cbuf, err := json.Marshal(content)
	assert.Nil(t, err)

	hex, err := HXTransaction(hxreq, HXMining, buf1, cbuf)
	fmt.Println("broadcast mortgage_send data:", string(hex))
	assert.Nil(t, err)
	ms2 := makeheader("mortgage_send", hxwalletid)
	bm2 := ms2["body"].(map[string]interface{})
	bm2["tranContent"] = content
	// bm2["data"] = string(hex)
	setBodyData(t, bm2, hex)

	if wallet {
		directPostWallet(t, "mortgage_send", hex)
		return
	}
	buf2 := doReq(t, ms2)
	fmt.Println("mortgage_send response:", string(buf2))
}

func testRewards(t *testing.T, wallet bool) {
	ms1 := makeheader("income_init", hxwalletid)
	bm := ms1["body"].(map[string]interface{})
	bm[fieldCoinAddr] = myhxaddr
	bm[fieldCoinType] = CoinHX
	buf1 := doReq(t, ms1)
	fmt.Println("income_init response:", string(buf1))

	content := map[string]string{
		"coinAddr": myhxaddr,
		"coinType": CoinHX,
		"tranAmt":  getFieldString(t, buf1, "tranAmt"),
		"tranFee":  getFieldString(t, buf1, "tranFee"),
	}

	cbuf, err := json.Marshal(content)
	assert.Nil(t, err)

	hex, err := HXTransaction(hxreq, HXRewards, buf1, cbuf)
	assert.Nil(t, err)

	ms2 := makeheader("income_send", hxwalletid)
	bm2 := ms2["body"].(map[string]interface{})
	bm2["tranContent"] = content
	// bm2["data"] = string(hex)
	setBodyData(t, bm2, hex)
	fmt.Println("broadcast income_send data:", string(hex))

	if wallet {
		directPostWallet(t, "bind", hex)
		return
	}
	buf2 := doReq(t, ms2)
	fmt.Println("income_send response:", string(buf2))
}

func testRedeem(t *testing.T, wallet bool) {
	citizens := mortgageList(t, myhxaddr)

	ms1 := makeheader("redeem_init", hxwalletid)
	bm := ms1["body"].(map[string]interface{})
	bm[fieldCoinAddr] = myhxaddr
	bm[fieldCoinType] = "HX"
	bm["citizenName"] = citizens[0]["citizenName"]
	fmt.Println("redeem citizen", citizens[0])

	buf1 := doReq(t, ms1)
	fmt.Println("redeem_init response:", string(buf1))

	content := map[string]interface{}{
		"coinAddr": myhxaddr,
		"citizen":  citizens[0],
		"coinType": CoinHX,
		"tranAmt":  fmt.Sprint(citizens[0]["amount"]),
		"tranFee":  getFieldString(t, buf1, "tranFee"),
	}

	cbuf, err := json.Marshal(content)
	assert.Nil(t, err)

	hex, err := HXTransaction(hxreq, HxRedeem, buf1, cbuf)
	assert.Nil(t, err)

	ms2 := makeheader("redeem_send", hxwalletid)
	bm2 := ms2["body"].(map[string]interface{})
	bm2["tranContent"] = content
	// bm2["data"] = string(hex)
	setBodyData(t, bm2, hex)
	fmt.Println("broadcast redeem_send data:", string(hex))

	if wallet {
		directPostWallet(t, "bind", hex)
		return
	}
	buf2 := doReq(t, ms2)
	fmt.Println("redeem_send response:", string(buf2))
}

func citizenList(t *testing.T, buf []byte) []map[string]string {
	item := getFieldIn(t, buf, "citizenList")
	ilist, ok := item.([]interface{})
	assert.True(t, ok)
	ml := []map[string]string{}

	for _, list := range ilist {
		v, ok := list.(map[string]interface{})
		assert.True(t, ok)
		ml = append(ml, map[string]string{
			"poolFee":           fmt.Sprint(v["poolFee"].(float64)),
			"participationRate": v["participationRate"].(string),
			"citizenId":         v["citizenId"].(string),
			"citizenAddress":    v["citizenAddress"].(string),
			"citizenName":       v["citizenName"].(string),
		})
	}

	return ml
}

func mortgageList(t *testing.T, addr string) []map[string]string {
	ms1 := makeheader("mortgage_list", hxwalletid)
	bm := ms1["body"].(map[string]interface{})
	bm[fieldCoinAddr] = addr

	buf1 := doReq(t, ms1)
	/*
		mortgageList	[]
		coinType
		logoUrl
		amount
		citizenId
		citizenAddress
		citizenName
	*/

	item := getFieldIn(t, buf1, "mortgageList")
	ilist, ok := item.([]interface{})
	assert.True(t, ok)
	ml := []map[string]string{}

	for _, list := range ilist {
		v, ok := list.(map[string]interface{})
		assert.True(t, ok)
		ml = append(ml, map[string]string{
			"coinType":       v["coinType"].(string),
			"logoUrl":        v["logoUrl"].(string),
			"amount":         fmt.Sprint(v["amount"]), //.(string),
			"citizenId":      v["citizenId"].(string),
			"citizenAddress": v["citizenAddress"].(string),
			"citizenName":    v["citizenName"].(string),
		})
	}

	return ml
}

// 2018-11-13 withdraw & rewards not test yet
// 2018-11-13 guarantee not test yet
func testHXTransaction(t *testing.T) {
	wallet := true

	hx.SetTestnetEthSig()
	// testHXContractInvoke(t, wallet)
	//testBindBTCUSDT(t, wallet)
	//testHXBind(t, wallet)
	//testUnbind(t, wallet)
	//testTransfer(t, mycointype, wallet)
	//testRegister(t, wallet, "a0000")
	//testMining(t, wallet)
	//testRedeem(t, wallet)
	//testRewards(t, wallet)
	//testWithdraw(t, wallet)
	_ = wallet
}

func testValidateAddress2(t *testing.T) {
	var addres = []string{
		"HXNZsWKfyfdTiDQhMtqxN2PhfTaLVw9AXutc",
		"HXNfHB9wQdJx7mhWDXam7gNsvyXd7j1yh1M3",
		"HXNWjKv1PUbZ6dgVoerkohdoCHRck9LAZh3Y",
		"HXNL3kJ4prkHUsHsnwW4HGSDUfYxncWcfgDn",
	}

	for _, addr := range addres {
		ok := hx.ValidateAddress(addr, "testnet")
		assert.True(t, ok)
	}
}

func testConvertHXSymbol(t *testing.T) {
	vs := []string{
		"Eth", "ETH", "ETH(HX)", "(HX)ETH",
	}

	for i, s := range vs {
		ret := convertHXSymbol(s)
		assert.Equal(t, "ETH", ret, "%s", vs[i])
	}
}

func testHXContractInvoke(t *testing.T, wallet bool) {
	ct := "XK1"
	ms1 := makeheader("hx_init", hxwalletid)
	bm := ms1["body"].(map[string]interface{})
	bm[fieldCoinAddr] = myhxaddr
	bm[fieldCoinType] = ct
	buf1 := doReq(t, ms1)
	fmt.Println("hx_init response:", string(buf1))

	content := map[string]string{
		"fromAddr": myhxaddr,
		"toAddr":   dsthxaddr,
		"coinType": ct,
		"tranAmt":  "50.05",
		//"tranFee":  getFieldString(t, buf1, "tranFee"),
		"bak": "",
	}

	cbuf, err := json.Marshal(content)
	assert.Nil(t, err)

	hex, err := HXTransaction(hxreq, HXContractInvoke, buf1, cbuf)
	assert.Nil(t, err)

	fmt.Println("contract invoke:", ghex.EncodeToString(hex))
	ms2 := makeheader("hx_send", hxwalletid)
	bm2 := ms2["body"].(map[string]interface{})
	bm2["tranContent"] = content
	// bm2["data"] = string(hex)
	setBodyData(t, bm2, hex)
	fmt.Println("broadcast hx_send data:", string(hex))

	if wallet {
		directPostWallet(t, "bind", hex)
		return
	}
	buf2 := doReq(t, ms2)
	fmt.Println("hx_send response:", string(buf2))
}
