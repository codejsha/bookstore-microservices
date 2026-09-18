package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.order.config.properties.TemporalWorkerConfig
import com.codejsha.bookstore.order.domain.workflow.OrderActivities
import io.temporal.worker.WorkerFactory
import org.junit.jupiter.api.Test
import org.mockito.Mockito.mock
import kotlin.test.assertEquals
import kotlin.test.assertTrue

class TemporalManagerTest {

    private fun manager(config: TemporalWorkerConfig): TemporalManager =
        TemporalManager(mock(WorkerFactory::class.java), mock(OrderActivities::class.java), config)

    @Test
    fun `workerOptions_configuredLimits_carryThemOntoTheWorker`() {
        val options = manager(
            TemporalWorkerConfig(
                maxConcurrentActivityExecutions = 8,
                maxConcurrentWorkflowTaskExecutions = 12,
                activityTaskPollers = 3,
                workflowTaskPollers = 4,
            ),
        ).workerOptions()

        assertEquals(8, options.maxConcurrentActivityExecutionSize)
        assertEquals(12, options.maxConcurrentWorkflowTaskExecutionSize)
        assertEquals(3, options.maxConcurrentActivityTaskPollers)
        assertEquals(4, options.maxConcurrentWorkflowTaskPollers)
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
