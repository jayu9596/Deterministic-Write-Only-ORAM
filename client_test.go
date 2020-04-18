package oram

import (
	"bytes"
	"fmt"
	"testing"

	oram "github.com/jaydeep/Path-Oram/helper"
)

// Golang has a very powerful routine for building tests.

// Run with "go test" to run the tests

// And "go test -v" to run verbosely so you see all the logging and
// what tests pass/fail individually.

// And "go test -cover" to check your code coverage in your tests

func TestORAM(t *testing.T) {
	//Generate Data
	data1 := oram.RandomBytes(8128)
	setMaxBucketSize(8)
	j := 0
	cnt := 0
	// valueFeeder := 0
	for j < 127 {
		k := 0
		for k < 8 {
			l := 0
			for l < 8 {
				data1[cnt] = byte(10 + j)
				cnt++
				l++
			}
			k++
		}
		j++
	}

	data2 := data1[0:8]
	_ = ServerUploadData(data1)
	// oldDB := GetMyDB()
	// fmt.Println("Old DB")
	// fmt.Println(oldDB)
	ret, _ := Access(0, data2)
	p := GetStash()
	fmt.Println(p)
	ret, _ = Access(0, nil)
	p = GetStash()
	fmt.Println(p)
	i := 1
	for i < 126 {
		ret, _ = Access(i, data2)
		ret, _ = Access(i, nil)
		p = GetStash()
		fmt.Println(p)
		if bytes.Compare(data2, ret) == 0 {
			fmt.Println("Pass")
		} else {
			fmt.Println("Fail")
		}
		i++
	}

	//Comparing Initial and Latest Stash and Database(server)
	// newDB := GetMyDB()
	// fmt.Println("New DB")
	// fmt.Println(newDB)
	fmt.Println("Stash")
	fmt.Println(GetStash())
	fmt.Println("Max Stash Size at Any point is")
	fmt.Println(GetMaxStashSize())
}
