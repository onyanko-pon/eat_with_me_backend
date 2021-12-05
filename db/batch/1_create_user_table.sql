CREATE TABLE users (
  id        SERIAL primary key,
  username  varchar(255) NOT NULL UNIQUE,
  image_url text
  twitter_screen_name varchar(255) NOT NULL DEFAULT '',
  twitter_username varchar(255) NOT NULL DEFAULT '',
  twitter_user_id SERIAL DEFAULT 0
);
