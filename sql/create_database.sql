CREATE TABLE IF NOT EXISTS team_stats (
  id serial PRIMARY KEY,
  name varchar(255),
  round_1_points int,
  round_2_points int,
  round_3_points int,
  round_4_points int
);

CREATE TABLE IF NOT EXISTS player_stats (
  id serial PRIMARY KEY,
  team_stat_id int REFERENCES team_stats (id),
  name varchar(255),

  reception_perfect int,
  reception_good int,
  reception_bad int,
  reception_fault int,

  attack_kill int,
  attack_not_kill int,
  attack_fault int,

  serve_ace int,
  serve_hard int,
  serve_easy int,
  serve_fault int,

  block_score int,
  block_damping int
);

CREATE TABLE IF NOT EXISTS matches (
  id serial PRIMARY KEY,
  date date,
  hall varchar(255),
  home_team int REFERENCES team_stats (id),
  guest_team int REFERENCES team_stats (id)
);
