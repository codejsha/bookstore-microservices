package com.codejsha.bookstore.settlement.application.port

import com.codejsha.bookstore.settlement.application.SettlementRunConflictException
import java.time.LocalDate
import java.time.LocalDateTime
import java.util.UUID

interface SettlementJobLauncher {

    fun runningTargetDates(): Set<LocalDate>

    fun launch(targetDate: LocalDate, runUid: UUID, rerun: Boolean): LaunchedRun
}

data class LaunchedRun(
    val status: String,
    val startedAt: LocalDateTime,
)
