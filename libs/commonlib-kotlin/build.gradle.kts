import java.util.*

plugins {
    kotlin("jvm") version "2.1.10"
    kotlin("plugin.spring") version "2.1.10"
    id("org.springframework.boot") version "3.4.3"
    id("io.spring.dependency-management") version "1.1.7"

    id("maven-publish")

    id("com.diffplug.spotless") version "7.2.1"
}

group = findProperty("group") as String
version = findProperty("version") as String

kotlin {
    jvmToolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

repositories {
    mavenLocal()
    mavenCentral()
    maven { url = uri("https://git.example.com/api/packages/example-corp/maven") }
}

dependencyManagement {
    imports {
        mavenBom("io.opentelemetry.instrumentation:opentelemetry-instrumentation-bom:2.16.0")
    }
}

dependencies {
    implementation("io.opentelemetry.instrumentation:opentelemetry-spring-boot-starter")

    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

publishing {
    val props = Properties()
    val propsFile = file("package.properties")
    if (propsFile.exists()) {
        propsFile.inputStream().use { props.load(it) }
    }
    repositories {
        maven {
            name = "commonlib-kotlin"
            url = uri("https://git.example.com/api/packages/example-corp/maven")
            credentials {
                username = props.getProperty("username") ?: System.getenv("MAVEN_USERNAME")
                password = props.getProperty("password") ?: System.getenv("MAVEN_PASSWORD")
            }
        }
    }
    publications {
        create<MavenPublication>("release") {
            from(components["kotlin"])
            artifactId = "commonlib-kotlin"
            groupId = project.group.toString()
            version = project.version.toString()

            pom {
                name.set("Common Kotlin Library")
                description.set("A common library for Kotlin projects.")
                url.set("https://git.example.com/example-corp/commonlib-kotlin.git")
                scm {
                    url.set("https://git.example.com/example-corp/commonlib-kotlin.git")
                    connection.set("scm:git:https://git.example.com/example-corp/commonlib-kotlin.git")
                    developerConnection.set("scm:git:ssh://git.example.com/example-corp/commonlib-kotlin.git")
                }
                developers {
                    developer {
                        id.set("admin")
                        name.set("admin")
                        email.set("admin@example.com")
                        organization.set("Example Corp")
                        organizationUrl.set("https://example.com")
                    }
                }
                licenses {
                    license {
                        name.set("Apache License, Version 2.0")
                        url.set("https://www.apache.org/licenses/LICENSE-2.0.txt")
                    }
                }
            }
        }
    }
}

tasks.withType<Test> {
    useJUnitPlatform()
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
