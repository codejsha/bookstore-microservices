package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.bookstore.payment.config.properties.TemporalWorkerConfig
import com.codejsha.bookstore.payment.domain.workflow.PaymentActivities
import io.temporal.worker.WorkerFactory
import org.junit.jupiter.api.Test
import org.mockito.Mockito.mock
import kotlin.test.assertEquals
import kotlin.test.assertTrue

class TemporalManagerTest {

    private fun manager(config: TemporalWorkerConfig): TemporalManager =
        TemporalManager(mock(WorkerFactory::class.java), mock(PaymentActivities::class.java), config)

    @Test
    fun `workerOptions_configuredLimits_carryThemOntoTheWorker`() {
        val options = manager(
            TemporalWorkerConfig(maxConcurrentActivityExecutions = 8, activityTaskPollers = 3),
        ).workerOptions()

        assertEquals(8, options.maxConcurrentActivityExecutionSize)
        assertEquals(3, options.maxConcurrentActivityTaskPollers)
    }

    @Test
    fun `workerOptions_defaults_stayWithinTheHikariPool`() {
        val options = manager(TemporalWorkerConfig()).workerOptions()

        assertEquals(16, options.maxConcurrentActivityExecutionSize)
        assertTrue(
            options.maxConcurrentActivityExecutionSize < HIKARI_MAXIMUM_POOL_SIZE,
            "activity concurrency ${options.maxConcurrentActivityExecutionSize} must leave room in the pool",
        )
    }

    private companion object {
        const val HIKARI_MAXIMUM_POOL_SIZE = 30
    }
}
