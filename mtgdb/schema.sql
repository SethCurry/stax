CREATE TABLE IF NOT EXISTS artist (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS set_ (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    code TEXT NOT NULL
);

CREATE TABLE card_ (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    oracle_id TEXT NOT NULL,
    color_identity TEXT
);

CREATE TABLE IF NOT EXISTS card_face (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    flavor_text TEXT NOT NULL,
    oracle_text TEXT NOT NULL,
    lang TEXT NOT NULL,
    cmc float,
    power_ TEXT,
    toughness TEXT,
    loyalty TEXT,
    mana_cost TEXT,
    type_line TEXT,
    colors TEXT
);

CREATE TABLE IF NOT EXISTS printing (
    id INTEGER PRIMARY KEY,
    artist_id INTEGER NOT NULL,
    set_id INTEGER NOT NULL,
    card_face_id INTEGER NOT NULL,
    rarity TEXT NOT NULL,
    name TEXT NOT NULL,
    FOREIGN KEY (artist_id) REFERENCES artist (id),
    FOREIGN KEY (set_id) REFERENCES set_ (id),
    FOREIGN KEY (card_face_id) REFERENCES card_face (id)
);

CREATE TABLE IF NOT EXISTS printing_image (
    id INTEGER PRIMARY KEY,
    printing_id INTEGER NOT NULL,
    image_url TEXT NOT NULL,
    FOREIGN KEY (printing_id) REFERENCES printing (id)
);

CREATE TABLE IF NOT EXISTS ruling (
    id INTEGER PRIMARY KEY,
    card_id INTEGER NOT NULL,
    ruling_text TEXT NOT NULL,
    ruling_date TEXT NOT NULL,
    FOREIGN KEY (card_id) REFERENCES card_ (id)
);