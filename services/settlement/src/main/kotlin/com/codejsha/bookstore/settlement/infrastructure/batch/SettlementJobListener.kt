package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.infrastructure.support.utils.uuidToBytes
import org.springframework.batch.core.BatchStatus
import org.springframework.batch.core.job.JobExecution
import org.springframework.batch.core.listener.JobExecutionListener
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Component
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.util.UUID
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid
import javax.sql.DataSource

@Component
class SettlementJobListener(
    dataSource: DataSource,
) : JobExecutionListener {

    private val jdbc = JdbcTemplate(dataSource)

    override fun afterJob(jobExecution: JobExecution) {
        val now = LocalDateTime.now(ZoneOffset.UTC)
        val completed = jobExecution.status == BatchStatus.COMPLETED
        val status = if (completed) "COMPLETED" else "FAILED"
        val failureMessage = jobExecution.allFailureExceptions
            .firstOrNull()
            ?.let { "${it.javaClass.simpleName}: ${it.message}" }
        val message = if (completed) "Settlement completed" else failureMessage

        val updated = jdbc.update(
            """
            UPDATE settlement_job_run
               SET status = ?, finished_at = ?, updated_at = ?,
                   message = COALESCE(message, ?)
             WHERE job_execution_id = ?
            """.trimIndent(),
            status, now, now, message, jobExecution.id,
        )

        if (updated == 0) {
            val targetDate = jobExecution.jobParameters.getLocalDate("targetDate")
            val uid = jobExecution.jobParameters.getString("runUid")?.let { UUID.fromString(it) } ?: Uuid.generateV7().toJavaUuid()
            jdbc.update(
                """
                INSERT INTO settlement_job_run
                    (uid, target_date, job_execution_id, status, started_at, finished_at, message, created_at)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                """.trimIndent(),
                uuidToBytes(uid), targetDate, jobExecution.id, status, now, now, message, now,
            )
        }
    }
}
