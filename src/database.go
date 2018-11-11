package main

import (
	"database/sql"
)

func saveMatch(db *sql.DB, match *Match) {
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
					default, $1, $2, $3, $4, $5, $6,
					$7, $8, $9, $10, $11, $12, $13,
					$14, $15
				) returning id;
			`,
				teamID,
				player.Name,
				player.Reception.Perfect,
				player.Reception.Good,
				player.Reception.Bad,
				player.Reception.Fault,
				player.Attack.Kill,
				player.Attack.NotKill,
				player.Attack.Fault,
				player.Serve.Ace,
				player.Serve.Hard,
				player.Serve.Easy,
				player.Serve.Fault,
				player.Block.Score,
				player.Block.Damping,
			).Scan(&insertedID)

			if err != nil {
				panic(err)
			}
		}
	}

	db.QueryRow(`
		INSERT INTO matches
		VALUES (default, $1, $2, $3, $4)
		returning id;
	`, match.Date, match.Hall, *teamIDs[0], *teamIDs[1])
}
