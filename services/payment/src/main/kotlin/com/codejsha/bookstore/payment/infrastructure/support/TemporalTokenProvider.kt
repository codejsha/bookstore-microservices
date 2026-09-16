package com.codejsha.bookstore.payment.infrastructure.support

import io.temporal.authorization.AuthorizationTokenSupplier
import org.slf4j.LoggerFactory
import java.time.Clock
import java.time.Duration
import java.time.Instant
import java.util.concurrent.Executors
import java.util.concurrent.RejectedExecutionException
import java.util.concurrent.ScheduledExecutorService
import java.util.concurrent.TimeUnit

data class TemporalAccessToken(
    val value: String,
    val expiresIn: Duration,
)

fun interface TemporalTokenFetcher {
    fun fetch(): TemporalAccessToken
}

class TemporalTokenProvider(
    private val fetcher: TemporalTokenFetcher,
    private val clock: Clock = Clock.systemUTC(),
    private val refreshRatio: Double = DEFAULT_REFRESH_RATIO,
    private val minRetryBackoff: Duration = DEFAULT_MIN_RETRY_BACKOFF,
    private val maxRetryBackoff: Duration = DEFAULT_MAX_RETRY_BACKOFF,
) : AuthorizationTokenSupplier {
    private val log = LoggerFactory.getLogger(TemporalTokenProvider::class.java)

    @Volatile
    private var current: CachedToken? = null
    private var retryBackoff: Duration = minRetryBackoff
    private var scheduler: ScheduledExecutorService? = null

    init {
        require(refreshRatio > 0.0 && refreshRatio < 1.0) { "refreshRatio must be between 0 and 1" }
    }

    @Synchronized
    fun initialize() {
        store(fetcher.fetch())
        retryBackoff = minRetryBackoff
    }

    @Synchronized
    fun start() {
        if (current == null) {
            initialize()
        }
        val executor = Executors.newSingleThreadScheduledExecutor(
            Thread.ofVirtual().name(REFRESH_THREAD_NAME).factory(),
        )
        scheduler = executor
        schedule(executor, nextRefreshDelay())
        log.info("temporal access token refresh started, next refresh in {} ms", nextRefreshDelay().toMillis())
    }

    @Synchronized
    fun stop() {
        scheduler?.shutdownNow()
        scheduler = null
    }

    override fun supply(): String {
        val token = checkNotNull(current) { "temporal access token is not initialized" }
        return BEARER_PREFIX + token.value
    }

    @Synchronized
    fun refresh(): Duration =
        try {
            store(fetcher.fetch())
            retryBackoff = minRetryBackoff
            nextRefreshDelay()
        } catch (e: Exception) {
            val delay = retryBackoff
            retryBackoff = minOf(retryBackoff.multipliedBy(2), maxRetryBackoff)
            log.warn(
                "temporal access token refresh failed, retrying in {} ms, current token expires at {}",
                delay.toMillis(),
                current?.expiresAt,
                e,
            )
            delay
        }

    fun nextRefreshDelay(): Duration {
        val refreshAt = current?.refreshAt ?: return Duration.ZERO
        val delay = Duration.between(clock.instant(), refreshAt)
        return if (delay.isNegative) Duration.ZERO else delay
    }

    private fun store(token: TemporalAccessToken) {
        require(token.value.isNotBlank()) { "temporal access token must not be blank" }
        require(!token.expiresIn.isNegative && !token.expiresIn.isZero) { "temporal access token lifetime must be positive" }
        val now = clock.instant()
        val refreshAfter = Duration.ofMillis((token.expiresIn.toMillis() * refreshRatio).toLong())
        current = CachedToken(token.value, now.plus(token.expiresIn), now.plus(refreshAfter))
    }

    private fun schedule(executor: ScheduledExecutorService, delay: Duration) {
        if (executor.isShutdown) {
            return
        }
        try {
            executor.schedule({ schedule(executor, refresh()) }, delay.toMillis(), TimeUnit.MILLISECONDS)
        } catch (_: RejectedExecutionException) {
            return
        }
    }

    private data class CachedToken(
        val value: String,
        val expiresAt: Instant,
        val refreshAt: Instant,
    )

    companion object {
        const val DEFAULT_REFRESH_RATIO = 0.75
        val DEFAULT_MIN_RETRY_BACKOFF: Duration = Duration.ofSeconds(1)
        val DEFAULT_MAX_RETRY_BACKOFF: Duration = Duration.ofSeconds(30)
        private const val BEARER_PREFIX = "Bearer "
        private const val REFRESH_THREAD_NAME = "temporal-token-refresh"
    }
}
