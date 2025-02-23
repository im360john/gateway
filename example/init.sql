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

INSERT INTO gachi_personas (name, strength_level, special_move, favorite_drink, battle_cry, team_id) VALUES
('Billy Herrington', 100, 'Anvil Drop', 'Protein Shake', 'Are you ready?', 2),
('Van Darkholme', 95, 'Whip of Submission', 'Red Wine', 'I like this kind of stuff', 1),
('Ricardo Milos', 90, 'Twerk of Power', 'Pina Colada', 'Let’s dance, boys!', 4),
('Mark Wolff', 85, 'Wolf Howl Slam', 'Whiskey on the Rocks', 'Awooo!', 2),
('Kazuhiko', 80, 'Smiling Slam', 'Green Tea', 'Good night, boy', 5),
('Dungeon Master', 99, 'Chains of Domination', 'Dark Ale', 'You have no choice', 1),
('Chad Thundercock', 98, 'Alpha Stomp', 'Pre-Workout Drink', 'Kneel before me!', 6),
('Big Boss', 97, 'Phantom Grip', 'Espresso', 'You’re pretty good', 6),
('Julius Belmont', 92, 'Holy Chains', 'Holy Water', 'Feel the power of discipline!', 5),
('Leather Baron', 91, 'Spanking of Justice', 'Black Coffee', 'You’ve been naughty!', 1),
('Hard Rod', 88, 'Steel Pipe Crush', 'Energy Drink', 'Let’s get HARD!', 2),
('Sweaty Steve', 84, 'Greased Lightning', 'Coconut Water', 'Dripping wet!', 3),
('Oil Overlord', 83, 'Slippery Escape', 'Olive Oil Shot', 'Too slick for you!', 3),
('Thicc Thunder', 81, 'Clap of Doom', 'Banana Smoothie', 'Feel the THICCNESS!', 4),
('Muscle Daddy', 79, 'Bear Hug Crush', 'Chocolate Milkshake', 'Come to daddy!', 2);

-- Вывод всех персонажей с командами для проверки
SELECT
    p.name,
    p.strength_level,
    p.special_move,
    p.favorite_drink,
    p.battle_cry,
    t.team_name,
    t.motto
FROM gachi_personas p
         LEFT JOIN gachi_teams t ON p.team_id = t.id
ORDER BY p.strength_level DESC;


-- Set the session parameter for date parsing in DMY format (e.g., "20-05-2017 14:56")
SET datestyle TO 'DMY';

-- Drop tables if they exist (order matters due to foreign keys)
DROP TABLE IF EXISTS fact_table;
DROP TABLE IF EXISTS payment_dim;
DROP TABLE IF EXISTS customer_dim;
DROP TABLE IF EXISTS item_dim;
DROP TABLE IF EXISTS store_dim;
DROP TABLE IF EXISTS time_dim;

-- Create the payment dimension table
CREATE TABLE IF NOT EXISTS payment_dim (
    payment_key TEXT PRIMARY KEY,
    trans_type  TEXT,
    bank_name   TEXT
);

-- Create the customer dimension table
CREATE TABLE IF NOT EXISTS customer_dim (
    customer_key TEXT PRIMARY KEY,
    name         TEXT,
    contact_no   TEXT,
    nid          TEXT
);

-- Create the item dimension table
CREATE TABLE IF NOT EXISTS item_dim (
    item_key    TEXT PRIMARY KEY,
    item_name   TEXT,
    description TEXT,
    unit_price  NUMERIC(10,2),
    man_country TEXT,
    supplier    TEXT,
    unit        TEXT
);

-- Create the store dimension table
CREATE TABLE IF NOT EXISTS store_dim (
    store_key TEXT PRIMARY KEY,
    division  TEXT,
    district  TEXT,
    upazila   TEXT
);

-- Create the time dimension table
CREATE TABLE IF NOT EXISTS time_dim (
    time_key TEXT PRIMARY KEY,
    date TIMESTAMP,       -- Example: "20-05-2017 14:56"
    hour INTEGER,         -- Hour
    day INTEGER,          -- Day of the month
    week TEXT,            -- Week (e.g., "3rd Week")
    month INTEGER,        -- Month (e.g., 5)
    quarter TEXT,         -- Quarter (e.g., "Q2")
    year INTEGER          -- Year
);

-- Create the fact table
CREATE TABLE IF NOT EXISTS fact_table (
    payment_key TEXT,
    customer_key TEXT,
    time_key TEXT,
    item_key TEXT,
    store_key TEXT,
    quantity INTEGER,
    unit TEXT,
    unit_price NUMERIC(10,2),
    total_price NUMERIC(10,2),
    FOREIGN KEY (payment_key) REFERENCES payment_dim(payment_key),
    FOREIGN KEY (customer_key) REFERENCES customer_dim(customer_key),
    FOREIGN KEY (time_key) REFERENCES time_dim(time_key),
    FOREIGN KEY (item_key) REFERENCES item_dim(item_key),
    FOREIGN KEY (store_key) REFERENCES store_dim(store_key)
);


-- Load data from CSV files with WIN1252 encoding

COPY payment_dim(payment_key, trans_type, bank_name)
    FROM '/var/lib/postgresql/csv/Trans_dim.csv'
    WITH (FORMAT csv, HEADER true, DELIMITER ',', ENCODING 'WIN1252');

COPY customer_dim(customer_key, name, contact_no, nid)
    FROM '/var/lib/postgresql/csv/customer_dim.csv'
    WITH (FORMAT csv, HEADER true, DELIMITER ',', ENCODING 'WIN1252');

COPY item_dim(item_key, item_name, description, unit_price, man_country, supplier, unit)
    FROM '/var/lib/postgresql/csv/item_dim.csv'
    WITH (FORMAT csv, HEADER true, DELIMITER ',', ENCODING 'WIN1252');

COPY store_dim(store_key, division, district, upazila)
    FROM '/var/lib/postgresql/csv/store_dim.csv'
    WITH (FORMAT csv, HEADER true, DELIMITER ',', ENCODING 'WIN1252');

COPY time_dim(time_key, date, hour, day, week, month, quarter, year)
    FROM '/var/lib/postgresql/csv/time_dim.csv'
    WITH (FORMAT csv, HEADER true, DELIMITER ',', ENCODING 'WIN1252');

COPY fact_table(payment_key, customer_key, time_key, item_key, store_key, quantity, unit, unit_price, total_price)
    FROM '/var/lib/postgresql/csv/fact_table.csv'
    WITH (FORMAT csv, HEADER true, DELIMITER ',', ENCODING 'WIN1252');
