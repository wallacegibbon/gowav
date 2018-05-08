package gowav

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type WavFile struct {
	AudioFormat   int
	NumChannels   int
	SampleRate    int
	ByteRate      int
	BlockAlign    int
	BitsPerSample int
}

func LoadWav(r io.Reader) (*WavFile, error) {
	var w WavFile
	err := w.fillFmt(r)
	if err != nil {
		return nil, err
	}
	//w.fillData(r)
	return &w, nil
}

func (w *WavFile) fillFmt(r io.Reader) error {
	buf := make([]byte, 40)
	n, err := r.Read(buf)
	if err != nil {
		return errors.New("[wav read error]" + err.Error())
	}
	if n < 40 {
		return errors.New("Invalid wav file(too short)")
	}
	if !bytes.Equal(buf[:4], []byte("RIFF")) {
		return errors.New("format: missing RIFF")
	}
	if !bytes.Equal(buf[8:12], []byte("WAVE")) {
		return errors.New("format: missing WAVE")
	}
	if !bytes.Equal(buf[12:16], []byte("fmt ")) {
		return errors.New("format: missing fmt ")
	}
	if !bytes.Equal(buf[36:40], []byte("data")) {
		return errors.New("format: missing data")
	}
	w.AudioFormat = toNum(buf[20:22])
	w.NumChannels = toNum(buf[22:24])
	w.SampleRate = toNum(buf[24:28])
	w.ByteRate = toNum(buf[28:32])
	w.BlockAlign = toNum(buf[32:34])
	w.BitsPerSample = toNum(buf[34:36])

	return nil
}

func (w *WavFile) fillData(r io.Reader) {
}

func (w *WavFile) String() string {
	return fmt.Sprintf(
		"[AudioFormat:%d,NumChannels:%d,SampleRate:%d,ByteRate:%d,BlockAlign:%d,BitsPerSample:%d]",
		w.AudioFormat,
		w.NumChannels,
		w.SampleRate,
		w.ByteRate,
		w.BlockAlign,
		w.BitsPerSample)
}

func toNum(p []byte) int {
	var r int
	for i := len(p) - 1; i >= 0; i-- {
		r <<= 8
		r += int(p[i])
	}
	return r
}
