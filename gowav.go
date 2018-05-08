package gowav

import (
	"io"
)

type WavFile struct {
	AudioFormat	int16
	NumChannels	int16
	SampleRate	int32
	ByteRate	int32
	BlockAlign	int16
	BitsPerSample	int16
}
