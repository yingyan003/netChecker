package encrypt

import (
	"fmt"
	"testing"
)

func Test_Encrypt(*testing.T) {
	e := NewEncryption("hello")
	fmt.Printf("%s\n", e.String())

	e = NewEncryption("world")
	fmt.Printf("%s\n", e.String())

	e = NewEncryption("hello")
	fmt.Printf("%s\n", e.String())
}
