package repository

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createTestLessons(t *testing.T, db *gorm.DB) []*Lesson {
	lessonData1, err := json.Marshal(map[string]interface{}{
		"title":       "Времена глаголов",
		"description": "Изучение времен глаголов в русском языке",
		"content":     "Содержание урока...",
	})
	require.NoError(t, err)

	lessonData2, err := json.Marshal(map[string]interface{}{
		"title":       "Части речи",
		"description": "Основные части речи",
		"content":     "Содержание урока...",
	})
	require.NoError(t, err)

	lessonData3, err := json.Marshal(map[string]interface{}{
		"title":       "Базовая лексика",
		"description": "Основные слова и фразы",
		"content":     "Содержание урока...",
	})
	require.NoError(t, err)

	lessons := []*Lesson{
		{
			Type:         "grammar",
			Level:        "beginner",
			Mode:         "theory",
			EnglishLevel: "A1",
			XP:           100,
			LessonData:   lessonData1,
		},
		{
			Type:         "grammar",
			Level:        "intermediate",
			Mode:         "practice",
			EnglishLevel: "B1",
			XP:           150,
			LessonData:   lessonData2,
		},
		{
			Type:         "vocabulary",
			Level:        "beginner",
			Mode:         "practice",
			EnglishLevel: "A2",
			XP:           80,
			LessonData:   lessonData3,
		},
	}

	for _, lesson := range lessons {
		err := db.Create(lesson).Error
		require.NoError(t, err)
	}

	return lessons
}

func TestLessonRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepository(db)

	lessonData, err := json.Marshal(map[string]interface{}{
		"title":       "Тестовый урок",
		"description": "Описание тестового урока",
		"content":     "Содержание урока...",
	})
	require.NoError(t, err)

	lesson := &Lesson{
		Type:         "test",
		Level:        "beginner",
		Mode:         "theory",
		EnglishLevel: "A1",
		XP:           100,
		LessonData:   lessonData,
	}

	// Тест создания урока
	err = repo.Create(lesson)
	assert.NoError(t, err)
	assert.NotZero(t, lesson.ID)
}

func TestLessonRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepository(db)
	lessons := createTestLessons(t, db)

	// Тест получения урока по ID
	lesson, err := repo.GetByID(lessons[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, lessons[0].ID, lesson.ID)
	assert.Equal(t, "grammar", lesson.Type)
	assert.Equal(t, "beginner", lesson.Level)
	assert.Equal(t, "A1", lesson.EnglishLevel) // Проверка уровня английского

	// Тест получения несуществующего урока
	lesson, err = repo.GetByID(9999)
	assert.Error(t, err)
	assert.Nil(t, lesson)
}

func TestLessonRepository_GetByTypeWithLimit(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepository(db)
	createTestLessons(t, db)

	// Тест получения уроков с лимитом
	lessons, err := repo.GetByTypeWithLimit("grammar", 1)
	assert.NoError(t, err)
	assert.Len(t, lessons, 1)
	assert.Equal(t, "grammar", lessons[0].Type)
}

func TestLessonRepository_GetCountByType(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepository(db)
	createTestLessons(t, db)

	// Тест получения количества уроков по типу "grammar"
	count, err := repo.GetCountByType("grammar")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// Тест получения количества уроков по типу "vocabulary"
	count, err = repo.GetCountByType("vocabulary")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Тест получения количества уроков по несуществующему типу
	count, err = repo.GetCountByType("nonexistent")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestLessonRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepository(db)
	lessons := createTestLessons(t, db)

	// Получаем урок для обновления
	lesson, err := repo.GetByID(lessons[0].ID)
	require.NoError(t, err)

	// Обновляем данные урока
	lesson.Level = "advanced"
	lesson.EnglishLevel = "C1" // Обновляем уровень английского
	lesson.XP = 200

	updatedLessonData, err := json.Marshal(map[string]interface{}{
		"title":       "Обновленный урок",
		"description": "Обновленное описание",
		"content":     "Обновленное содержание...",
	})
	require.NoError(t, err)
	lesson.LessonData = updatedLessonData

	// Тест обновления урока
	err = repo.Update(lesson)
	assert.NoError(t, err)

	// Проверяем, что урок обновился
	updatedLesson, err := repo.GetByID(lesson.ID)
	assert.NoError(t, err)
	assert.Equal(t, "advanced", updatedLesson.Level)
	assert.Equal(t, "C1", updatedLesson.EnglishLevel) // Проверка уровня английского
	assert.Equal(t, 200, updatedLesson.XP)

	var lessonData map[string]interface{}
	err = json.Unmarshal(updatedLesson.LessonData, &lessonData)
	assert.NoError(t, err)
	assert.Equal(t, "Обновленный урок", lessonData["title"])
}

func TestLessonRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLessonRepository(db)
	lessons := createTestLessons(t, db)

	// Проверяем, что урок существует
	_, err := repo.GetByID(lessons[0].ID)
	require.NoError(t, err)

	// Тест удаления урока
	err = repo.Delete(lessons[0].ID)
	assert.NoError(t, err)

	// Проверяем, что урок удален
	_, err = repo.GetByID(lessons[0].ID)
	assert.Error(t, err)
}
