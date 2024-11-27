CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR(180) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    confirmation_token VARCHAR(120) NOT NULL,
    confirmed INTEGER NOT NULL,
    roles VARCHAR(255) NOT NULL,
    last_login DATETIME NOT NULL,
    date_created  DATETIME NOT NULL,
    date_modified DATETIME NOT NULL
)
