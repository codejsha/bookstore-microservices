package com.codejsha.bookstore.payment.infrastructure.support

import org.junit.jupiter.api.Test
import java.time.Clock
import java.time.Duration
import java.time.Instant
import java.time.ZoneId
import java.time.ZoneOffset
import java.util.ArrayDeque
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class TemporalTokenProviderTest {

    private val clock = MutableClock(Instant.parse("2026-09-15T00:00:00Z"))

    @Test
    fun initialize_300sLifetime_schedulesRefreshAt75PercentOfLifetime() {
        val provider = provider(ScriptedFetcher(token("first", 300)))

        provider.initialize()

        assertEquals(Duration.ofSeconds(225), provider.nextRefreshDelay())
    }

    @Test
    fun supply_repeatedCalls_reusesCachedTokenWithoutFetching() {
        val fetcher = ScriptedFetcher(token("first", 300))
        val provider = provider(fetcher)
        provider.initialize()

        repeat(5) { assertEquals("Bearer first", provider.supply()) }

        assertEquals(1, fetcher.calls)
    }

    @Test
    fun supply_beforeInitialize_throwsIllegalState() {
        val provider = provider(ScriptedFetcher())

        assertFailsWith<IllegalStateException> { provider.supply() }
    }

    @Test
    fun initialize_fetchFails_propagatesError() {
        val provider = provider(ScriptedFetcher(failure()))

        assertFailsWith<IllegalStateException> { provider.initialize() }
    }

    @Test
    fun refresh_fetchSucceeds_replacesTokenAndReschedulesAheadOfExpiry() {
        val provider = provider(ScriptedFetcher(token("first", 300), token("second", 120)))
        provider.initialize()
        clock.advance(Duration.ofSeconds(225))

        val delay = provider.refresh()

        assertEquals("Bearer second", provider.supply())
        assertEquals(Duration.ofSeconds(90), delay)
    }

    @Test
    fun refresh_fetchFails_keepsPreviousToken() {
        val provider = provider(ScriptedFetcher(token("first", 300), failure()))
        provider.initialize()
        clock.advance(Duration.ofSeconds(225))

        val delay = provider.refresh()

        assertEquals("Bearer first", provider.supply())
        assertEquals(Duration.ofSeconds(1), delay)
    }

    @Test
    fun refresh_consecutiveFailures_backsOffExponentiallyUpToMax() {
        val provider = provider(
            ScriptedFetcher(token("first", 300), failure(), failure(), failure(), failure(), failure()),
        )
        provider.initialize()

        val delays = (1..5).map { provider.refresh() }

        assertEquals(listOf(1L, 2L, 4L, 8L, 8L).map(Duration::ofSeconds), delays)
    }

    @Test
    fun refresh_successAfterFailures_resetsBackoff() {
        val provider = provider(
            ScriptedFetcher(token("first", 300), failure(), failure(), token("second", 300), failure()),
        )
        provider.initialize()
        provider.refresh()
        provider.refresh()
        provider.refresh()

        val delay = provider.refresh()

        assertEquals("Bearer second", provider.supply())
        assertEquals(Duration.ofSeconds(1), delay)
    }

    @Test
    fun nextRefreshDelay_refreshPointPassed_isZero() {
        val provider = provider(ScriptedFetcher(token("first", 300)))
        provider.initialize()
        clock.advance(Duration.ofSeconds(400))

        assertEquals(Duration.ZERO, provider.nextRefreshDelay())
    }

    private fun provider(fetcher: TemporalTokenFetcher) =
        TemporalTokenProvider(
            fetcher = fetcher,
            clock = clock,
            refreshRatio = 0.75,
            minRetryBackoff = Duration.ofSeconds(1),
            maxRetryBackoff = Duration.ofSeconds(8),
        )

    private fun token(value: String, seconds: Long): () -> TemporalAccessToken =
        { TemporalAccessToken(value, Duration.ofSeconds(seconds)) }

    private fun failure(): () -> TemporalAccessToken =
        { throw IllegalStateException("token endpoint unavailable") }

    private class ScriptedFetcher(vararg steps: () -> TemporalAccessToken) : TemporalTokenFetcher {
        private val steps = ArrayDeque(steps.toList())
        var calls = 0
            private set

        override fun fetch(): TemporalAccessToken {
            calls++
            return steps.removeFirst()()
        }
    }

    private class MutableClock(private var now: Instant) : Clock() {
        fun advance(duration: Duration) {
            now = now.plus(duration)
        }

        override fun getZone(): ZoneId = ZoneOffset.UTC

        override fun withZone(zone: ZoneId?): Clock = this

        override fun instant(): Instant = now
    }
}
