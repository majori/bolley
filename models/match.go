package models

import (
	"errors"
	"fmt"
	
	"github.com/majori/bolley/parser"
	"github.com/majori/bolley/db"
)

const (
	DBErrDublicate = "dublicate"
)

type Match struct {}

func (h Match) Create(match *parser.Match) error {
	db := db.GetDB()
	var err error
	var teamIDs [2]*int
	fmt.Println(match)

	// Check if match already exists
	rows, err := db.Query(`
		SELECT id
		FROM matches
		WHERE id = $1
	`, match.ID)

	if err != nil {
		panic(err)
	}

	if rows.Next() {
		return errors.New(DBErrDublicate)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	for i, team := range []parser.Team{match.Home, match.Guest} {
		var teamID int
		err = tx.QueryRow(`
			INSERT INTO team_stats
			VALUES(default, $1, $2, $3, $4, $5)
			returning id;
		`,
			team.Name,
			team.Scores[0], team.Scores[1], team.Scores[2], team.Scores[3],
		).Scan(&teamID)

		if err != nil {
			tx.Rollback()
			panic(err)
		}

		teamIDs[i] = &teamID

		for _, player := range team.Players {
			var insertedID int
			err = tx.QueryRow(`
				INSERT INTO player_stats VALUES (
					default, $1, $2, $3, $4, $5, $6, $7,
					$8, $9, $10, $11, $12, $13, $14, $15, $16
				) returning id;
			`,
				teamID,
				player.Number,
				player.Name,
				player.Reception.Excellent,
				player.Reception.Positive,
				player.Reception.Negative,
				player.Reception.Error,
				player.Attack.Killed,
				player.Attack.NotKilled,
				player.Attack.Error,
				player.Serve.Ace,
				player.Serve.Positive,
				player.Serve.Negative,
				player.Serve.Fault,
				player.Block.Killed,
				player.Block.Positive,
			).Scan(&insertedID)

			if err != nil {
				tx.Rollback()
				panic(err)
			}
		}
	}

	var insertedID int
	err = tx.QueryRow(`
		INSERT INTO matches
		VALUES ($1, $2, $3, $4, $5)
		returning id;
	`, match.ID, match.Date, match.Hall, *teamIDs[0], *teamIDs[1],
	).Scan(&insertedID)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
	return nil
}
