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
  number int,
  name varchar(255),

  reception_excellent int,
  reception_positive int,
  reception_negative int,
  reception_error int,

  attack_killed int,
  attack_not_killed int,
  attack_error int,

  serve_ace int,
  serve_positive int,
  serve_negative int,
  serve_error int,

  block_killed int,
  block_positive int
);

CREATE TABLE IF NOT EXISTS matches (
  id int PRIMARY KEY,
  date date,
  hall varchar(255),
  home_team int REFERENCES team_stats (id),
  guest_team int REFERENCES team_stats (id)
);
