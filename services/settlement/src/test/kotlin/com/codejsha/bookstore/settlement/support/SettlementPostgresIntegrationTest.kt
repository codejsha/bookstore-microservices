package com.codejsha.bookstore.settlement.support

import org.flywaydb.core.Flyway
import org.junit.jupiter.api.BeforeEach
import org.springframework.batch.test.context.SpringBatchTest
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.jdbc.datasource.DriverManagerDataSource
import org.springframework.test.context.DynamicPropertyRegistry
import org.springframework.test.context.DynamicPropertySource
import org.testcontainers.containers.PostgreSQLContainer
import org.testcontainers.junit.jupiter.Testcontainers
import javax.sql.DataSource

@SpringBootTest(
    properties = [
        "spring.flyway.enabled=false",
        "spring.batch.job.enabled=false",
        "management.otlp.metrics.export.enabled=false",
    ],
)
@SpringBatchTest
@Testcontainers(disabledWithoutDocker = true)
abstract class SettlementPostgresIntegrationTest {

    @Autowired
    protected lateinit var dataSource: DataSource

    protected val settlementJdbc: JdbcTemplate by lazy { JdbcTemplate(dataSource) }

    protected val paymentJdbc: JdbcTemplate get() = paymentJdbcTemplate

    @BeforeEach
    fun cleanTables() {
        listOf("settlement_detail", "daily_settlement", "settlement_job_run", "settlement_run_lock")
            .forEach { settlementJdbc.execute("DELETE FROM $it") }
        listOf("payments", "refunds").forEach { paymentJdbc.execute("DELETE FROM $it") }
    }

    companion object {
        private const val PAYMENT_DB = "payment_db"

        private const val DB_USER = "postgres"
        private const val DB_PASSWORD = "test"
        private const val DRIVER = "org.postgresql.Driver"

        private val postgres: PostgreSQLContainer<*> by lazy {
            PostgreSQLContainer("postgres:16")
                .withDatabaseName("settlement_db")
                .withUsername(DB_USER)
                .withPassword(DB_PASSWORD)
                .also { it.start() }
        }

        private val paymentJdbcUrl: String
            get() = "jdbc:postgresql://${postgres.host}:${postgres.getMappedPort(5432)}/$PAYMENT_DB?sslmode=disable"

        private val paymentJdbcTemplate: JdbcTemplate by lazy {
            JdbcTemplate(
                DriverManagerDataSource(paymentJdbcUrl, DB_USER, DB_PASSWORD).apply { setDriverClassName(DRIVER) },
            )
        }

        private var migrated = false

        @Synchronized
        private fun migrateOnce() {
            if (migrated) return
            JdbcTemplate(
                DriverManagerDataSource(postgres.jdbcUrl, DB_USER, DB_PASSWORD).apply { setDriverClassName(DRIVER) },
            ).execute("CREATE DATABASE $PAYMENT_DB")

            Flyway.configure()
                .dataSource(postgres.jdbcUrl, DB_USER, DB_PASSWORD)
                .locations("filesystem:db/migrations")
                .load()
                .migrate()

            Flyway.configure()
                .dataSource(paymentJdbcUrl, DB_USER, DB_PASSWORD)
                .locations("classpath:payment-fixture")
                .load()
                .migrate()

            migrated = true
        }

        @JvmStatic
        @DynamicPropertySource
        fun properties(registry: DynamicPropertyRegistry) {
            migrateOnce()

            registry.add("spring.datasource.url") { postgres.jdbcUrl }
            registry.add("spring.datasource.username") { DB_USER }
            registry.add("spring.datasource.password") { DB_PASSWORD }
            registry.add("spring.datasource.driver-class-name") { DRIVER }

            registry.add("settlement.timezone") { "Asia/Seoul" }
            registry.add("settlement.chunk-size") { 2 }
            registry.add("settlement.stale-execution-timeout") { "PT5M" }
            registry.add("settlement.lock.ttl") { "PT5M" }
            registry.add("settlement.lock.heartbeat-interval") { "PT1H" }
            registry.add("settlement.payment-db.host") { postgres.host }
            registry.add("settlement.payment-db.port") { postgres.getMappedPort(5432) }
            registry.add("settlement.payment-db.params") { "sslmode=disable" }
            registry.add("payment.db.username") { DB_USER }
            registry.add("payment.db.password") { DB_PASSWORD }
            registry.add("payment.db.name") { PAYMENT_DB }
        }
    }
}
