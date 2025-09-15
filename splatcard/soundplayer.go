package main

import "github.com/hajimehoshi/ebiten/v2/audio"

var (
	audioContext *audio.Context
	soundPlayer  *audio.Player
)

func PlayMusic() {
	const sampleRate = 44100

	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}
	if soundPlayer == nil {
		var err error
		soundPlayer, err = audioContext.NewPlayer(Music)
		if err != nil {
			return
		}
	}

	if soundPlayer != nil && !soundPlayer.IsPlaying() {
		soundPlayer.Rewind()
		soundPlayer.Play()
	}
}
