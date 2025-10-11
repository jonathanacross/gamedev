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
	ArenaView          *ArenaView
	ArenaTimer         *Timer
	HumanController1   *core.HumanController
	HumanController2   *core.HumanController
	CharacterCards     []*CharacterFrame
	WaitingForStart    bool
	WaitingForNewRound bool
	EndRoundTimer      *Timer
	Round              int
	PreviousIsAlive    []bool // Stores player.IsAlive status from the *last* update
	RemainingRanks     []int  // Stores the score values of remaining ranks [8, 6, 4, 2]
	RoundScores        map[*CharData]int
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
	var arena = core.NewArenaFromGrid(GetGrid(round), players)

	initialStatus := make([]bool, len(players))
	for i := range players {
		initialStatus[i] = players[i].IsAlive
	}

	roundScores := make(map[*CharData]int)
	for _, char := range characters {
		roundScores[char] = 0
	}

	// Initialize RemainingRanks with the top scores (1st, 2nd, 3rd, 4th)
	initialRanks := []int{ScoreFirstPlace, ScoreSecondPlace, ScoreThirdPlace, ScoreFourthPlace}
	// Trim the rank array to match the number of players
	numPlayers := len(characters)
	if numPlayers < len(initialRanks) {
		initialRanks = initialRanks[:numPlayers]
	}

	return &GamePlayState{
		ArenaView:          NewArenaView(arena, characters),
		ArenaTimer:         NewTimer(GameUpdateSpeedMillis * time.Millisecond),
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
	}
}

func (gs *GamePlayState) scoreDiedPlayers(justDiedPlayers []*CharData) {
	numDiedThisTick := len(justDiedPlayers)
	if numDiedThisTick == 0 {
		return
	}

	// The ranks we are scoring are the WORST available ranks (the end of the RemainingRanks slice).
	numRemaining := len(gs.RemainingRanks)

	// Ensure we don't try to score more ranks than are available
	if numDiedThisTick > numRemaining {
		numDiedThisTick = numRemaining
	}

	// Identify the score values for the ranks involved in the tie (the last elements)
	ranksInTie := gs.RemainingRanks[numRemaining-numDiedThisTick:]

	scoreSum := 0
	for _, score := range ranksInTie {
		scoreSum += score
	}

	// Calculate the integer-averaged score for the tie group.
	// Since we used doubled scores, this results in a guaranteed integer.
	avgScore := scoreSum / numDiedThisTick

	// Assign the score and immediately update total score.
	for _, char := range justDiedPlayers {
		gs.RoundScores[char] = avgScore
		char.Score += avgScore
	}

	// Remove the assigned ranks from the pool.
	gs.RemainingRanks = gs.RemainingRanks[:numRemaining-numDiedThisTick]
}

func (gs *GamePlayState) handleArenaUpdate() int {
	gs.ArenaView.Update()

	currentActivePlayers := 0
	justDiedPlayers := []*CharData{}

	// Identify players who just died in this time step
	for i, player := range gs.ArenaView.Arena.Players {
		char := gs.ArenaView.Characters[i]

		if gs.PreviousIsAlive[i] && !player.IsAlive {
			justDiedPlayers = append(justDiedPlayers, char)
		}
		if player.IsAlive {
			currentActivePlayers++
		}
	}

	// Score players who just died and update RemainingRanks
	gs.scoreDiedPlayers(justDiedPlayers)

	// Prepare for the next frame's comparison
	gs.PreviousIsAlive = make([]bool, len(gs.ArenaView.Arena.Players))
	for i, player := range gs.ArenaView.Arena.Players {
		gs.PreviousIsAlive[i] = player.IsAlive
	}

	return currentActivePlayers
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
			if nextRound < NumRounds {
				g.State = NewGamePlayState(gs.ArenaView.Characters, nextRound)
			} else {
				g.State = &TitleScreenState{}
			}
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
			numWinners := len(gs.RemainingRanks)

			if numWinners > 0 { // If 0, everyone died in the last tick, and they were already scored.
				scoreSum := 0
				for _, score := range gs.RemainingRanks {
					scoreSum += score
				}

				avgScore := scoreSum / numWinners

				// Identify the remaining alive players and assign the final score
				for i, player := range gs.ArenaView.Arena.Players {
					if player.IsAlive {
						charData := gs.ArenaView.Characters[i]
						gs.RoundScores[charData] = avgScore
						charData.Score += avgScore
					}
				}
			}

			// Start the clock to delay before the next round
			gs.WaitingForNewRound = true
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
		if roundScore > 0 {
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
