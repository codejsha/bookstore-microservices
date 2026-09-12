package com.codejsha.bookstore.payment.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.api.RefundApi
import com.codejsha.bookstore.generated.application.port.openapi.model.*
import com.codejsha.bookstore.payment.application.usecase.RefundUseCase
import com.codejsha.bookstore.payment.domain.aggregate.RefundAggregate
import com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.payment.infrastructure.support.auth.assertManager
import com.codejsha.bookstore.payment.infrastructure.support.auth.assertStaff
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.springframework.data.domain.Pageable
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController
import java.net.URI
import java.time.ZoneOffset
import java.util.*

@RestController
class RefundController(
    private val refundUseCase: RefundUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : RefundApi {

    override fun refundsGetAll(
        paymentId: String?,
        status: RefundStatus?,
        pageable: Pageable?
    ): ResponseEntity<RefundFindAllResponse> = runBlocking {
        principalResolver.require().assertStaff()
        val option = RefundQueryOption(paymentId = paymentId, status = status?.value)
        val context = buildContext()

        val result = refundUseCase.findAllRefunds(option, pageable ?: Pageable.unpaged(), context)
        val response = RefundFindAllResponse(
            total = result.totalElements,
            items = result.content.map { toRefundFindResponse(it) }
        )
        ResponseEntity.ok(response)
    }

    override fun refundsCreate(requestBody: RefundCreateRequest): ResponseEntity<Unit> = runBlocking {
        principalResolver.require().assertManager()
        val context = buildContext()
        val command = RefundCreateCommand(
            paymentId = requestBody.paymentId,
            amount = requestBody.amount,
            currency = requestBody.currency,
            reason = requestBody.reason,
            refundType = requestBody.refundType?.value,
            metadata = null,
            idempotencyKey = "${requestBody.paymentId}:${requestBody.idempotencyKey}",
        )
        val refund = refundUseCase.createRefund(command, context)
        ResponseEntity.created(URI.create("/api/v1/refunds/${refund.uid}")).build()
    }

    override fun refundsRead(uid: String): ResponseEntity<RefundFindResponse> = runBlocking {
        principalResolver.require().assertStaff()
        val context = buildContext()
        val refund = refundUseCase.findRefund(UUID.fromString(uid), context)
        ResponseEntity.ok(toRefundFindResponse(refund))
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun toRefundFindResponse(agg: RefundAggregate) = RefundFindResponse(
        uid = agg.uid.toString(),
        refundId = agg.refundId,
        paymentId = agg.paymentId,
        connector = agg.connector,
        connectorRefundId = agg.connectorRefundId,
        amount = agg.amount.amount,
        currency = agg.amount.currency,
        status = RefundStatus.fromValue(agg.status.value),
        reason = agg.reason,
        refundType = RefundType.fromValue(agg.refundType.value),
        errorCode = agg.errorCode,
        errorMessage = agg.errorMessage,
        createdAt = agg.createdAt.atOffset(ZoneOffset.UTC),
        updatedAt = agg.updatedAt?.atOffset(ZoneOffset.UTC),
    )

    private fun buildContext() = ActorContext(actorId = 0L, ActorType.USER)
}
