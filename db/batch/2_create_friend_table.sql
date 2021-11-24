CREATE TYPE FRIEND_STATUS AS ENUM ('applying', 'accepted');

CREATE TABLE friends (
  user_id        INT NOT NULL REFERENCES users(id),
  friend_user_id INT NOT NULL REFERENCES users(id),
  status         FRIEND_STATUS NOT NULL,
  unique (user_id, friend_user_id)
);
