CREATE TABLE IF NOT EXISTS basic_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOt NULL,
    api_token VARCHAR(80) NOT NULL,
    popup_id INTEGER NOT NULL,
    os VARCHAR(255),
    browser VARCHAR(255),
    country VARCHAR(255),
    area VARCHAR(255),
    city VARCHAR(255),
    date_created DATETIME NOT NULL
)