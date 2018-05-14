package gowav

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func oneFrame(w *WavFile) error {
	frm, err := w.GetFrame()
	if err != nil {
		return errors.New("GetFrame failed " + err.Error())
	}
	fmt.Println("The 1st frame:", frm)
	return nil
}

func allFrames(w *WavFile) error {
	frms, err := w.GetAllFrames()
	if err != nil {
		return errors.New("GetAllFrames failed " + err.Error())
	}
	fmt.Println("frame bytes length:", len(frms))
	fmt.Println("The first 10 frames:", frms[:10])

	return nil
}

func Test_BasicInfo(t *testing.T) {
	w, err := NewWavFile("./frog.wav")
	if err != nil {
		t.Error("Failed", err)
		return
	}
	defer w.Close()
	fmt.Println("The wav:", w)
	t.Log("Succeed")
}

func Test_Frame_1(t *testing.T) {
	w, err := NewWavFile("./frog.wav")
	if err != nil {
		t.Error("Failed", err)
		return
	}
	defer w.Close()
	oneFrame(w)
	t.Log("Succeed")
}

func Test_Frame_2(t *testing.T) {
	w, err := NewWavFile("./frog.wav")
	if err != nil {
		t.Error("Failed", err)
		return
	}
	defer w.Close()
	allFrames(w)
	t.Log("Succeed")
}

func loadThenWrite(infile, outfile string) error {
	w, err := NewWavFile(infile)
	if err != nil {
		return err
	}
	defer w.Close()
	out, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer out.Close()

	w.WriteParams(out)
	for {
		frm, err := w.GetFrame()
		if err != nil {
			return err
		}
		if frm != nil {
			out.Write(frm)
		} else {
			break
		}
	}
	return nil
}

func Test_LoadDump_1(t *testing.T) {
	err := loadThenWrite("./frog.wav", "./out.wav")
	if err != nil {
		t.Error("Failed", err)
	} else {
		t.Log("Succeed")
	}
}
