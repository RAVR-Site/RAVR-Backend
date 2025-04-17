package ru.itis.fpsbackend.controller;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;
import ru.itis.fpsbackend.dto.*;
import ru.itis.fpsbackend.security.UserDetailsImpl;
import ru.itis.fpsbackend.service.TokenService;
import ru.itis.fpsbackend.service.UserService;

@RestController
@RequestMapping("/api/auth")
@Tag(name = "Аутентификация", description = "API для авторизации, регистрации и управления токенами")
public class AuthController {

    private final AuthenticationManager authenticationManager;
    private final UserService userService;
    private final TokenService tokenService;

    public AuthController(AuthenticationManager authenticationManager,
                          UserService userService,
                          TokenService tokenService) {
        this.authenticationManager = authenticationManager;
        this.userService = userService;
        this.tokenService = tokenService;
    }

    @Operation(
            summary = "Аутентификация пользователя", 
            description = "Позволяет пользователю войти в систему, используя имя пользователя и пароль. " +
                    "Возвращает JWT токены для авторизации запросов."
    )
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Успешная аутентификация",
                    content = { @Content(mediaType = "application/json", 
                            schema = @Schema(implementation = ru.itis.fpsbackend.dto.ApiResponse.class)) }),
            @ApiResponse(responseCode = "401", description = "Неверные учетные данные", 
                    content = @Content)
    })
    @PostMapping("/login")
    public ResponseEntity<?> authenticateUser(@Valid @RequestBody LoginRequest loginRequest) {
        Authentication authentication = authenticationManager.authenticate(
                new UsernamePasswordAuthenticationToken(loginRequest.getUsername(), loginRequest.getPassword()));

        SecurityContextHolder.getContext().setAuthentication(authentication);

        UserDetailsImpl userDetails = (UserDetailsImpl) authentication.getPrincipal();

        JwtResponse jwtResponse = tokenService.generateTokens(userDetails);

        return ResponseEntity.ok(ru.itis.fpsbackend.dto.ApiResponse.success("Вход выполнен успешно", jwtResponse));
    }

    @Operation(
            summary = "Регистрация нового пользователя", 
            description = "Создает новую учетную запись пользователя в системе."
    )
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Пользователь успешно зарегистрирован",
                    content = { @Content(mediaType = "application/json", 
                            schema = @Schema(implementation = ru.itis.fpsbackend.dto.ApiResponse.class)) }),
            @ApiResponse(responseCode = "400", description = "Ошибка валидации данных или пользователь уже существует", 
                    content = { @Content(mediaType = "application/json", 
                            schema = @Schema(implementation = ru.itis.fpsbackend.dto.ApiResponse.class)) })
    })
    @PostMapping("/register")
    public ResponseEntity<?> registerUser(@Valid @RequestBody UserRegisterRequest registerRequest) {

        UserResponse userResponse = userService.registerUser(registerRequest);

        return ResponseEntity.ok(
                ru.itis.fpsbackend.dto.ApiResponse.success("Пользователь успешно зарегистрирован",
                        userResponse)
        );
    }

    @Operation(
            summary = "Обновление JWT токена", 
            description = "Получает новый access токен, используя действительный refresh токен."
    )
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Токен успешно обновлен",
                    content = { @Content(mediaType = "application/json", 
                            schema = @Schema(implementation = ru.itis.fpsbackend.dto.ApiResponse.class)) }),
            @ApiResponse(responseCode = "400", description = "Недействительный refresh токен", 
                    content = { @Content(mediaType = "application/json", 
                            schema = @Schema(implementation = ru.itis.fpsbackend.dto.ApiResponse.class)) })
    })
    @PostMapping("/refresh")
    public ResponseEntity<?> refreshToken(@Valid @RequestBody TokenRefreshRequest request) {
        String refreshToken = request.getRefreshToken();

        try {
            JwtResponse jwtResponse = tokenService.refreshToken(refreshToken);
            return ResponseEntity.ok(ru.itis.fpsbackend.dto.ApiResponse.success("Токен успешно обновлен", jwtResponse));
        } catch (Exception e) {
            return ResponseEntity.badRequest()
                    .body(ru.itis.fpsbackend.dto.ApiResponse.error(e.getMessage()));
        }
    }
}