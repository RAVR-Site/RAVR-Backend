package repository

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestUserRepository_CreateAndFindByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &User{Username: "testuser", Password: "secret"}
	err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	found, err := repo.FindByUsername("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.Username, found.Username)
	assert.Equal(t, user.Password, found.Password)
}

func TestUserRepository_FindByUsername_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	found, err := repo.FindByUsername("nouser")
	assert.NoError(t, err)
	assert.Nil(t, found)
}
