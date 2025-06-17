-- Добавляем поле experience в таблицу пользователей
ALTER TABLE users ADD COLUMN IF NOT EXISTS experience BIGINT DEFAULT 0 NOT NULL;

-- Создаем таблицу для хранения истории рейтинга пользователей
CREATE TABLE IF NOT EXISTS user_rankings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    experience BIGINT NOT NULL,
    period VARCHAR(20) NOT NULL,
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
