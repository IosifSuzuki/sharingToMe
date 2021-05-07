# sharingToMe

# Documentation

## Instruction to create database to run the apllication

```CREATE DATABASE shareme;

\c shareme;

CREATE TABLE publisher (
    id SERIAL PRIMARY KEY,
    nickname CHAR(32) NOT NULL,
    email text,
    ip INET NOT NULL,
    country_flag_url TEXT NOT NULL,
    latitude real NOT NULL,
    longitude real NOT NULL
);

CREATE TABLE consumer (
    id SERIAL PRIMARY KEY,
    nickname CHAR(32) NOT NULL,
    phone_number CHAR(17) NOT NULL,
    password_hash TEXT NOT NULL,
    birth_date TIMESTAMP NOT NULL,
    reference TEXT,
    ip INET NOT NULL,
    country_flag_url TEXT NOT NULL,
    latitude real NOT NULL,
    longitude real NOT NULL
);

CREATE TABLE post (
    id SERIAL PRIMARY KEY,
    publisher_id INT NOT NULL,
    description TEXT,
    file_path TEXT NOT NULL,
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX post_id_index ON post(post_id);
ALTER TABLE post ADD FOREIGN KEY (publisher_id) REFERENCES publisher(id) ON DELETE CASCADE;
```