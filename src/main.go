package main

import (
	"fmt"
)

func main() {
	filename := "" // TODO: Read filename somehow

	game, err := parseSpreadsheet(filename)
	if err == nil {
		fmt.Println(game)
	}
}
