package gowav

import (
	"errors"
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
	w, err := NewWav(f)
	if err != nil {
		return err
	}
	fmt.Println(w)
	frm, err := w.GetFrame()
	if err != nil {
		return errors.New("GetFrame failed " + err.Error())
	}
	fmt.Println(frm)
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

	w, err := NewWav(in)
	if err != nil {
		return err
	}

	w.WriteParams(out)
	for {
		frm, err := w.GetFrame()
		if err != nil {
			return nil
		}
		out.Write(frm)
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
