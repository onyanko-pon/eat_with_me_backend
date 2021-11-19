CREATE TABLE events (
  id          SERIAL PRIMARY key,
  title       varchar(255) NOT NULL,
  description text NOT NULL,
  latitude    DECIMAL NOT NULL,
  longitude   DECIMAL NOT NULL,
  organize_user_id INT NOT NULL REFERENCES users(id),
  start_datetime TIMESTAMP,
  end_datetime TIMESTAMP
)