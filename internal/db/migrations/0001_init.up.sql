CREATE TABLE news (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    link TEXT UNIQUE NOT NULL,
    description TEXT,
    full_text TEXT,
    published_at TIMESTAMP,
    publisher TEXT,
    vector FLOAT[] 
);
