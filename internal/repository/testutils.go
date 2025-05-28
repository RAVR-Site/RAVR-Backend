package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB создает экземпляр БД в памяти для тестирования.
// Эта функция используется в различных тестах репозитория для настройки тестового окружения.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Автоматически мигрируем все необходимые модели
	err = db.AutoMigrate(&User{}, &Lesson{})
	require.NoError(t, err)

	return db
}
