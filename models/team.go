package models

import (
	"database/sql"

	"github.com/majori/bolley/db"
)

type Team struct{}

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

func (t Team) GetAll() []string {
	db := db.GetDB()
	rows, err := db.Query(`
		SELECT name
		FROM team_stats
		GROUP BY name
	`)

	if err != nil {
		panic(err)
	}

	var teams []string
	for rows.Next() {
		var team string
		rows.Scan(&team)
		teams = append(teams, team)
	}

	return teams
}

func (t Team) GetCumulativeStats(name string) []CumulativeStats {
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
