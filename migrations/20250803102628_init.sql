-- +goose Up
CREATE EXTENSION IF NOT EXISTS "vector";
CREATE TABLE news (
                      id TEXT PRIMARY KEY,
                      title TEXT NOT NULL,
                      link TEXT UNIQUE NOT NULL,
                      description TEXT,
                      full_text TEXT,
                      published_at TIMESTAMPTZ,
                      publisher TEXT,
                      vector VECTOR(1024)
);
CREATE TABLE source_last_parsed (
                      source_id TEXT PRIMARY KEY,
                      last_parsed TIMESTAMPTZ NOT NULL
);


-- +goose Down
DROP TABLE IF EXISTS news;
DROP TABLE IF EXISTS source_last_parsed;
