CREATE TABLE IF NOT EXISTS popup_message (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    api_token VARCHAR(80) NOT NULL,
    popup_type INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    show_time INTEGER NOT NULL,
    close_time INTEGER NOT NULL,
    popup_pos INTEGER NOT NULL,
    image_name VARCHAR(255) NULL,
    enabled INTEGER NOT NULL,
    date_created DATETIME NOT NULL,
    date_modified DATETIME NOT NULL
) 