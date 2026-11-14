package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.CART
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.CART_ITEM
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDERS
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDER_ADJUSTMENT
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDER_ITEM
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.ORDER_SHIPPING
import com.codejsha.bookstore.order.domain.model.command.OrderCreateCommand
import com.codejsha.bookstore.order.infrastructure.adapter.pgsql.CartItemRepoImpl
import com.codejsha.bookstore.order.infrastructure.adapter.pgsql.CartRepoImpl
import com.codejsha.bookstore.order.infrastructure.adapter.pgsql.OrderAdjustmentRepoImpl
import com.codejsha.bookstore.order.infrastructure.adapter.pgsql.OrderItemRepoImpl
import com.codejsha.bookstore.order.infrastructure.adapter.pgsql.OrderRepoImpl
import com.codejsha.bookstore.order.infrastructure.adapter.pgsql.OrderShippingRepoImpl
import com.codejsha.bookstore.order.support.OrderTestFixtures
import org.jooq.DSLContext
import org.junit.jupiter.api.BeforeEach
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.jooq.test.autoconfigure.JooqTest
import org.springframework.boot.test.util.TestPropertyValues
import org.springframework.context.ApplicationContextInitializer
import org.springframework.context.ConfigurableApplicationContext
import org.springframework.context.annotation.Import
import org.springframework.core.io.FileSystemResource
import org.springframework.jdbc.datasource.init.ScriptUtils
import org.springframework.test.context.ContextConfiguration
import org.springframework.transaction.annotation.Propagation
import org.springframework.transaction.annotation.Transactional
import org.testcontainers.containers.PostgreSQLContainer
import org.testcontainers.junit.jupiter.Testcontainers
import java.io.File
import java.math.BigDecimal
import java.sql.DriverManager
import java.sql.SQLException
import java.util.UUID

@JooqTest(
    properties = [
        "spring.cloud.config.enabled=false",
        "spring.cloud.config.import-check.enabled=false",
    ],
)
@Import(
    PersistenceManager::class,
    JooqTransactionRunner::class,
    OrderRepoImpl::class,
    OrderItemRepoImpl::class,
    OrderShippingRepoImpl::class,
    OrderAdjustmentRepoImpl::class,
    CartRepoImpl::class,
    CartItemRepoImpl::class,
)
@Transactional(propagation = Propagation.NOT_SUPPORTED)
@Testcontainers(disabledWithoutDocker = true)
@ContextConfiguration(initializers = [AbstractPostgresTestSupport.Initializer::class])
abstract class AbstractPostgresTestSupport {

    @Autowired
    protected lateinit var dsl: DSLContext

    @Autowired
    protected lateinit var orderRepo: OrderRepoImpl

    protected val ctx = OrderTestFixtures.DEFAULT_CONTEXT

    @BeforeEach
    fun cleanTables() {
        dsl.deleteFrom(ORDER_ITEM).execute()
        dsl.deleteFrom(ORDER_ADJUSTMENT).execute()
        dsl.deleteFrom(ORDER_SHIPPING).execute()
        dsl.deleteFrom(ORDERS).execute()
        dsl.deleteFrom(CART_ITEM).execute()
        dsl.deleteFrom(CART).execute()
    }

    protected fun orderCommand(
        idempotencyKey: String = "idem-${UUID.randomUUID()}",
        userUid: UUID = OrderTestFixtures.USER_UID,
        totalAmount: BigDecimal = BigDecimal("10000.00"),
    ): OrderCreateCommand = OrderCreateCommand(
        userUid = userUid,
        currency = "KRW",
        itemsAmount = totalAmount,
        discountAmount = BigDecimal.ZERO,
        shippingAmount = BigDecimal.ZERO,
        taxAmount = BigDecimal.ZERO,
        totalAmount = totalAmount,
        idempotencyKey = idempotencyKey,
    )

    protected fun sqlStateOf(throwable: Throwable): String? =
        generateSequence(throwable) { it.cause }
            .filterIsInstance<SQLException>()
            .firstOrNull()
            ?.sqlState

    companion object {
        const val UNIQUE_VIOLATION = "23505"

        private val MIGRATION_FILE = Regex("""V(\d+)__.+\.sql""")

        val container = PostgreSQLContainer("postgres:16").apply {
            withDatabaseName("testdb")
            withUsername("postgres")
            withPassword("test")
        }

        private var migrated = false

        @Synchronized
        fun startAndMigrate() {
            container.start()
            if (migrated) return
            val scripts = File("db/migrations").listFiles().orEmpty()
                .mapNotNull { file -> MIGRATION_FILE.matchEntire(file.name)?.let { it.groupValues[1].toInt() to file } }
                .sortedBy { it.first }
            check(scripts.isNotEmpty()) { "no migration scripts under db/migrations" }
            DriverManager.getConnection(container.jdbcUrl, container.username, container.password).use { connection ->
                scripts.forEach { (_, file) -> ScriptUtils.executeSqlScript(connection, FileSystemResource(file)) }
            }
            migrated = true
        }
    }

    internal class Initializer : ApplicationContextInitializer<ConfigurableApplicationContext> {
        override fun initialize(configurableApplicationContext: ConfigurableApplicationContext) {
            startAndMigrate()
            TestPropertyValues.of(
                "spring.datasource.driver-class-name=org.postgresql.Driver",
                "spring.datasource.url=${container.jdbcUrl}",
                "spring.datasource.username=${container.username}",
                "spring.datasource.password=${container.password}",
            ).applyTo(configurableApplicationContext.environment)
        }
    }
}
