package ru.itis.fpsbackend.dto;

import jakarta.validation.constraints.NotBlank;
import lombok.*;

@AllArgsConstructor
@NoArgsConstructor
@Builder
@Getter
@Setter
public class TokenRefreshRequest {
    @NotBlank(message = "Refresh token cannot be blank")
    private String refreshToken;
}