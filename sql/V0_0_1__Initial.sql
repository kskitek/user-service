CREATE TABLE Users (
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL UNIQUE,
  email VARCHAR NOT NULL UNIQUE,
  password VARCHAR NOT NULL,
  creationDate DATE NOT NULL
);
