
CREATE TABLE friend_applys (
  user_id        INT NOT NULL REFERENCES users(id),
  friend_user_id INT NOT NULL REFERENCES users(id),
  accepted_at    timestamp with time zone DEFAULT NULL,
  unique (user_id, friend_user_id)
);