package ru.itis.fpsbackend.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotBlank;
import lombok.Getter;
import lombok.Setter;

@Setter
@Getter
@Schema(description = "Запрос на аутентификацию пользователя")
public class LoginRequest {
    // Геттеры и сеттеры
    @NotBlank(message = "Имя пользователя не может быть пустым")
    @Schema(description = "Имя пользователя", example = "user1", required = true)
    private String username;

    @NotBlank(message = "Пароль не может быть пустым")
    @Schema(description = "Пароль пользователя", example = "password123", required = true)
    private String password;

}