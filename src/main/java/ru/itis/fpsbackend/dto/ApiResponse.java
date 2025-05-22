package ru.itis.fpsbackend.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.Getter;
import lombok.Setter;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;

@Setter
@Getter
@Schema(description = "Стандартный формат ответа API")
public class ApiResponse<T> {
    @Schema(description = "Статус успешности операции", example = "true")
    private boolean success;
    
    @Schema(description = "Сообщение о результате операции", example = "Операция выполнена успешно")
    private String message;
    
    @Schema(description = "Метка времени ответа", example = "2025-04-21T14:30:15.123456")
    private String timestamp;
    
    @Schema(description = "Данные ответа (могут отличаться в зависимости от endpoint)")
    private T data;

    public ApiResponse() {
        this.timestamp = LocalDateTime.now().format(DateTimeFormatter.ISO_DATE_TIME);
    }

    public ApiResponse(boolean success, String message) {
        this();
        this.success = success;
        this.message = message;
    }

    public ApiResponse(boolean success, String message, T data) {
        this(success, message);
        this.data = data;
    }

    public static <T> ApiResponse<T> success(String message) {
        return new ApiResponse<>(true, message);
    }

    public static <T> ApiResponse<T> success(String message, T data) {
        return new ApiResponse<>(true, message, data);
    }

    public static <T> ApiResponse<T> error(String message) {
        return new ApiResponse<>(false, message);
    }

}