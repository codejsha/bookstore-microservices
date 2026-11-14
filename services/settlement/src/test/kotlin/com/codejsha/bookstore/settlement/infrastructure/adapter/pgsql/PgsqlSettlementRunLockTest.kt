package com.codejsha.bookstore.settlement.infrastructure.adapter.pgsql

import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import com.codejsha.bookstore.settlement.config.properties.SettlementRunLockProperties
import com.codejsha.bookstore.settlement.support.SettlementPostgresIntegrationTest
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import java.time.Duration
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.util.concurrent.Callable
import java.util.concurrent.CountDownLatch
import java.util.concurrent.Executors
import java.util.concurrent.TimeUnit
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

class PgsqlSettlementRunLockTest : SettlementPostgresIntegrationTest() {

    private val ttl: Duration = Duration.ofMinutes(5)
    private val locks = mutableListOf<PgsqlSettlementRunLock>()

    private val date: LocalDate = LocalDate.of(2026, 10, 1)

    @AfterEach
    fun shutdownLocks() {
        locks.forEach { it.shutdown() }
        locks.clear()
    }

    @Test
    fun `tryAcquire_freeDate_recordsALeaseForTheOwner`() {
        val lock = newLock()
        val before = LocalDateTime.now(ZoneOffset.UTC)

        assertTrue(lock.tryAcquire(date, "run-1"))

        val lease = lease(date)
        assertEquals("run-1", lease.ownerToken)
        assertNull(lease.releasedAt)
        assertEquals(lease.acquiredAt, lease.heartbeatAt)
        assertEquals(lease.acquiredAt.plus(ttl), lease.expiresAt)
        assertFalse(lease.acquiredAt.isBefore(before.minusSeconds(1)))
        assertTrue(lock.isHeldBy(date, "run-1"))
        assertFalse(lock.isHeldBy(date, "run-2"))
        assertFalse(lock.isHeldBy(date.plusDays(1), "run-1"))
    }

    @Test
    fun `tryAcquire_dateHeldByAnotherRun_isRejectedAndKeepsTheHolder`() {
        val lock = newLock()
        assertTrue(lock.tryAcquire(date, "run-1"))
        val held = lease(date)

        assertFalse(lock.tryAcquire(date, "run-2"))
        assertFalse(lock.tryAcquire(date, "run-1"))

        assertEquals(held, lease(date))
        assertTrue(lock.isHeldBy(date, "run-1"))
    }

    @Test
    fun `tryAcquire_concurrentContenders_exactlyOneWins`() {
        val contenders = 8
        val contenderLocks = (1..contenders).map { newLock() }
        val pool = Executors.newFixedThreadPool(contenders)
        val start = CountDownLatch(1)
        try {
            val futures = contenderLocks.mapIndexed { index, lock ->
                pool.submit(
                    Callable {
                        start.await()
                        "run-$index" to lock.tryAcquire(date, "run-$index")
                    },
                )
            }
            start.countDown()
            val outcomes = futures.map { it.get(30, TimeUnit.SECONDS) }

            val winners = outcomes.filter { it.second }.map { it.first }
            assertEquals(1, winners.size)
            assertEquals(winners.single(), lease(date).ownerToken)
            assertEquals(1, leaseRowCount(date))
        } finally {
            pool.shutdownNow()
        }
    }

    @Test
    fun `release_owner_freesTheDateForTheNextRun`() {
        val lock = newLock()
        assertTrue(lock.tryAcquire(date, "run-1"))

        assertTrue(lock.release(date, "run-1"))

        assertFalse(lock.isHeldBy(date, "run-1"))
        assertNotNull(lease(date).releasedAt)
        assertFalse(lock.release(date, "run-1"))
        assertTrue(lock.tryAcquire(date, "run-2"))
        assertTrue(lock.isHeldBy(date, "run-2"))
        assertNull(lease(date).releasedAt)
        assertEquals(1, leaseRowCount(date))
    }

    @Test
    fun `release_nonOwner_leavesTheHolderInPlace`() {
        val lock = newLock()
        assertTrue(lock.tryAcquire(date, "run-1"))

        assertFalse(lock.release(date, "run-2"))
        assertFalse(lock.release(date.plusDays(1), "run-1"))

        assertTrue(lock.isHeldBy(date, "run-1"))
    }

    @Test
    fun `tryAcquire_softDeletedLease_isTakenOver`() {
        val lock = newLock()
        assertTrue(lock.tryAcquire(date, "run-1"))
        settlementJdbc.update(
            "UPDATE settlement_run_lock SET deleted_at = ? WHERE target_date = ?",
            LocalDateTime.now(ZoneOffset.UTC), date,
        )

        assertFalse(lock.isHeldBy(date, "run-1"))
        assertTrue(lock.tryAcquire(date, "run-2"))

        assertTrue(lock.isHeldBy(date, "run-2"))
        assertNull(
            settlementJdbc.queryForObject(
                "SELECT deleted_at FROM settlement_run_lock WHERE target_date = ?", LocalDateTime::class.java, date,
            ),
        )
    }

    @Test
    fun `heartbeat_heldLease_extendsTheExpiry`() {
        val lock = newLock()
        assertTrue(lock.tryAcquire(date, "run-1"))
        val acquired = lease(date)

        val renewed = awaitLease(date) { it.heartbeatAt.isAfter(acquired.heartbeatAt) }

        assertTrue(renewed.expiresAt.isAfter(acquired.expiresAt))
        assertEquals(renewed.heartbeatAt.plus(ttl), renewed.expiresAt)
        assertEquals("run-1", renewed.ownerToken)
        assertTrue(lock.isHeldBy(date, "run-1"))
    }

    @Test
    fun `heartbeat_leaseLostToAnotherRun_stopsRenewing`() {
        val lock = newLock()
        assertTrue(lock.tryAcquire(date, "run-1"))
        settlementJdbc.update("UPDATE settlement_run_lock SET owner_token = 'run-2' WHERE target_date = ?", date)
        Thread.sleep(HEARTBEAT_SETTLE_MILLIS)
        val frozen = LocalDateTime.of(2026, 10, 1, 0, 0)
        settlementJdbc.update(
            "UPDATE settlement_run_lock SET owner_token = 'run-1', heartbeat_at = ? WHERE target_date = ?",
            frozen, date,
        )

        Thread.sleep(HEARTBEAT_SETTLE_MILLIS)

        assertEquals(frozen, lease(date).heartbeatAt)
    }

    @Test
    fun `release_heldLease_stopsTheHeartbeat`() {
        val lock = newLock()
        assertTrue(lock.tryAcquire(date, "run-1"))
        assertTrue(lock.release(date, "run-1"))
        val frozen = LocalDateTime.of(2026, 10, 1, 0, 0)
        settlementJdbc.update(
            "UPDATE settlement_run_lock SET released_at = NULL, heartbeat_at = ? WHERE target_date = ?",
            frozen, date,
        )

        Thread.sleep(HEARTBEAT_SETTLE_MILLIS)

        assertEquals(frozen, lease(date).heartbeatAt)
    }

    @Test
    fun `shutdown_heldLease_stopsTheHeartbeatAndLetsTheLeaseLapse`() {
        val lock = newLock()
        assertTrue(lock.tryAcquire(date, "run-1"))
        val acquired = lease(date)

        lock.shutdown()
        Thread.sleep(HEARTBEAT_SETTLE_MILLIS)

        assertEquals(acquired, lease(date))
    }

    private fun newLock(): PgsqlSettlementRunLock =
        PgsqlSettlementRunLock(
            dataSource,
            SettlementBatchProperties(
                lock = SettlementRunLockProperties(ttl = ttl, heartbeatInterval = Duration.ofSeconds(1)),
            ),
        ).also { locks += it }

    private fun awaitLease(targetDate: LocalDate, condition: (Lease) -> Boolean): Lease {
        val deadline = System.nanoTime() + TimeUnit.SECONDS.toNanos(10)
        while (System.nanoTime() < deadline) {
            val current = lease(targetDate)
            if (condition(current)) return current
            Thread.sleep(100)
        }
        throw AssertionError("the lease for $targetDate never reached the expected state: ${lease(targetDate)}")
    }

    private fun leaseRowCount(targetDate: LocalDate): Int =
        settlementJdbc.queryForObject(
            "SELECT COUNT(*) FROM settlement_run_lock WHERE target_date = ?", Int::class.java, targetDate,
        )!!

    private fun lease(targetDate: LocalDate): Lease =
        settlementJdbc.queryForObject(
            """
            SELECT owner_token, acquired_at, heartbeat_at, expires_at, released_at
              FROM settlement_run_lock WHERE target_date = ?
            """.trimIndent(),
            { rs, _ ->
                Lease(
                    ownerToken = rs.getString("owner_token"),
                    acquiredAt = rs.getObject("acquired_at", LocalDateTime::class.java),
                    heartbeatAt = rs.getObject("heartbeat_at", LocalDateTime::class.java),
                    expiresAt = rs.getObject("expires_at", LocalDateTime::class.java),
                    releasedAt = rs.getObject("released_at", LocalDateTime::class.java),
                )
            },
            targetDate,
        )!!

    private data class Lease(
        val ownerToken: String,
        val acquiredAt: LocalDateTime,
        val heartbeatAt: LocalDateTime,
        val expiresAt: LocalDateTime,
        val releasedAt: LocalDateTime?,
    )

    private companion object {
        const val HEARTBEAT_SETTLE_MILLIS = 2_500L
    }
}
