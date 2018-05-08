package gowav

import (
	"fmt"
	"os"
	"testing"
)

func Show(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	w, err := LoadWav(f)
	if err != nil {
		return err
	}
	fmt.Println(w)
	return nil
}

func Test_Show_1(t *testing.T) {
	err := Show("./frog.wav")
	if err != nil {
		t.Error("Failed", err)
	} else {
		t.Log("Succeed")
	}
}

