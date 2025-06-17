-- Проверяем существование колонки xp, и если она существует, делаем её nullable
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
              WHERE table_name='results' AND column_name='xp') THEN
        ALTER TABLE results ALTER COLUMN xp DROP NOT NULL;

        -- Добавляем комментарий о том, что поле устарело и будет удалено
        COMMENT ON COLUMN results.xp IS 'Устаревшее поле, используйте added_experience';

        -- Копируем значения из added_experience в xp для обратной совместимости
        UPDATE results SET xp = added_experience WHERE xp IS NULL;
    END IF;
END
$$;

