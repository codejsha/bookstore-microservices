package com.codejsha.bookstore.payment.infrastructure.support

import org.flywaydb.core.Flyway
import org.springframework.boot.test.util.TestPropertyValues
import org.springframework.context.ApplicationContextInitializer
import org.springframework.context.ConfigurableApplicationContext
import org.springframework.test.context.ContextConfiguration
import org.testcontainers.containers.PostgreSQLContainer
import org.testcontainers.junit.jupiter.Testcontainers

@Testcontainers(disabledWithoutDocker = true)
@ContextConfiguration(initializers = [AbstractPostgresTestSupport.Initializer::class])
abstract class AbstractPostgresTestSupport {
    companion object {
        val container = PostgreSQLContainer("postgres:16").apply {
            withDatabaseName("testdb")
            withUsername("postgres")
            withPassword("test")
        }
    }

    internal class Initializer : ApplicationContextInitializer<ConfigurableApplicationContext> {
        override fun initialize(configurableApplicationContext: ConfigurableApplicationContext) {
            container.start()
            Flyway.configure()
                .dataSource(container.jdbcUrl, container.username, container.password)
                .locations("filesystem:db/migrations")
                .load()
                .migrate()
            TestPropertyValues.of(
                "spring.datasource.driver-class-name=org.postgresql.Driver",
                "spring.datasource.url=${container.jdbcUrl}",
                "spring.datasource.username=${container.username}",
                "spring.datasource.password=${container.password}",
            ).applyTo(configurableApplicationContext.environment)
        }
    }
}
