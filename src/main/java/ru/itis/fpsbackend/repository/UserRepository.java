package ru.itis.fpsbackend.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import ru.itis.fpsbackend.model.User;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {
}