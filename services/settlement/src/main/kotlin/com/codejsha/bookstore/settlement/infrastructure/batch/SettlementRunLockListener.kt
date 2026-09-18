package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.application.port.SettlementRunLock
import org.slf4j.LoggerFactory
import org.springframework.batch.core.job.JobExecution
import org.springframework.batch.core.listener.JobExecutionListener
import org.springframework.stereotype.Component

@Component
class SettlementRunLockListener(
    private val runLock: SettlementRunLock,
) : JobExecutionListener {

    private val log = LoggerFactory.getLogger(javaClass)

    override fun afterJob(jobExecution: JobExecution) {
        val targetDate = jobExecution.jobParameters.getLocalDate("targetDate") ?: return
        val ownerToken = jobExecution.jobParameters.getString("runUid") ?: return
        if (runLock.release(targetDate, ownerToken)) {
            log.info("Released the settlement run lock for {} held by {}", targetDate, ownerToken)
        }
    }
}
