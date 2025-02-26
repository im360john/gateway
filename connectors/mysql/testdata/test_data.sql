CREATE TABLE gachi_teams (
    id INT AUTO_INCREMENT PRIMARY KEY,
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
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    strength_level INT NOT NULL,
    special_move VARCHAR(100) NOT NULL,
    favorite_drink VARCHAR(50) NOT NULL,
    battle_cry VARCHAR(100) NOT NULL,
    team_id INT,
    FOREIGN KEY (team_id) REFERENCES gachi_teams(id)
);

INSERT INTO gachi_personas (name, strength_level, special_move, favorite_drink, battle_cry, team_id) VALUES
    ('Billy Herrington', 100, 'Anvil Drop', 'Protein Shake', 'Are you ready?', 2),
    ('Van Darkholme', 95, 'Whip of Submission', 'Red Wine', 'I like this kind of stuff', 1); 