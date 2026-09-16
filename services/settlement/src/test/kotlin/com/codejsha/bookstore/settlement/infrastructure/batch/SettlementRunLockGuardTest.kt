package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.application.SettlementRunLockLostException
import com.codejsha.bookstore.settlement.support.FakeSettlementRunLock
import java.time.LocalDate
import kotlin.test.Test
import kotlin.test.assertFailsWith

class SettlementRunLockGuardTest {

    private val targetDate: LocalDate = LocalDate.of(2026, 7, 10)

    @Test
    fun `verifyOwnership_runStillHoldsTheLock_passes`() {
        val guard = SettlementRunLockGuard(FakeSettlementRunLock(heldBy = mapOf(targetDate to "run-1")))

        guard.verifyOwnership(targetDate, "run-1")
    }

    @Test
    fun `verifyOwnership_abandonedRunWithoutTheLock_throwsLockLost`() {
        val runLock = FakeSettlementRunLock(heldBy = mapOf(targetDate to "run-1"))
        runLock.release(targetDate, "run-1")
        val guard = SettlementRunLockGuard(runLock)

        assertFailsWith<SettlementRunLockLostException> { guard.verifyOwnership(targetDate, "run-1") }
    }

    @Test
    fun `verifyOwnership_lockTakenOverByAnotherRun_throwsLockLost`() {
        val guard = SettlementRunLockGuard(FakeSettlementRunLock(heldBy = mapOf(targetDate to "run-2")))

        assertFailsWith<SettlementRunLockLostException> { guard.verifyOwnership(targetDate, "run-1") }
    }

    @Test
    fun `verifyOwnership_executionWithoutRunUid_throwsLockLost`() {
        val guard = SettlementRunLockGuard(FakeSettlementRunLock(heldBy = mapOf(targetDate to "run-1")))

        assertFailsWith<SettlementRunLockLostException> { guard.verifyOwnership(targetDate, null) }
        assertFailsWith<SettlementRunLockLostException> { guard.verifyOwnership(targetDate, "  ") }
    }
}
