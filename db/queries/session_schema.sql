    CREATE TABLE IF NOT EXISTS sessions (
    session_id INTEGER PRIMARY KEY, 
    access_jwt TEXT UNIQUE NOT NULL, 
    refresh_jwt TEXT UNIQUE NOT NULL, 
    session_user_handle TEXT NOT NULL, 
    did TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    last_updated DATE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE);
