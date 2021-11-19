CREATE TABLE friends (
  user_id        INT NOT NULL,
  friend_user_id INT NOT NULL,
  unique (user_id, friend_user_id)
);
