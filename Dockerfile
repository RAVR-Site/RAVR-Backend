FROM maven:3.9.9-eclipse-temurin-24-alpine
LABEL authors="daniilstudenikin"

COPY src /app/src
COPY pom.xml /app/pom.xml
COPY mvnw /app/mvnw
COPY mvnw.cmd /app/mvnw.cmd

WORKDIR /app

RUN mvn clean package

ENTRYPOINT ["ls", "/app/target"]