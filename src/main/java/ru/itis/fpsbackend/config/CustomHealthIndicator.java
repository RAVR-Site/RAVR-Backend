package ru.itis.fpsbackend.config;

import org.springframework.boot.actuate.health.Health;
import org.springframework.boot.actuate.health.HealthIndicator;
import org.springframework.stereotype.Component;

import java.util.Random;

/**
 * Пользовательский индикатор здоровья для мониторинга дополнительных компонентов системы
 */
@Component
public class CustomHealthIndicator implements HealthIndicator {

    private final Random random = new Random();
    
    @Override
    public Health health() {
        // Проверяем доступную память в JVM
        long freeMemory = Runtime.getRuntime().freeMemory();
        long totalMemory = Runtime.getRuntime().totalMemory();
        double memoryUsagePercent = 100 - ((double) freeMemory / totalMemory) * 100;
        
        // Имитируем проверку внешних сервисов (например, сторонние API)
        boolean externalServiceStatus = checkExternalService();
        
        // Если общее состояние нормальное - возвращаем UP с дополнительными данными
        if (memoryUsagePercent < 90 && externalServiceStatus) {
            return Health.up()
                    .withDetail("memory", String.format("%.2f%%", memoryUsagePercent))
                    .withDetail("externalServices", "available")
                    .build();
        }
        
        // Если есть проблемы - возвращаем состояние DOWN с информацией о проблеме
        return Health.down()
                .withDetail("memory", String.format("%.2f%%", memoryUsagePercent))
                .withDetail("externalServices", externalServiceStatus ? "available" : "unavailable")
                .withDetail("message", "System resources are under pressure")
                .build();
    }
    
    /**
     * Имитация проверки внешнего сервиса
     * В реальном приложении здесь будет фактическая проверка доступности внешних сервисов
     */
    private boolean checkExternalService() {
        // 95% вероятность, что сервис доступен
        return random.nextDouble() < 0.95;
    }
}