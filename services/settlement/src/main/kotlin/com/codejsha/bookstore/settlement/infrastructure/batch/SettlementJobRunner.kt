package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import org.slf4j.LoggerFactory
import org.springframework.batch.core.BatchStatus
import org.springframework.batch.core.job.Job
import org.springframework.batch.core.job.parameters.JobParameters
import org.springframework.batch.core.job.parameters.JobParametersBuilder
import org.springframework.batch.core.launch.JobExecutionAlreadyRunningException
import org.springframework.batch.core.launch.JobInstanceAlreadyCompleteException
import org.springframework.batch.core.launch.JobOperator
import org.springframework.batch.core.repository.JobRepository
import org.springframework.boot.ApplicationArguments
import org.springframework.boot.ApplicationRunner
import org.springframework.boot.ExitCodeGenerator
import org.springframework.boot.SpringApplication
import org.springframework.context.ConfigurableApplicationContext
import org.springframework.context.annotation.Profile
import org.springframework.stereotype.Component
import java.time.LocalDate
import java.time.ZoneId
import java.util.UUID
import kotlin.system.exitProcess

@Component
@Profile("batch")
class SettlementJobRunner(
    private val jobOperator: JobOperator,
    private val dailySettlementJob: Job,
    private val jobRepository: JobRepository,
    private val staleExecutionRecovery: StaleExecutionRecovery,
    private val applicationContext: ConfigurableApplicationContext,
    private val properties: SettlementBatchProperties,
) : ApplicationRunner {

    private val log = LoggerFactory.getLogger(javaClass)

    override fun run(args: ApplicationArguments) {
        val targetDate = args.getOptionValues("targetDate")
            ?.firstOrNull()
            ?.let { LocalDate.parse(it) }
            ?: LocalDate.now(ZoneId.of(properties.timezone)).minusDays(1)

        val builder = JobParametersBuilder().addLocalDate("targetDate", targetDate, true)
        args.getOptionValues("rerun.id")?.firstOrNull()?.let { builder.addString("rerun.id", it, true) }

        val code = try {
            launch(targetDate, builder.toJobParameters())
        } catch (e: Exception) {
            log.error("Settlement run for {} failed to launch", targetDate, e)
            1
        }
        val exitCode = SpringApplication.exit(applicationContext, ExitCodeGenerator { code })
        exitProcess(exitCode)
    }

    private fun launch(targetDate: LocalDate, baseParameters: JobParameters): Int {
        staleExecutionRecovery.abandonStaleFor(dailySettlementJob.name, targetDate)

        if (staleExecutionRecovery.hasLiveExecution(dailySettlementJob.name, targetDate)) {
            log.warn("A live settlement execution for {} is already running; exiting", targetDate)
            return 1
        }

        val last = jobRepository.getLastJobExecution(dailySettlementJob.name, baseParameters)
        val parameters = when {
            last == null -> baseParameters
            last.status == BatchStatus.COMPLETED -> {
                log.info("Settlement for {} is already complete; nothing to do", targetDate)
                return 0
            }
            else -> {
                log.warn(
                    "Prior settlement execution for {} ended {}; rebuilding as a new instance",
                    targetDate, last.status,
                )
                JobParametersBuilder(baseParameters)
                    .addString("recovery.id", UUID.randomUUID().toString(), true)
                    .toJobParameters()
            }
        }

        return try {
            val execution = jobOperator.start(dailySettlementJob, parameters)
            if (execution.status == BatchStatus.COMPLETED) 0 else 1
        } catch (e: JobInstanceAlreadyCompleteException) {
            0
        } catch (e: JobExecutionAlreadyRunningException) {
            log.warn("Settlement for {} started concurrently elsewhere; exiting", targetDate)
            1
        }
    }
}
