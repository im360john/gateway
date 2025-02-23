CREATE TABLE testdb.gachi_teams (
    id UInt32,
    team_name String,
    motto String
) ENGINE = MergeTree()
ORDER BY id;

INSERT INTO testdb.gachi_teams (id, team_name, motto) VALUES (1, 'Dungeon Lords', 'Pain is pleasure');

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
