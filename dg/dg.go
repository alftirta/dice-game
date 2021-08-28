package dg

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	ID, Score    int
	IsEliminated bool
	Dice         []int
}

type Setting struct {
	TotalPlayers, TotalDice int
	TimeDelay               time.Duration
}

type Game struct {
	Setting
	TotalRemainingPlayers int
	Players, Winners      []Player
}

// rollDice lets the player to roll its dice.
func (p *Player) rollDice() {
	for i := 0; i < len(p.Dice); i++ {
		p.Dice[i] = rand.Intn(6) + 1
	}
}

// removeDie removes one of the player's dice. *Note: "die" is the singular form of "dice".
func (p *Player) removeDie(i int) {
	copy(p.Dice[i:], p.Dice[i+1:])
	p.Dice = p.Dice[:len(p.Dice)-1]
}

// evaluate evaluates the rolled dice according to the dice-game rules.
func (g *Game) evaluate() {
	oneCounter := make(map[int]int) // key: Player.ID; value: count

	// If it's 1, then increment the oneCounter and delete the die.
	// If it's 6, then delete the die.
	for i := range g.Players {
		oneCounter[g.Players[i].ID] = 0
		for j := 0; j < len(g.Players[i].Dice); j++ {
			if g.Players[i].Dice[j] == 1 {
				oneCounter[g.Players[i].ID]++
				g.Players[i].removeDie(j)
				j--
			} else if g.Players[i].Dice[j] == 6 {
				g.Players[i].Score += 6
				g.Players[i].removeDie(j)
				j--
			}
		}
	}

	// Give the numbered 1 dice to the next player.
	for i, v := range oneCounter {
		if v > 0 {
			if i < len(oneCounter) {
				for j := 0; j < v; j++ {
					if g.Players[i].IsEliminated == false {
						g.Players[i].Dice = append(g.Players[i].Dice, 1)
					}
				}
			} else {
				for j := 0; j < v; j++ {
					if g.Players[0].IsEliminated == false {
						g.Players[0].Dice = append(g.Players[0].Dice, 1)
					}
				}
			}
		}
	}

	// Eliminate players who don't have any die left.
	for i, p := range g.Players {
		if len(p.Dice) == 0 && p.IsEliminated == false {
			g.Players[i].IsEliminated = true
			g.TotalRemainingPlayers--
		}
	}
}

// convertDiceToString convert a dice slice into string.
func (g *Game) convertDiceToString(i int) (result string) {
	for j := 0; j < len(g.Players[i].Dice); j++ {
		if j == len(g.Players[i].Dice)-1 {
			result += strconv.Itoa(g.Players[i].Dice[j])
			break
		}
		result += strconv.Itoa(g.Players[i].Dice[j]) + ", "
	}
	if result == "" {
		result = "_"
	}
	return
}

// getWinners get game winners based on players' score.
func (g *Game) getWinners() {
	for _, p := range g.Players {
		if p.Score > g.Winners[0].Score {
			g.Winners = []Player{p}
		} else if p.Score == g.Winners[0].Score && p.ID != g.Players[0].ID {
			g.Winners = append(g.Winners, p)
		}
	}
}

// announceTheWinners announce the winners on the terminal.
func (g *Game) announceTheWinners() {
	if len(g.Winners) == len(g.Players) && len(g.Winners) > 1 {
		fmt.Printf("Seri! Setiap pemain mendapatkan %d poin.\n", g.Winners[0].Score)
	} else if len(g.Winners) == 1 {
		fmt.Printf("Selamat! Pemenangnya adalah Pemain #%d dengan %d poin.\n", g.Winners[0].ID, g.Winners[0].Score)
	} else if len(g.Winners) == 2 {
		fmt.Printf(
			"Seri! Pemain #%d dan Pemain #%d masing-masing mendapatkan %d poin.\n",
			g.Winners[0].ID,
			g.Winners[1].ID,
			g.Winners[1].Score,
		)
	} else if len(g.Winners) > 2 {
		str := "Seri! "
		for i, w := range g.Winners {
			if i < len(g.Winners)-1 {
				str += fmt.Sprintf("Pemain #%d, ", w.ID)
			} else {
				str += fmt.Sprintf("and Pemain #%d masing-masing mendapatkan %d poin.\n", w.ID, w.Score)
			}
		}
		fmt.Println(str)
	}
}

// Play starts the game.
func (g *Game) Play() {
	fmt.Println("====================")
	fmt.Printf("Pemain = %d, Dadu = %d\n", g.Setting.TotalPlayers, g.Setting.TotalDice)
	fmt.Println("====================")
	time.Sleep(g.Setting.TimeDelay)

	turn := 1
	for {
		// All players roll their dice.
		fmt.Printf("Giliran %d lempar dadu:\n", turn)
		for i := 0; i < len(g.Players); i++ {
			// Check if each player still has dice.
			if len(g.Players[i].Dice) > 0 && g.Players[i].IsEliminated == false {
				g.Players[i].rollDice()
			}
			fmt.Printf("    Pemain #%d (%d): %s\n", g.Players[i].ID, g.Players[i].Score, g.convertDiceToString(i))
		}

		// Evaluate according to the game rules.
		g.evaluate()
		fmt.Println("Setelah evaluasi:")
		for i := 0; i < len(g.Players); i++ {
			fmt.Printf("    Pemain #%d (%d): %s\n", g.Players[i].ID, g.Players[i].Score, g.convertDiceToString(i))
		}
		fmt.Println("====================")
		time.Sleep(g.Setting.TimeDelay)

		// Check if the game is over.
		if g.TotalRemainingPlayers <= 1 {
			g.getWinners()
			g.announceTheWinners()
			fmt.Printf("Permainan telah berakhir dalam %d giliran.\n", turn)
			break
		}

		turn++
	}
}

// CreatePlayers creates players.
func CreatePlayers(totalPlayers, totalDice int) ([]Player, error) {
	if totalPlayers < 0 {
		msg := "totalPlayers cannot be negative"
		return nil, errors.New(msg)
	}
	players := make([]Player, totalPlayers)
	for i := 0; i < totalPlayers; i++ {
		players[i] = Player{
			ID: i+1,
			Dice: make([]int, totalDice),
		}
	}
	return players, nil
}

// CreateGame creates a dice game.
func CreateGame(ps []Player, td time.Duration) (Game, error) {
	if ps == nil {
		msg := "players cannot be nil"
		return Game{}, errors.New(msg)
	}
	game := Game{
		Setting: Setting{
			TotalPlayers: len(ps),
			TotalDice: len(ps[0].Dice),
			TimeDelay: td * time.Second,
		},
		TotalRemainingPlayers: len(ps),
		Players: ps,
		Winners: []Player{ps[0]},
	}
	return game, nil
}
