package main

import "github.com/hajimehoshi/ebiten/v2/audio"

var (
	audioContext *audio.Context
	soundPlayer  *audio.Player
)

func PlayWinSound() {
	const sampleRate = 44100

	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}
	if soundPlayer == nil {
		var err error
		soundPlayer, err = audioContext.NewPlayer(WinSound)
		if err != nil {
			return
		}
	}

	if soundPlayer != nil && !soundPlayer.IsPlaying() {
		soundPlayer.Rewind()
		soundPlayer.Play()
	}
}
