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
	FileSize      int
	NumChannels   int
	SampleRate    int
	ByteRate      int
	BlockAlign    int
	BitsPerSample int
	DataSize      int
	SampleCount   int
	stream        io.Reader
}

func (w *WavFile) GetParams() error {
	buf, err := w.read(36)
	if err != nil {
		return err
	}
	if !eql(buf[:4], "RIFF") {
		return errors.New("missing RIFF")
	}
	w.FileSize = toNum(buf[4:8])

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

	w.DataSize, err = w.getDataChunk()
	if err != nil {
		return err
	}
	w.SampleCount = w.DataSize / w.BlockAlign

	return nil
}

func (w *WavFile) getDataChunk() (int, error) {
	var dataSize int
	for {
		buf, err := w.read(8)
		if err != nil {
			return 0, err
		}
		dataSize = toNum(buf[4:8])
		if !eql(buf[0:4], "data") {
			_, err := w.read(dataSize)
			if err != nil {
				return 0, err
			}
		} else {
			break
		}
	}
	return dataSize, nil
}

func (w *WavFile) read(size int) ([]byte, error) {
	buf := make([]byte, size)
	n, err := w.stream.Read(buf)
	if err != nil {
		return nil, err
	}
	if n < size {
		return nil, errors.New("incomplete stream")
	}
	return buf, nil
}

func (w *WavFile) GetFrame() ([]byte, error) {
	d, err := w.read(w.BlockAlign)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return d, nil
}

func (w *WavFile) String() string {
	return fmt.Sprintf(
		"[AudioFormat:%d,NumChannels:%d,SampleRate:%d,"+
			"ByteRate:%d,BlockAlign:%d,BitsPerSample:%d,"+
			"SampleCount:%d]",
		w.AudioFormat,
		w.NumChannels,
		w.SampleRate,
		w.ByteRate,
		w.BlockAlign,
		w.BitsPerSample,
		w.SampleCount)
}

func (w *WavFile) WriteParams(out io.Writer) error {
	buf := make([]byte, 44)
	copy(buf[0:4], []byte("RIFF"))
	copy(buf[4:8], toBytes(w.FileSize, 4))
	copy(buf[8:16], []byte("WAVEfmt "))
	copy(buf[16:20], []byte{0x10, 0, 0, 0})
	copy(buf[20:22], toBytes(w.AudioFormat, 2))
	copy(buf[22:24], toBytes(w.NumChannels, 2))
	copy(buf[24:28], toBytes(w.SampleRate, 2))
	copy(buf[28:32], toBytes(w.ByteRate, 2))
	copy(buf[32:34], toBytes(w.BlockAlign, 2))
	copy(buf[34:36], toBytes(w.BitsPerSample, 2))
	copy(buf[36:40], []byte("data"))
	copy(buf[40:44], toBytes(w.DataSize, 4))
	_, err := out.Write(buf)

	return err
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

func eql(raw []byte, cmp string) bool {
	return bytes.Equal(raw, []byte(cmp))
}

func NewWav(r io.Reader) (*WavFile, error) {
	w := WavFile{stream: r}
	err := w.GetParams()
	if err != nil {
		return nil, err
	}

	return &w, nil
}
