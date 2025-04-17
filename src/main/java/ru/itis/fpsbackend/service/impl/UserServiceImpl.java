package ru.itis.fpsbackend.service.impl;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;
import ru.itis.fpsbackend.dto.UserRegisterRequest;
import ru.itis.fpsbackend.dto.UserResponse;
import ru.itis.fpsbackend.exception.BusinessException;
import ru.itis.fpsbackend.model.User;
import ru.itis.fpsbackend.repository.UserRepository;
import ru.itis.fpsbackend.service.UserService;

@Service
public class UserServiceImpl implements UserService {

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private BCryptPasswordEncoder passwordEncoder;

    @Override
    public UserResponse registerUser(UserRegisterRequest request) {
        // Проверка существования пользователя с таким именем
        if (userRepository.existsByUsername(request.getUsername())) {
            throw new BusinessException("Пользователь с таким именем уже существует");
        }

        // Проверка существования пользователя с таким email
        if (userRepository.existsByEmail(request.getEmail())) {
            throw new BusinessException("Пользователь с таким email уже существует");
        }

        // Хеширование пароля
        String hashedPassword = passwordEncoder.encode(request.getPassword());

        // Создание нового пользователя
        User user = new User();
        user.setUsername(request.getUsername());
        user.setEmail(request.getEmail());
        user.setPassword(hashedPassword);

        // Сохранение пользователя в базу данных
        User savedUser = userRepository.save(user);

        // Формирование ответа
        return mapUserToResponse(savedUser);
    }

    @Override
    public UserResponse mapUserToResponse(User user) {
        return new UserResponse(user.getId(), user.getUsername(), user.getEmail());
    }
}