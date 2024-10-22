package common

import (
	"log"
	"testing"
)

type test struct {
	Name string `toml:"naym"`
}

func TestStore(t *testing.T) {
	var tst test

	tst.Name = "123"

	if err := Save("./test.toml", &tst); err != nil {
		log.Println("failed to save:", err)
		return
	}
	if err := Load("test.toml", &tst); err != nil {
		log.Println("failed to load:", err)
		return
	}

}
