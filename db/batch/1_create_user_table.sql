CREATE TABLE users (
  id        SERIAL primary key,
  username  varchar(255) NOT NULL UNIQUE,
  image_url text
);
