/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: trx_test.go
 * Date: 2018-09-05
 *
 */

package hx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"strconv"
	"testing"
	"time"
)

func testTimeConvert(t *testing.T) {

	tmp_time := time.Now().Unix()

	fmt.Println(Time2Str(tmp_time))

}

func testJson(t *testing.T) {

	transferOp := DefaultTransferOperation()

	transferTrx := Transaction{
		1,
		2,
		"2018-09-04T08:16:25",
		[][]interface{}{{0, transferOp}},
		make([]interface{}, 0),
		[]string{"2018-09-04T08:16:25"},
		3,
		nil,
	}

	b, err := json.Marshal(transferTrx)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)

}

func testSignature(t *testing.T) {

	chainidHex := "fe70279c1d9850d4ddb6ca1f00c577bc2e86bf33d54fafd4c606a6937b89ae32"

	seed := MnemonicToSeed("venture lazy digital aware plug hire acquire abuse chunk know gloom snow much employ glow rich exclude allow", "123")

	//addrKey, _ := GetAddressKey(seed, 0, 0)

	wif, err := ExportWif(seed, 0, 0)
	assert.Nil(t, err)

	bin, err := hex.DecodeString(chainidHex)
	assert.Nil(t, err)

	sig, err := GetSignature(wif, bin)
	assert.Nil(t, err)

	fmt.Println("bts sign by C:", len(sig), hex.EncodeToString(sig))

	sig2, err := btsSign(wif, bin)
	assert.Nil(t, err)

	fmt.Println("bts sign by go:", len(sig2), hex.EncodeToString(sig2))
}

func testPack(t *testing.T) {
	//out, err := BuildBindAccountTransaction("6413,2521010061", "5KbxwmcNhUQe7oVN5oMpC3BiYmGpNDf8u3W1EPn3qzfrGxVahyq", "HXNRhTbDKiw2ut91BJEc5zy49HKWQnRw9as7", 0, "TshDfDSPRhV2BCDAJFAGAjs5K2TJGfrCPua",
	//	"HC", "PtWVJdsvDGYsidC9igf6h2KRBFCdUb7k6Phx9DoZPeXzHirhL2yAM", "", "9f3b24c962226c1cb775144e73ba7bb177f9ed0b72fac69cd38764093ab530bd")

	//BuildUnBindAccountTransaction("6413,2521010061", "5KbxwmcNhUQe7oVN5oMpC3BiYmGpNDf8u3W1EPn3qzfrGxVahyq","HXNRhTbDKiw2ut91BJEc5zy49HKWQnRw9as7", 0, "TshDfDSPRhV2BCDAJFAGAjs5K2TJGfrCPua",
	//	"HC", "PtWVJdsvDGYsidC9igf6h2KRBFCdUb7k6Phx9DoZPeXzHirhL2yAM", "9f3b24c962226c1cb775144e73ba7bb177f9ed0b72fac69cd38764093ab530bd")

	//BuildWithdrawCrosschainTransaction("39618,358453409", "5KbxwmcNhUQe7oVN5oMpC3BiYmGpNDf8u3W1EPn3qzfrGxVahyq","HXNRhTbDKiw2ut91BJEc5zy49HKWQnRw9as7", 0, "TshDfDSPRhV2BCDAJFAGAjs5K2TJGfrCPua",
	//	"HC", "1.2", "9f3b24c962226c1cb775144e73ba7bb177f9ed0b72fac69cd38764093ab530bd")

	//BuildRegisterAccountTransaction("39618,358453409", "5KbxwmcNhUQe7oVN5oMpC3BiYmGpNDf8u3W1EPn3qzfrGxVahyq","HXNRhTbDKiw2ut91BJEc5zy49HKWQnRw9as7", "HX77DEz5FFdsbyM4P4XMyZ5Xm2DHPph4o3GjLXcyc8Eq62s84SMw",500000, "", "wens", "9f3b24c962226c1cb775144e73ba7bb177f9ed0b72fac69cd38764093ab530bd")

	out, err := BuildLockBalanceTransaction("33081,1272682105", "5KR6ocp5eUdWWYPX7mYp4XLGBcZ2xHVHVsNaco6K2YZSWQTqES7", "HXNcikaxB2rsK26JCiwvzse9AqFPkBGyAynG", "1.2.60", "1.3.0", 100000, 0,
		"1.6.11", "InvalidAddress", "07c870b857439cc298de0f7747d475c57320ddfdd6f28357f7bed2a7ff41e821")

	//BuildRedeemBalanceTransaction("3379,2788439324", "5KbxwmcNhUQe7oVN5oMpC3BiYmGpNDf8u3W1EPn3qzfrGxVahyq","HXNRhTbDKiw2ut91BJEc5zy49HKWQnRw9as7", "1.2.78", "1.3.0", 400000,0,
	//	                    "1.6.11", "HXNKuyBkoGdZZSLyPbJEetheRhMjezkaXk2J","9f3b24c962226c1cb775144e73ba7bb177f9ed0b72fac69cd38764093ab530bd")

	asset_arr := []string{"citizen10,54459861,1.3.0", "citizen9,39886,1.3.0"}
	out, err = BuildObtainPaybackTransaction("5595,4227186882", "5KbxwmcNhUQe7oVN5oMpC3BiYmGpNDf8u3W1EPn3qzfrGxVahyq", "HXNRhTbDKiw2ut91BJEc5zy49HKWQnRw9as7", 200, asset_arr,
		"2.22.15", "9f3b24c962226c1cb775144e73ba7bb177f9ed0b72fac69cd38764093ab530bd")

	if err != nil {
		fmt.Println("error")
	}
	_ = out

	fmt.Println(string(out))
}

func testPrecision(t *testing.T) {
	s := "5.00001"

	f, err := strconv.ParseFloat(s, 64)
	assert.Nil(t, err)

	i := int64(math.Round(f * 100000.00))
	fmt.Printf("s=%v  f=%v  i=%d\n", s, f, i)
}
