import org.gradle.api.tasks.testing.logging.TestExceptionFormat
import org.gradle.api.tasks.testing.logging.TestLogEvent
import org.springframework.boot.gradle.tasks.bundling.BootJar
import java.util.*

plugins {
    kotlin("jvm") version "2.1.10"
    kotlin("plugin.spring") version "2.1.10"
    id("org.springframework.boot") version "3.4.3"
    id("io.spring.dependency-management") version "1.1.7"

    id("com.google.devtools.ksp") version "2.1.10-1.0.31"
    kotlin("plugin.serialization") version "2.1.10"
    kotlin("plugin.jpa") version "2.1.10"

    id("org.flywaydb.flyway") version "11.8.0"
}

group = findProperty("group") as String
version = findProperty("version") as String

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
    sourceSets {
        main {
            java {
                srcDirs(
                    "src/main/kotlin",
                    "src/main/generated",
                    "build/generated/ksp/main/kotlin"
                )
            }
            resources {
                srcDir("src/main/resources")
            }
        }
        test {
            java {
                srcDir("src/test/kotlin")
            }
            resources {
                srcDir("src/test/resources")
            }
        }
    }
}

repositories {
    mavenLocal()
    mavenCentral()
}

dependencyManagement {
    imports {
        mavenBom("org.springframework.cloud:spring-cloud-dependencies:2024.0.1")
        mavenBom("io.github.openfeign.querydsl:querydsl-bom:6.10.1")
        mavenBom("io.opentelemetry.instrumentation:opentelemetry-instrumentation-bom:2.16.0")
        mavenBom("io.micrometer:micrometer-bom:1.15.0")
    }
}

dependencies {
    implementation("org.jetbrains.kotlin:kotlin-reflect")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core")

    // web mvc
    implementation("org.springframework.boot:spring-boot-starter-web")

    // jpa
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    runtimeOnly("com.mysql:mysql-connector-j")

    // querydsl
    implementation("jakarta.persistence:jakarta.persistence-api")
    implementation(group = "com.querydsl", name = "querydsl-jpa", classifier = "jakarta")
    // querydsl openfeign
    implementation("io.github.openfeign.querydsl:querydsl-core")
    ksp("io.github.openfeign.querydsl:querydsl-ksp-codegen")
    testImplementation("io.github.openfeign.querydsl:querydsl-jpa")

    // grpc
    implementation("io.grpc:grpc-kotlin-stub:1.4.1")
    implementation("io.grpc:grpc-protobuf:1.70.0")
    implementation("com.google.protobuf:protobuf-kotlin:4.29.3")

    // conductor
    implementation("org.conductoross:conductor-client-spring:4.0.6")

    // observability
    implementation("io.opentelemetry.instrumentation:opentelemetry-spring-boot-starter")
    implementation("io.opentelemetry.instrumentation:opentelemetry-logback-appender-1.0:2.16.0-alpha")
    implementation("io.opentelemetry.instrumentation:opentelemetry-micrometer-1.5:2.16.0-alpha")
    implementation("io.micrometer:micrometer-registry-otlp")
    implementation("org.springframework.boot:spring-boot-starter-actuator")

    // custom
    implementation("com.codejsha.bookstore:commonlib-kotlin:0.1.0")

    implementation("org.springframework.boot:spring-boot-starter-validation")
    implementation("org.springframework.cloud:spring-cloud-starter-config")

    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5")
    testImplementation("org.jetbrains.kotlinx:kotlinx-coroutines-test")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

kotlin {
    compilerOptions {
        freeCompilerArgs.addAll("-Xjsr305=strict")
    }
}

val dbProps = Properties().apply {
    val secretFile = file("/vault/secrets/db.properties")
    if (secretFile.exists()) {
        load(secretFile.inputStream())
    }
}

flyway {
    url = dbProps.getProperty("db.url", "")
    user = dbProps.getProperty("db.user", "")
    password = dbProps.getProperty("db.password", "")
    locations = arrayOf("classpath:db/migrations")
}

tasks.withType<Test> {
    useJUnitPlatform()
    testLogging {
        events(
            TestLogEvent.FAILED,
            TestLogEvent.PASSED,
            TestLogEvent.SKIPPED
        )
        debug {
            events(
                TestLogEvent.FAILED,
                TestLogEvent.PASSED,
                TestLogEvent.SKIPPED,
                TestLogEvent.STANDARD_OUT,
                TestLogEvent.STANDARD_ERROR
            )
            showStackTraces = true
            exceptionFormat = TestExceptionFormat.FULL
        }
    }
}

tasks.withType<BootJar> {
    archiveAppendix.set("app")
    archiveVersion.set("")
}

tasks.withType<Jar> {
    archiveVersion.set("")
}

tasks.processResources {
    duplicatesStrategy = DuplicatesStrategy.INCLUDE
}
