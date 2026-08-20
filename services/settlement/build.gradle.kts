import org.gradle.api.tasks.testing.logging.TestExceptionFormat
import org.gradle.api.tasks.testing.logging.TestLogEvent
import org.jooq.meta.jaxb.*
import org.springframework.boot.gradle.tasks.bundling.BootJar

import java.util.*

plugins {
    kotlin("jvm") version "2.3.10"
    kotlin("plugin.spring") version "2.3.10"
    kotlin("plugin.serialization") version "2.3.10"
    id("org.springframework.boot") version "4.0.2"
    id("io.spring.dependency-management") version "1.1.7"

    id("com.codejsha.platform.jooq-codegen-plugin") version "0.1.0"
    id("org.jooq.jooq-codegen-gradle") version "3.19.29"
    id("com.diffplug.spotless") version "8.3.0"
}

group = findProperty("group") as String
version = findProperty("version") as String

kotlin {
    jvmToolchain {
        languageVersion = JavaLanguageVersion.of(25)
    }
    compilerOptions {
        freeCompilerArgs.addAll("-Xjsr305=strict")
        optIn.add("kotlin.uuid.ExperimentalUuidApi")
    }
}

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(25)
    }
}

sourceSets {
    main {
        kotlin {
            srcDir("src/main/kotlin")
        }
        java {
            srcDir("src/main/generated")
        }
        resources {
            srcDir("src/main/resources")
        }
    }
    test {
        kotlin {
            srcDir("src/test/kotlin")
        }
        resources {
            srcDir("src/test/resources")
        }
    }
}

repositories {
    mavenLocal()
    mavenCentral()
    maven {
        url = uri("https://maven.pkg.github.com/codejsha/package-repo")
    }
}

dependencyManagement {
    imports {
        mavenBom("org.springframework.cloud:spring-cloud-dependencies:2025.1.0")
        mavenBom("org.testcontainers:testcontainers-bom:2.0.3")
    }
}

dependencies {
    implementation("org.jetbrains.kotlin:kotlin-reflect")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core")
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5")
    testImplementation("org.jetbrains.kotlinx:kotlinx-coroutines-test")

    // web
    implementation("org.springframework.boot:spring-boot-starter-webmvc")
    implementation("org.springframework.boot:spring-boot-starter-validation")
    implementation("org.springframework.boot:spring-boot-starter-security")

    // data
    implementation("org.springframework.data:spring-data-commons")
    implementation("org.springframework.boot:spring-boot-data-commons")
    implementation("org.springframework.boot:spring-boot-starter-jooq")
    testImplementation("org.springframework.boot:spring-boot-starter-jooq-test")
    implementation("com.mysql:mysql-connector-j")
    testImplementation("org.flywaydb:flyway-core")
    testImplementation("org.flywaydb:flyway-mysql")

    // batch
    implementation("org.springframework.boot:spring-boot-starter-batch")
    testImplementation("org.springframework.batch:spring-batch-test")

    // integration
    implementation("org.springframework.cloud:spring-cloud-starter-config")
    compileOnly("io.opentelemetry.instrumentation:opentelemetry-instrumentation-annotations:2.24.0")
    implementation("org.springframework.boot:spring-boot-starter-actuator")
    implementation("io.micrometer:micrometer-registry-otlp")
    implementation("org.springframework.boot:spring-boot-opentelemetry")

    // custom
    implementation("com.codejsha.platform:shared-library-kotlin:0.1.0-SNAPSHOT")

    // test
    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
    testImplementation("org.testcontainers:testcontainers-junit-jupiter")
    testImplementation("org.testcontainers:testcontainers-mysql")
}

configurations.all {
    exclude(group = "io.opentelemetry.instrumentation", module = "opentelemetry-spring-boot-starter")
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

tasks.processTestResources {
    duplicatesStrategy = DuplicatesStrategy.INCLUDE
}

val dbProps =
    Properties().apply {
        val secretFile = file("/vault/secrets/db.properties")
        if (secretFile.exists() && secretFile.isFile) {
            secretFile.inputStream().use { fis ->
                load(fis)
            }
        }
    }
val dbName: String = dbProps.getProperty("db.name", "${project.name}_db")

val jooqGenerator =
    Generator().apply {
        name = "org.jooq.codegen.KotlinGenerator"
        database =
            Database().apply {
                name = "org.jooq.meta.mysql.MySQLDatabase"
                schemata = listOf(SchemaMappingType().apply { inputSchema = dbName })
                includes = "$dbName.*"
                excludes = "$dbName.flyway_schema_history|$dbName.BATCH_.*"
                forcedTypes =
                    listOf(
                        ForcedType().apply {
                            name = "BOOLEAN"
                            includeTypes = "TINYINT\\(1\\)"
                            includeExpression = ".*\\.deleted"
                        }
                    )
            }
        strategy =
            Strategy().apply {
                name = "org.jooq.codegen.DefaultGeneratorStrategy"
                matchers =
                    Matchers().apply {
                        catalogs =
                            listOf(
                                MatchersCatalogType().apply {
                                }
                            )
                        tables =
                            listOf(
                                MatchersTableType().apply {
                                    tableClass =
                                        MatcherRule().apply {
                                            transform = MatcherTransformType.PASCAL
                                            expression = "$0_table"
                                        }
                                    interfaceClass =
                                        MatcherRule().apply {
                                            transform = MatcherTransformType.PASCAL
                                            expression = "i_$0_entity"
                                        }
                                    pojoClass =
                                        MatcherRule().apply {
                                            transform = MatcherTransformType.PASCAL
                                            expression = "$0_entity"
                                        }
                                }
                            )
                        enums =
                            listOf(
                                MatchersEnumType().apply {
                                    enumClass =
                                        MatcherRule().apply {
                                            transform = MatcherTransformType.PASCAL
                                            expression = "$0"
                                        }
                                }
                            )
                    }
            }
        generate =
            Generate().apply {
                withTables(true)
                withRecords(true)
                withPojos(true)
                withInterfaces(true)
                withDaos(true)
                withSpringAnnotations(true)
                withSequences(true)
                withRoutines(true)
                withIndexes(true)
            }
        target =
            org.jooq.meta.jaxb.Target().apply {
                packageName = "com.codejsha.bookstore.generated.infrastructure.adapter.jooq"
                directory = "$projectDir/src/main/generated"
            }
    }

jooq {
    configuration {
        jdbc {
            driver = "com.mysql.cj.jdbc.Driver"
            url = dbProps.getProperty("db.url", "")
            user = dbProps.getProperty("db.user", "")
            password = dbProps.getProperty("db.password", "")
        }
        generator = jooqGenerator
        logging = org.jooq.meta.jaxb.Logging.INFO
        onError = OnError.FAIL
        onUnused = OnError.LOG
    }
}

jooqContainer {
    databaseType = "mysql"
    databaseName = dbName
    location = "filesystem:db/migrations"
    generator = jooqGenerator
}
