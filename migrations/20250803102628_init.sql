-- +goose Up
CREATE EXTENSION IF NOT EXISTS "vector";
CREATE TABLE news (
                      id TEXT PRIMARY KEY,
                      title TEXT NOT NULL,
                      link TEXT UNIQUE NOT NULL,
                      description TEXT,
                      full_text TEXT,
                      published_at TIMESTAMP,
                      publisher TEXT,
                      vector VECTOR(1024)
);


-- +goose Down
DROP TABLE IF EXISTS news;
