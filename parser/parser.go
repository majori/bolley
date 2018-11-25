package parser

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/extrame/xls"
)

// Reception stats
type Reception struct {
	Excellent int
	Positive  int
	Negative  int
	Error     int
}

// Attack stats
type Attack struct {
	Killed    int
	NotKilled int
	Error     int
}

// Serve stats
type Serve struct {
	Ace      int
	Positive int
	Negative int
	Fault    int
}

// Block stats
type Block struct {
	Killed   int
	Positive int
}

// Player stats
type Player struct {
	Number    int
	Name      string
	Reception Reception
	Attack    Attack
	Serve     Serve
	Block     Block
}

// Team stats from the game
type Team struct {
	Name    string
	Scores  [4]int
	Players []Player
}

// Match contains information about one game
type Match struct {
	ID    string
	Date  time.Time
	Hall  string
	Home  Team
	Guest Team
}

func Parse(reader io.ReadSeeker) (*Match, error) {
	var match Match
	if xlFile, err := xls.OpenReader(reader, "utf-8"); err == nil {
		if sheet := xlFile.GetSheet(0); sheet != nil {
			teams := readTeams(sheet)

			if sheet.Row(3).Col(8) == "x" {
				match.Home, match.Guest = teams[0], teams[1]
			} else {
				match.Home, match.Guest = teams[1], teams[0]
			}

			match.ID = sheet.Row(2).Col(12)
			match.Hall = sheet.Row(4).Col(15)
			// TODO: Weird problem with date parsing
			// match.Date, err = time.Parse(time.RFC3339, sheet.Row(3).Col(12))
			// if err != nil {
			// 	return nil, err
			// }
			return &match, nil

		}

		return nil, errors.New("file doen't contain any sheets")
	}

	return nil, errors.New("can't open file")
}

func readTeams(sheet *xls.WorkSheet) [2]Team {
	var teams [2]Team
	for r := 3; r <= 4; r++ {
		name := sheet.Row(r).Col(1)
		var scores [4]int
		for c := 3; c <= 6; c++ {
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

		reception := collectPoints(row, 2, 4)
		attack := collectPoints(row, 8, 3)
		serve := collectPoints(row, 13, 4)
		block := collectPoints(row, 19, 2)

		players = append(players, Player{
			Number:    number,
			Name:      strings.Join(nameSplit[1:3], " "),
			Reception: Reception{reception[0], reception[1], reception[2], reception[3]},
			Attack:    Attack{attack[0], attack[1], attack[2]},
			Serve:     Serve{serve[0], serve[1], serve[2], serve[3]},
			Block:     Block{block[0], block[1]},
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
