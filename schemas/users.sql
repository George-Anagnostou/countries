CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username                TEXT NOT NULL UNIQUE,
    password                TEXT NOT NULL,
    longest_country_score   TEXT NOT NULL DEFAULT 0,
    current_country_score   TEXT NOT NULL DEFAULT 0,
    longest_capital_score   TEXT NOT NULL DEFAULT 0,
    current_capital_score   TEXT NOT NULL DEFAULT 0,
    current_country         TEXT NOT NULL,
    current_capital         TEXT NOT NULL,
    created_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
