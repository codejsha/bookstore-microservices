package com.codejsha.bookstore.settlement.domain.aggregate

import java.time.LocalDate
import java.time.LocalDateTime
import java.util.UUID

data class SettlementRun(
    val uid: UUID,
    val targetDate: LocalDate,
    val status: String,
    val startedAt: LocalDateTime,
)
