package responses

// Response базовая структура для всех ответов API
type Response struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo содержит детали ошибки
type ErrorInfo struct {
	Code    string `json:"code" example:"INVALID_CREDENTIALS"`
	Message string `json:"message" example:"Неверный логин или пароль"`
}

// DataResponse создает успешный ответ с данными
func DataResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// ErrorResponse создает ответ с ошибкой
func ErrorResponse(code, message string) Response {
	return Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
}

// MessageResponse создает успешный ответ с сообщением
func MessageResponse(message string) Response {
	return DataResponse(map[string]string{
		"message": message,
	})
}

// Empty структура для пустых ответов
type Empty struct{}

// EmptyResponse создает успешный ответ без данных
func EmptyResponse() Response {
	return DataResponse(Empty{})
}

// Стандартные модели ответов для документации Swagger

// @Description Ответ с токеном авторизации
type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// @Description Ответ с информацией о пользователе
type UserResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"johndoe"`
}
