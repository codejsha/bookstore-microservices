package com.codejsha.bookstore.payment.infrastructure.support

import org.junit.jupiter.api.Test
import org.mockito.ArgumentMatchers.anyLong
import org.mockito.ArgumentMatchers.eq
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.mockingDetails
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoInteractions
import org.mockito.Mockito.verifyNoMoreInteractions
import org.redisson.api.RFuture
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
import kotlin.test.assertSame
import kotlin.test.assertTrue

class RedissonDistributedLockTest {

    private val redisson: RedissonClient = mock(RedissonClient::class.java)
    private val rLock: RLock = mock(RLock::class.java)
    private val lock = RedissonDistributedLock(redisson)

    private fun acquired(value: Boolean): RFuture<Boolean> = CompletableFutureWrapper(CompletableFuture.completedFuture(value))

    private fun released(): RFuture<Void> = CompletableFutureWrapper(CompletableFuture.completedFuture<Void>(null))

    private fun releaseFailed(cause: Throwable): RFuture<Void> = CompletableFutureWrapper(CompletableFuture.failedFuture<Void>(cause))

    @Test
    fun `tryLock_noTtlAndLockFree_acquiresPrefixedKeyAndYieldsThreadIdToken`() {
        given(redisson.getLock("lock:order-1")).willReturn(rLock)
        given(rLock.tryLockAsync(anyLong())).willReturn(acquired(true))

        val token = assertNotNull(lock.tryLock("order-1"))

        val threadId = token.toLong()
        assertTrue(threadId >= 1L, "thread id must be positive: $threadId")
        verify(redisson).getLock("lock:order-1")
        verify(rLock).tryLockAsync(threadId)
        verifyNoMoreInteractions(redisson, rLock)
    }

    @Test
    fun `tryLock_ttlGiven_requestsZeroWaitLeaseInMilliseconds`() {
        given(redisson.getLock("lock:order-2")).willReturn(rLock)
        given(rLock.tryLockAsync(anyLong(), anyLong(), eq(TimeUnit.MILLISECONDS), anyLong()))
            .willReturn(acquired(true))

        val token = assertNotNull(lock.tryLock("order-2", Duration.ofSeconds(5)))

        verify(rLock).tryLockAsync(0L, 5_000L, TimeUnit.MILLISECONDS, token.toLong())
        verifyNoMoreInteractions(rLock)
    }

    @Test
    fun `tryLock_lockHeldByAnotherOwner_yieldsNull`() {
        given(redisson.getLock("lock:order-3")).willReturn(rLock)
        given(rLock.tryLockAsync(anyLong())).willReturn(acquired(false))

        assertNull(lock.tryLock("order-3"))
    }

    @Test
    fun `tryLock_ttlGivenAndLockHeldByAnotherOwner_yieldsNull`() {
        given(redisson.getLock("lock:order-4")).willReturn(rLock)
        given(rLock.tryLockAsync(anyLong(), anyLong(), eq(TimeUnit.MILLISECONDS), anyLong()))
            .willReturn(acquired(false))

        assertNull(lock.tryLock("order-4", Duration.ofMillis(250)))
    }

    @Test
    fun `unlock_tokenFromTryLock_releasesWithSameThreadId`() {
        given(redisson.getLock("lock:order-5")).willReturn(rLock)
        given(rLock.unlockAsync(42L)).willReturn(released())

        lock.unlock("order-5", "42")

        verify(redisson).getLock("lock:order-5")
        verify(rLock).unlockAsync(42L)
        verifyNoMoreInteractions(redisson, rLock)
    }

    @Test
    fun `unlock_malformedToken_ignoresWithoutTouchingRedisson`() {
        lock.unlock("order-6", "not-a-number")

        verifyNoInteractions(redisson)
    }

    @Test
    fun `unlock_leaseAlreadyExpired_swallowsIllegalMonitorState`() {
        given(redisson.getLock("lock:order-7")).willReturn(rLock)
        given(rLock.unlockAsync(7L)).willReturn(releaseFailed(IllegalMonitorStateException("not locked by current thread")))

        lock.unlock("order-7", "7")

        verify(rLock).unlockAsync(7L)
    }

    @Test
    fun `unlock_redisFailure_rethrowsUnderlyingCause`() {
        val failure = IllegalArgumentException("redis down")
        given(redisson.getLock("lock:order-8")).willReturn(rLock)
        given(rLock.unlockAsync(8L)).willReturn(releaseFailed(failure))

        val thrown = assertFailsWith<IllegalArgumentException> { lock.unlock("order-8", "8") }

        assertSame(failure, thrown)
    }

    @Test
    fun `withLock_lockFree_runsActionThenReleasesWithAcquiredToken`() {
        given(redisson.getLock("lock:order-9")).willReturn(rLock)
        given(rLock.tryLockAsync(anyLong())).willReturn(acquired(true))
        given(rLock.unlockAsync(anyLong())).willReturn(released())
        var heldDuringAction = false

        val result = lock.withLock("order-9") {
            heldDuringAction = mockingDetails(rLock).invocations.none { it.method.name == "unlockAsync" }
            "done"
        }

        assertEquals("done", result)
        assertTrue(heldDuringAction, "lock must still be held while the action runs")
        val threadId = acquiredThreadId()
        verify(rLock).tryLockAsync(threadId)
        verify(rLock).unlockAsync(threadId)
        verifyNoMoreInteractions(rLock)
    }

    @Test
    fun `withLock_ttlGiven_acquiresWithLeaseAndReleases`() {
        given(redisson.getLock("lock:order-10")).willReturn(rLock)
        given(rLock.tryLockAsync(anyLong(), anyLong(), eq(TimeUnit.MILLISECONDS), anyLong()))
            .willReturn(acquired(true))
        given(rLock.unlockAsync(anyLong())).willReturn(released())

        val result = lock.withLock("order-10", Duration.ofSeconds(30)) { 7 }

        assertEquals(7, result)
        val threadId = acquiredThreadId()
        verify(rLock).tryLockAsync(0L, 30_000L, TimeUnit.MILLISECONDS, threadId)
        verify(rLock).unlockAsync(threadId)
        verifyNoMoreInteractions(rLock)
    }

    @Test
    fun `withLock_actionThrows_releasesLockAndPropagatesFailure`() {
        given(redisson.getLock("lock:order-11")).willReturn(rLock)
        given(rLock.tryLockAsync(anyLong())).willReturn(acquired(true))
        given(rLock.unlockAsync(anyLong())).willReturn(released())
        val failure = IllegalStateException("payment failed")

        val thrown = assertFailsWith<IllegalStateException> { lock.withLock<Unit>("order-11") { throw failure } }

        assertSame(failure, thrown)
        verify(rLock).unlockAsync(acquiredThreadId())
    }

    @Test
    fun `withLock_lockHeldByAnotherOwner_failsWithoutRunningActionOrReleasing`() {
        given(redisson.getLock("lock:order-12")).willReturn(rLock)
        given(rLock.tryLockAsync(anyLong())).willReturn(acquired(false))
        var ran = false

        val thrown = assertFailsWith<IllegalStateException> { lock.withLock("order-12") { ran = true } }

        assertEquals("Failed to acquire lock for key: order-12", thrown.message)
        assertFalse(ran)
        verify(rLock).tryLockAsync(acquiredThreadId())
        verifyNoMoreInteractions(rLock)
    }

    private fun acquiredThreadId(): Long {
        val invocation = mockingDetails(rLock).invocations.first { it.method.name == "tryLockAsync" }
        return invocation.arguments.last() as Long
    }
}
