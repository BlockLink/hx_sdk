package hx

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"

	// "github.com/HcashOrg/hcd/chaincfg/chainec"
	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/HcashOrg/hcd/hcec/secp256k1"
	"github.com/HcashOrg/hcd/hcutil"
	"github.com/HcashOrg/hcd/wire"
	"github.com/stretchr/testify/assert"
)

// hx hc bind signature

func testHxBindSign(t *testing.T) {
	var (
		wif = "PtWVNgWCnPpiM9rC3RL7KyhzNs1RoeqT3wCU71WpWrP8Sv9tvxqB1"
		//pubkey = "037fb3a233ad44d892dc8d1505670eb87d0d000dc3c3ece55ed1e055523ab813d1"
		addr = "TsZZqfAWSsyZ2631bi1xKYvMEx8CLk9fSkJ"
		exp  = "H8l+f8Dq+Z65ofbrXOwgxTaMniNJ73Z1TOmZOJgcCQP7NFCCQ7TjzO573UOhno5k3OVvrFk5xrqXO7V7q6wakHM="
	)

	_ = exp
	w, err := hcutil.DecodeWIF(wif)
	assert.Nil(t, err, "decode wif failed")

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
	assert.Nil(t, err)
	fmt.Println("sig hex:", hex.EncodeToString(res))
	ret := base64.StdEncoding.EncodeToString(res)
	fmt.Println("sig:", ret)

	// assert.Equal(t, ret, exp)
	/*
		r, s, err := chainec.Secp256k1.Sign(w.PrivKey, []byte("Hc Signed Message:\n" + addr))
		assert.Nil(t, err)

		sig := chainec.Secp256k1.NewSignature(r, s)
		res := sig.Serialize()

		fmt.Println("sig hex:", hex.EncodeToString(res))
		ret := base64.StdEncoding.EncodeToString(res)
		fmt.Println("sig:", ret)
	*/
}

/*
{
  "addr": "mftonGbWCGFmXyrKYTxufMctJ71vzSEZsK",
  "pubkey": "02f4707a30d245fb85e05b34c59d0946bbe793541fe884d9fd02e6d50be721bf26",
  "wif_key": "cQ4i3HjrVbXyp863E9SHmt93eFJ2q6kcupxZW5HmoKD9ayu9pFY6"
}
unlocked >>> create_crosschain_symbol ETH
create_crosschain_symbol ETH
2250138ms th_a       wallet.cpp:1112               save_wallet_file     ] saving wallet to file wallet.json
{
  "addr": "0x7b487d0e51d33e6434f98985abf4c0f04503547b",
  "pubkey": "0x0465b419454c0c196eb07b7504054b447b5c27f164349e852a4d25f787de7485719edc4163bb79243f792640efba355f23800060e981494bab37769a59e84cdd7f",
  "wif_key": "0d6f7ad765debb2a7fb99d87deaa1b65084919f6d0f22b308c8edbaa58578b9a"
}

unlocked >>> wallet_create_account guotie
wallet_create_account guotie
2407087ms th_a       wallet.cpp:1112               save_wallet_file     ] saving wallet to file wallet.json
"HXNRy59BQPvKdKGhh52oa7LiYkr69judc7oe"

unlocked >>>  bind_tunnel_account guotie 0x7b487d0e51d33e6434f98985abf4c0f04503547b ETH true
 bind_tunnel_account guotie 0x7b487d0e51d33e6434f98985abf4c0f04503547b ETH true
2662347ms th_a       wallet.cpp:700                get_account          ] my account id 0.0.0 different from blockchain id 1.2.53
{
  "ref_block_num": 52728,
  "ref_block_prefix": 3769530690,
  "expiration": "2018-11-02T03:54:20",
  "operations": [[
      10,{
        "fee": {
          "amount": 0,
          "asset_id": "1.3.0"
        },
        "crosschain_type": "ETH",
        "addr": "HXNRy59BQPvKdKGhh52oa7LiYkr69judc7oe",
        "account_signature": "1f11fcd3714236b3ad9786751c90763e4371b548ee4cdb40629112d9ae6d4a443d7b447eacf1dc0ae2bcf24309748ff2bae565649f5967308db97a377a24b97e99",
        "tunnel_address": "0x7b487d0e51d33e6434f98985abf4c0f04503547b",
        "tunnel_signature": "0xeb1cd9f5c05d82d488e91bf5e78a1dc9c0b4517a8bdb22b9db8e73c1bf303dd62b57b2963369127d5ad0f1c54f0989b9ea2571f9ce4fe3c8389f9138f47549f81b"
      }
    ]
  ],
  "extensions": [],
  "signatures": [
    "1f6a536ba0f11fadfc3298a74d69eebbd7f36395197ad565284ceb8e126ccf02e523d6a691588046366073a78a3fdb860380edf8be39db028170726a4185e672d4"
  ],
  "block_num": 0,
  "trxid": "dafb8e53c87ea4287becea18a983a5e2bade122a"
}
*/
func testEthSignAddr(t *testing.T) {
	var ethAddres = []struct {
		wif  string
		addr string
		exp  string
	}{{
		wif: "3f2153c638e857ae4b5ef132c1ee09c24bb48484d2dea91a5071b202be2e2a90",
		//pubkey = "037fb3a233ad44d892dc8d1505670eb87d0d000dc3c3ece55ed1e055523ab813d1"
		addr: "1891025831596418915523e786334b2b44985272",
		exp:  "0xa4770569de58c0ffaff3ec741c57db40cd9bc37d2de60ea8af91d4e73668257d42bc19e6c7bcb4ea8099fbeb8a0d8e4258e573839cea46b2f2219917de8c9ef61b",
	},
		{
			wif: "0d6f7ad765debb2a7fb99d87deaa1b65084919f6d0f22b308c8edbaa58578b9a",
			//pubkey = "037fb3a233ad44d892dc8d1505670eb87d0d000dc3c3ece55ed1e055523ab813d1"
			addr: "7b487d0e51d33e6434f98985abf4c0f04503547b",
			exp:  "0xeb1cd9f5c05d82d488e91bf5e78a1dc9c0b4517a8bdb22b9db8e73c1bf303dd62b57b2963369127d5ad0f1c54f0989b9ea2571f9ce4fe3c8389f9138f47549f81b",
		},
	}

	SetTestnetEthSig()

	/*
		ret, err := ethSignAddress(wif, addr)
		assert.Nil(t, err)

		_=exp
		fmt.Println("eth sign:", ret)
		assert.Equal(t, exp, ret)
	*/

	for i, eth := range ethAddres {
		wif := eth.wif
		addr := eth.addr
		exp := eth.exp

		ret2, err := ethSignAddress2(wif, addr)
		assert.Nil(t, err)

		_ = exp
		fmt.Println("eth sign:", i, ret2)
		assert.Equal(t, exp, ret2)
	}
}

func testBtcSignAddr(t *testing.T) {
	var (
		wif = "cQNNUDuVz1qBEVhuhHtjgY5Cb8FieZXDLc9e45DWhd2RzZw4fT5W"
		//pubkey = "037fb3a233ad44d892dc8d1505670eb87d0d000dc3c3ece55ed1e055523ab813d1"
		addr = "msxWoneJBLcBWeZDHH89PDLgYyGmtCpGQe"
		exp  = "H/wTWGMMH6F687Zjj2N48YAwMYmMR9ZN2EDk/kmqEd/lBszu8Hxijg8lakUQy6EChkTTRPh6JksY7V/+VG7ykuc="
	)

	ret, err := btcSignAddress(wif, addr)
	assert.Nil(t, err)

	_ = exp
	fmt.Println("btc sign:", ret)
	assert.Equal(t, exp, ret)
}

func testLtcSignAddr(t *testing.T) {
	var ltcAddres = []struct {
		wif  string
		addr string
		exp  string
	}{{
		wif: "cV3KHPBxArGaQcBSLaEvWTqorAV4FJZmbvXXvP41BqmzZAnGCP1b",
		//pubkey = "037fb3a233ad44d892dc8d1505670eb87d0d000dc3c3ece55ed1e055523ab813d1"
		addr: "mmhLVJ6WW42BbGoCfgs6zpg4j6WCGP4ffB",
		exp:  "IEe4kdlDKOhUj+g1fHuRX2MbjX4m9O/ScuHOx810LHjdXV7qSRmOruAX4Bc3qXOwMq0cKgSWJqWXcHHUnzKpEHE=",
	},
		{
			wif:  "T5LN3q4o3CUtKpoi3TtiL4JMPNjqaBhkoWzQJ2F4HLWuu94Y6G8y",
			addr: "LP2JKjy9WmSygMdoe2CzEHabXrPSPXspNF",
			exp:  "H6zfPRpOp/on4gS/2JcpbHo9Up8cj4I+1zIYng64/RUzOqo0allqC7J+hxStQY+E7+Ze9fZdac+KuMwBmXjvCww=",
		},
	}

	for _, laddr := range ltcAddres {
		wif := laddr.wif
		addr := laddr.addr
		exp := laddr.exp

		ret, err := ltcSignAddress(wif, addr)
		assert.Nil(t, err)

		_ = exp
		fmt.Println("ltc sign:", ret)
		assert.Equal(t, exp, ret)
	}
}

/*
func TestHash(t *testing.T) {
	msg := []byte("1234567890abcdef")
	h1:=crypto.Keccak256(msg)
	h2:=Keccak256(msg)
	assert.Equal(t, h1, h2)
}

func testSignv2(t *testing.T) {
	var (
		wif = "0d6f7ad765debb2a7fb99d87deaa1b65084919f6d0f22b308c8edbaa58578b9a"
		msg = []byte("1234567890abcdef1234567890abcdef")
		)

	buf, err := hex.DecodeString(wif)
	if err != nil {
		fmt.Println("decode wif failed: ", err)
		return
	}

	key, err := crypto.ToECDSA(buf)
	fmt.Println("key1:", key.PublicKey, "D:", *key.D)

	s1, err := crypto.Sign(msg, key)
	assert.Nil(t, err)


	key2, pk := btcec.PrivKeyFromBytes(btcec.S256(), buf)
	_=pk
	fmt.Println("key2:", key2.PublicKey, "D:", *key2.D)

	s2, err := btcec.SignCompact(btcec.S256(), key2, msg, true)
	assert.Nil(t, err)

	fmt.Println("s1:", hex.EncodeToString(s1))
	fmt.Println("s2:", hex.EncodeToString(s2))


	s3, err := btcec.SignCompact(btcec.S256(), key2, msg, false)
	assert.Nil(t, err)
	fmt.Println("s3:", hex.EncodeToString(s3))
}
*/

func randomMsg(l int) string {
	chars := "01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	msg := ""
	for i := 0; i < l; i++ {
		msg += string(chars[rand.Intn(len(chars))])
	}

	return msg
}

/*
func testEthSign(t *testing.T) {
	var (
		wif = "0d6f7ad765debb2a7fb99d87deaa1b65084919f6d0f22b308c8edbaa58578b9a"
		msgs = []string {}
	)

	for i := 0; i < 100; i ++ {
		msgs = append(msgs, randomMsg(32))
	}

	buf, err := hex.DecodeString(wif)
	if err != nil {
		fmt.Println("decode wif failed: ", err)
		return
	}

	key, err := crypto.ToECDSA(buf)
	fmt.Println("key1:", key.PublicKey, "D:", *key.D)

	for _, msg := range msgs {
		b1, err := crypto.Sign([]byte(msg), key)
		assert.Nil(t, err)

		b2, err := Sign2(wif, []byte(msg))
		assert.Nil(t, err)

		s1 := hex.EncodeToString(b1)
		s2 := hex.EncodeToString(b2)
		assert.Equal(t, s1[0:128], s2[2:])
		// fmt.Printf("sign: %v %v\n", s1, s2)
	}
}
*/
