package com.codejsha.bookstore.order.infrastructure.support

import com.zaxxer.hikari.HikariDataSource
import org.springframework.boot.jdbc.autoconfigure.DataSourceProperties
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@Configuration
class ReadinessDataSourceConfig {

    @Bean(destroyMethod = "close")
    fun readinessDataSource(dataSourceProperties: DataSourceProperties): ReadinessDataSource {
        val pool =
            dataSourceProperties
                .initializeDataSourceBuilder()
                .type(HikariDataSource::class.java)
                .build()
                .apply {
                    poolName = POOL_NAME
                    maximumPoolSize = MAXIMUM_POOL_SIZE
                    minimumIdle = MINIMUM_IDLE
                    connectionTimeout = CONNECTION_TIMEOUT_MILLIS
                    validationTimeout = VALIDATION_TIMEOUT_MILLIS
                    maxLifetime = MAX_LIFETIME_MILLIS
                    isAutoCommit = true
                }
        return ReadinessDataSource(pool)
    }

    companion object {
        const val POOL_NAME = "order-readiness"
        const val MAXIMUM_POOL_SIZE = 2
        const val MINIMUM_IDLE = 1
        const val CONNECTION_TIMEOUT_MILLIS = 10_000L
        const val VALIDATION_TIMEOUT_MILLIS = 5_000L
        const val MAX_LIFETIME_MILLIS = 1_800_000L
    }
}
