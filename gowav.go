// Use of this source code is governed by a MIT license that can be found
// in the LICENSE file.

// Package gowav implements wav audio file reading and writing.

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
	DataSize      int
	SampleCount   int
	Samples       [][]int
}

func (w *WavFile) fillFmt(r io.Reader) error {
	buf, err := fullRead(r, 36)
	if err != nil {
		return err
	}
	if !eql(buf[:4], "RIFF") {
		return errors.New("missing RIFF")
	}
	if !eql(buf[8:12], "WAVE") {
		return errors.New("missing WAVE")
	}
	if !eql(buf[12:16], "fmt ") {
		return errors.New("missing fmt")
	}
	if toNum(buf[16:20]) != 16 {
		return errors.New("unsupported fmt size")
	}
	w.AudioFormat = toNum(buf[20:22])
	if w.AudioFormat != 1 {
		return errors.New("unsupported audio format")
	}
	w.NumChannels = toNum(buf[22:24])
	w.SampleRate = toNum(buf[24:28])
	w.ByteRate = toNum(buf[28:32])
	w.BlockAlign = toNum(buf[32:34])
	w.BitsPerSample = toNum(buf[34:36])

	w.DataSize, err = readUntilFindData(r)
	if err != nil {
		return err
	}
	w.SampleCount = w.DataSize / w.BlockAlign

	return nil
}

func readUntilFindData(r io.Reader) (int, error) {
	var dataSize int
	for {
		buf, err := fullRead(r, 8)
		if err != nil {
			return 0, err
		}
		dataSize = toNum(buf[4:8])
		if !eql(buf[0:4], "data") {
			_, err := fullRead(r, dataSize)
			if err != nil {
				return 0, err
			}
		} else {
			break
		}
	}
	return dataSize, nil
}

func eql(raw []byte, cmp string) bool {
	return bytes.Equal(raw, []byte(cmp))
}

func (w *WavFile) fillData(r io.Reader) error {
	buf, err := fullRead(r, w.DataSize)
	if err != nil {
		return err
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
		w.Samples[:10])
}

func (w *WavFile) Dump(out io.Writer) error {
	dsize := w.BlockAlign * w.SampleCount
	fsize := dsize + 44 - 8
	buf := make([]byte, 44)
	copy(buf[0:4], []byte("RIFF"))
	copy(buf[4:8], toBytes(fsize, 4))
	copy(buf[8:16], []byte("WAVEfmt "))
	copy(buf[16:20], []byte{0x10, 0, 0, 0})
	copy(buf[20:22], toBytes(w.AudioFormat, 2))
	copy(buf[22:24], toBytes(w.NumChannels, 2))
	copy(buf[24:28], toBytes(w.SampleRate, 2))
	copy(buf[28:32], toBytes(w.ByteRate, 2))
	copy(buf[32:34], toBytes(w.BlockAlign, 2))
	copy(buf[34:36], toBytes(w.BitsPerSample, 2))
	copy(buf[36:40], []byte("data"))
	copy(buf[40:44], toBytes(dsize, 4))
	_, err := out.Write(buf)
	if err != nil {
		return err
	}

	if len(w.Samples) != w.SampleCount {
		return errors.New("Samples length and SampleCount mismatch")
	}
	for i := 0; i < w.SampleCount; i++ {
		for n := 0; n < w.NumChannels; n++ {
			sampleBytes := w.BitsPerSample / 8
			if len(w.Samples[i]) != w.NumChannels {
				return errors.New("NumChannels error")
			}
			_, err := out.Write(
				toBytes(w.Samples[i][n], sampleBytes))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func toBytes(num, size int) []byte {
	r := make([]byte, size)
	for i := 0; i < size; i++ {
		r[i] = byte(num & 0xff)
		num >>= 8
	}
	return r
}

func toNum(p []byte) int {
	var r int
	for i := len(p) - 1; i >= 0; i-- {
		r <<= 8
		r += int(p[i])
	}
	return r
}

func fullRead(r io.Reader, size int) ([]byte, error) {
	buf := make([]byte, size)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	if n < size {
		return nil, errors.New("incomplete stream")
	}
	return buf, nil
}

func LoadWav(r io.Reader) (*WavFile, error) {
	var w WavFile
	err := w.fillFmt(r)
	if err != nil {
		return nil, err
	}

	err = w.fillData(r)
	if err != nil {
		return nil, err
	}
	return &w, nil
}
