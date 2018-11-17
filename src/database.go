package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func saveMatch(match *Match) {
	var err error
	var teamIDs [2]*int
	for i, team := range []Team{match.Home, match.Guest} {
		var teamID int
		err = db.QueryRow(`
			INSERT INTO team_stats
			VALUES(default, $1, $2, $3, $4, $5)
			returning id;
		`,
			team.Name,
			team.Scores[0], team.Scores[1], team.Scores[2], team.Scores[3],
		).Scan(&teamID)

		if err != nil {
			panic(err)
		}

		teamIDs[i] = &teamID

		for _, player := range team.Players {
			var insertedID int
			err = db.QueryRow(`
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
				panic(err)
			}
		}
	}

	var insertedID int
	err = db.QueryRow(`
		INSERT INTO matches
		VALUES ($1, $2, $3, $4, $5)
		returning id;
	`, match.ID, match.Date, match.Hall, *teamIDs[0], *teamIDs[1],
	).Scan(&insertedID)

	if err != nil {
		panic(err)
	}
}
