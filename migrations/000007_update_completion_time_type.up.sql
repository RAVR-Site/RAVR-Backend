-- Изменяем тип поля completion_time на BIGINT для хранения времени в секундах
ALTER TABLE results ALTER COLUMN completion_time TYPE BIGINT USING (completion_time::bigint);

-- Обновляем комментарий, чтобы уточнить, что время хранится в секундах
COMMENT ON COLUMN results.completion_time IS 'Время прохождения урока в секундах';
