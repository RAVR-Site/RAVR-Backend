package service

import (
	"encoding/json"
	"fmt"
	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"go.uber.org/zap"
	"os"
)

// LessonService интерфейс сервиса для работы с уроками
type LessonService interface {
	GetLesson(id uint) (*repository.Lesson, error)
	GetLessonWithParsedData(id uint) (*LessonDTO, error)
	GetLessonsByTypeWithLimit(lessonType string, limit int) ([]*repository.Lesson, error)
	GetLessonsCountByType(lessonType string) (int64, error)
	LoadLessonsFromFile(filePath string) error
}

// LessonDTO представляет данные урока с распарсенными JSON данными
type LessonDTO struct {
	ID           uint                   `json:"id"`
	Type         string                 `json:"type"`
	Level        string                 `json:"level"`
	Mode         string                 `json:"mode"`
	EnglishLevel string                 `json:"english_level"`
	XP           int                    `json:"xp"`
	LessonData   map[string]interface{} `json:"lesson_data"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
}

// FileLessonData представляет структуру для загрузки уроков из файла
type FileLessonData struct {
	Type         string                 `json:"type"`
	Level        string                 `json:"level"`
	Mode         string                 `json:"mode"`
	EnglishLevel string                 `json:"english_level"`
	XP           int                    `json:"xp"`
	LessonData   map[string]interface{} `json:"lesson_data"`
}

// lessonService имплементация LessonService
type lessonService struct {
	repo   repository.LessonRepository
	logger *zap.Logger
}

// NewLessonService создает новый экземпляр LessonService
func NewLessonService(repo repository.LessonRepository, logger *zap.Logger) LessonService {
	return &lessonService{
		repo:   repo,
		logger: logger,
	}
}

// GetLesson возвращает урок по ID
func (s *lessonService) GetLesson(id uint) (*repository.Lesson, error) {
	lesson, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Ошибка получения урока", zap.Uint("id", id), zap.Error(err))
		return nil, fmt.Errorf("урок не найден: %w", err)
	}
	return lesson, nil
}

// GetLessonWithParsedData возвращает урок с распарсенными данными
func (s *lessonService) GetLessonWithParsedData(id uint) (*LessonDTO, error) {
	lesson, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Ошибка получения урока", zap.Uint("id", id), zap.Error(err))
		return nil, fmt.Errorf("урок не найден: %w", err)
	}

	// Распарсим JSON данные урока
	var lessonData map[string]interface{}
	if err := json.Unmarshal(lesson.LessonData, &lessonData); err != nil {
		s.logger.Error("Ошибка десериализации данных урока", zap.Error(err))
		return nil, fmt.Errorf("ошибка десериализации данных урока: %w", err)
	}

	return &LessonDTO{
		ID:           lesson.ID,
		Type:         lesson.Type,
		Level:        lesson.Level,
		Mode:         lesson.Mode,
		EnglishLevel: lesson.EnglishLevel,
		XP:           lesson.XP,
		LessonData:   lessonData,
		CreatedAt:    lesson.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    lesson.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetLessonsByTypeWithLimit возвращает ограниченное количество уроков определенного типа
func (s *lessonService) GetLessonsByTypeWithLimit(lessonType string, limit int) ([]*repository.Lesson, error) {
	lessons, err := s.repo.GetByTypeWithLimit(lessonType, limit)
	if err != nil {
		s.logger.Error("Ошибка получения уроков по типу с лимитом", zap.String("type", lessonType), zap.Int("limit", limit), zap.Error(err))
		return nil, fmt.Errorf("ошибка получения уроков по типу: %w", err)
	}
	return lessons, nil
}

// GetLessonsCountByType возвращает количество уроков определенного типа
func (s *lessonService) GetLessonsCountByType(lessonType string) (int64, error) {
	count, err := s.repo.GetCountByType(lessonType)
	if err != nil {
		s.logger.Error("Ошибка получения количества уроков по типу", zap.String("type", lessonType), zap.Error(err))
		return 0, fmt.Errorf("ошибка получения количества уроков по типу: %w", err)
	}
	return count, nil
}

// LoadLessonsFromFile загружает уроки из JSON файла
func (s *lessonService) LoadLessonsFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		s.logger.Error("Ошибка чтения файла с уроками", zap.String("path", filePath), zap.Error(err))
		return fmt.Errorf("ошибка чтения файла с уроками: %w", err)
	}

	var fileLessons []FileLessonData
	if err := json.Unmarshal(data, &fileLessons); err != nil {
		s.logger.Error("Ошибка разбора JSON файла с уроками", zap.Error(err))
		return fmt.Errorf("ошибка разбора файла с уроками: %w", err)
	}

	for _, fileLesson := range fileLessons {
		// Проверяем, существует ли урок с такими же параметрами
		lessons, err := s.repo.GetByTypeWithLimit(fileLesson.Type, 1000) // Используем большой лимит для получения всех уроков
		if err != nil {
			s.logger.Error("Ошибка при проверке существующих уроков", zap.Error(err))
			return err
		}

		lessonExists := false
		lessonDataJSON, err := json.Marshal(fileLesson.LessonData)
		if err != nil {
			s.logger.Error("Ошибка сериализации данных урока", zap.Error(err))
			return err
		}

		// Проверка на дубликаты
		for _, existingLesson := range lessons {
			if existingLesson.Type == fileLesson.Type &&
				existingLesson.Level == fileLesson.Level &&
				existingLesson.Mode == fileLesson.Mode &&
				existingLesson.EnglishLevel == fileLesson.EnglishLevel {
				// Проверяем содержимое lesson_data
				var existingData map[string]interface{}
				if err := json.Unmarshal(existingLesson.LessonData, &existingData); err != nil {
					s.logger.Error("Ошибка десериализации данных существующего урока", zap.Error(err))
					continue
				}

				// Сравниваем содержимое уроков
				existingDataJSON, err := json.Marshal(existingData)
				if err != nil {
					s.logger.Error("Ошибка сериализации данных существующего урока", zap.Error(err))
					continue
				}

				// Здесь мы проверяем, совпадают ли данные уроков
				if string(existingDataJSON) == string(lessonDataJSON) {
					lessonExists = true
					s.logger.Info("Урок уже существует, пропускаем",
						zap.String("type", fileLesson.Type),
						zap.String("level", fileLesson.Level),
						zap.String("english_level", fileLesson.EnglishLevel))
					break
				}
			}
		}

		// Если урок не существует, создаем его
		if !lessonExists {
			newLesson := &repository.Lesson{
				Type:         fileLesson.Type,
				Level:        fileLesson.Level,
				Mode:         fileLesson.Mode,
				EnglishLevel: fileLesson.EnglishLevel,
				XP:           fileLesson.XP,
				LessonData:   lessonDataJSON,
			}

			if err := s.repo.Create(newLesson); err != nil {
				s.logger.Error("Ошибка создания урока из файла", zap.Error(err))
				return err
			}

			s.logger.Info("Создан новый урок",
				zap.String("type", fileLesson.Type),
				zap.String("level", fileLesson.Level),
				zap.String("english_level", fileLesson.EnglishLevel))
		}
	}

	s.logger.Info("Загрузка уроков из файла завершена успешно", zap.Int("count", len(fileLessons)))
	return nil
}
