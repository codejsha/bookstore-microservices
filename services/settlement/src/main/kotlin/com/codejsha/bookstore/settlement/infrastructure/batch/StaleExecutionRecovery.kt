package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import org.slf4j.LoggerFactory
import org.springframework.batch.core.BatchStatus
import org.springframework.batch.core.ExitStatus
import org.springframework.batch.core.job.JobExecution
import org.springframework.batch.core.repository.JobRepository
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Component
import java.time.Duration
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.ZoneOffset
import javax.sql.DataSource

@Component
class StaleExecutionRecovery(
    private val jobRepository: JobRepository,
    private val properties: SettlementBatchProperties,
    dataSource: DataSource,
) {
    private val log = LoggerFactory.getLogger(javaClass)
    private val jdbc = JdbcTemplate(dataSource)

    fun liveRunningDates(jobName: String): Set<LocalDate> =
        jobRepository.findRunningJobExecutions(jobName)
            .filterNot { isStale(it) }
            .mapNotNull { it.jobParameters.getLocalDate("targetDate") }
            .toSet()

    fun hasLiveExecution(jobName: String, targetDate: LocalDate): Boolean =
        targetDate in liveRunningDates(jobName)

    fun abandonStaleFor(jobName: String, targetDate: LocalDate): Int {
        val stale = jobRepository.findRunningJobExecutions(jobName)
            .filter { it.jobParameters.getLocalDate("targetDate") == targetDate }
            .filter { isStale(it) }
        stale.forEach { abandon(it) }
        return stale.size
    }

    fun isStale(execution: JobExecution): Boolean {
        val lastActivity = (
            execution.stepExecutions.mapNotNull { it.lastUpdated } +
                listOfNotNull(execution.lastUpdated, execution.startTime)
            ).maxOrNull() ?: return true
        return Duration.between(lastActivity, LocalDateTime.now()) > properties.staleExecutionTimeout
    }

    private fun abandon(execution: JobExecution) {
        val now = LocalDateTime.now()
        val description = "abandoned: no progress for ${properties.staleExecutionTimeout} " +
            "(batch pod presumed killed)"
        execution.stepExecutions
            .filter { it.status.isRunning }
            .forEach { step ->
                step.status = BatchStatus.FAILED
                step.exitStatus = ExitStatus.FAILED.addExitDescription(description)
                step.setEndTime(now)
                jobRepository.update(step)
            }
        execution.status = BatchStatus.FAILED
        execution.exitStatus = ExitStatus.FAILED.addExitDescription(description)
        execution.setEndTime(now)
        jobRepository.update(execution)

        val nowUtc = LocalDateTime.now(ZoneOffset.UTC)
        jdbc.update(
            """
            UPDATE settlement_job_run
               SET status = 'FAILED', finished_at = ?, updated_at = ?,
                   message = COALESCE(message, ?)
             WHERE job_execution_id = ? AND status = 'STARTED'
            """.trimIndent(),
            nowUtc, nowUtc, description, execution.id,
        )
        log.warn(
            "Abandoned stale settlement execution {} (targetDate={}): {}",
            execution.id, execution.jobParameters.getLocalDate("targetDate"), description,
        )
    }
}
