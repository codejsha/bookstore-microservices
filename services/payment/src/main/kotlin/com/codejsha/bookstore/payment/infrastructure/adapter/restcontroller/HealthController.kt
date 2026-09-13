package com.codejsha.bookstore.payment.infrastructure.adapter.restcontroller

import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController
import java.util.concurrent.ExecutorService
import java.util.concurrent.Executors
import java.util.concurrent.TimeUnit
import javax.sql.DataSource

@RestController
class HealthController(private val dataSources: List<DataSource>) {

    private val readinessExecutor: ExecutorService =
        Executors.newCachedThreadPool { r -> Thread(r, "health-ready").apply { isDaemon = true } }

    @GetMapping("/health")
    fun health(): ResponseEntity<Void> = ResponseEntity.ok().build()

    @GetMapping("/health/ready")
    fun ready(): ResponseEntity<Void> =
        if (dataSources.all { isReachable(it) }) {
            ResponseEntity.ok().build()
        } else {
            ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE).build()
        }

    private fun isReachable(dataSource: DataSource): Boolean {
        val probe = readinessExecutor.submit<Boolean> {
            dataSource.connection.use { it.isValid(VALIDATION_TIMEOUT_SECONDS) }
        }
        return try {
            probe.get(READINESS_TIMEOUT_SECONDS, TimeUnit.SECONDS)
        } catch (_: Exception) {
            probe.cancel(true)
            false
        }
    }

    companion object {
        private const val VALIDATION_TIMEOUT_SECONDS = 2
        private const val READINESS_TIMEOUT_SECONDS = 3L
    }
}
