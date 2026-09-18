package com.codejsha.bookstore.settlement.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.settlement.infrastructure.support.ReadinessDataSource
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController
import java.util.concurrent.ExecutorService
import java.util.concurrent.Executors
import java.util.concurrent.TimeUnit

@RestController
class HealthController(private val readinessDataSource: ReadinessDataSource) {

    private val readinessExecutor: ExecutorService =
        Executors.newCachedThreadPool { r -> Thread(r, "health-ready").apply { isDaemon = true } }

    @GetMapping("/health")
    fun health(): ResponseEntity<Void> = ResponseEntity.ok().build()

    @GetMapping("/health/ready")
    fun ready(): ResponseEntity<Void> =
        if (isReachable()) {
            ResponseEntity.ok().build()
        } else {
            ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE).build()
        }

    private fun isReachable(): Boolean {
        val probe = readinessExecutor.submit<Boolean> {
            readinessDataSource.isReachable(VALIDATION_TIMEOUT_SECONDS)
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
