package main

import (
	"fmt"
	"gocycle/core"
	"image/color"
	"math/rand/v2"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameState interface {
	Update(g *Game) error
	Draw(g *Game, screen *ebiten.Image)
}

type Game struct {
	State GameState
}

func NewGame() *Game {
	return &Game{
		State: &TitleScreenState{},
	}
}

func (g *Game) Update() error {
	g.State.Update(g)
	PlayMusic()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.State.Draw(g, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func drawTextAt(screen *ebiten.Image, message string, x float64, y float64, align text.Align, c color.Color) {
	fontSize := float64(16)
	fontFace := &text.GoTextFace{
		Source: MainFaceSource,
		Size:   fontSize,
	}

	// Manually handle alignment to ensure pixel-perfect rendering
	textWidth, _ := text.Measure(message, fontFace, 1.0)
	if align == text.AlignCenter {
		x -= float64(textWidth) / 2
	} else if align == text.AlignEnd {
		x -= float64(textWidth)
	}
	x = float64(int(x))
	y = float64(int(y))

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignStart

	text.Draw(screen, message, fontFace, op)
}

// ------------------- Title Screen State

type TitleScreenState struct{}

func (gs *TitleScreenState) Update(g *Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.State = NewCharacterPickerState()
	}
	return nil
}

func (gs *TitleScreenState) Draw(g *Game, screen *ebiten.Image) {
	drawTextAt(screen, "GoCycle", ScreenWidth/2, ScreenHeight/2, text.AlignCenter, color.White)
	drawTextAt(screen, "Press Space", ScreenWidth/2, 3*ScreenHeight/4, text.AlignCenter, color.White)
}

// ------------------- Character Picker State

type CharacterPickerState struct {
	Picker *CharacterPicker
}

func NewCharacterPickerState() *CharacterPickerState {
	return &CharacterPickerState{
		Picker: NewCharacterPicker(),
	}
}

func (gs *CharacterPickerState) Update(g *Game) error {
	gs.Picker.Update()

	if gs.Picker.IsValid() && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		selectedChars := gs.Picker.GetSelectedCharacters()
		// Shuffle so that chars don't always start in the same place.
		rand.Shuffle(len(selectedChars), func(i, j int) {
			selectedChars[i], selectedChars[j] = selectedChars[j], selectedChars[i]
		})

		initialScores := make(map[int]int)
		for _, char := range selectedChars {
			initialScores[char.ID] = 0
		}

		g.State = NewGamePlayState(selectedChars, 0, initialScores)
	}

	return nil
}

func (gs *CharacterPickerState) Draw(g *Game, screen *ebiten.Image) {
	gs.Picker.Draw(screen)
}

// ------------------- Scrore Screen State

type ScoreScreenState struct {
	CharacterCards []*CharacterFrame
	Scores         map[int]int // Key: character ID, Value: total score across rounds
}

func NewScoreScreenState(chars []*CharData, scores map[int]int) *ScoreScreenState {
	sort.Slice(chars, func(i, j int) bool {
		return scores[chars[i].ID] > scores[chars[j].ID]
	})

	winnerScore := scores[chars[0].ID]
	loserScore := scores[chars[len(chars)-1].ID]
	cards := []*CharacterFrame{}
	spaceX := float64(CharPortraitWidth + 20)
	startX := (ScreenWidth - (spaceX*float64(len(chars)-1) + CharPortraitWidth)) / 2
	y := float64(ScreenHeight/3 - CharPortraitBigHeight/2)
	for i, char := range chars {
		x := startX + float64(i)*spaceX

		mood := CharacterNeutral
		if scores[char.ID] == winnerScore {
			mood = CharacterHappy
		} else if scores[char.ID] == loserScore {
			mood = CharacterSad
		}

		cards = append(cards, NewCharacterFrame(char, x, y, mood, false))
	}

	return &ScoreScreenState{
		CharacterCards: cards,
		Scores:         scores,
	}
}

func (gs *ScoreScreenState) Update(g *Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.State = &TitleScreenState{}
	}

	return nil
}

func (gs *ScoreScreenState) Draw(g *Game, screen *ebiten.Image) {
	for _, card := range gs.CharacterCards {
		card.Draw(screen)
		// Draw the total scores below each card
		score := gs.Scores[card.CharData.ID]
		scoreText := fmt.Sprintf("%d", score)
		scoreX := card.X + card.HitBox().Width()/2
		scoreY := card.Y + card.HitBox().Height() + 5
		drawTextAt(screen, scoreText, scoreX, scoreY, text.AlignCenter, color.White)
	}

	drawTextAt(screen, "Press Space", ScreenWidth/2, 3*ScreenHeight/4, text.AlignCenter, color.White)
}

// ------------------- Game Play State

type GamePlayState struct {
	ArenaView          *ArenaView
	ArenaTimer         *Timer
	ArenaTimeSpeedMs   int
	HumanController1   *core.HumanController
	HumanController2   *core.HumanController
	CharacterCards     []*CharacterFrame
	WaitingForStart    bool
	WaitingForNewRound bool
	EndRoundTimer      *Timer
	Round              int
	PreviousIsAlive    []bool      // Tracks status for score calculation
	RemainingRanks     []int       // Scores pool for tie-breaking
	TotalScores        map[int]int // Key: character ID, Value: total score across rounds
	RoundScores        map[int]int // Key: character ID, Value: score for this round
}

func NewGamePlayState(characters []*CharData, round int, prevTotalScores map[int]int) *GamePlayState {
	var human1 *core.HumanController
	var human2 *core.HumanController

	// Get starting Locations
	numPlayers := len(characters)
	arenaLocs := core.GetStartVectors(numPlayers)
	positionData := PositionDataByNumChars[numPlayers]

	cards := []*CharacterFrame{}
	for i, char := range characters {
		cards = append(cards, NewCharacterFrame(char,
			positionData[i].CardX, positionData[i].CardY, CharacterNeutral, false))
	}

	players := []*core.Player{}

	// Initial directions correspond to the GetStartVectors order: Top-Left=Right, Bottom-Right=Left, Bottom-Left=Up, Top-Right=Down
	initialDirections := []core.Vector{core.Right, core.Left, core.Up, core.Down}

	for i, char := range characters {
		controllerInstance := char.NewController()

		// Use the core.Vector from core.GetStartVectors
		players = append(players, core.NewPlayer(i+1,
			arenaLocs[i], initialDirections[i], controllerInstance))

		// Check and assign Human controllers for input handling using the fresh instance
		switch char.ControllerType {
		case HumanFirstPlayer:
			human1 = controllerInstance.(*core.HumanController)
		case HumanSecondPlayer:
			human2 = controllerInstance.(*core.HumanController)
		}
	}
	var arena = core.NewArenaFromGrid(core.GetGrid(round), players)

	initialStatus := make([]bool, numPlayers)
	for i := range numPlayers {
		initialStatus[i] = players[i].IsAlive
	}

	roundScores := make(map[int]int)
	totalScores := make(map[int]int)
	for _, char := range characters {
		roundScores[char.ID] = -1 // not yet scored
		totalScores[char.ID] = prevTotalScores[char.ID]
	}

	initialRanks := []int{}
	for rank := range numPlayers {
		initialRanks = append(initialRanks, rank*2)
	}

	return &GamePlayState{
		ArenaView:          NewArenaView(arena, characters),
		ArenaTimer:         NewTimer(GameUpdateSpeedMillis * time.Millisecond),
		ArenaTimeSpeedMs:   GameUpdateSpeedMillis,
		HumanController1:   human1,
		HumanController2:   human2,
		CharacterCards:     cards,
		WaitingForStart:    true,
		WaitingForNewRound: false,
		EndRoundTimer:      NewTimer(2 * time.Second),
		Round:              round,
		PreviousIsAlive:    initialStatus,
		RemainingRanks:     initialRanks,
		RoundScores:        roundScores,
		TotalScores:        totalScores,
	}
}

func (gs *GamePlayState) Update(g *Game) error {
	// wait to press space to start the first time
	if gs.WaitingForStart {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			gs.WaitingForStart = false
		}
		return nil
	}

	if gs.WaitingForNewRound {
		gs.EndRoundTimer.Update()
		if gs.EndRoundTimer.IsReady() {
			gs.EndRoundTimer.Reset()
			gs.WaitingForNewRound = false

			// Prepare the next round
			nextRound := gs.Round + 1
			newTotals := make(map[int]int)
			for k, v := range gs.TotalScores {
				newTotals[k] = v + gs.RoundScores[k]
			}
			if nextRound < NumRounds {
				g.State = NewGamePlayState(gs.ArenaView.Characters, nextRound, newTotals)
			} else {
				g.State = NewScoreScreenState(gs.ArenaView.Characters, newTotals)
			}
		}
		return nil
	}

	// Main game play

	if gs.HumanController1 != nil {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			gs.HumanController1.EnqueueDirection(core.Up)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			gs.HumanController1.EnqueueDirection(core.Down)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			gs.HumanController1.EnqueueDirection(core.Left)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			gs.HumanController1.EnqueueDirection(core.Right)
		}
	}
	if gs.HumanController2 != nil {
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			gs.HumanController2.EnqueueDirection(core.Up)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			gs.HumanController2.EnqueueDirection(core.Down)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			gs.HumanController2.EnqueueDirection(core.Left)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			gs.HumanController2.EnqueueDirection(core.Right)
		}
	}

	// Allow player to speed up the game if there are no active human players
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && gs.AllHumanPlayersDead() {
		gs.ArenaTimeSpeedMs /= 2
		gs.ArenaTimer = NewTimer(time.Duration(gs.ArenaTimeSpeedMs) * time.Millisecond)
	}

	gs.ArenaTimer.Update()
	if gs.ArenaTimer.IsReady() {
		gs.ArenaTimer.Reset()
		gs.handleArenaUpdate()

		// Check end of round
		numActivePlayers := 0
		for _, p := range gs.ArenaView.Arena.Players {
			if p.IsAlive {
				numActivePlayers++
			}
		}

		if numActivePlayers <= 1 {
			// Score remaining player(s) and end the round
			gs.scoreRemainingPlayers()
			gs.WaitingForNewRound = true
			gs.EndRoundTimer.Reset()
		}
	}

	return nil
}

func (gs *GamePlayState) AllHumanPlayersDead() bool {
	hasHuman := false
	humansAlive := false
	for i, p := range gs.ArenaView.Arena.Players {
		char := gs.ArenaView.Characters[i]
		if char.ControllerType == HumanFirstPlayer || char.ControllerType == HumanSecondPlayer {
			hasHuman = true
			if p.IsAlive {
				humansAlive = true
			}
		}
	}

	return !hasHuman || !humansAlive
}

// handleArenaUpdate runs one game tick and updates scoring for dead players.
func (gs *GamePlayState) handleArenaUpdate() {
	gs.ArenaView.Arena.Update()

	// 1. Score players who died *this tick* using the core logic.
	gs.scoreDiedPlayers()

	// 2. Update previous status for the next tick.
	for i, p := range gs.ArenaView.Arena.Players {
		gs.PreviousIsAlive[i] = p.IsAlive
	}
}

// scoreDiedPlayers calls the core scoring function and translates the results back
// to the UI's CharData-keyed map.
func (gs *GamePlayState) scoreDiedPlayers() {
	players := gs.ArenaView.Arena.Players

	// Create a temporary map to hold scores by Player ID for the core function
	roundScoresByID := make(map[int]int)
	for i, p := range players {
		char := gs.ArenaView.Characters[i]
		roundScoresByID[p.ID] = gs.RoundScores[char.ID]
	}

	// Update roundScoresByID and gs.RemainingRanks
	core.HandleScoreUpdate(players, gs.PreviousIsAlive, roundScoresByID, &gs.RemainingRanks)

	// Apply the new scores back to the local UI-specific map and update card status
	for i, p := range players {
		char := gs.ArenaView.Characters[i]
		newScore := roundScoresByID[p.ID]
		gs.RoundScores[char.ID] = newScore
	}
}

// scoreRemainingPlayers scores the final winner(s) or the last tie group.
func (gs *GamePlayState) scoreRemainingPlayers() {
	players := gs.ArenaView.Arena.Players

	// Create a temporary map to hold scores by Player ID for the core function
	roundScoresByID := make(map[int]int)
	for i, p := range players {
		char := gs.ArenaView.Characters[i]
		roundScoresByID[p.ID] = gs.RoundScores[char.ID]
	}

	core.ScoreRemainingPlayers(players, roundScoresByID, gs.RemainingRanks)

	// Apply the final scores back to the local UI-specific map and update card status
	for i, p := range players {
		char := gs.ArenaView.Characters[i]
		finalScore := roundScoresByID[p.ID]
		if finalScore >= 0 && gs.RoundScores[char.ID] == -1 {
			gs.RoundScores[char.ID] = finalScore
		}
	}
}

func (gs *GamePlayState) Draw(g *Game, screen *ebiten.Image) {
	roundString := fmt.Sprintf("Round %d/%d", gs.Round+1, NumRounds)
	drawTextAt(screen, roundString, 90, 10, text.AlignStart, color.White)
	gs.ArenaView.Draw(screen)
	for _, card := range gs.CharacterCards {
		card.Draw(screen)

		// Draw the scores below each card
		roundScore := gs.RoundScores[card.CharData.ID]
		if roundScore >= 0 {
			scoreText := fmt.Sprintf("%d", roundScore)
			scoreX := card.X + card.HitBox().Width()/2
			scoreY := card.Y + card.HitBox().Height() + 5
			drawTextAt(screen, scoreText, scoreX, scoreY, text.AlignCenter, color.White)
		}
	}
}

var PositionDataByNumChars = getPositionData()

type PositionData struct {
	ArenaLoc core.Vector
	CardX    float64
	CardY    float64
}

// Positions of where to put players in the arena and where to draw the player
// cards.  This depends on the number of players.
func getPositionData() [][]PositionData {
	lo := 12
	// TODO: this should be 37, but there's a bug when players don't all
	// start with even coordinates.
	hi := 38
	return [][]PositionData{
		{},
		{
			{ArenaLoc: core.Vector{X: lo, Y: lo}, CardX: 10, CardY: 10},
		},
		{
			{ArenaLoc: core.Vector{X: lo, Y: lo}, CardX: 10, CardY: 10},
			{ArenaLoc: core.Vector{X: hi, Y: hi}, CardX: 300, CardY: 10},
		},
		{
			{ArenaLoc: core.Vector{X: lo, Y: lo}, CardX: 10, CardY: 10},
			{ArenaLoc: core.Vector{X: hi, Y: hi}, CardX: 300, CardY: 10},
			{ArenaLoc: core.Vector{X: lo, Y: hi}, CardX: 10, CardY: 120},
		},
		{
			{ArenaLoc: core.Vector{X: lo, Y: lo}, CardX: 10, CardY: 10},
			{ArenaLoc: core.Vector{X: hi, Y: hi}, CardX: 300, CardY: 120},
			{ArenaLoc: core.Vector{X: lo, Y: hi}, CardX: 10, CardY: 120},
			{ArenaLoc: core.Vector{X: hi, Y: lo}, CardX: 300, CardY: 10},
		},
	}
}
