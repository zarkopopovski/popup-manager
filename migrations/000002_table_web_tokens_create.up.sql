CREATE TABLE IF NOT EXISTS web_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title VARCHAR(120) NOT NULL,
    description TEXT NULL,
    web_url VARCHAR(255),
    api_token VARCHAR(120),
    is_valid INTEGER,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL
)