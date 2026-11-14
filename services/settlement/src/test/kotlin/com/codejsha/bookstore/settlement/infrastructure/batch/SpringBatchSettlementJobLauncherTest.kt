package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.application.SettlementRunConflictException
import com.codejsha.bookstore.settlement.application.SettlementRunValidationException
import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import com.codejsha.bookstore.settlement.support.FakeSettlementRunLock
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.mockito.ArgumentMatchers.any
import org.mockito.ArgumentMatchers.anyString
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.never
import org.mockito.Mockito.verify
import org.springframework.batch.core.BatchStatus
import org.springframework.batch.core.job.Job
import org.springframework.batch.core.job.JobExecution
import org.springframework.batch.core.job.JobInstance
import org.springframework.batch.core.job.parameters.InvalidJobParametersException
import org.springframework.batch.core.job.parameters.JobParameters
import org.springframework.batch.core.job.parameters.JobParametersBuilder
import org.springframework.batch.core.job.parameters.JobParametersValidator
import org.springframework.batch.core.repository.JobRepository
import org.springframework.batch.infrastructure.item.ExecutionContext
import java.time.LocalDate
import java.util.UUID
import javax.sql.DataSource
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

class SpringBatchSettlementJobLauncherTest {

    private val date: LocalDate = LocalDate.of(2026, 7, 1)
    private val jobRepository: JobRepository = mock(JobRepository::class.java)
    private val job: Job = mock(Job::class.java).also { given(it.name).willReturn(JOB_NAME) }
    private val launcher = SpringBatchSettlementJobLauncher(
        jobRepository,
        job,
        StaleExecutionRecovery(
            jobRepository, SettlementBatchProperties(), FakeSettlementRunLock(), mock(DataSource::class.java),
        ),
    )

    @AfterEach
    fun shutdownLauncher() {
        launcher.shutdown()
    }

    @Test
    fun `launch_lastExecutionEndedUnknown_throwsConflictWithoutCreatingAnExecution`() {
        val instance = JobInstance(1L, JOB_NAME)
        val unknown = JobExecution(1L, instance, identifying()).apply { status = BatchStatus.UNKNOWN }
        given(jobRepository.getJobInstance(anyString(), any(JobParameters::class.java))).willReturn(instance)
        given(jobRepository.getJobExecutions(instance)).willReturn(listOf(unknown))
        given(jobRepository.getLastJobExecution(instance)).willReturn(unknown)

        val error = assertFailsWith<SettlementRunConflictException> {
            launcher.launch(date, UUID.randomUUID(), rerun = true)
        }

        assertTrue(error.message!!.startsWith("Settlement run for $date cannot be restarted"))
        verifyNoExecutionCreated()
    }

    @Test
    fun `launch_jobRejectsTheParameters_throwsValidationWithoutCreatingAnExecution`() {
        given(job.jobParametersValidator).willReturn(
            JobParametersValidator { throw InvalidJobParametersException("targetDate is out of range") },
        )
        given(jobRepository.createJobInstance(anyString(), any(JobParameters::class.java)))
            .willReturn(JobInstance(1L, JOB_NAME))

        val error = assertFailsWith<SettlementRunValidationException> {
            launcher.launch(date, UUID.randomUUID(), rerun = false)
        }

        assertTrue(error.message!!.contains("targetDate is out of range"))
        verifyNoExecutionCreated()
    }

    private fun verifyNoExecutionCreated() {
        verify(jobRepository, never()).createJobExecution(
            any(JobInstance::class.java), any(JobParameters::class.java), any(ExecutionContext::class.java),
        )
    }

    private fun identifying(): JobParameters =
        JobParametersBuilder().addLocalDate("targetDate", date, true).toJobParameters()

    private companion object {
        const val JOB_NAME = "dailySettlementJob"
    }
}
