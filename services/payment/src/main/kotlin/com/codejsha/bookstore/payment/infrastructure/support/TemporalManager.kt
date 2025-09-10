package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.bookstore.payment.domain.workflow.PaymentActivities
import io.temporal.worker.WorkerFactory
import jakarta.annotation.PostConstruct
import jakarta.annotation.PreDestroy
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Component
import java.util.concurrent.TimeUnit

@Component
class TemporalManager(
    private val workerFactory: WorkerFactory,
    private val paymentActivities: PaymentActivities,
) {
    private val log = LoggerFactory.getLogger(TemporalManager::class.java)

    @PostConstruct
    fun startWorkers() {
        val worker = workerFactory.newWorker("payment-task-queue")
        worker.registerActivitiesImplementations(paymentActivities)
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
