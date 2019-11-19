/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: address_test
 * Date: 2018-09-04
 *
 */

package hx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNewPrivate(t *testing.T) {
	priv, pub, addr, err := GetNewPrivate()
	assert.Nil(t, err)
	fmt.Println("private:", priv)
	fmt.Println("pubkey:", pub)
	fmt.Println("address:", addr)
}
func TestGetAddressBytes(t *testing.T) {

}
