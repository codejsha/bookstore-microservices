package com.codejsha.bookstore.settlement.application.port.repo

import java.time.LocalDate
import java.time.LocalDateTime
import java.util.UUID

data class DailySettlementResult(
    val id: Long,
    val uid: UUID,
    val settlementDate: LocalDate,
    val currency: String,
    val paymentMethod: String?,
    val grossAmount: Long,
    val refundAmount: Long,
    val feeAmount: Long,
    val netAmount: Long,
    val paymentCount: Int,
    val refundCount: Int,
    val status: String,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class SettlementDetailResult(
    val id: Long,
    val uid: UUID,
    val settlementDate: LocalDate,
    val sourceType: String,
    val sourceId: String,
    val paymentId: String,
    val amount: Long,
    val currency: String,
    val paymentMethod: String?,
    val occurredAt: LocalDateTime,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)
