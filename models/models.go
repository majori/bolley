package models

import (
	"database/sql"
	"errors"

	"github.com/majori/bolley/db"
	"github.com/majori/bolley/parser"
)

type CumulativeStats struct {
	name              string
	points_scored     int
	points_per_match  float32
	attacks           int
	blocks            int
	blocks_per_match  float32
	aces              int
	aces_per_match    float32
	attack_precent    sql.NullFloat64
	reception_precent sql.NullFloat64
	won_lost          int
}

const (
	DBErrDublicate = "dublicate"
)

func CreateMatch(match *parser.Match) error {
	db := db.GetDB()
	var err error
	var teamIDs [2]*int

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

func GetTeams() {

}

func GetTeamStats(name string) []CumulativeStats {
	db := db.GetDB()
	rows, err := db.Query(`
		WITH cumulative_stats AS (
			SELECT
				ps.name,
		  	COUNT(*) AS matches_played,
				SUM(ps.reception_excellent) AS re_ex,
				SUM(ps.reception_positive) AS re_po,
				SUM(ps.reception_negative) AS re_ne,
				SUM(ps.reception_error) AS re_er,
				SUM(ps.attack_killed) AS at_ki,
				SUM(ps.attack_not_killed) AS at_nk,
				SUM(ps.attack_error) AS at_er,
				SUM(ps.serve_ace) AS se_ac,
				SUM(ps.serve_positive) AS se_po,
				SUM(ps.serve_negative) AS se_ne,
				SUM(ps.serve_error) AS se_er,
				SUM(ps.block_killed) AS bl_ki,
				SUM(ps.block_positive) AS bl_po
			FROM player_stats AS ps
			JOIN team_stats ON ps.team_stat_id = team_stats.id
			WHERE team_stats.name = $1
			GROUP BY ps.name
			ORDER BY ps.name
		)
		SELECT
			cs.name,
			(cs.at_ki + cs.se_ac + cs.bl_ki) AS points_scored,
			ROUND((cs.at_ki + cs.se_ac + cs.bl_ki)::numeric / cs.matches_played, 1) AS points_per_match,
			cs.at_ki AS attacks,
			cs.bl_ki AS blocks,
			ROUND(cs.bl_ki::numeric / cs.matches_played, 1) AS blocks_per_match,
			cs.se_ac AS aces,
			ROUND(cs.se_ac::numeric / cs.matches_played, 1) AS aces_per_match,
			ROUND(100 * (cs.at_ki::numeric / NULLIF((cs.at_ki + cs.at_nk + cs.at_er)::numeric, 0)), 2) AS attack_precent,
			ROUND(100 * ((cs.re_ex + cs.re_po)::numeric / NULLIF((cs.re_ex + cs.re_po + cs.re_ne + cs.re_er)::numeric, 0)),2) AS reception_precent,
			cs.at_ki + cs.se_ac + cs.bl_ki - cs.re_er - cs.at_er - cs.se_er AS won_lost
		FROM cumulative_stats AS cs
	`, name)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var playerStats []CumulativeStats
	for rows.Next() {
		stats := CumulativeStats{}
		err = rows.Scan(
			&stats.name, &stats.points_scored, &stats.points_per_match,
			&stats.attacks, &stats.blocks, &stats.blocks_per_match, &stats.aces,
			&stats.aces_per_match, &stats.attack_precent, &stats.reception_precent,
			&stats.won_lost,
		)
		if err != nil {
			panic(err)
		}

		playerStats = append(playerStats, stats)
	}

	return playerStats
}
