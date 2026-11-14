package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import com.codejsha.bookstore.settlement.support.FakeSettlementRunLock
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.io.TempDir
import org.mockito.ArgumentMatchers.any
import org.mockito.ArgumentMatchers.anyString
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.springframework.batch.core.BatchStatus
import org.springframework.batch.core.job.Job
import org.springframework.batch.core.job.JobExecution
import org.springframework.batch.core.job.JobInstance
import org.springframework.batch.core.job.parameters.JobParameters
import org.springframework.batch.core.job.parameters.JobParametersBuilder
import org.springframework.batch.core.launch.JobExecutionAlreadyRunningException
import org.springframework.batch.core.launch.JobInstanceAlreadyCompleteException
import org.springframework.batch.core.launch.JobOperator
import org.springframework.batch.core.repository.JobRepository
import org.springframework.boot.DefaultApplicationArguments
import org.springframework.context.support.GenericApplicationContext
import java.lang.management.ManagementFactory
import java.nio.file.Files
import java.nio.file.Path
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.ZoneId
import java.util.Properties
import java.util.UUID
import java.util.concurrent.TimeUnit
import javax.sql.DataSource
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.system.exitProcess
import kotlin.test.assertNull

class SettlementJobRunnerTest {

    @TempDir
    lateinit var workDir: Path

    @Test
    fun `run_explicitTargetDate_launchesThatDateUnderItsRunLockAndExitsZero`() {
        val outcome = runInChildJvm(Scenario.JOB_COMPLETES, "--targetDate=2026-07-01")

        val runUid = assertNotNull(outcome["start.runUid"])
        assertEquals(0, outcome.exitCode)
        assertEquals("1", outcome["start.count"])
        assertEquals("2026-07-01", outcome["start.targetDate"])
        assertEquals("true", outcome["start.targetDate.identifying"])
        assertEquals(runUid, UUID.fromString(runUid).toString())
        assertEquals("false", outcome["start.runUid.identifying"])
        assertNull(outcome["start.rerunId"])
        assertNull(outcome["start.recoveryId"])
        assertEquals(runUid, outcome["start.lockHolder"])
        assertNull(outcome["exit.lockHolder"])
        assertEquals("false", outcome["exit.contextActive"])
    }

    @Test
    fun `run_noTargetDate_settlesYesterdayInTheConfiguredZone`() {
        val zone = ZoneId.of(TIMEZONE)
        val before = LocalDate.now(zone).minusDays(1)

        val outcome = runInChildJvm(Scenario.JOB_COMPLETES)

        val after = LocalDate.now(zone).minusDays(1)
        assertEquals(0, outcome.exitCode)
        assertEquals(true, outcome["start.targetDate"] in setOf(before.toString(), after.toString()))
    }

    @Test
    fun `run_rerunIdOption_addsAnIdentifyingRerunId`() {
        val outcome = runInChildJvm(Scenario.JOB_COMPLETES, "--targetDate=2026-07-01", "--rerun.id=manual-2")

        assertEquals(0, outcome.exitCode)
        assertEquals("manual-2", outcome["start.rerunId"])
        assertEquals("true", outcome["start.rerunId.identifying"])
    }

    @Test
    fun `run_malformedTargetDate_failsBeforeTouchingTheLockOrTheJob`() {
        val outcome = runInChildJvm(Scenario.JOB_COMPLETES, "--targetDate=2026-7-1")

        assertEquals(UNCAUGHT_EXIT_CODE, outcome.exitCode)
        assertEquals("java.time.format.DateTimeParseException", outcome["uncaught"])
        assertEquals("0", outcome["start.count"])
        assertNull(outcome["exit.lockHolder"])
    }

    @Test
    fun `run_dateLockedByAnotherRun_exitsOneWithoutStarting`() {
        val outcome = runInChildJvm(Scenario.DATE_LOCKED_BY_ANOTHER_RUN, "--targetDate=2026-07-01")

        assertEquals(1, outcome.exitCode)
        assertEquals("0", outcome["start.count"])
        assertEquals(OTHER_RUN, outcome["exit.lockHolder"])
        assertEquals("false", outcome["exit.contextActive"])
    }

    @Test
    fun `run_liveExecutionForTheDate_exitsOneAndReleasesTheLock`() {
        val outcome = runInChildJvm(Scenario.LIVE_EXECUTION_FOR_THE_DATE, "--targetDate=2026-07-01")

        assertEquals(1, outcome.exitCode)
        assertEquals("0", outcome["start.count"])
        assertNull(outcome["exit.lockHolder"])
    }

    @Test
    fun `run_liveExecutionForAnotherDate_stillLaunches`() {
        val outcome = runInChildJvm(Scenario.LIVE_EXECUTION_FOR_ANOTHER_DATE, "--targetDate=2026-07-01")

        assertEquals(0, outcome.exitCode)
        assertEquals("2026-07-01", outcome["start.targetDate"])
    }

    @Test
    fun `run_priorExecutionCompleted_exitsZeroWithoutStarting`() {
        val outcome = runInChildJvm(Scenario.PRIOR_EXECUTION_COMPLETED, "--targetDate=2026-07-01")

        assertEquals(0, outcome.exitCode)
        assertEquals("0", outcome["start.count"])
        assertNull(outcome["exit.lockHolder"])
    }

    @Test
    fun `run_priorExecutionFailed_launchesANewRecoveryInstance`() {
        val outcome = runInChildJvm(Scenario.PRIOR_EXECUTION_FAILED, "--targetDate=2026-07-01")

        assertEquals(0, outcome.exitCode)
        assertEquals("1", outcome["start.count"])
        assertEquals("2026-07-01", outcome["start.targetDate"])
        assertNotNull(outcome["start.recoveryId"])
        assertEquals("true", outcome["start.recoveryId.identifying"])
        assertNotNull(outcome["start.runUid"])
    }

    @Test
    fun `run_jobEndsFailed_exitsOneAndReleasesTheLock`() {
        val outcome = runInChildJvm(Scenario.JOB_FAILS, "--targetDate=2026-07-01")

        assertEquals(1, outcome.exitCode)
        assertEquals("1", outcome["start.count"])
        assertNull(outcome["exit.lockHolder"])
        assertEquals("false", outcome["exit.contextActive"])
    }

    @Test
    fun `run_instanceAlreadyComplete_exitsZero`() {
        val outcome = runInChildJvm(Scenario.INSTANCE_ALREADY_COMPLETE, "--targetDate=2026-07-01")

        assertEquals(0, outcome.exitCode)
        assertEquals("1", outcome["start.count"])
        assertNull(outcome["exit.lockHolder"])
    }

    @Test
    fun `run_startedConcurrentlyElsewhere_exitsOne`() {
        val outcome = runInChildJvm(Scenario.STARTED_CONCURRENTLY_ELSEWHERE, "--targetDate=2026-07-01")

        assertEquals(1, outcome.exitCode)
        assertEquals("1", outcome["start.count"])
        assertNull(outcome["exit.lockHolder"])
    }

    @Test
    fun `run_launchFailsUnexpectedly_exitsOneAndReleasesTheLock`() {
        val outcome = runInChildJvm(Scenario.LAUNCH_FAILS_UNEXPECTEDLY, "--targetDate=2026-07-01")

        assertEquals(1, outcome.exitCode)
        assertNull(outcome["uncaught"])
        assertNull(outcome["exit.lockHolder"])
        assertEquals("false", outcome["exit.contextActive"])
    }

    private class Outcome(val exitCode: Int, private val report: Properties) {
        operator fun get(key: String): String? = report.getProperty(key)
    }

    private fun runInChildJvm(scenario: Scenario, vararg runnerArgs: String): Outcome {
        val reportFile = workDir.resolve("report.properties")
        val outputFile = workDir.resolve("child-output.log")
        val agents = ManagementFactory.getRuntimeMXBean().inputArguments.filter { it.startsWith("-javaagent:") }
        val command = listOf(Path.of(System.getProperty("java.home"), "bin", "java").toString()) +
            agents +
            listOf(SettlementJobRunnerProcess::class.java.name, scenario.name, reportFile.toString()) +
            runnerArgs
        val process = ProcessBuilder(command)
            .redirectErrorStream(true)
            .redirectOutput(outputFile.toFile())
            .apply { environment()["CLASSPATH"] = System.getProperty("java.class.path") }
            .start()
        if (!process.waitFor(120, TimeUnit.SECONDS)) {
            process.destroyForcibly()
            throw AssertionError("the settlement runner process did not exit:\n${Files.readString(outputFile)}")
        }
        if (!Files.exists(reportFile)) {
            throw AssertionError(
                "the settlement runner process exited ${process.exitValue()} without a report:\n" +
                    Files.readString(outputFile),
            )
        }
        val report = Properties().apply { Files.newBufferedReader(reportFile).use { load(it) } }
        return Outcome(process.exitValue(), report)
    }
}

internal enum class Scenario {
    JOB_COMPLETES,
    JOB_FAILS,
    DATE_LOCKED_BY_ANOTHER_RUN,
    LIVE_EXECUTION_FOR_THE_DATE,
    LIVE_EXECUTION_FOR_ANOTHER_DATE,
    PRIOR_EXECUTION_COMPLETED,
    PRIOR_EXECUTION_FAILED,
    INSTANCE_ALREADY_COMPLETE,
    STARTED_CONCURRENTLY_ELSEWHERE,
    LAUNCH_FAILS_UNEXPECTEDLY,
}

internal const val UNCAUGHT_EXIT_CODE = 70
internal const val OTHER_RUN = "other-run"
internal const val TIMEZONE = "Asia/Seoul"
private const val JOB_NAME = "dailySettlementJob"
private val DATE: LocalDate = LocalDate.of(2026, 7, 1)

internal object SettlementJobRunnerProcess {

    @JvmStatic
    fun main(args: Array<String>) {
        val scenario = Scenario.valueOf(args[0])
        val reportFile = Path.of(args[1])
        val runnerArgs = args.drop(2).toTypedArray()

        val report = Properties()
        var startCount = 0
        val runLock = FakeSettlementRunLock(
            heldBy = if (scenario == Scenario.DATE_LOCKED_BY_ANOTHER_RUN) mapOf(DATE to OTHER_RUN) else emptyMap(),
        )
        val properties = SettlementBatchProperties(timezone = TIMEZONE)
        val applicationContext = GenericApplicationContext().apply { refresh() }
        val job = mock(Job::class.java)
        given(job.name).willReturn(JOB_NAME)
        val jobRepository = mock(JobRepository::class.java)
        val jobOperator = mock(JobOperator::class.java)

        when (scenario) {
            Scenario.LIVE_EXECUTION_FOR_THE_DATE ->
                given(jobRepository.findRunningJobExecutions(JOB_NAME)).willReturn(setOf(liveExecution(DATE)))
            Scenario.LIVE_EXECUTION_FOR_ANOTHER_DATE ->
                given(jobRepository.findRunningJobExecutions(JOB_NAME))
                    .willReturn(setOf(liveExecution(DATE.minusDays(1))))
            Scenario.PRIOR_EXECUTION_COMPLETED ->
                given(jobRepository.getLastJobExecution(anyString(), any(JobParameters::class.java)))
                    .willReturn(execution(BatchStatus.COMPLETED, identifying(DATE)))
            Scenario.PRIOR_EXECUTION_FAILED ->
                given(jobRepository.getLastJobExecution(anyString(), any(JobParameters::class.java)))
                    .willReturn(execution(BatchStatus.FAILED, identifying(DATE)))
            else -> Unit
        }

        given(jobOperator.start(any(Job::class.java), any(JobParameters::class.java))).willAnswer { invocation ->
            val parameters = invocation.getArgument<JobParameters>(1)
            val targetDate = parameters.getLocalDate("targetDate")!!
            startCount++
            report.setProperty("start.targetDate", targetDate.toString())
            report.setProperty("start.targetDate.identifying", parameters.getParameter("targetDate")!!.identifying().toString())
            parameters.getParameter("runUid")?.let {
                report.setProperty("start.runUid", it.value().toString())
                report.setProperty("start.runUid.identifying", it.identifying().toString())
            }
            parameters.getParameter("rerun.id")?.let {
                report.setProperty("start.rerunId", it.value().toString())
                report.setProperty("start.rerunId.identifying", it.identifying().toString())
            }
            parameters.getParameter("recovery.id")?.let {
                report.setProperty("start.recoveryId", it.value().toString())
                report.setProperty("start.recoveryId.identifying", it.identifying().toString())
            }
            runLock.holderOf(targetDate)?.let { report.setProperty("start.lockHolder", it) }
            when (scenario) {
                Scenario.JOB_FAILS -> execution(BatchStatus.FAILED, parameters)
                Scenario.INSTANCE_ALREADY_COMPLETE -> throw JobInstanceAlreadyCompleteException("already complete")
                Scenario.STARTED_CONCURRENTLY_ELSEWHERE -> throw JobExecutionAlreadyRunningException("already running")
                Scenario.LAUNCH_FAILS_UNEXPECTEDLY -> throw IllegalStateException("job repository unavailable")
                else -> execution(BatchStatus.COMPLETED, parameters)
            }
        }

        Runtime.getRuntime().addShutdownHook(
            Thread {
                report.setProperty("start.count", startCount.toString())
                report.setProperty("exit.contextActive", applicationContext.isActive.toString())
                runLock.holderOf(DATE)?.let { report.setProperty("exit.lockHolder", it) }
                Files.newBufferedWriter(reportFile).use { report.store(it, null) }
            },
        )

        val runner = SettlementJobRunner(
            jobOperator = jobOperator,
            dailySettlementJob = job,
            jobRepository = jobRepository,
            staleExecutionRecovery = StaleExecutionRecovery(
                jobRepository, properties, runLock, mock(DataSource::class.java),
            ),
            runLock = runLock,
            applicationContext = applicationContext,
            properties = properties,
        )
        try {
            runner.run(DefaultApplicationArguments(*runnerArgs))
        } catch (e: Throwable) {
            report.setProperty("uncaught", e.javaClass.name)
            exitProcess(UNCAUGHT_EXIT_CODE)
        }
    }

    private fun identifying(targetDate: LocalDate): JobParameters =
        JobParametersBuilder().addLocalDate("targetDate", targetDate, true).toJobParameters()

    private fun liveExecution(targetDate: LocalDate): JobExecution =
        execution(BatchStatus.STARTED, identifying(targetDate)).apply { setStartTime(LocalDateTime.now()) }

    private fun execution(status: BatchStatus, parameters: JobParameters): JobExecution =
        JobExecution(1L, JobInstance(1L, JOB_NAME), parameters).apply { this.status = status }
}
