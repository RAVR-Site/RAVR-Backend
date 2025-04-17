package ru.itis.fpsbackend.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;
import ru.itis.fpsbackend.model.Token;
import ru.itis.fpsbackend.model.User;

import java.util.List;
import java.util.Optional;

@Repository
public interface TokenRepository extends JpaRepository<Token, Long> {

    Optional<Token> findByAccessToken(String accessToken);

    Optional<Token> findByRefreshToken(String refreshToken);

    List<Token> findAllByUser(User user);

    void deleteAllByUser(User user);

    @Query("SELECT t FROM Token t WHERE t.user = ?1 AND t.refreshTokenExpiresAt > CURRENT_TIMESTAMP")
    List<Token> findAllValidTokensByUser(User user);
}