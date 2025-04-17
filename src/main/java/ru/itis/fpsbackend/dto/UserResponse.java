package ru.itis.fpsbackend.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.*;
import ru.itis.fpsbackend.model.User;


@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
@Schema(description = "Ответ с информацией о пользователе")
public class UserResponse {
    @Schema(description = "ID пользователя", example = "1")
    private Long id;
    
    @Schema(description = "Имя пользователя", example = "user1")
    private String username;
    
    @Schema(description = "Email пользователя", example = "user@example.com")
    private String email;

    public static UserResponse fromUser(User user) {
        return new UserResponse(
                user.getId(),
                user.getUsername(),
                user.getEmail()
        );
    }
}