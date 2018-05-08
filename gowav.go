package gowav

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const dataOffset = 44

type WavFile struct {
	AudioFormat   int
	NumChannels   int
	SampleRate    int
	ByteRate      int
	BlockAlign    int
	BitsPerSample int
	DataSize      int
	SampleCount   int
	Samples       [][]int
}

func (w *WavFile) fillFmt(r io.Reader) error {
	buf := make([]byte, dataOffset)
	n, err := r.Read(buf)
	if err != nil {
		return errors.New("[wav head read error]" + err.Error())
	}
	if n < dataOffset {
		return errors.New("Invalid wav file(header incomplete)")
	}
	if !eql(buf[:4], "RIFF") {
		return errors.New("format: missing RIFF")
	}
	if !eql(buf[8:12], "WAVE") {
		return errors.New("format: missing WAVE")
	}
	if !eql(buf[12:16], "fmt ") {
		return errors.New("format: missing fmt ")
	}
	if !eql(buf[36:40], "data") {
		return errors.New("format: missing data")
	}
	w.AudioFormat = toNum(buf[20:22])
	w.NumChannels = toNum(buf[22:24])
	w.SampleRate = toNum(buf[24:28])
	w.ByteRate = toNum(buf[28:32])
	w.BlockAlign = toNum(buf[32:34])
	w.BitsPerSample = toNum(buf[34:36])
	w.DataSize = toNum(buf[40:44])
	w.SampleCount = w.DataSize / w.BlockAlign

	return nil
}

func (w *WavFile) fillData(r io.Reader) error {
	buf := make([]byte, w.DataSize)
	n, err := r.Read(buf)
	if err != nil {
		return errors.New("[wav data read error]" + err.Error())
	}
	if n < w.DataSize {
		return errors.New("Invalid wav file(data incomplete)")
	}
	for i := 0; i < w.SampleCount; i++ {
		tmp := make([]int, w.NumChannels)
		for n := 0; n < w.NumChannels; n++ {
			sampleBytes := w.BitsPerSample / 8
			offset := i*w.BlockAlign + n*sampleBytes
			tmp[n] = toNum(buf[offset : offset+sampleBytes])
		}
		w.Samples = append(w.Samples, tmp)
	}
	return nil
}

func (w *WavFile) String() string {
	return fmt.Sprintf("[AudioFormat:%d,NumChannels:%d,SampleRate:%d,"+
		"ByteRate:%d,BlockAlign:%d,BitsPerSample:%d,"+
		"SampleCount:%d]%v",
		w.AudioFormat,
		w.NumChannels,
		w.SampleRate,
		w.ByteRate,
		w.BlockAlign,
		w.BitsPerSample,
		w.SampleCount,
		w.Samples[0:10])
}

func toNum(p []byte) int {
	var r int
	for i := len(p) - 1; i >= 0; i-- {
		r <<= 8
		r += int(p[i])
	}
	return r
}

func eql(raw []byte, cmp string) bool {
	return bytes.Equal(raw, []byte(cmp))
}

func LoadWav(r io.Reader) (*WavFile, error) {
	var w WavFile
	err := w.fillFmt(r)
	if err != nil {
		return nil, err
	}

	w.fillData(r)
	if err != nil {
		return nil, err
	}

	return &w, nil
}
