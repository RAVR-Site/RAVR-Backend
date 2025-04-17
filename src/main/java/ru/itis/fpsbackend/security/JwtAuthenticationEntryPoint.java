package ru.itis.fpsbackend.security;

import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.InsufficientAuthenticationException;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.AuthenticationEntryPoint;
import org.springframework.stereotype.Component;
import ru.itis.fpsbackend.dto.ApiResponse;

import java.io.IOException;
import java.io.Serial;
import java.io.Serializable;

@Component
public class JwtAuthenticationEntryPoint implements AuthenticationEntryPoint, Serializable {

    @Serial
    private static final long serialVersionUID = 1L;
    
    private final ObjectMapper objectMapper = new ObjectMapper();

    @Override
    public void commence(HttpServletRequest request, HttpServletResponse response,
                         AuthenticationException authException) throws IOException {
        
        String errorMessage = "Ошибка аутентификации";
        
        // Определение типа ошибки аутентификации и соответствующего сообщения
        if (authException instanceof BadCredentialsException) {
            errorMessage = "Неверное имя пользователя или пароль";
        } else if (authException instanceof InsufficientAuthenticationException) {
            errorMessage = "Для доступа к этому ресурсу требуется аутентификация";
        } else if (authException.getCause() != null) {
            errorMessage = authException.getCause().getMessage();
        } else if (authException.getMessage() != null) {
            errorMessage = authException.getMessage();
        }
        
        // Создание объекта ApiResponse для сериализации в JSON
        ApiResponse<?> apiResponse = ApiResponse.error(errorMessage);
        
        // Настройка HTTP-ответа
        response.setStatus(HttpStatus.UNAUTHORIZED.value());
        response.setContentType(MediaType.APPLICATION_JSON_VALUE);
        
        // Запись JSON-ответа в поток вывода
        objectMapper.writeValue(response.getOutputStream(), apiResponse);
    }
}