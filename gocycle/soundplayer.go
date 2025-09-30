package main

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

var (
	audioContext *audio.Context
	musicPlayer  *audio.Player // Dedicated player for music
)

func init() {
	const sampleRate = 44100
	audioContext = audio.NewContext(sampleRate)

	// Create a dedicated player for music from the byte slice
	musicStream, err := mp3.DecodeWithoutResampling(bytes.NewReader(MusicBytes))
	if err != nil {
		panic(err)
	}

	musicPlayer, err = audioContext.NewPlayer(musicStream)
	if err != nil {
		panic(err)
	}
}

// PlayMusic starts playing the background music
func PlayMusic() {
	if !musicPlayer.IsPlaying() {
		musicPlayer.Rewind()
		musicPlayer.Play()
	}
}

// PlaySound creates and plays a new player for the given sound stream
func PlaySound(s []byte) {
	stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(s))
	if err != nil {
		return
	}
	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		return
	}
	//	player.SetVolume(0.25)
	player.Play()
}
