package com.codejsha.bookstore.settlement.domain.service

import com.codejsha.bookstore.settlement.application.SettlementRunConflictException
import com.codejsha.bookstore.settlement.application.SettlementRunValidationException
import com.codejsha.bookstore.settlement.application.port.SettlementJobLauncher
import com.codejsha.bookstore.settlement.application.port.SettlementRunLock
import com.codejsha.bookstore.settlement.application.usecase.TriggerSettlementRunUseCase
import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import com.codejsha.bookstore.settlement.domain.aggregate.SettlementRun
import com.codejsha.bookstore.settlement.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.platform.shared.data.ActorContext
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.stereotype.Service
import java.time.LocalDate
import java.time.ZoneId
import java.util.UUID
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid

@Service
class SettlementRunService(
    private val launcher: SettlementJobLauncher,
    private val runLock: SettlementRunLock,
    private val properties: SettlementBatchProperties,
) : TriggerSettlementRunUseCase {

    private val zone: ZoneId get() = ZoneId.of(properties.timezone)

    @WithSpan
    override fun triggerSettlementRun(
        command: TriggerSettlementRunCommand,
        context: ActorContext,
    ): SettlementRun {
        val today = LocalDate.now(zone)
        if (!command.targetDate.isBefore(today)) {
            throw SettlementRunValidationException(
                "targetDate ${command.targetDate} must be a closed past day (before $today in $zone)",
            )
        }

        if (command.targetDate in launcher.runningTargetDates()) {
            throw SettlementRunConflictException(
                "A settlement run for ${command.targetDate} is already running",
            )
        }

        val runUid = Uuid.generateV7().toJavaUuid()
        val ownerToken = runUid.toString()
        if (!runLock.tryAcquire(command.targetDate, ownerToken)) {
            throw SettlementRunConflictException(
                "A settlement run for ${command.targetDate} is already running",
            )
        }

        val launched = try {
            launcher.launch(command.targetDate, runUid, command.rerun)
        } catch (e: Throwable) {
            runLock.release(command.targetDate, ownerToken)
            throw e
        }
        return SettlementRun(
            uid = runUid,
            targetDate = command.targetDate,
            status = launched.status,
            startedAt = launched.startedAt,
        )
    }
}
