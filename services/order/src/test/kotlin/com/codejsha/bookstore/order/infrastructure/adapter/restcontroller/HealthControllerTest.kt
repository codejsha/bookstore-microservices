package com.codejsha.bookstore.order.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.order.infrastructure.support.ReadinessDataSource
import org.junit.jupiter.api.Test
import org.mockito.ArgumentMatchers.anyInt
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.springframework.http.HttpStatus
import java.sql.Connection
import java.sql.SQLException
import javax.sql.DataSource
import kotlin.test.assertEquals
import kotlin.test.assertTrue

class HealthControllerTest {

    private fun readinessDataSource(valid: Boolean): ReadinessDataSource {
        val connection = mock(Connection::class.java)
        given(connection.isValid(anyInt())).willReturn(valid)
        val dataSource = mock(DataSource::class.java)
        given(dataSource.connection).willReturn(connection)
        return ReadinessDataSource(dataSource)
    }

    @Test
    fun `health_always_returns200`() {
        val controller = HealthController(readinessDataSource(valid = true))

        assertEquals(HttpStatus.OK, controller.health().statusCode)
    }

    @Test
    fun `ready_readinessPoolValid_returns200`() {
        val controller = HealthController(readinessDataSource(valid = true))

        assertEquals(HttpStatus.OK, controller.ready().statusCode)
    }

    @Test
    fun `ready_readinessPoolInvalid_returns503`() {
        val controller = HealthController(readinessDataSource(valid = false))

        assertEquals(HttpStatus.SERVICE_UNAVAILABLE, controller.ready().statusCode)
    }

    @Test
    fun `ready_acquisitionFails_returns503`() {
        val dataSource = mock(DataSource::class.java)
        given(dataSource.connection).willThrow(SQLException("pool exhausted"))
        val controller = HealthController(ReadinessDataSource(dataSource))

        assertEquals(HttpStatus.SERVICE_UNAVAILABLE, controller.ready().statusCode)
    }

    @Test
    fun `ready_acquisitionBlocks_returns503WithinTimeout`() {
        val dataSource = mock(DataSource::class.java)
        given(dataSource.connection).willAnswer {
            Thread.sleep(10_000)
            mock(Connection::class.java)
        }
        val controller = HealthController(ReadinessDataSource(dataSource))

        val start = System.nanoTime()
        val status = controller.ready().statusCode
        val elapsedMillis = (System.nanoTime() - start) / 1_000_000

        assertEquals(HttpStatus.SERVICE_UNAVAILABLE, status)
        assertTrue(elapsedMillis < 6_000, "readiness took ${elapsedMillis}ms")
    }
}
