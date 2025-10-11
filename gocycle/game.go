package main

import (
	"fmt"
	"gocycle/core"
	"image/color"
	"math/rand/v2"
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
	Selector *CharacterSelector
}

func NewCharacterPickerState() *CharacterPickerState {
	return &CharacterPickerState{
		Selector: NewCharacterSelector(16, 30, 74, 90, 2, 5),
	}
}

func (gs *CharacterPickerState) Update(g *Game) error {
	gs.Selector.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		selectedChars := gs.Selector.GetSelectedCharacters()
		// Shuffle so that chars don't always start in the same place.
		rand.Shuffle(len(selectedChars), func(i, j int) {
			selectedChars[i], selectedChars[j] = selectedChars[j], selectedChars[i]
		})

		g.State = NewGamePlayState(selectedChars, 0)
	}

	return nil
}

func (gs *CharacterPickerState) Draw(g *Game, screen *ebiten.Image) {
	gs.Selector.Draw(screen)
}

// ------------------- Game Play State

type GamePlayState struct {
	ArenaView        *ArenaView
	ArenaTimer       *Timer
	HumanController1 *core.HumanController
	HumanController2 *core.HumanController
	CharacterCards   []*CharacterFrame
	WaitingForStart  bool
	Round            int
	PreviousIsAlive  []bool      // Stores player.IsAlive status from the *last* update
	DeathOrder       []*CharData // Stores players as they die
	RoundScores      map[*CharData]int
}

func NewGamePlayState(characters []*CharData, round int) *GamePlayState {
	var human1 *core.HumanController
	var human2 *core.HumanController
	for _, char := range characters {
		switch char.ControllerType {
		case HumanFirstPlayer:
			human1 = char.Controller.(*core.HumanController)
		case HumanSecondPlayer:
			human2 = char.Controller.(*core.HumanController)
		}
	}

	positionData := PositionDataByNumChars[len(characters)]

	cards := []*CharacterFrame{}
	for i, char := range characters {
		cards = append(cards, NewCharacterFrame(char,
			positionData[i].CardX, positionData[i].CardY, CharacterNeutral, false))
	}

	players := []*core.Player{}
	initialDirections := []core.Vector{core.Right, core.Left, core.Down, core.Up}
	for i, char := range characters {
		players = append(players, core.NewPlayer(i+1,
			positionData[i].ArenaLoc, initialDirections[i], char.Controller))
	}
	var arena = core.NewArena(ArenaWidth, ArenaHeight, players)

	initialStatus := make([]bool, len(players))
	for i := range players {
		initialStatus[i] = players[i].IsAlive
	}

	roundScores := make(map[*CharData]int)
	for _, char := range characters {
		roundScores[char] = 0
	}

	return &GamePlayState{
		ArenaView:        NewArenaView(arena, characters),
		ArenaTimer:       NewTimer(GameUpdateSpeedMillis * time.Millisecond),
		HumanController1: human1,
		HumanController2: human2,
		CharacterCards:   cards,
		WaitingForStart:  true,
		Round:            round,
		PreviousIsAlive:  initialStatus,
		DeathOrder:       []*CharData{},
		RoundScores:      roundScores,
	}
}

func (gs *GamePlayState) handleArenaUpdate() int {
	gs.ArenaView.Update() // This calls gs.Arena.Update()

	currentActivePlayers := 0
	justDiedPlayers := []*core.Player{}

	// 1. Identify players who just died in this time step and update DeathOrder
	for i, player := range gs.ArenaView.Arena.Players {
		// Check if player just died this frame (was alive, now dead)
		if gs.PreviousIsAlive[i] && !player.IsAlive {
			justDiedPlayers = append(justDiedPlayers, player)
			// Add the character to the death order list
			gs.DeathOrder = append(gs.DeathOrder, gs.ArenaView.Characters[i])
		}
		if player.IsAlive {
			currentActivePlayers++
		}
	}

	// 2. Prepare for the next frame's comparison
	gs.PreviousIsAlive = make([]bool, len(gs.ArenaView.Arena.Players))
	for i, player := range gs.ArenaView.Arena.Players {
		gs.PreviousIsAlive[i] = player.IsAlive
	}

	// --- Integrated Scoring Logic ---
	dyingCharScores := make(map[*CharData]int)

	if len(justDiedPlayers) > 0 {
		groupSize := len(justDiedPlayers)

		// The number of players who have already died. This determines the starting score slot.
		baseIndex := len(gs.DeathOrder) - groupSize

		// Look up the scores using the base index and the group size.
		scores, ok := ScoreLookup[baseIndex][groupSize]
		if !ok {
			scores = make([]int, groupSize)
		}

		// Map the dying characters to the scores they receive for fast lookup
		for i, player := range justDiedPlayers {
			charData := gs.ArenaView.Characters[player.ID-1]
			dyingCharScores[charData] = scores[i]
		}
	}

	// 3. Iterate over ALL characters in the round to update scores
	// This ensures that the loop runs for all players, fulfilling your request.
	for _, charData := range gs.ArenaView.Characters {
		// If the character is in the dying map, add the calculated score.
		// Otherwise, the lookup returns 0 (which is correct for alive/already dead players).
		if score, found := dyingCharScores[charData]; found {
			gs.RoundScores[charData] += score
		}
	}
	// --- End Integrated Scoring Logic ---

	return currentActivePlayers
}

var ScoreLookup = map[int]map[int][]int{
	0: { // First death.
		1: {0},          // 1 player died
		2: {1, 1},       // 2 players died simultaneously. Both get 1 point.
		3: {2, 2, 2},    // 3 players died. All get 2 points.
		4: {3, 3, 3, 3}, // All 4 players died. All get 3 points.
	},
	1: { // Second death.
		1: {2},       // 1 player died. Score 2.
		2: {3, 3},    // 2 players died. Both get 3 points.
		3: {4, 4, 4}, // 3 players died. All get 4 points.
	},
	2: { // Third death.
		1: {4},    // 1 player died. Score 4.
		2: {5, 5}, // 2 players died. Both get 5 points.
	},
	3: { // 1 player remaining
		1: {6}, // Last player gets all 6 points.
	},
}

func (gs *GamePlayState) Update(g *Game) error {
	// wait to press space to start the first time
	if gs.WaitingForStart {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			gs.WaitingForStart = false
		}
		return nil
	}

	if gs.HumanController1 != nil {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			gs.HumanController1.EnqueueDirection(core.Left)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			gs.HumanController1.EnqueueDirection(core.Right)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			gs.HumanController1.EnqueueDirection(core.Up)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			gs.HumanController1.EnqueueDirection(core.Down)
		}
	}
	if gs.HumanController2 != nil {
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			gs.HumanController2.EnqueueDirection(core.Left)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			gs.HumanController2.EnqueueDirection(core.Right)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			gs.HumanController2.EnqueueDirection(core.Up)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			gs.HumanController2.EnqueueDirection(core.Down)
		}
	}

	gs.ArenaTimer.Update()
	if gs.ArenaTimer.IsReady() {
		gs.ArenaTimer.Reset()

		numActivePlayers := gs.handleArenaUpdate()

		// check for end of round
		if numActivePlayers <= 1 {
			for charData, roundScore := range gs.RoundScores {
				charData.Score += roundScore
			}

			nextRound := gs.Round + 1
			if nextRound < NumRounds {
				g.State = NewGamePlayState(gs.ArenaView.Characters, nextRound)
			} else {
				g.State = &TitleScreenState{}
			}
		}
	}
	return nil
}

func (gs *GamePlayState) Draw(g *Game, screen *ebiten.Image) {
	roundString := fmt.Sprintf("Round %d/%d", gs.Round+1, NumRounds)
	drawTextAt(screen, roundString, 90, 10, text.AlignStart, color.White)
	gs.ArenaView.Draw(screen)
	for _, card := range gs.CharacterCards {
		card.Draw(screen)

		// Draw the scores below each card
		roundScore := gs.RoundScores[card.CharData]
		scoreText := fmt.Sprintf("%d", roundScore)
		scoreX := card.X + card.HitBox().Width()/2
		scoreY := card.Y + card.HitBox().Height() + 5 // +5 for a small offset below the card
		drawTextAt(screen, scoreText, scoreX, scoreY, text.AlignCenter, color.White)
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
	return [][]PositionData{
		{},
		{
			{ArenaLoc: core.Vector{X: 10, Y: 10}, CardX: 10, CardY: 10},
		},
		{
			{ArenaLoc: core.Vector{X: 10, Y: 10}, CardX: 10, CardY: 10},
			{ArenaLoc: core.Vector{X: 30, Y: 30}, CardX: 300, CardY: 10},
		},
		{
			{ArenaLoc: core.Vector{X: 10, Y: 10}, CardX: 10, CardY: 10},
			{ArenaLoc: core.Vector{X: 30, Y: 30}, CardX: 300, CardY: 10},
			{ArenaLoc: core.Vector{X: 10, Y: 30}, CardX: 10, CardY: 120},
		},
		{
			{ArenaLoc: core.Vector{X: 10, Y: 10}, CardX: 10, CardY: 10},
			{ArenaLoc: core.Vector{X: 30, Y: 30}, CardX: 300, CardY: 120},
			{ArenaLoc: core.Vector{X: 10, Y: 30}, CardX: 10, CardY: 120},
			{ArenaLoc: core.Vector{X: 30, Y: 10}, CardX: 300, CardY: 10},
		},
	}
}
