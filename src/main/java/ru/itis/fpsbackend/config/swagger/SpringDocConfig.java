package ru.itis.fpsbackend.config.swagger;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import org.springdoc.core.customizers.OpenApiCustomizer;
import org.springdoc.core.customizers.OperationCustomizer;
import org.springframework.web.method.HandlerMethod;
import io.swagger.v3.oas.models.Operation;

/**
 * Конфигурация для настройки сериализации объектов в OpenAPI документации
 */
@Configuration
public class SpringDocConfig {

    /**
     * Кастомизатор OpenAPI для исключения проблемных схем или полей, 
     * которые могут вызывать ошибки при генерации документации
     */
    @Bean
    public OpenApiCustomizer openApiCustomizer() {
        return openApi -> {
            // Здесь можно настроить какие схемы нужно исключить или модифицировать
            // Это поможет избежать ошибок с циклическими ссылками
            
            // Пример: если нужно удалить определенную схему
            // openApi.getComponents().getSchemas().remove("ProblemSchemaName");
            
            // Для отладки можно добавить информацию о версии API
            openApi.getInfo().setDescription(openApi.getInfo().getDescription() + 
                    "\n\nAPI Version: 1.0.0" + 
                    "\nБолее подробная документация доступна на отдельных вкладках по группам API.");
        };
    }

    /**
     * Кастомизатор операций для добавления дополнительной информации к методам контроллера
     */
    @Bean
    public OperationCustomizer operationCustomizer() {
        return (Operation operation, HandlerMethod handlerMethod) -> {
            // Здесь можно настроить документацию для отдельных операций
            // Например, добавить заметки или пометки
            return operation;
        };
    }

    /**
     * Настроенный ObjectMapper для OpenAPI
     */
    @Bean
    @Primary
    public ObjectMapper objectMapper() {
        ObjectMapper objectMapper = new ObjectMapper();
        // Отключаем сериализацию для пустых значений, это может помочь с некоторыми проблемами
        objectMapper.setSerializationInclusion(com.fasterxml.jackson.annotation.JsonInclude.Include.NON_NULL);
        return objectMapper;
    }
}