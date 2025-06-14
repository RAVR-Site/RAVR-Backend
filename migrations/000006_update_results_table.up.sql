-- Сначала удаляем ограничение внешнего ключа
ALTER TABLE results DROP CONSTRAINT IF EXISTS results_lesson_id_fkey;

-- Меняем тип столбца lesson_id на VARCHAR
ALTER TABLE results ALTER COLUMN lesson_id TYPE VARCHAR(255);

-- Добавляем новые столбцы
ALTER TABLE results ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE results ADD COLUMN IF NOT EXISTS completion_time VARCHAR(20);
ALTER TABLE results ADD COLUMN IF NOT EXISTS added_experience BIGINT DEFAULT 0;

-- Переименовываем столбец time_taken в score
ALTER TABLE results RENAME COLUMN time_taken TO score;

-- Меняем тип столбца score на BIGINT
ALTER TABLE results ALTER COLUMN score TYPE BIGINT;

-- Примечание: не создаем новое ограничение внешнего ключа,
-- так как теперь lesson_id будет строкой и не может ссылаться на числовой столбец id
