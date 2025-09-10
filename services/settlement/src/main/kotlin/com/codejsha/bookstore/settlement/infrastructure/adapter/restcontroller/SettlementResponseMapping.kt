package com.codejsha.bookstore.settlement.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementDetailResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementFindResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementRunResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementSourceType
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementStatus
import com.codejsha.bookstore.settlement.domain.aggregate.DailySettlement
import com.codejsha.bookstore.settlement.domain.aggregate.SettlementDetailEntity
import com.codejsha.bookstore.settlement.domain.aggregate.SettlementRun
import java.time.ZoneOffset

internal fun toSettlementFindResponse(agg: DailySettlement) = SettlementFindResponse(
    uid = agg.uid.toString(),
    settlementDate = agg.settlementDate,
    currency = agg.currency,
    paymentMethod = agg.paymentMethod,
    grossAmount = agg.grossAmount,
    refundAmount = agg.refundAmount,
    feeAmount = agg.feeAmount,
    netAmount = agg.netAmount,
    paymentCount = agg.paymentCount,
    refundCount = agg.refundCount,
    status = SettlementStatus.fromValue(agg.status.value),
    details = agg.details.takeIf { it.isNotEmpty() }?.map { toSettlementDetailResponse(it) },
    createdAt = agg.createdAt.atOffset(ZoneOffset.UTC),
    updatedAt = agg.updatedAt?.atOffset(ZoneOffset.UTC),
)

internal fun toSettlementRunResponse(run: SettlementRun) = SettlementRunResponse(
    uid = run.uid.toString(),
    targetDate = run.targetDate,
    status = run.status,
    startedAt = run.startedAt.atOffset(ZoneOffset.UTC),
)

internal fun toSettlementDetailResponse(entity: SettlementDetailEntity) = SettlementDetailResponse(
    uid = entity.uid.toString(),
    sourceType = SettlementSourceType.fromValue(entity.sourceType.value),
    sourceId = entity.sourceId,
    paymentId = entity.paymentId,
    amount = entity.amount,
    currency = entity.currency,
    paymentMethod = entity.paymentMethod,
    occurredAt = entity.occurredAt.atOffset(ZoneOffset.UTC),
)
