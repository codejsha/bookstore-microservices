package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.application.SettlementRunConflictException
import com.codejsha.bookstore.settlement.application.port.SettlementJobLauncher
import com.codejsha.bookstore.settlement.application.port.SettlementRunLock
import com.codejsha.bookstore.settlement.support.SettlementPostgresIntegrationTest
import org.junit.jupiter.api.Test
import org.springframework.batch.core.BatchStatus
import org.springframework.batch.core.job.parameters.JobParametersBuilder
import org.springframework.batch.core.repository.JobRepository
import org.springframework.batch.infrastructure.item.ExecutionContext
import org.springframework.beans.factory.annotation.Autowired
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.util.UUID
import java.util.concurrent.TimeUnit
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

class SpringBatchSettlementJobLauncherIntegrationTest : SettlementPostgresIntegrationTest() {

    @Autowired
    private lateinit var launcher: SettlementJobLauncher

    @Autowired
    private lateinit var runLock: SettlementRunLock

    @Autowired
    private lateinit var jobRepository: JobRepository

    @Test
    fun `launch_unsettledDate_runsTheJobInTheBackgroundUnderTheRunUid`() {
        val date = LocalDate.of(2026, 8, 1)
        insertPayment("pay_launch_1", 4_000, date)
        val runUid = acquire(date)
        val before = LocalDateTime.now(ZoneOffset.UTC).minusSeconds(5)

        val launched = launcher.launch(date, runUid, rerun = false)

        assertTrue(launched.status in setOf("STARTING", "STARTED", "COMPLETED"), launched.status)
        assertTrue(launched.startedAt.isAfter(before))
        val runs = awaitFinishedRuns(date, expected = 1)
        assertEquals(listOf(JobRun(runUid, "COMPLETED")), runs.map { JobRun(it.uid, it.status) })
        assertEquals(4_000L, grossAmount(date))
        assertNull(parametersOf(runs.single()).getParameter("rerun.id"))
        assertNull(parametersOf(runs.single()).getParameter("recovery.id"))
        assertFalse(runLock.isHeldBy(date, runUid.toString()))
    }

    @Test
    fun `launch_settledDateWithoutRerun_throwsConflict`() {
        val date = LocalDate.of(2026, 8, 2)
        insertPayment("pay_launch_2", 4_000, date)
        launcher.launch(date, acquire(date), rerun = false)
        awaitFinishedRuns(date, expected = 1)

        val error = assertFailsWith<SettlementRunConflictException> {
            launcher.launch(date, acquire(date), rerun = false)
        }

        assertEquals("$date is already settled; set rerun=true to re-run it", error.message)
        assertEquals(1, finishedRuns(date).size)
    }

    @Test
    fun `launch_settledDateWithRerun_settlesTheDateAgainAsANewInstance`() {
        val date = LocalDate.of(2026, 8, 3)
        insertPayment("pay_launch_3", 4_000, date)
        launcher.launch(date, acquire(date), rerun = false)
        val first = awaitFinishedRuns(date, expected = 1).single()
        insertPayment("pay_launch_3b", 1_500, date)
        val rerunUid = acquire(date)

        launcher.launch(date, rerunUid, rerun = true)

        val runs = awaitFinishedRuns(date, expected = 2)
        val second = runs.single { it.uid == rerunUid }
        assertEquals("COMPLETED", second.status)
        assertNotNull(parametersOf(second).getParameter("rerun.id"))
        assertTrue(instanceIdOf(first) != instanceIdOf(second))
        assertEquals(5_500L, grossAmount(date))
    }

    @Test
    fun `launch_priorRunFailed_startsARecoveryInstanceThatSettlesTheDate`() {
        val date = LocalDate.of(2026, 8, 4)
        insertPayment("pay_launch_4", 4_000, date)
        val lostUid = UUID.randomUUID()
        launcher.launch(date, lostUid, rerun = false)
        val failed = awaitFinishedRuns(date, expected = 1).single()
        assertEquals("FAILED", failed.status)
        val recoveryUid = acquire(date)

        launcher.launch(date, recoveryUid, rerun = false)

        val recovered = awaitFinishedRuns(date, expected = 2).single { it.uid == recoveryUid }
        assertEquals("COMPLETED", recovered.status)
        assertNotNull(parametersOf(recovered).getParameter("recovery.id"))
        assertNull(parametersOf(recovered).getParameter("rerun.id"))
        assertTrue(instanceIdOf(failed) != instanceIdOf(recovered))
        assertEquals(4_000L, grossAmount(date))
    }

    @Test
    fun `launch_dateAlreadyRunning_throwsConflict`() {
        val date = LocalDate.of(2026, 8, 5)
        val runningId = createStartedExecution(date)
        try {
            assertTrue(date in launcher.runningTargetDates())

            val error = assertFailsWith<SettlementRunConflictException> {
                launcher.launch(date, acquire(date), rerun = false)
            }

            assertEquals("A settlement run for $date is already running", error.message)
            assertEquals(BatchStatus.STARTED, jobRepository.getJobExecution(runningId)!!.status)
            assertTrue(finishedRuns(date).isEmpty())
        } finally {
            finishExecution(runningId)
        }
        assertFalse(date in launcher.runningTargetDates())
    }

    private data class JobRun(val uid: UUID, val status: String)

    private data class FinishedRun(val uid: UUID, val status: String, val jobExecutionId: Long)

    private fun acquire(date: LocalDate): UUID {
        val runUid = UUID.randomUUID()
        check(runLock.tryAcquire(date, runUid.toString())) { "the settlement run lock for $date must be free" }
        return runUid
    }

    private fun finishedRuns(date: LocalDate): List<FinishedRun> =
        settlementJdbc.query(
            """
            SELECT uid, status, job_execution_id FROM settlement_job_run
             WHERE target_date = ? AND finished_at IS NOT NULL ORDER BY id
            """.trimIndent(),
            { rs, _ ->
                FinishedRun(
                    rs.getObject("uid", UUID::class.java), rs.getString("status"), rs.getLong("job_execution_id"),
                )
            },
            date,
        )

    private fun awaitFinishedRuns(date: LocalDate, expected: Int): List<FinishedRun> {
        val deadline = System.nanoTime() + TimeUnit.SECONDS.toNanos(60)
        while (System.nanoTime() < deadline) {
            val runs = finishedRuns(date)
            val settled = runs.size >= expected &&
                runs.none { jobRepository.getJobExecution(it.jobExecutionId)!!.status.isRunning } &&
                date !in launcher.runningTargetDates()
            if (settled) return runs
            Thread.sleep(100)
        }
        throw AssertionError("expected $expected finished settlement runs for $date, found ${finishedRuns(date)}")
    }

    private fun parametersOf(run: FinishedRun) = jobRepository.getJobExecution(run.jobExecutionId)!!.jobParameters

    private fun instanceIdOf(run: FinishedRun): Long = jobRepository.getJobExecution(run.jobExecutionId)!!.jobInstanceId

    private fun grossAmount(date: LocalDate): Long =
        settlementJdbc.queryForObject(
            "SELECT gross_amount FROM daily_settlement WHERE settlement_date = ? AND currency = 'KRW'",
            Long::class.java, date,
        )!!

    private fun insertPayment(paymentId: String, amount: Long, date: LocalDate) {
        paymentJdbc.update(
            """
            INSERT INTO payments
                (uid, payment_id, amount, amount_captured, currency, status, payment_method, captured_at, created_at)
            VALUES (?, ?, ?, ?, 'KRW', 'succeeded', 'card', ?, ?)
            """.trimIndent(),
            UUID.randomUUID(), paymentId, amount, amount, date.atTime(3, 0), LocalDateTime.now(ZoneOffset.UTC),
        )
    }

    private fun createStartedExecution(date: LocalDate): Long {
        val parameters = JobParametersBuilder().addLocalDate("targetDate", date, true).toJobParameters()
        val instance = jobRepository.createJobInstance(JOB_NAME, parameters)
        val execution = jobRepository.createJobExecution(instance, parameters, ExecutionContext())
        val now = LocalDateTime.now()
        settlementJdbc.update(
            "UPDATE BATCH_JOB_EXECUTION SET STATUS = 'STARTED', START_TIME = ?, LAST_UPDATED = ? WHERE JOB_EXECUTION_ID = ?",
            now, now, execution.id,
        )
        return execution.id
    }

    private fun finishExecution(executionId: Long) {
        val execution = jobRepository.getJobExecution(executionId)!!
        execution.status = BatchStatus.FAILED
        execution.setEndTime(LocalDateTime.now())
        jobRepository.update(execution)
    }

    private companion object {
        const val JOB_NAME = "dailySettlementJob"
    }
}
