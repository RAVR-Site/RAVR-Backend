package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	// Генерируем 32 байта (256 бит) случайных данных
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal("Error generating random bytes:", err)
	}

	// Кодируем в base64 для удобства использования
	secret := base64.URLEncoding.EncodeToString(bytes)

	fmt.Println("Generated JWT Secret:")
	fmt.Println(secret)
	fmt.Println("\nAdd this to your .env file:")
	fmt.Printf("JWT_SECRET=%s\n", secret)
}
