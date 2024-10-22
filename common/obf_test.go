package common

import (
	"fmt"
	"testing"
)

func TestObfuscateText(t *testing.T) {
	obf := ObfuscateText("Pier")
	fmt.Println(obf)
}
