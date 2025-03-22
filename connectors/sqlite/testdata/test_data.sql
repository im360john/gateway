-- Create test tables
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    age INTEGER,
    email TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    title TEXT NOT NULL,
    content TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Insert test data
INSERT INTO users (name, age, email) VALUES
    ('John Doe', 30, 'john@example.com'),
    ('Jane Smith', 25, 'jane@example.com'),
    ('Bob Johnson', 35, 'bob@example.com');

INSERT INTO posts (user_id, title, content) VALUES
    (1, 'First Post', 'This is the content of the first post'),
    (1, 'Second Post', 'This is the content of the second post'),
    (2, 'Hello World', 'This is Jane''s first post'); 