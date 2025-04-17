package ru.itis.fpsbackend.service;

import ru.itis.fpsbackend.dto.UserRegisterRequest;
import ru.itis.fpsbackend.dto.UserResponse;
import ru.itis.fpsbackend.model.User;

public interface UserService {
    
    /**
     * Регистрирует нового пользователя
     * 
     * @param request данные пользователя для регистрации
     * @return информация о зарегистрированном пользователе
     */
    UserResponse registerUser(UserRegisterRequest request);
    
    /**
     * Преобразует сущность User в DTO-ответ
     * 
     * @param user сущность пользователя
     * @return DTO с данными пользователя
     */
    UserResponse mapUserToResponse(User user);
}