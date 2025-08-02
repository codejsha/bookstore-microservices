import org.springframework.boot.gradle.tasks.bundling.BootJar

plugins {
    kotlin("jvm") version "2.1.10"
    kotlin("plugin.spring") version "2.1.10"
    id("org.springframework.boot") version "3.4.3"
    id("io.spring.dependency-management") version "1.1.7"

    id("com.diffplug.spotless") version "7.2.1"
}

group = findProperty("group") as String
version = findProperty("version") as String

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

repositories {
    mavenCentral()
    maven { url = uri("https://git.example.com/api/packages/example-corp/maven") }
}

dependencyManagement {
    imports {
        mavenBom("org.springframework.cloud:spring-cloud-dependencies:2024.0.1")
    }
}

dependencies {
    implementation("org.jetbrains.kotlin:kotlin-reflect")
    implementation("org.springframework.cloud:spring-cloud-config-server")
    implementation("org.springframework.boot:spring-boot-starter-actuator")

    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

kotlin {
    compilerOptions {
        freeCompilerArgs.addAll("-Xjsr305=strict")
    }
}

tasks.withType<Test> {
    useJUnitPlatform()
}

tasks.withType<BootJar> {
    archiveVersion.set("")
}

tasks.withType<Jar> {
    archiveVersion.set("")
}

tasks.processResources {
    duplicatesStrategy = DuplicatesStrategy.INCLUDE
}

spotless {
    java {
        importOrder("com.codejsha.**", "|", "*", "|", "java.**", "javax.**", "|", "\$*")
        palantirJavaFormat()
        removeUnusedImports()
        formatAnnotations()
    }
    kotlin {
        ktlint()
            .setEditorConfigPath("$projectDir/.editorconfig")
        suppressLintsFor {
            step = "ktlint"
            shortCode = "standard:no-wildcard-imports"
        }
    }
    kotlinGradle {
        target("*.gradle.kts")
        ktlint()
            .setEditorConfigPath("$projectDir/.editorconfig")
        suppressLintsFor {
            step = "ktlint"
            shortCode = "standard:no-wildcard-imports"
        }
    }
}
