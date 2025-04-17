package ru.itis.fpsbackend.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotBlank;
import lombok.*;

@AllArgsConstructor
@NoArgsConstructor
@Builder
@Getter
@Setter
@Schema(description = "Запрос на обновление токена доступа")
public class TokenRefreshRequest {
    @NotBlank(message = "Refresh token cannot be blank")
    @Schema(description = "Refresh токен для обновления", example = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", required = true)
    private String refreshToken;
}