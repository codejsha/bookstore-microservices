package com.codejsha.bookstore.settlement.infrastructure.batch

import java.time.LocalDate
import java.time.LocalDateTime
import java.util.UUID

data class PaymentRow(
    val paymentId: String,
    val status: String,
    val amount: Long,
    val amountCaptured: Long?,
    val currency: String,
    val paymentMethod: String?,
    val capturedAt: LocalDateTime,
)

data class RefundRow(
    val refundId: String,
    val paymentId: String,
    val amount: Long,
    val currency: String,
    val paymentMethod: String?,
    val createdAt: LocalDateTime,
)

data class SettlementDetailInsert(
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
)
