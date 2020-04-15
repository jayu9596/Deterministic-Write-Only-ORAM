package oram

import (
	"fmt"
	"testing"

	oram "github.com/jaydeep/Path-Oram/helper"
)

// Golang has a very powerful routine for building tests.

// Run with "go test" to run the tests

// And "go test -v" to run verbosely so you see all the logging and
// what tests pass/fail individually.

// And "go test -cover" to check your code coverage in your tests

func TestDatastore(t *testing.T) {
	data1 := oram.RandomBytes(4096)
	err := Write(0, data1)
	err = Write(1, data1)
	err = Write(2, data1)
	ret, err2 := Read(1)
	fmt.Println(err2)
	fmt.Println(ret)
	fmt.Println(err)
}
