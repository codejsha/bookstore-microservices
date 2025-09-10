package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.order.application.port.support.DistributedLock
import org.redisson.api.RedissonClient
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Component
import java.time.Duration
import java.util.concurrent.ExecutionException
import java.util.concurrent.ThreadLocalRandom
import java.util.concurrent.TimeUnit

@Component
class RedissonDistributedLock(
    private val redissonClient: RedissonClient,
) : DistributedLock {

    override fun tryLock(key: String, ttl: Duration?): String? {
        val threadId = newThreadId()
        val lock = redissonClient.getLock(lockKey(key))
        val acquired = if (ttl == null) {
            lock.tryLockAsync(threadId)
        } else {
            lock.tryLockAsync(NO_WAIT_MILLIS, ttl.toMillis(), TimeUnit.MILLISECONDS, threadId)
        }.toCompletableFuture().get()
        return if (acquired == true) threadId.toString() else null
    }

    override fun unlock(key: String, token: String) {
        val threadId = token.toLongOrNull()
        if (threadId == null) {
            log.warn("Ignoring unlock for key {} with malformed token: {}", key, token)
            return
        }
        try {
            redissonClient.getLock(lockKey(key))
                .unlockAsync(threadId)
                .toCompletableFuture()
                .get()
        } catch (e: ExecutionException) {
            val cause = e.cause
            if (cause is IllegalMonitorStateException) {
                log.warn("Lock for key {} is no longer held by token {}; it likely expired", key, token)
            } else {
                throw cause ?: e
            }
        }
    }

    override fun <T> withLock(key: String, ttl: Duration?, action: () -> T): T {
        val token = tryLock(key, ttl) ?: error("Failed to acquire lock for key: $key")
        try {
            return action()
        } finally {
            unlock(key, token)
        }
    }

    private fun lockKey(key: String): String = "$KEY_PREFIX$key"

    private fun newThreadId(): Long = ThreadLocalRandom.current().nextLong(1L, Long.MAX_VALUE)

    companion object {
        private const val KEY_PREFIX = "lock:"
        private const val NO_WAIT_MILLIS = 0L
        private val log = LoggerFactory.getLogger(RedissonDistributedLock::class.java)
    }
}
