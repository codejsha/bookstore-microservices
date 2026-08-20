package com.codejsha.bookstore.settlement.application.port

import java.time.LocalDate
import java.time.LocalDateTime

interface PgReconciliationPort {

    fun fetchDailyTotals(
        targetDate: LocalDate,
        startUtc: LocalDateTime,
        endUtc: LocalDateTime,
    ): List<CurrencyTotal>

    data class CurrencyTotal(
        val currency: String,
        val count: Long,
        val grossAmount: Long,
    )
}
