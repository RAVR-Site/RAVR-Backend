package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/Ravr-Site/Ravr-Backend/config"
	"github.com/Ravr-Site/Ravr-Backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	NUM_USERS      = 50  // Количество фейковых пользователей
	NUM_LESSONS    = 10  // Количество уроков, которые будут пройдены
	MIN_SCORE      = 60  // Минимальное время прохождения урока (в секундах)
	MAX_SCORE      = 300 // Максимальное время прохождения урока (в секундах)
	EXPERIENCE_MIN = 50  // Минимальный опыт за урок
	EXPERIENCE_MAX = 200 // Максимальный опыт за урок
)

var (
	// Имена и фамилии для генерации пользователей
	firstNames = []string{"Александр", "Иван", "Максим", "Дмитрий", "Андрей", "Артем", "Сергей", "Владимир", "Никита", "Михаил"}
	lastNames  = []string{"Смирнов", "Иванов", "Кузнецов", "Соколов", "Попов", "Лебедев", "Козлов", "Новиков", "Морозов", "Петров"}
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Подключаемся к базе данных
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Создаем репозитории
	userRepo := repository.NewUserRepository(db)
	resultRepo := repository.NewResultRepository(db)

	// Устанавливаем seed для генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Creating fake users and results...")

	// Создаем фейковых пользователей
	createdUsers := createFakeUsers(userRepo)

	// Заполняем результаты уроков для каждого пользователя
	createFakeResults(resultRepo, createdUsers)

	fmt.Printf("Successfully created %d fake users with lesson results!\n", NUM_USERS)
}

// Создает фейковых пользователей
func createFakeUsers(userRepo repository.UserRepository) []*repository.User {
	users := make([]*repository.User, 0, NUM_USERS)

	for i := 0; i < NUM_USERS; i++ {
		// Генерируем случайное имя и фамилию
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]

		// Создаем уникальное имя пользователя
		username := fmt.Sprintf("user%d", i+1)

		// Хэшируем пароль
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		// Создаем пользователя
		user := &repository.User{
			Username:   username,
			Password:   string(hashedPassword),
			FirstName:  firstName,
			LastName:   lastName,
			Experience: 0, // Опыт будет накапливаться через результаты
		}

		err := userRepo.Create(user)
		if err != nil {
			log.Printf("Error creating user %s: %v", username, err)
			continue
		}

		users = append(users, user)
		fmt.Printf("Created user: %s (%s %s)\n", username, firstName, lastName)
	}

	return users
}

// Создает фейковые результаты уроков
func createFakeResults(resultRepo repository.ResultRepository, users []*repository.User) {
	// Для каждого пользователя
	for _, user := range users {
		totalExperience := uint64(0)

		// Генерируем несколько результатов для разных уроков
		for lessonID := 1; lessonID <= NUM_LESSONS; lessonID++ {
			// Случайное время прохождения
			score := rand.Intn(MAX_SCORE-MIN_SCORE+1) + MIN_SCORE // Время в секундах

			// Случайный опыт
			experience := uint64(rand.Intn(EXPERIENCE_MAX-EXPERIENCE_MIN+1) + EXPERIENCE_MIN)

			// Используем время в секундах напрямую
			completionTime := uint64(score)

			// Случайная дата завершения в пределах последних 30 дней
			daysAgo := rand.Intn(30)
			completedAt := time.Now().AddDate(0, 0, -daysAgo)

			// Создаем запись о результате
			result := &repository.Result{
				UserID:          user.ID,
				LessonID:        strconv.Itoa(lessonID),
				Score:           uint64(score),
				CompletionTime:  completionTime,
				CompletedAt:     completedAt,
				AddedExperience: experience,
			}

			err := resultRepo.Create(result)
			if err != nil {
				log.Printf("Error creating result for user %s, lesson %d: %v", user.Username, lessonID, err)
				continue
			}

			totalExperience += experience
		}

		// Обновляем общий опыт пользователя
		user.Experience = totalExperience
		if err := userRepo.Update(user.Username, map[string]interface{}{"experience": totalExperience}); err != nil {
			log.Printf("Error updating experience for user %s: %v", user.Username, err)
		}
	}
}

// userRepo - глобальная переменная для доступа к репозиторию пользователей
var userRepo repository.UserRepository
