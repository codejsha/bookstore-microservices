package com.codejsha.bookstore.settlement.infrastructure.batch

import java.time.LocalDate
import java.time.LocalDateTime
import java.time.ZoneId
import java.time.ZoneOffset

data class SettlementWindow(
    val startUtc: LocalDateTime,
    val endUtc: LocalDateTime,
) {
    companion object {
        fun of(targetDate: LocalDate, zone: ZoneId): SettlementWindow {
            val start = targetDate.atStartOfDay(zone).toInstant().atZone(ZoneOffset.UTC).toLocalDateTime()
            val end = targetDate.plusDays(1).atStartOfDay(zone).toInstant().atZone(ZoneOffset.UTC).toLocalDateTime()
            return SettlementWindow(start, end)
        }
    }
}
