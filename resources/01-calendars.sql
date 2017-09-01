CREATE TABLE calendars
(
  id      STRING PRIMARY KEY NOT NULL,
  user_id INT                NOT NULL,
  CONSTRAINT calendars_users_strava_id_fk FOREIGN KEY (user_id) REFERENCES users (strava_id)
);
CREATE UNIQUE INDEX calendars_id_uindex
  ON calendars (id);