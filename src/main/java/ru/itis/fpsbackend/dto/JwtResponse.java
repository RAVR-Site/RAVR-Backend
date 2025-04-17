package ru.itis.fpsbackend.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@Schema(description = "Ответ с JWT токенами после успешной аутентификации")
public class JwtResponse {
    @Schema(description = "Токен доступа (access token)", example = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
    private String accessToken;
    
    @Schema(description = "Токен обновления (refresh token)", example = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
    private String refreshToken;
    
    @Schema(description = "Тип токена", example = "Bearer", defaultValue = "Bearer")
    private final String tokenType = "Bearer";
    
    @Schema(description = "ID пользователя", example = "1")
    private Long id;
    
    @Schema(description = "Имя пользователя", example = "user1")
    private String username;
    
    @Schema(description = "Email пользователя", example = "user@example.com")
    private String email;

    public JwtResponse(String accessToken, String refreshToken, Long id, String username, String email) {
        this.accessToken = accessToken;
        this.refreshToken = refreshToken;
        this.id = id;
        this.username = username;
        this.email = email;
    }
}