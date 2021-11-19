CREATE TABLE friends (
  user_id        INT NOT NULL REFERENCES users(id),
  friend_user_id INT NOT NULL REFERENCES users(id),
  unique (user_id, friend_user_id)
);
