-- CREATE TYPE FRIEND_STATUS AS ENUM ('applying', 'accepted', 'muted');

-- CREATE TABLE friends (
--   user_id        INT NOT NULL REFERENCES users(id),
--   friend_user_id INT NOT NULL REFERENCES users(id),
--   status         FRIEND_STATUS NOT NULL,
--   unique (user_id, friend_user_id)
-- );

CREATE TABLE user_relations (
  user_id        INT NOT NULL REFERENCES users(id),
  friend_user_id INT NOT NULL REFERENCES users(id),
  blinding BOOLEAN NOT NULL DEFAULT FALSE,
  blinded BOOLEAN NOT NULL DEFAULT FALSE,
  unique (user_id, friend_user_id)
);