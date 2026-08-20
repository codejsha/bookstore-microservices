package com.codejsha.bookstore.settlement.domain.aggregate

import com.codejsha.bookstore.settlement.application.port.repo.DailySettlementResult
import com.codejsha.bookstore.settlement.application.port.repo.SettlementDetailResult
import com.codejsha.bookstore.settlement.domain.constant.SettlementSourceType
import com.codejsha.bookstore.settlement.domain.constant.SettlementStatus
import java.time.LocalDate
import java.time.LocalDateTime
import java.util.UUID

data class DailySettlement(
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
    val status: SettlementStatus,
    val details: List<SettlementDetailEntity>,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class SettlementDetailEntity(
    val id: Long,
    val uid: UUID,
    val settlementDate: LocalDate,
    val sourceType: SettlementSourceType,
    val sourceId: String,
    val paymentId: String,
    val amount: Long,
    val currency: String,
    val paymentMethod: String?,
    val occurredAt: LocalDateTime,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

fun DailySettlementResult.toAggregate() = DailySettlement(
    id = id,
    uid = uid,
    settlementDate = settlementDate,
    currency = currency,
    paymentMethod = paymentMethod,
    grossAmount = grossAmount,
    refundAmount = refundAmount,
    feeAmount = feeAmount,
    netAmount = netAmount,
    paymentCount = paymentCount,
    refundCount = refundCount,
    status = SettlementStatus.fromValue(status),
    details = emptyList(),
    createdAt = createdAt,
    updatedAt = updatedAt,
)

fun SettlementDetailResult.toEntity() = SettlementDetailEntity(
    id = id,
    uid = uid,
    settlementDate = settlementDate,
    sourceType = SettlementSourceType.fromValue(sourceType),
    sourceId = sourceId,
    paymentId = paymentId,
    amount = amount,
    currency = currency,
    paymentMethod = paymentMethod,
    occurredAt = occurredAt,
    createdAt = createdAt,
    updatedAt = updatedAt,
)
