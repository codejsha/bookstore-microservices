package com.codejsha.bookstore.payment.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.api.MandateApi
import com.codejsha.bookstore.generated.application.port.openapi.model.*
import com.codejsha.bookstore.payment.application.usecase.MandateUseCase
import com.codejsha.bookstore.payment.domain.aggregate.MandateAggregate
import com.codejsha.bookstore.payment.domain.model.command.MandateSetupCommand
import com.codejsha.bookstore.payment.domain.model.option.MandateQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.payment.infrastructure.support.auth.Principal
import com.codejsha.bookstore.payment.infrastructure.support.auth.isStaff
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController
import org.springframework.web.server.ResponseStatusException
import java.time.ZoneOffset
import java.util.*

@RestController
class MandateController(
    private val mandateUseCase: MandateUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : MandateApi {

    override fun mandatesGetAll(
        customerId: String?,
        mandateStatus: MandateStatus?,
        pageable: Pageable?
    ): ResponseEntity<MandateFindAllResponse> = runBlocking {
        val principal = principalResolver.require()
        val ownerFilter = if (principal.isStaff()) customerId else principal.sub
        val option = MandateQueryOption(customerId = ownerFilter, mandateStatus = mandateStatus?.value)
        val context = buildContext()

        val result = mandateUseCase.findAllMandates(option, pageable ?: Pageable.unpaged(), context)
        val response = MandateFindAllResponse(
            total = result.totalElements,
            items = result.content.map { toMandateFindResponse(it) }
        )
        ResponseEntity.ok(response)
    }

    override fun mandatesSetup(body: MandateSetupRequest): ResponseEntity<MandateFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val command = MandateSetupCommand(
            customerId = principal.sub,
            paymentMethodToken = body.paymentMethodToken,
            currency = body.mandateCurrency,
            mandateAmountMinor = body.mandateAmount,
        )
        val mandate = mandateUseCase.setupMandate(command, buildContext())
        ResponseEntity.ok(toMandateFindResponse(mandate))
    }

    override fun mandatesRead(uid: String): ResponseEntity<MandateFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val mandate = mandateUseCase.findMandate(UUID.fromString(uid), buildContext())
        assertOwnerOrStaff(principal, mandate)
        ResponseEntity.ok(toMandateFindResponse(mandate))
    }

    override fun mandatesRevoke(uid: String): ResponseEntity<MandateFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        assertOwnerOrStaff(principal, mandateUseCase.findMandate(UUID.fromString(uid), context))
        val mandate = mandateUseCase.revokeMandate(UUID.fromString(uid), context)
        ResponseEntity.ok(toMandateFindResponse(mandate))
    }

    // ─── Authorization ──────────────────────────────────────────────────────

    private fun assertOwnerOrStaff(principal: Principal, mandate: MandateAggregate) {
        if (principal.isStaff() || mandate.customerId == principal.sub) return
        throw ResponseStatusException(HttpStatus.NOT_FOUND, "Mandate not found")
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun toMandateFindResponse(agg: MandateAggregate) = MandateFindResponse(
        uid = agg.uid.toString(),
        mandateId = agg.mandateId,
        customerId = agg.customerId,
        paymentMethodId = agg.paymentMethodId,
        mandateType = MandateType.fromValue(agg.mandateType.value),
        mandateStatus = MandateStatus.fromValue(agg.mandateStatus.value),
        mandateAmount = agg.mandateAmount,
        mandateCurrency = agg.mandateCurrency,
        startDate = agg.startDate?.atOffset(ZoneOffset.UTC),
        endDate = agg.endDate?.atOffset(ZoneOffset.UTC),
        setupFutureUsage = agg.setupFutureUsage?.let { FutureUsage.fromValue(it.value) },
        customerAcceptanceType = agg.customerAcceptanceType,
        customerAcceptedAt = agg.customerAcceptedAt?.atOffset(ZoneOffset.UTC),
        createdAt = agg.createdAt.atOffset(ZoneOffset.UTC),
        updatedAt = agg.updatedAt?.atOffset(ZoneOffset.UTC),
    )

    private fun buildContext() = ActorContext(actorId = 0L, ActorType.USER)
}
