FROM alpine/java:21-jre
LABEL authors="daniilstudenikin"

WORKDIR /app

COPY target/*.jar /app/app.jar

ENTRYPOINT ["java", "/app/app.jar"]