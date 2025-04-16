import org.gradle.api.tasks.testing.logging.TestExceptionFormat
import org.gradle.api.tasks.testing.logging.TestLogEvent
import org.springframework.boot.gradle.tasks.bundling.BootJar

plugins {
    kotlin("jvm") version "2.1.10"
    kotlin("plugin.spring") version "2.1.10"
    id("org.springframework.boot") version "3.4.3"
    id("io.spring.dependency-management") version "1.1.7"

    id("com.google.devtools.ksp") version "2.1.10-1.0.31"
    kotlin("plugin.serialization") version "2.1.10"
    kotlin("plugin.jpa") version "2.1.10"
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
    mavenCentral()
}

dependencyManagement {
    imports {
        mavenBom("org.springframework.cloud:spring-cloud-dependencies:2024.0.0")
        mavenBom("io.github.openfeign.querydsl:querydsl-bom:6.10.1")
        mavenBom("io.opentelemetry.instrumentation:opentelemetry-instrumentation-bom:2.14.0")
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

    // opentelemetry
    implementation("io.opentelemetry:opentelemetry-sdk")
    implementation("io.opentelemetry:opentelemetry-exporter-otlp")
    // opentelemetry instrumentation
    implementation("io.opentelemetry.instrumentation:opentelemetry-instrumentation-annotations")
    implementation("io.opentelemetry.semconv:opentelemetry-semconv")

    implementation("org.springframework.boot:spring-boot-starter-validation")
    implementation("org.springframework.boot:spring-boot-starter-actuator")
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
