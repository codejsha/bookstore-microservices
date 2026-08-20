package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.order.domain.workflow.OrderActivities
import com.codejsha.bookstore.order.domain.workflow.OrderCancellationWorkflowImpl
import com.codejsha.bookstore.order.domain.workflow.OrderFulfillmentWorkflowImpl
import com.codejsha.bookstore.order.domain.workflow.OrderPlacementWorkflowImpl
import io.temporal.worker.WorkerFactory
import jakarta.annotation.PostConstruct
import jakarta.annotation.PreDestroy
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Component
import java.util.concurrent.TimeUnit

@Component
class TemporalManager(
    private val workerFactory: WorkerFactory,
    private val orderActivities: OrderActivities,
) {
    private val log = LoggerFactory.getLogger(TemporalManager::class.java)

    @PostConstruct
    fun startWorkers() {
        val worker = workerFactory.newWorker("order-task-queue")
        worker.registerWorkflowImplementationTypes(
            OrderPlacementWorkflowImpl::class.java,
            OrderCancellationWorkflowImpl::class.java,
            OrderFulfillmentWorkflowImpl::class.java,
        )
        worker.registerActivitiesImplementations(orderActivities)
        workerFactory.start()
    }

    @PreDestroy
    fun stopWorkers() {
        workerFactory.shutdown()
        workerFactory.awaitTermination(SHUTDOWN_TIMEOUT_SECONDS, TimeUnit.SECONDS)
        log.info("temporal worker factory shutdown complete")
    }

    companion object {
        private const val SHUTDOWN_TIMEOUT_SECONDS = 30L
    }
}
