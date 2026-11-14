package com.codejsha.bookstore.order.infrastructure.support

import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.mockito.ArgumentCaptor
import org.mockito.ArgumentMatchers.anyLong
import org.mockito.ArgumentMatchers.eq
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoInteractions
import org.mockito.Mockito.verifyNoMoreInteractions
import org.redisson.api.RLock
import org.redisson.api.RedissonClient
import org.redisson.misc.CompletableFutureWrapper
import java.time.Duration
import java.util.concurrent.CompletableFuture
import java.util.concurrent.TimeUnit
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

class RedissonDistributedLockTest {

    private val client = mock(RedissonClient::class.java)
    private val rLock = mock(RLock::class.java)
    private val lock = RedissonDistributedLock(client)

    @BeforeEach
    fun routeLockKey() {
        given(client.getLock("lock:order:1")).willReturn(rLock)
    }

    @Test
    fun `tryLock_noTtlAndLockFree_yieldsTokenUsedAsRedissonThreadId`() {
        given(rLock.tryLockAsync(anyLong())).willReturn(CompletableFutureWrapper(true))

        val token = assertNotNull(lock.tryLock("order:1"))

        assertTrue(token.toLong() > 0L)
        verify(rLock).tryLockAsync(token.toLong())
        verifyNoMoreInteractions(rLock)
    }

    @Test
    fun `tryLock_ttlGiven_requestsLeaseWithoutWaiting`() {
        given(rLock.tryLockAsync(eq(0L), eq(5000L), eq(TimeUnit.MILLISECONDS), anyLong()))
            .willReturn(CompletableFutureWrapper(true))

        val token = assertNotNull(lock.tryLock("order:1", Duration.ofSeconds(5)))

        verify(rLock).tryLockAsync(0L, 5000L, TimeUnit.MILLISECONDS, token.toLong())
        verifyNoMoreInteractions(rLock)
    }

    @Test
    fun `tryLock_lockHeldElsewhere_yieldsNull`() {
        given(rLock.tryLockAsync(anyLong())).willReturn(CompletableFutureWrapper(false))

        assertNull(lock.tryLock("order:1"))
    }

    @Test
    fun `tryLock_ttlGivenAndLockHeldElsewhere_yieldsNull`() {
        given(rLock.tryLockAsync(eq(0L), eq(1000L), eq(TimeUnit.MILLISECONDS), anyLong()))
            .willReturn(CompletableFutureWrapper(false))

        assertNull(lock.tryLock("order:1", Duration.ofSeconds(1)))
    }

    @Test
    fun `tryLock_calledTwice_issuesDistinctTokens`() {
        given(rLock.tryLockAsync(anyLong())).willReturn(CompletableFutureWrapper(true))

        assertFalse(lock.tryLock("order:1") == lock.tryLock("order:1"))
    }

    @Test
    fun `unlock_tokenFromTryLock_releasesWithThatThreadId`() {
        given(rLock.unlockAsync(42L)).willReturn(completed())

        lock.unlock("order:1", "42")

        verify(rLock).unlockAsync(42L)
        verifyNoMoreInteractions(rLock)
    }

    @Test
    fun `unlock_malformedToken_skipsRedisson`() {
        lock.unlock("order:1", "test-token")

        verifyNoInteractions(rLock)
    }

    @Test
    fun `unlock_lockNoLongerHeld_completesQuietly`() {
        given(rLock.unlockAsync(42L)).willReturn(failed(IllegalMonitorStateException("not locked by thread 42")))

        lock.unlock("order:1", "42")

        verify(rLock).unlockAsync(42L)
    }

    @Test
    fun `unlock_redisFails_rethrowsUnderlyingCause`() {
        given(rLock.unlockAsync(42L)).willReturn(failed(IllegalStateException("redis unreachable")))

        val ex = assertFailsWith<IllegalStateException> { lock.unlock("order:1", "42") }

        assertEquals("redis unreachable", ex.message)
    }

    @Test
    fun `withLock_lockAcquired_runsActionAndReleasesWithSameToken`() {
        given(rLock.tryLockAsync(eq(0L), eq(3000L), eq(TimeUnit.MILLISECONDS), anyLong()))
            .willReturn(CompletableFutureWrapper(true))
        given(rLock.unlockAsync(anyLong())).willReturn(completed())

        val result = lock.withLock("order:1", Duration.ofSeconds(3)) { "done" }

        assertEquals("done", result)
        val threadId = ArgumentCaptor.forClass(Long::class.javaObjectType)
        verify(rLock).tryLockAsync(eq(0L), eq(3000L), eq(TimeUnit.MILLISECONDS), threadId.capture() ?: 0L)
        verify(rLock).unlockAsync(threadId.value)
        verifyNoMoreInteractions(rLock)
    }

    @Test
    fun `withLock_actionThrows_releasesLockAndPropagates`() {
        given(rLock.tryLockAsync(anyLong())).willReturn(CompletableFutureWrapper(true))
        given(rLock.unlockAsync(anyLong())).willReturn(completed())

        val ex = assertFailsWith<IllegalArgumentException> {
            lock.withLock<Unit>("order:1") { throw IllegalArgumentException("action failed") }
        }

        assertEquals("action failed", ex.message)
        val threadId = ArgumentCaptor.forClass(Long::class.javaObjectType)
        verify(rLock).tryLockAsync(threadId.capture() ?: 0L)
        verify(rLock).unlockAsync(threadId.value)
    }

    @Test
    fun `withLock_lockUnavailable_failsWithoutRunningActionOrUnlocking`() {
        given(rLock.tryLockAsync(anyLong())).willReturn(CompletableFutureWrapper(false))
        var ran = false

        val ex = assertFailsWith<IllegalStateException> { lock.withLock("order:1") { ran = true } }

        assertEquals("Failed to acquire lock for key: order:1", ex.message)
        assertFalse(ran)
        val threadId = ArgumentCaptor.forClass(Long::class.javaObjectType)
        verify(rLock).tryLockAsync(threadId.capture() ?: 0L)
        verifyNoMoreInteractions(rLock)
    }

    private fun completed(): CompletableFutureWrapper<Void> =
        CompletableFutureWrapper(CompletableFuture.completedFuture<Void>(null))

    private fun failed(cause: Throwable): CompletableFutureWrapper<Void> =
        CompletableFutureWrapper(CompletableFuture.failedFuture<Void>(cause))
}
