package com.codejsha.bookstore.payment.infrastructure.support

import com.zaxxer.hikari.HikariDataSource
import javax.sql.DataSource

class ReadinessDataSource(private val dataSource: DataSource) : AutoCloseable {

    fun isReachable(validationTimeoutSeconds: Int): Boolean =
        dataSource.connection.use { it.isValid(validationTimeoutSeconds) }

    fun unwrapHikari(): HikariDataSource? =
        when (dataSource) {
            is HikariDataSource -> dataSource
            else -> runCatching { dataSource.unwrap(HikariDataSource::class.java) }.getOrNull()
        }

    override fun close() {
        unwrapHikari()?.close()
    }
}
