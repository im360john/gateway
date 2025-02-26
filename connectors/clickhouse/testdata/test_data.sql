CREATE TABLE testdb.gachi_teams (
    id UInt32,
    team_name String,
    motto String
) ENGINE = MergeTree()
ORDER BY id;

INSERT INTO testdb.gachi_teams (id, team_name, motto) VALUES (1, 'Dungeon Lords', 'Pain is pleasure');
INSERT INTO testdb.gachi_teams (id, team_name, motto) VALUES (2, 'Muscle Brothers', 'Strength and honor');

CREATE TABLE testdb.gachi_personas (
    id UInt32,
    name String,
    strength_level Int32,
    special_move String,
    favorite_drink String,
    battle_cry String,
    team_id UInt32
) ENGINE = MergeTree()
ORDER BY id;

INSERT INTO testdb.gachi_personas (id, name, strength_level, special_move, favorite_drink, battle_cry, team_id) VALUES (1, 'Billy Herrington', 100, 'Anvil Drop', 'Protein Shake', 'Are you ready?', 2);
INSERT INTO testdb.gachi_personas (id, name, strength_level, special_move, favorite_drink, battle_cry, team_id) VALUES (2, 'Van Darkholme', 95, 'Whip of Submission', 'Red Wine', 'I like this kind of stuff', 1);
INSERT INTO testdb.gachi_personas (id, name, strength_level, special_move, favorite_drink, battle_cry, team_id) VALUES (3, 'Ricardo Milos', 90, 'Twerk of Power', 'Pina Colada', 'Let''s dance, boys!', 4);
