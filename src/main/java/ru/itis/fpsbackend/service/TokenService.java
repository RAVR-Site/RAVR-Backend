package ru.itis.fpsbackend.service;

import ru.itis.fpsbackend.dto.JwtResponse;
import ru.itis.fpsbackend.model.User;
import ru.itis.fpsbackend.security.UserDetailsImpl;

public interface TokenService {
    JwtResponse generateTokens(UserDetailsImpl userDetails);
    JwtResponse refreshToken(String refreshToken);
    void invalidateAllUserTokens(User user);
}