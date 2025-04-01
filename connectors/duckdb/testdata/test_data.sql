-- Create test tables
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    name VARCHAR NOT NULL,
    age INTEGER,
    email VARCHAR UNIQUE
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    title VARCHAR NOT NULL,
    content VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Insert test data
INSERT INTO users (id, name, age, email) VALUES
    (1, 'John Doe', 30, 'john@example.com'),
    (2, 'Jane Smith', 25, 'jane@example.com'),
    (3, 'Bob Johnson', 35, 'bob@example.com');

INSERT INTO posts (id, user_id, title, content) VALUES
    (1, 1, 'First Post', 'This is the content of the first post'),
    (2, 1, 'Second Post', 'This is the content of the second post'),
    (3, 2, 'Hello World', 'This is Jane''s first post'); 