package responses

// SuccessResponse обобщенная структура для успешных ответов
type SuccessResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

// ErrorResponse структура для ответов с ошибкой
type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   ErrorInfo `json:"error"`
}

// ErrorInfo содержит детали ошибки
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success создает успешный ответ с данными указанного типа T
func Success[T any](data T) SuccessResponse[T] {
	return SuccessResponse[T]{
		Success: true,
		Data:    data,
	}
}

// Error создает ответ с ошибкой
func Error(code, message string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error: ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
}

// MessageResponse создает успешный ответ с сообщением
func MessageResponse(message string) SuccessResponse[map[string]string] {
	return Success(map[string]string{
		"message": message,
	})
}
