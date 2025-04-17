# Этап сборки
FROM maven:3.9.6-eclipse-temurin-21 AS build

WORKDIR /app
COPY pom.xml .
# Кеширование зависимостей для ускорения последующих сборок
RUN mvn dependency:go-offline -B

COPY src ./src
COPY checkstyle.xml ./checkstyle.xml
RUN mvn package -B -DskipTests

# Этап запуска
FROM eclipse-temurin:21-jre-alpine AS runtime

LABEL maintainer="fastrapier"
LABEL application="ravr-backend"

# Установка wget для HEALTHCHECK
RUN apk add --no-cache wget

# Создание непривилегированного пользователя
RUN addgroup --system --gid 1001 appuser && \
    adduser --system --uid 1001 --ingroup appuser appuser

WORKDIR /app

# Предоставление переменных среды для настройки JVM
ENV JAVA_OPTS="-XX:+UseContainerSupport -XX:MaxRAMPercentage=75.0"

# Копирование JAR-файла из этапа сборки
COPY --from=build /app/target/*.jar /app/app.jar

# Установка прав на файлы
RUN chown -R appuser:appuser /app
USER appuser

# Указание порта, который использует приложение
EXPOSE 8080

# Проверка работоспособности приложения
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD wget -q --spider http://localhost:8080/actuator/health || exit 1

# Запуск приложения
ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar /app/app.jar"]