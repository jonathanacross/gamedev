package main

import (
	"bufio"
	"fmt"
	"os"
)

func DrawBoard(b *Board) string {
	result := ""
	result += "     "
	for c := range BoardSize {
		result += fmt.Sprintf("%2d", c)
	}
	result += "\n"
	result += "    +"
	for range BoardSize {
		result += fmt.Sprintf("--")
	}
	result += "\n"
	for r := range BoardSize {
		rowidx := GetIndex(r, 0)
		result += fmt.Sprintf("%3d | ", rowidx)
		for c := range BoardSize {
			idx := GetIndex(r, c)
			if b.white.Get(idx) {
				result += "o "
			} else if b.black.Get(idx) {
				result += "x "
			} else {
				result += ". "
			}
		}
		result += "\n"
	}

	return result
}

func IsValidMove(b *Board, m Move) (bool, string) {
	if !b.empty.Get(m.to) {
		return false, "target square is not empty"
	}
	if b.playerToMove == White && !b.white.Get(m.from) {
		return false, "from square is not occupied by white"
	}
	if b.playerToMove == Black && !b.black.Get(m.from) {
		return false, "from square is not occupied by black"
	}
	return true, ""
}

func ParseMove(input string) (Move, error) {
	var m Move
	_, err := fmt.Sscanf(input, "%d %d", &m.from, &m.to)
	if err != nil {
		return m, fmt.Errorf("invalid move format: %v", err)
	}
	if m.from < 0 || m.from >= NumSquares {
		return m, fmt.Errorf("from index %v out of bounds", m.from)
	}
	if m.to < 0 || m.to >= NumSquares {
		return m, fmt.Errorf("to index %v out of bounds", m.to)
	}
	if adjacentBitboards[m.from].Get(m.to) {
		m.jump = false
		return m, nil
	}
	if jumpBitboards[m.from].Get(m.to) {
		m.jump = true
		return m, nil
	}
	return m, fmt.Errorf("invalid move: from %d to %d is neither a step nor a jump", m.from, m.to)
}

func getMoveFromUser(b *Board) Move {
	scanner := bufio.NewScanner(os.Stdin)
	var move Move
	var err error
	for {
		fmt.Printf("> ")
		scanner.Scan()
		input := scanner.Text()
		if input == "q" || input == "quit" {
			os.Exit(1)
		}

		move, err = ParseMove(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		valid, msg := IsValidMove(b, move)
		if !valid {
			fmt.Println(msg)
			continue
		}

		break
	}
	return move
}

// func main() {
// 	engines := map[Player]Engine{
// 		White: &Human{},
// 		Black: &GreedyEngine{},
// 	}
// 	b := NewBoard()
// 	ctr := 0
// 	for ctr < 49 {
// 		ctr++
// 		fmt.Println(DrawBoard(b))
// 		if b.playerToMove == White {
// 			fmt.Printf("White (o) move ")
// 		} else {
// 			fmt.Printf("Black (x) move ")
// 		}
// 		move := engines[b.playerToMove].GenMove(b)
// 		fmt.Println(move)
// 		b.Move(move)
// 	}
// }
