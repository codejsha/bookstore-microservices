package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.platform.shared.vault.VaultCredentialWatcher
import com.zaxxer.hikari.HikariDataSource
import org.junit.jupiter.api.Test
import org.springframework.boot.autoconfigure.AutoConfigurations
import org.springframework.boot.jdbc.autoconfigure.DataSourceAutoConfiguration
import org.springframework.boot.test.context.runner.ApplicationContextRunner
import javax.sql.DataSource
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNotSame

class ReadinessDataSourceAutoConfigurationTest {

    private val contextRunner =
        ApplicationContextRunner()
            .withConfiguration(AutoConfigurations.of(DataSourceAutoConfiguration::class.java))
            .withUserConfiguration(ReadinessDataSourceConfig::class.java)
            .withPropertyValues(
                "spring.datasource.driver-class-name=com.mysql.cj.jdbc.Driver",
                "spring.datasource.url=jdbc:mysql://localhost:3306/payment_db",
                "spring.datasource.username=app",
                "spring.datasource.password=app",
                "spring.datasource.hikari.maximum-pool-size=$APPLICATION_POOL_MAXIMUM_SIZE",
                "spring.datasource.hikari.minimum-idle=$APPLICATION_POOL_MAXIMUM_SIZE",
            )

    @Test
    fun `context_readinessPoolRegistered_keepsSingleApplicationDataSource`() {
        contextRunner.run { context ->
            assertEquals(1, context.getBeanNamesForType(DataSource::class.java).size)
            val applicationPool = context.getBean(DataSource::class.java) as HikariDataSource

            assertEquals(APPLICATION_POOL_MAXIMUM_SIZE, applicationPool.maximumPoolSize)
        }
    }

    @Test
    fun `context_readinessPoolRegistered_isolatesReadinessFromApplicationPool`() {
        contextRunner.run { context ->
            val applicationPool = context.getBean(DataSource::class.java) as HikariDataSource
            val readinessPool = assertNotNull(context.getBean(ReadinessDataSource::class.java).unwrapHikari())

            assertNotSame(applicationPool, readinessPool)
            assertEquals(ReadinessDataSourceConfig.MAXIMUM_POOL_SIZE, readinessPool.maximumPoolSize)
            assertEquals(ReadinessDataSourceConfig.POOL_NAME, readinessPool.poolName)
            assertEquals(APPLICATION_POOL_MAXIMUM_SIZE, applicationPool.maximumPoolSize)
        }
    }

    @Test
    fun `context_vaultWatchersRegistered_rotatesBothPools`() {
        contextRunner.withUserConfiguration(VaultCredentialConfig::class.java).run { context ->
            assertEquals(2, context.getBeanNamesForType(VaultCredentialWatcher::class.java).size)
        }
    }

    companion object {
        private const val APPLICATION_POOL_MAXIMUM_SIZE = 30
    }
}
