package ru.itis.fpsbackend.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
@Schema(description = "Запрос на регистрацию нового пользователя")
public class UserRegisterRequest {
    @NotBlank(message = "Имя пользователя не может быть пустым")
    @Size(min = 3, max = 50, message = "Имя пользователя должно содержать от 3 до 50 символов")
    @Schema(description = "Имя пользователя", example = "newuser", required = true, minLength = 3, maxLength = 50)
    private String username;

    @NotBlank(message = "Email не может быть пустым")
    @Email(message = "Неверный формат email")
    @Schema(description = "Email пользователя", example = "user@example.com", required = true)
    private String email;

    @NotBlank(message = "Пароль не может быть пустым")
    @Size(min = 6, message = "Пароль должен содержать минимум 6 символов")
    @Schema(description = "Пароль пользователя", example = "securePassword123", required = true, minLength = 6)
    private String password;
}