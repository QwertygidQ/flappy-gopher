package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type Decoded struct {
	streamer beep.StreamSeekCloser
	format   beep.Format
}

type SoundHandler struct {
	decoded map[string]*Decoded
}

func initBeep(sampleRate beep.SampleRate) *SoundHandler {
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	return &SoundHandler{decoded: make(map[string]*Decoded)}
}

func (sh *SoundHandler) loadSound(path, key string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	streamer, format, err := wav.Decode(file)
	if err != nil {
		panic(err)
	}

	sh.decoded[key] = &Decoded{streamer: streamer, format: format}
}

func (sh *SoundHandler) playSound(key string, loop bool) {
	streamer := sh.decoded[key].streamer
	if loop {
		loopStreamer := beep.Loop(-1, streamer)
		speaker.Play(loopStreamer)
	} else {
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			streamer.Seek(0)
		})))
	}
}
