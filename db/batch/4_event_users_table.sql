CREATE TABLE event_users (
  user_id  INT NOT NULL REFERENCES users(id),
  event_id INT NOT NULL REFERENCES events(id),
  unique (user_id, event_id)
);