CREATE TABLE IF NOT EXISTS results
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    time_taken INTEGER NOT NULL,
    xp         INTEGER NOT NULL,
    lesson_id  INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (lesson_id) REFERENCES lessons (id) ON DELETE CASCADE
);