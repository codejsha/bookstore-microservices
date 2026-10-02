package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.application.port.restclient.SettlementClient
import com.codejsha.bookstore.admin.domain.model.external.SettlementBucket
import com.codejsha.bookstore.admin.domain.model.external.SettlementDetailLine
import com.codejsha.bookstore.admin.domain.model.external.SettlementRunAck
import com.codejsha.bookstore.admin.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.admin.domain.model.option.SettlementQueryOption
import com.codejsha.bookstore.generated.application.port.restclient.settlement.api.SettlementApi
import com.codejsha.bookstore.generated.application.port.restclient.settlement.model.SettlementDetailResponse
import com.codejsha.bookstore.generated.application.port.restclient.settlement.model.SettlementFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.settlement.model.SettlementRunRequest
import com.codejsha.bookstore.generated.application.port.restclient.settlement.model.SettlementRunResponse
import com.codejsha.bookstore.generated.application.port.restclient.settlement.model.SettlementStatus
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Component

@Component
class SettlementRestClient(
    private val settlementApi: SettlementApi,
) : SettlementClient {

    override fun findAllSettlements(option: SettlementQueryOption, pageable: Pageable): Page<SettlementBucket> {
        val response = settlementApi.settlementsGetAll(
            option.settlementDate,
            option.status?.toSettlementStatus(),
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toBucket() }, pageable, response.total)
    }

    override fun findSettlement(uid: String): SettlementBucket =
        settlementApi.settlementsRead(uid).toBucket()

    override fun triggerSettlementRun(command: TriggerSettlementRunCommand): SettlementRunAck =
        settlementApi.settlementRunsCreate(
            SettlementRunRequest(targetDate = command.targetDate, rerun = command.rerun),
        ).toAck()

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun String.toSettlementStatus(): SettlementStatus =
        parseEnumValue("settlement status", this, SettlementStatus.Companion::fromValue)

    private fun SettlementFindResponse.toBucket() = SettlementBucket(
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
        status = status.value,
        details = details.orEmpty().map { it.toLine() },
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun SettlementDetailResponse.toLine() = SettlementDetailLine(
        uid = uid,
        sourceType = sourceType.value,
        sourceId = sourceId,
        paymentId = paymentId,
        amount = amount,
        currency = currency,
        paymentMethod = paymentMethod,
        occurredAt = occurredAt,
    )

    private fun SettlementRunResponse.toAck() = SettlementRunAck(
        uid = uid,
        targetDate = targetDate,
        status = status,
        startedAt = startedAt,
    )
}
