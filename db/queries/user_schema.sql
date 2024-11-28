    CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY,
    handle TEXT NOT NULL UNIQUE,
    app_pass TEXT NOT NULL);
