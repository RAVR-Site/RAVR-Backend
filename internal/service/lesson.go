package service

import (
	"encoding/json"
	"fmt"
	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"go.uber.org/zap"
	"os"
	"sort"
	"strings"
	"time"
)

// LessonService интерфейс сервиса для работы с уроками
type LessonService interface {
	GetLesson(id uint) (*Lesson, error)
	GetLessonByType(lessonType string) ([]LessonByType, error)
	LoadLessonsFromFile(filePath string) error
}

// Lesson представляет данные урока с распарсенными JSON данными
type Lesson struct {
	ID         uint                   `json:"id"`
	Type       string                 `json:"type"`
	Level      uint                   `json:"level"`
	Mode       string                 `json:"mode"`
	LessonData map[string]interface{} `json:"lesson_data,omitempty"`
}

// FileLessonData представляет структуру для загрузки уроков из файла
type FileLessonData struct {
	Type       string                   `json:"type"`
	Mode       string                   `json:"mode"`
	LessonData []map[string]interface{} `json:"lesson_data"`
}

type lessonService struct {
	repo   repository.LessonRepository
	logger *zap.Logger
}

func NewLessonService(repo repository.LessonRepository, logger *zap.Logger) LessonService {
	return &lessonService{
		repo:   repo,
		logger: logger,
	}
}

type LessonByType struct {
	Level  uint `json:"level"`
	EasyId uint `json:"easyId"`
	HardId uint `json:"hardId"`
}

func (s lessonService) GetLessonByType(lessonType string) ([]LessonByType, error) {
	lessons, err := s.repo.GetByType(lessonType)
	if err != nil {
		s.logger.Error("Ошибка получения уроков по типу", zap.String("type", lessonType), zap.Error(err))
		return nil, fmt.Errorf("ошибка получения уроков по типу %s: %w", lessonType, err)
	}

	result := make(map[uint]*LessonByType)
	for _, lesson := range lessons {
		if result[lesson.Level] == nil {
			result[lesson.Level] = &LessonByType{
				Level: lesson.Level,
			}
		}

		if lesson.Mode == "easy" {
			result[lesson.Level].EasyId = lesson.ID
		} else if lesson.Mode == "hard" {
			result[lesson.Level].HardId = lesson.ID
		} else {
			s.logger.Warn("Неизвестный режим урока", zap.String("mode", lesson.Mode), zap.Uint("id", lesson.ID))
			continue
		}
	}

	var resultSlice []LessonByType
	for _, lessonByType := range result {
		if lessonByType.EasyId == 0 || lessonByType.HardId == 0 {
			s.logger.Warn("Урок не содержит обоих режимов", zap.Uint("level", lessonByType.Level), zap.Uint("easyId", lessonByType.EasyId), zap.Uint("hardId", lessonByType.HardId))
			continue
		}
		resultSlice = append(resultSlice, *lessonByType)
	}
	sort.Slice(resultSlice, func(i, j int) bool {
		return resultSlice[i].Level < resultSlice[j].Level
	})
	return resultSlice, nil
}

func (s *lessonService) GetLesson(id uint) (*Lesson, error) {
	lesson, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Ошибка получения урока", zap.Uint("id", id), zap.Error(err))
		return nil, fmt.Errorf("урок не найден: %w", err)
	}
	var lessonData map[string]interface{}
	if err := json.Unmarshal(lesson.LessonData, &lessonData); err != nil {
		s.logger.Error("Ошибка разбора данных урока", zap.Uint("id", lesson.ID), zap.Error(err))
		return nil, fmt.Errorf("ошибка разбора данных урока: %w", err)
	}
	return &Lesson{
		ID:         lesson.ID,
		Type:       lesson.Type,
		Level:      lesson.Level,
		Mode:       lesson.Mode,
		LessonData: lessonData,
	}, nil
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

	var lessons, lessonsFromDb []*repository.Lesson
	lessonsFromDb, err = s.repo.GetAll()
	if err != nil {
		s.logger.Error("Ошибка получения всех уроков из базы данных", zap.Error(err))
		return fmt.Errorf("ошибка получения уроков из базы данных: %w", err)
	}
	for _, fileLesson := range fileLessons {
		for _, lessonData := range fileLesson.LessonData {
			level, ok := lessonData["level"].(float64)
			if !ok {
				s.logger.Error("Некорректный уровень урока", zap.String("type", fileLesson.Type))
				return fmt.Errorf("некорректный уровень урока: %s", fileLesson.Type)
			}
			delete(lessonData, "level")
			ld, err := json.Marshal(lessonData)
			if err != nil {
				s.logger.Error("Ошибка сериализации данных урока", zap.String("type", fileLesson.Type), zap.Error(err))
				return fmt.Errorf("ошибка сериализации данных урока: %w", err)
			}
			lesson := &repository.Lesson{
				Type:       fileLesson.Type,
				Mode:       fileLesson.Mode,
				Level:      uint(level),
				LessonData: ld,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			lessons = append(lessons, lesson)
		}
	}

	var lessonsToCreate []*repository.Lesson
	seen := make(map[uint]*repository.Lesson)
	// Проверяем, есть ли уже такие уроки в базе данных
	for _, lesson := range lessons {
		exists := false
		for _, dbLesson := range lessonsFromDb {
			if dbLesson.Type == lesson.Type && dbLesson.Mode == lesson.Mode && dbLesson.Level == lesson.Level {
				err := validateLessonData(dbLesson.LessonData, lesson.LessonData)
				if err != nil {
					s.logger.Warn("Данные урока отличаются от данных в базе", zap.Uint("id", dbLesson.ID), zap.String("type", lesson.Type), zap.String("mode", lesson.Mode), zap.Uint("level", lesson.Level), zap.Error(err))
				}
				exists = true
				seen[dbLesson.ID] = dbLesson
				break
			}
		}
		if !exists {
			lessonsToCreate = append(lessonsToCreate, lesson)
		} else {
			s.logger.Info("Урок уже существует в базе данных", zap.String("type", lesson.Type), zap.String("mode", lesson.Mode), zap.Uint("level", lesson.Level))
		}
	}
	for _, lesson := range lessonsToCreate {
		if err := s.repo.Create(lesson); err != nil {
			s.logger.Error("Ошибка создания урока в базе данных", zap.String("type", lesson.Type), zap.String("mode", lesson.Mode), zap.Uint("level", lesson.Level), zap.Error(err))
			return fmt.Errorf("ошибка создания урока в базе данных: %w", err)
		}
	}

	for _, dbLesson := range lessonsFromDb {
		if seen[dbLesson.ID] == nil {
			if err := s.repo.Delete(dbLesson.ID); err != nil {
				s.logger.Error("Ошибка удаления урока в базе данных", zap.Uint("id", dbLesson.ID), zap.Error(err))
				return fmt.Errorf("ошибка удаления урока в базе данных: %w", err)
			}
		}
	}

	s.logger.Info("Уроки успешно загружены из файла", zap.String("filePath", filePath), zap.Int("count", len(lessonsToCreate)))
	return nil
}

func validateLessonData(dbLessonData, lessonData []byte) error {
	var dbLessonMap, lessonMap map[string]interface{}
	if err := json.Unmarshal(dbLessonData, &dbLessonMap); err != nil {
		return fmt.Errorf("ошибка разбора данных урока из базы данных: %w", err)
	}
	if err := json.Unmarshal(lessonData, &lessonMap); err != nil {
		return fmt.Errorf("ошибка разбора данных урока: %w", err)
	}

	for key, value := range lessonMap {
		if dbValue, exists := dbLessonMap[key]; exists {
			if !strings.EqualFold(fmt.Sprintf("%v", dbValue), fmt.Sprintf("%v", value)) {
				return fmt.Errorf("данные урока отличаются по ключу %s: %v != %v", key, dbValue, value)
			}
		} else {
			return fmt.Errorf("ключ %s отсутствует в данных урока из базы данных", key)
		}
	}

	return nil

}
