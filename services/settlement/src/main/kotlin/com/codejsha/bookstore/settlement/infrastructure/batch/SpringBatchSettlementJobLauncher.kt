package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.application.SettlementRunConflictException
import com.codejsha.bookstore.settlement.application.SettlementRunValidationException
import com.codejsha.bookstore.settlement.application.port.LaunchedRun
import com.codejsha.bookstore.settlement.application.port.SettlementJobLauncher
import jakarta.annotation.PreDestroy
import org.springframework.batch.core.BatchStatus
import org.springframework.batch.core.configuration.support.MapJobRegistry
import org.springframework.batch.core.job.Job
import org.springframework.batch.core.job.parameters.InvalidJobParametersException
import org.springframework.batch.core.job.parameters.JobParametersBuilder
import org.springframework.batch.core.launch.JobExecutionAlreadyRunningException
import org.springframework.batch.core.launch.JobInstanceAlreadyCompleteException
import org.springframework.batch.core.launch.JobRestartException
import org.springframework.batch.core.launch.support.TaskExecutorJobOperator
import org.springframework.batch.core.repository.JobRepository
import org.springframework.scheduling.concurrent.ThreadPoolTaskExecutor
import org.springframework.stereotype.Component
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.util.UUID

@Component
class SpringBatchSettlementJobLauncher(
    private val jobRepository: JobRepository,
    private val dailySettlementJob: Job,
    private val staleExecutionRecovery: StaleExecutionRecovery,
) : SettlementJobLauncher {

    private val taskExecutor = ThreadPoolTaskExecutor().apply {
        corePoolSize = 1
        maxPoolSize = 1
        queueCapacity = 8
        setThreadNamePrefix("settlement-trigger-")
        setWaitForTasksToCompleteOnShutdown(true)
        setAwaitTerminationSeconds(60)
        initialize()
    }

    private val jobRegistry = MapJobRegistry().apply { register(dailySettlementJob) }

    private val operator = TaskExecutorJobOperator().apply {
        setJobRepository(jobRepository)
        setJobRegistry(jobRegistry)
        setTaskExecutor(this@SpringBatchSettlementJobLauncher.taskExecutor)
        afterPropertiesSet()
    }

    override fun runningTargetDates(): Set<LocalDate> =
        staleExecutionRecovery.liveRunningDates(dailySettlementJob.name)

    override fun launch(targetDate: LocalDate, runUid: UUID, rerun: Boolean): LaunchedRun {
        staleExecutionRecovery.abandonStaleFor(dailySettlementJob.name, targetDate)

        val builder = JobParametersBuilder()
            .addLocalDate("targetDate", targetDate, true)
            .addString("runUid", runUid.toString(), false)
        if (rerun) {
            builder.addString("rerun.id", UUID.randomUUID().toString(), true)
        } else {
            val last = jobRepository.getLastJobExecution(
                dailySettlementJob.name,
                JobParametersBuilder().addLocalDate("targetDate", targetDate, true).toJobParameters(),
            )
            if (last != null && last.status != BatchStatus.COMPLETED && !last.status.isRunning) {
                builder.addString("recovery.id", UUID.randomUUID().toString(), true)
            }
        }

        val execution = try {
            operator.start(dailySettlementJob, builder.toJobParameters())
        } catch (e: JobExecutionAlreadyRunningException) {
            throw SettlementRunConflictException("A settlement run for $targetDate is already running")
        } catch (e: JobInstanceAlreadyCompleteException) {
            throw SettlementRunConflictException("$targetDate is already settled; set rerun=true to re-run it")
        } catch (e: JobRestartException) {
            throw SettlementRunConflictException("Settlement run for $targetDate cannot be restarted: ${e.message}")
        } catch (e: InvalidJobParametersException) {
            throw SettlementRunValidationException("Invalid settlement run parameters: ${e.message}")
        }

        return LaunchedRun(
            status = execution.status.name,
            startedAt = execution.startTime ?: LocalDateTime.now(ZoneOffset.UTC),
        )
    }

    @PreDestroy
    fun shutdown() {
        taskExecutor.shutdown()
    }
}
