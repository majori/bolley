package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/extrame/xls"
)

// Defence stats
type Defence struct {
	Perfect int
	Good    int
	Bad     int
	Fault   int
}

// Attack stats
type Attack struct {
	Kill    int
	NotKill int
	Fault   int
}

// Serve stats
type Serve struct {
	Ace   int
	Hard  int
	Easy  int
	Fault int
}

// Block stats
type Block struct {
	Score   int
	Damping int
}

// Player stats
type Player struct {
	Number  int
	Name    string
	Defence Defence
	Attack  Attack
	Serve   Serve
	Block   Block
}

// Team stats from the game
type Team struct {
	Name    string
	Scores  [5]int
	Players []Player
}

// Game contains information about one game
type Game struct {
	Home  Team
	Quest Team
}

func main() {
	excelFileName := "" // TODO: Read file somehow
	if xlFile, err := xls.Open(excelFileName, "utf-8"); err == nil {
		if sheet := xlFile.GetSheet(0); sheet != nil {
			teams := readTeams(sheet)
			fmt.Println(teams)
		}
	}
}

func readTeams(sheet *xls.WorkSheet) [2]Team {
	var teams [2]Team
	for r := 3; r <= 4; r++ {
		name := sheet.Row(r).Col(1)
		var scores [5]int
		for c := 3; c <= 7; c++ {
			scores[c-3], _ = strconv.Atoi(sheet.Row(r).Col(c))
		}

		// Find the row where team table starts
		i := 5
		for sheet.Row(i).Col(0) != name {
			i++
		}

		teams[r-3] = Team{
			Name:    sheet.Row(r).Col(1),
			Scores:  scores,
			Players: readPlayers(sheet, i),
		}
	}
	return teams
}

func readPlayers(sheet *xls.WorkSheet, startRow int) []Player {
	var players []Player
	for r := startRow + 2; sheet.Row(r).Col(0) != ""; r++ {
		row := sheet.Row(r)
		// Skip rows which doesn't have name
		if strings.Contains(row.Col(1), "?") {
			continue
		}

		nameSplit := strings.Split(row.Col(1), " ")
		// Skip rows which have more than just number, first name and last name
		if len(nameSplit) > 3 {
			continue
		}

		number, _ := strconv.Atoi(nameSplit[0])

		defence := collectPoints(row, 2, 4)
		attack := collectPoints(row, 8, 3)
		serve := collectPoints(row, 13, 4)
		block := collectPoints(row, 19, 2)

		players = append(players, Player{
			Number:  number,
			Name:    strings.Join(nameSplit[1:3], " "),
			Defence: Defence{defence[0], defence[1], defence[2], defence[3]},
			Attack:  Attack{attack[0], attack[1], attack[2]},
			Serve:   Serve{serve[0], serve[1], serve[2], serve[3]},
			Block:   Block{block[0], block[1]},
		})
	}
	return players
}

func collectPoints(row *xls.Row, start int, length int) []int {
	arr := make([]int, length)
	for c := start; c < start+length; c++ {
		val, _ := strconv.Atoi(row.Col(c))
		arr[c-start] = val
	}
	return arr
}
