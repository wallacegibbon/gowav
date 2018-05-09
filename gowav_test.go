package gowav

import (
	"fmt"
	"os"
	"testing"
)

func Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	w, err := LoadWav(f)
	if err != nil {
		return err
	}
	fmt.Println(w)
	return nil
}

func Test_Load_1(t *testing.T) {
	err := Load("./frog.wav")
	if err != nil {
		t.Error("Failed", err)
	} else {
		t.Log("Succeed")
	}
}

func LoadDump(infile, outfile string) error {
	in, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer out.Close()

	w, err := LoadWav(in)
	if err != nil {
		return err
	}
	err = w.Dump(out)
	if err != nil {
		return err
	}
	return nil
}

func Test_LoadDump_1(t *testing.T) {
	err := LoadDump("./frog.wav", "./out.wav")
	if err != nil {
		t.Error("Failed", err)
	} else {
		t.Log("Succeed")
	}
}
