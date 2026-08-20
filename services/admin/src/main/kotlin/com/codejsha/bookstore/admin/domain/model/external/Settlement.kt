package com.codejsha.bookstore.admin.domain.model.external

import java.time.LocalDate
import java.time.OffsetDateTime

data class SettlementBucket(
    val uid: String,
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
    val details: List<SettlementDetailLine>,
    val createdAt: OffsetDateTime,
    val updatedAt: OffsetDateTime?,
)

data class SettlementDetailLine(
    val uid: String,
    val sourceType: String,
    val sourceId: String,
    val paymentId: String,
    val amount: Long,
    val currency: String,
    val paymentMethod: String?,
    val occurredAt: OffsetDateTime,
)

data class SettlementRunAck(
    val uid: String,
    val targetDate: LocalDate,
    val status: String,
    val startedAt: OffsetDateTime,
)
