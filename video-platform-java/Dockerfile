FROM eclipse-temurin:21-jre-alpine

WORKDIR /app

RUN addgroup -S videoplatform && adduser -S videoplatform -G videoplatform

COPY target/video-platform-java-1.0.0-SNAPSHOT.jar app.jar

USER videoplatform

EXPOSE 8081

ENTRYPOINT ["java", "-XX:+UseContainerSupport", "-XX:MaxRAMPercentage=75.0", "-jar", "app.jar"]
