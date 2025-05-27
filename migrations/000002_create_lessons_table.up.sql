CREATE TABLE IF NOT EXISTS lessons
(
    id            SERIAL PRIMARY KEY,
    type          VARCHAR(50)  NOT NULL,
    level         VARCHAR(255) NOT NULL,
    mode          VARCHAR(50)  NOT NULL,
    english_level VARCHAR(10)  NOT NULL,
    xp            INT          NOT NULL,
    lesson_data   JSONB        NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

