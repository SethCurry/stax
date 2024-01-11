CREATE TABLE artist (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE set_ (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    code TEXT NOT NULL
);


CREATE TABLE card_ (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    oracle_text TEXT,
    cmc float,
    power TEXT,
    toughness TEXT,
    loyalty TEXT,
    mana_cost TEXT,
    type_line TEXT,
    colors TEXT,
    color_identity TEXT
);

CREATE TABLE printing (
    id INTEGER PRIMARY KEY,
    artist_id INTEGER NOT NULL,
    set_id INTEGER NOT NULL,
    card_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    FOREIGN KEY (artist_id) REFERENCES artist (id),
    FOREIGN KEY (set_id) REFERENCES set_ (id),
    FOREIGN KEY (card_id) REFERENCES card_ (id)
);

CREATE TABLE ruling (
    id INTEGER PRIMARY KEY,
    card_id INTEGER NOT NULL,
    ruling_text TEXT NOT NULL,
    ruling_date TEXT NOT NULL,
    FOREIGN KEY (card_id) REFERENCES card_ (id)
);