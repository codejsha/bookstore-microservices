package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.support.FakeSettlementRunLock
import org.springframework.batch.core.BatchStatus
import org.springframework.batch.core.job.JobExecution
import org.springframework.batch.core.job.JobInstance
import org.springframework.batch.core.job.parameters.JobParametersBuilder
import java.time.LocalDate
import kotlin.test.Test
import kotlin.test.assertFalse
import kotlin.test.assertTrue

class SettlementRunLockListenerTest {

    private val targetDate: LocalDate = LocalDate.of(2026, 7, 10)

    @Test
    fun `afterJob_completedRun_releasesTheDateLock`() {
        val runLock = FakeSettlementRunLock(heldBy = mapOf(targetDate to "run-1"))
        val listener = SettlementRunLockListener(runLock)

        listener.afterJob(jobExecution(BatchStatus.COMPLETED, "run-1"))

        assertFalse(runLock.isHeld(targetDate))
    }

    @Test
    fun `afterJob_failedRun_releasesTheDateLock`() {
        val runLock = FakeSettlementRunLock(heldBy = mapOf(targetDate to "run-1"))
        val listener = SettlementRunLockListener(runLock)

        listener.afterJob(jobExecution(BatchStatus.FAILED, "run-1"))

        assertFalse(runLock.isHeld(targetDate))
    }

    @Test
    fun `afterJob_runThatNeverOwnedTheLock_leavesTheHolderInPlace`() {
        val runLock = FakeSettlementRunLock(heldBy = mapOf(targetDate to "run-2"))
        val listener = SettlementRunLockListener(runLock)

        listener.afterJob(jobExecution(BatchStatus.FAILED, "run-1"))

        assertTrue(runLock.isHeldBy(targetDate, "run-2"))
    }

    @Test
    fun `afterJob_executionWithoutRunUid_leavesTheHolderInPlace`() {
        val runLock = FakeSettlementRunLock(heldBy = mapOf(targetDate to "run-1"))
        val listener = SettlementRunLockListener(runLock)

        listener.afterJob(jobExecution(BatchStatus.COMPLETED, null))

        assertTrue(runLock.isHeldBy(targetDate, "run-1"))
    }

    private fun jobExecution(status: BatchStatus, runUid: String?): JobExecution {
        val builder = JobParametersBuilder().addLocalDate("targetDate", targetDate, true)
        runUid?.let { builder.addString("runUid", it, false) }
        val execution = JobExecution(1L, JobInstance(1L, "dailySettlementJob"), builder.toJobParameters())
        execution.status = status
        return execution
    }
}
