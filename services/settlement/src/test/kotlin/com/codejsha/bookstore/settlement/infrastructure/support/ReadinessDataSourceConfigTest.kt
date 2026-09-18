package com.codejsha.bookstore.settlement.infrastructure.support

import org.junit.jupiter.api.Test
import org.springframework.boot.jdbc.autoconfigure.DataSourceProperties
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

class ReadinessDataSourceConfigTest {

    private fun readinessDataSource(): ReadinessDataSource {
        val properties = DataSourceProperties().apply {
            driverClassName = "com.mysql.cj.jdbc.Driver"
            url = "jdbc:mysql://localhost:3306/settlement_db"
            username = "readiness"
            password = "readiness"
        }
        return ReadinessDataSourceConfig().readinessDataSource(properties)
    }

    @Test
    fun `readinessDataSource_built_capsPoolAtTwoConnections`() {
        val pool = assertNotNull(readinessDataSource().unwrapHikari())

        assertEquals(ReadinessDataSourceConfig.MAXIMUM_POOL_SIZE, pool.maximumPoolSize)
        assertEquals(ReadinessDataSourceConfig.MINIMUM_IDLE, pool.minimumIdle)
        assertEquals(ReadinessDataSourceConfig.POOL_NAME, pool.poolName)
    }

    @Test
    fun `readinessDataSource_built_staysFarBelowApplicationPool`() {
        val pool = assertNotNull(readinessDataSource().unwrapHikari())

        assertTrue(
            pool.maximumPoolSize < APPLICATION_POOL_MAXIMUM_SIZE,
            "readiness pool ${pool.maximumPoolSize} must stay below application pool $APPLICATION_POOL_MAXIMUM_SIZE",
        )
    }

    @Test
    fun `readinessDataSource_closed_releasesPool`() {
        val readiness = readinessDataSource()
        val pool = assertNotNull(readiness.unwrapHikari())

        readiness.close()

        assertTrue(pool.isClosed)
    }

    companion object {
        private const val APPLICATION_POOL_MAXIMUM_SIZE = 30
    }
}
