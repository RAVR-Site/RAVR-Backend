package ru.itis.fpsbackend.service.impl;

import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import ru.itis.fpsbackend.dto.JwtResponse;
import ru.itis.fpsbackend.model.Token;
import ru.itis.fpsbackend.model.User;
import ru.itis.fpsbackend.repository.TokenRepository;
import ru.itis.fpsbackend.repository.UserRepository;
import ru.itis.fpsbackend.security.JwtTokenProvider;
import ru.itis.fpsbackend.security.UserDetailsImpl;
import ru.itis.fpsbackend.service.TokenService;

import java.time.LocalDateTime;

@Service
public class TokenServiceImpl implements TokenService {

    private final TokenRepository tokenRepository;
    private final UserRepository userRepository;
    private final JwtTokenProvider jwtTokenProvider;

    public TokenServiceImpl(TokenRepository tokenRepository,
                            UserRepository userRepository,
                            JwtTokenProvider jwtTokenProvider) {
        this.tokenRepository = tokenRepository;
        this.userRepository = userRepository;
        this.jwtTokenProvider = jwtTokenProvider;
    }

    @Override
    @Transactional
    public JwtResponse generateTokens(UserDetailsImpl userDetails) {
        // Генерация токенов
        String accessToken = jwtTokenProvider.generateAccessToken(userDetails);
        String refreshToken = jwtTokenProvider.generateRefreshToken(userDetails);

        // Получаем пользователя и сохраняем токены в базу
        User user = userRepository.findById(userDetails.getId())
                .orElseThrow(() -> new RuntimeException("Пользователь не найден"));

        LocalDateTime accessTokenExpiresAt = jwtTokenProvider.getExpirationDateFromToken(accessToken);
        LocalDateTime refreshTokenExpiresAt = jwtTokenProvider.getExpirationDateFromToken(refreshToken);

        Token token = new Token(user, accessToken, refreshToken, accessTokenExpiresAt, refreshTokenExpiresAt);
        tokenRepository.save(token);

        // Создание ответа без информации о ролях
        return new JwtResponse(
                accessToken,
                refreshToken,
                userDetails.getId(),
                userDetails.getUsername(),
                userDetails.getEmail()
        );
    }

    @Override
    @Transactional
    public JwtResponse refreshToken(String refreshToken) {
        // Проверяем валидность токена
        if (!jwtTokenProvider.validateToken(refreshToken)) {
            throw new RuntimeException("Refresh токен недействителен");
        }

        // Ищем токен в базе
        Token storedToken = tokenRepository.findByRefreshToken(refreshToken)
                .orElseThrow(() -> new RuntimeException("Refresh токен не найден"));

        // Проверяем, что токен не истек
        if (storedToken.isRefreshTokenExpired()) {
            tokenRepository.delete(storedToken);
            throw new RuntimeException("Refresh токен истек");
        }

        // Получаем пользователя
        User user = storedToken.getUser();

        // Создаем UserDetails для генерации новых токенов
        UserDetailsImpl userDetails = UserDetailsImpl.build(user);

        // Генерируем новую пару токенов
        String newAccessToken = jwtTokenProvider.generateAccessToken(userDetails);
        String newRefreshToken = jwtTokenProvider.generateRefreshToken(userDetails);

        // Обновляем токены в базе
        LocalDateTime accessTokenExpiresAt = jwtTokenProvider.getExpirationDateFromToken(newAccessToken);
        LocalDateTime refreshTokenExpiresAt = jwtTokenProvider.getExpirationDateFromToken(newRefreshToken);

        storedToken.setAccessToken(newAccessToken);
        storedToken.setRefreshToken(newRefreshToken);
        storedToken.setAccessTokenExpiresAt(accessTokenExpiresAt);
        storedToken.setRefreshTokenExpiresAt(refreshTokenExpiresAt);

        tokenRepository.save(storedToken);

        // Создаем ответ без информации о ролях
        return new JwtResponse(
                newAccessToken,
                newRefreshToken,
                userDetails.getId(),
                userDetails.getUsername(),
                userDetails.getEmail()
        );
    }

    @Override
    @Transactional
    public void invalidateAllUserTokens(User user) {
        tokenRepository.deleteAllByUser(user);
    }
}
