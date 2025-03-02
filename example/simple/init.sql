-- Create a demo table for personas
CREATE TABLE gachi_teams (
    id SERIAL PRIMARY KEY,
    team_name VARCHAR(50) NOT NULL,
    motto VARCHAR(100) NOT NULL
);

INSERT INTO gachi_teams (team_name, motto) VALUES
    ('Dungeon Lords', 'Pain is pleasure'),
    ('Muscle Brothers', 'Strength and honor'),
    ('Oil Masters', 'Slip into submission'),
    ('Thicc Squad', 'The heavier, the better'),
    ('Holy Disciples', 'Divine domination'),
    ('The Alpha Pack', 'Only the strongest survive');

CREATE TABLE gachi_personas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    strength_level INT NOT NULL,
    special_move VARCHAR(100) NOT NULL,
    favorite_drink VARCHAR(50) NOT NULL,
    battle_cry VARCHAR(100) NOT NULL,
    team_id INT REFERENCES gachi_teams(id)
);
