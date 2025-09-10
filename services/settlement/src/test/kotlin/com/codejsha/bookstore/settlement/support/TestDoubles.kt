package com.codejsha.bookstore.settlement.support

import com.codejsha.bookstore.settlement.application.SettlementRunConflictException
import com.codejsha.bookstore.settlement.application.port.LaunchedRun
import com.codejsha.bookstore.settlement.application.port.SettlementJobLauncher
import com.codejsha.bookstore.settlement.application.port.repo.DailySettlementRepo
import com.codejsha.bookstore.settlement.application.port.repo.DailySettlementResult
import com.codejsha.bookstore.settlement.application.port.repo.SettlementDetailRepo
import com.codejsha.bookstore.settlement.application.port.repo.SettlementDetailResult
import com.codejsha.bookstore.settlement.application.port.support.TransactionRunner
import com.codejsha.bookstore.settlement.domain.model.option.SettlementQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.util.UUID

class FakeTransactionRunner : TransactionRunner {
    override suspend fun <T> tx(block: () -> T): T = block()
}

class FakeDailySettlementRepo(
    private val rows: List<DailySettlementResult>,
) : DailySettlementRepo {
    override fun findAll(
        option: SettlementQueryOption, pageable: Pageable, context: ActorContext
    ): Page<DailySettlementResult> {
        val filtered = rows.filter { row ->
            (option.settlementDate == null || row.settlementDate == option.settlementDate) &&
                (option.status == null || row.status == option.status)
        }
        return PageImpl(filtered, pageable, filtered.size.toLong())
    }

    override fun findOne(uid: UUID, context: ActorContext): DailySettlementResult =
        rows.firstOrNull { it.uid == uid }
            ?: throw NoSuchElementException("Settlement with uid $uid not found")
}

class FakeSettlementJobLauncher(
    private val running: Set<LocalDate> = emptySet(),
    private val launchError: RuntimeException? = null,
) : SettlementJobLauncher {

    var lastLaunch: LaunchArgs? = null
        private set

    data class LaunchArgs(val targetDate: LocalDate, val runUid: UUID, val rerun: Boolean)

    override fun runningTargetDates(): Set<LocalDate> = running

    override fun launch(targetDate: LocalDate, runUid: UUID, rerun: Boolean): LaunchedRun {
        launchError?.let { throw it }
        lastLaunch = LaunchArgs(targetDate, runUid, rerun)
        return LaunchedRun(status = "STARTED", startedAt = LocalDateTime.now(ZoneOffset.UTC))
    }
}

fun conflictingLauncher(message: String = "already running"): FakeSettlementJobLauncher =
    FakeSettlementJobLauncher(launchError = SettlementRunConflictException(message))

class FakeSettlementDetailRepo(
    private val rows: List<SettlementDetailResult>,
) : SettlementDetailRepo {
    override fun findAllByBucket(
        settlementDate: LocalDate,
        currency: String,
        paymentMethod: String?,
        context: ActorContext,
    ): List<SettlementDetailResult> =
        rows.filter {
            it.settlementDate == settlementDate &&
                it.currency == currency &&
                it.paymentMethod == paymentMethod
        }
}
