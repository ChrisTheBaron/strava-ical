CREATE TABLE users
(
  strava_id           INT PRIMARY KEY NOT NULL,
  firstname           STRING          NOT NULL,
  lastname            STRING          NOT NULL,
  email               STRING          NOT NULL,
  strava_access_token STRING          NOT NULL
);
CREATE UNIQUE INDEX users_strava_id_uindex
  ON users (strava_id);
CREATE UNIQUE INDEX users_strava_access_token_uindex
  ON users (strava_access_token);