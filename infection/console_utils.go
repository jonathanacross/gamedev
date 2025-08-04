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

func parseMove(input string) (Move, error) {
	var from, to int
	_, err := fmt.Sscanf(input, "%d %d", from, to)
	if err != nil {
		return Move{}, fmt.Errorf("invalid move format: %v", err)
	}
	return CreateMove(from, to)
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

		move, err = parseMove(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		valid, msg := move.IsValid(b)
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
