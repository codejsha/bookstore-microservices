package com.codejsha.bookstore.payment.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.api.PaymentApi
import com.codejsha.bookstore.generated.application.port.openapi.api.PaymentAttemptApi
import com.codejsha.bookstore.generated.application.port.openapi.model.*
import com.codejsha.bookstore.payment.application.usecase.PaymentUseCase
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAggregate
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAttemptEntity
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.auth.BadRequestException
import com.codejsha.bookstore.payment.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.payment.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.payment.infrastructure.support.auth.Principal
import com.codejsha.bookstore.payment.infrastructure.support.auth.ROLE_ADMIN
import com.codejsha.bookstore.payment.infrastructure.support.auth.assertPaymentOwner
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
class PaymentController(
    private val paymentUseCase: PaymentUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : PaymentApi, PaymentAttemptApi {

    // ─── PaymentApi ─────────────────────────────────────────────────────────

    override fun paymentsGetAll(
        customerId: String?,
        status: PaymentStatus?,
        connector: String?,
        pageable: Pageable?
    ): ResponseEntity<PaymentFindAllResponse> = runBlocking {
        val principal = principalResolver.require()
        val ownerFilter = if (principal.hasRole(ROLE_ADMIN)) customerId else principal.sub
        val option = PaymentQueryOption(customerId = ownerFilter, status = status?.value, connector = connector)
        val context = buildContext()

        val result = paymentUseCase.findAllPayments(option, pageable ?: Pageable.unpaged(), context)
        val response = PaymentFindAllResponse(
            total = result.totalElements,
            items = result.content.map { toPaymentFindResponse(it) }
        )
        ResponseEntity.ok(response)
    }

    override fun paymentsCreate(requestBody: PaymentCreateRequest): ResponseEntity<Unit> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        val ownerId = createOwnerOf(principal, requestBody.customerId)
        val command = PaymentCreateCommand(
            amount = requestBody.amount,
            currency = requestBody.currency,
            customerId = ownerId,
            paymentMethod = requestBody.paymentMethod?.value,
            paymentMethodType = requestBody.paymentMethodType,
            authenticationType = requestBody.authenticationType?.value,
            setupFutureUsage = requestBody.setupFutureUsage?.value,
            description = requestBody.description,
            returnUrl = requestBody.returnUrl,
            billingAddress = null,
            shippingAddress = null,
            metadata = null,
            idempotencyKey = "$ownerId:${requestBody.idempotencyKey}",
        )
        val payment = paymentUseCase.createPayment(command, context)
        ResponseEntity.created(URI.create("/api/v1/payments/${payment.uid}")).build()
    }

    override fun paymentsRead(uid: String): ResponseEntity<PaymentFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        val payment = paymentUseCase.findPayment(UUID.fromString(uid), context)
        principal.assertPaymentOwner(payment.customerId)
        ResponseEntity.ok(toPaymentFindResponse(payment))
    }

    override fun paymentsUpdate(
        uid: String, requestBody: PaymentUpdateRequest
    ): ResponseEntity<PaymentUpdateResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        principal.assertPaymentOwner(paymentUseCase.findPayment(UUID.fromString(uid), context).customerId)
        val command = PaymentUpdateCommand(
            amount = requestBody.amount,
            currency = requestBody.currency,
            customerId = updateOwnerOf(principal, requestBody.customerId),
            paymentMethod = requestBody.paymentMethod?.value,
            paymentMethodType = requestBody.paymentMethodType,
            authenticationType = requestBody.authenticationType?.value,
            setupFutureUsage = requestBody.setupFutureUsage?.value,
            description = requestBody.description,
            returnUrl = requestBody.returnUrl,
            billingAddress = null,
            shippingAddress = null,
            metadata = null,
        )
        val payment = paymentUseCase.updatePayment(UUID.fromString(uid), command, context)
        ResponseEntity.ok(toPaymentUpdateResponse(payment))
    }

    // ─── PaymentAttemptApi ──────────────────────────────────────────────────

    override fun paymentAttemptsGetAll(
        paymentUid: String,
        status: PaymentStatus?,
        pageable: Pageable?
    ): ResponseEntity<PaymentAttemptFindAllResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        authorizePayment(paymentUid, principal, context)

        val result =
            paymentUseCase.findAllPaymentAttempts(UUID.fromString(paymentUid), pageable ?: Pageable.unpaged(), context)
        val response = PaymentAttemptFindAllResponse(
            total = result.totalElements,
            items = result.content.map { toPaymentAttemptFindResponse(it) }
        )
        ResponseEntity.ok(response)
    }

    override fun paymentAttemptsRead(
        paymentUid: String, uid: String
    ): ResponseEntity<PaymentAttemptFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext()
        authorizePayment(paymentUid, principal, context)
        val attempt = paymentUseCase.findPaymentAttempt(UUID.fromString(paymentUid), UUID.fromString(uid), context)
        ResponseEntity.ok(toPaymentAttemptFindResponse(attempt))
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun toPaymentFindResponse(agg: PaymentAggregate) = PaymentFindResponse(
        uid = agg.uid.toString(),
        paymentId = agg.paymentId,
        merchantId = agg.merchantId,
        profileId = agg.profileId,
        customerId = agg.customerId,
        paymentMethodId = agg.paymentMethodId,
        mandateId = agg.mandateId,
        connector = agg.connector?.connector,
        connectorTransactionId = agg.connector?.transactionId,
        amount = agg.amount.amount,
        currency = agg.amount.currency,
        amountCapturable = agg.amountCapturable,
        amountCaptured = agg.amountCaptured,
        surchargeAmount = agg.surchargeAmount,
        taxAmount = agg.taxAmount,
        status = PaymentStatus.fromValue(agg.status.value),
        captureMethod = CaptureMethod.fromValue(agg.captureMethod.value),
        authenticationType = agg.authenticationType?.let { AuthenticationType.fromValue(it.value) },
        paymentMethod = agg.paymentMethod?.let { PaymentMethodEnum.fromValue(it.value) },
        paymentMethodType = agg.paymentMethodType,
        clientSecret = agg.clientSecret,
        setupFutureUsage = agg.setupFutureUsage?.let { FutureUsage.fromValue(it.value) },
        offSession = agg.offSession,
        description = agg.description,
        returnUrl = agg.returnUrl,
        statementDescriptor = agg.statementDescriptor,
        errorCode = agg.errorCode,
        errorMessage = agg.errorMessage,
        confirmedAt = agg.confirmedAt?.atOffset(ZoneOffset.UTC),
        capturedAt = agg.capturedAt?.atOffset(ZoneOffset.UTC),
        cancelledAt = agg.cancelledAt?.atOffset(ZoneOffset.UTC),
        createdAt = agg.createdAt.atOffset(ZoneOffset.UTC),
        updatedAt = agg.updatedAt?.atOffset(ZoneOffset.UTC),
        attempts = agg.attempts.map { toPaymentAttemptFindResponse(it) },
    )

    private fun toPaymentUpdateResponse(agg: PaymentAggregate) = PaymentUpdateResponse(
        uid = agg.uid.toString(),
        paymentId = agg.paymentId,
        amount = agg.amount.amount,
        currency = agg.amount.currency,
        status = PaymentStatus.fromValue(agg.status.value),
        captureMethod = CaptureMethod.fromValue(agg.captureMethod.value),
        offSession = agg.offSession,
        createdAt = agg.createdAt.atOffset(ZoneOffset.UTC),
        updatedAt = agg.updatedAt?.atOffset(ZoneOffset.UTC),
    )

    private fun toPaymentAttemptFindResponse(entity: PaymentAttemptEntity) = PaymentAttemptFindResponse(
        uid = entity.uid.toString(),
        attemptId = entity.attemptId,
        paymentId = entity.paymentId,
        connector = entity.connector?.connector,
        connectorTransactionId = entity.connector?.transactionId,
        amount = entity.amount.amount,
        currency = entity.amount.currency,
        status = PaymentStatus.fromValue(entity.status.value),
        authenticationType = entity.authenticationType?.let { AuthenticationType.fromValue(it.value) },
        paymentMethod = entity.paymentMethod?.let { PaymentMethodEnum.fromValue(it.value) },
        paymentMethodType = entity.paymentMethodType,
        errorCode = entity.errorCode,
        errorMessage = entity.errorMessage,
        createdAt = entity.createdAt.atOffset(ZoneOffset.UTC),
    )

    private suspend fun authorizePayment(paymentUid: String, principal: Principal, context: ActorContext) {
        val payment = paymentUseCase.findPayment(UUID.fromString(paymentUid), context)
        principal.assertPaymentOwner(payment.customerId)
    }

    private fun createOwnerOf(principal: Principal, requestedCustomerId: String?): String =
        if (principal.hasRole(ROLE_ADMIN)) {
            requestedCustomerId ?: throw BadRequestException("customer_id is required when acting as admin")
        } else {
            principal.sub
        }

    private fun updateOwnerOf(principal: Principal, requestedCustomerId: String?): String? {
        if (principal.hasRole(ROLE_ADMIN)) return requestedCustomerId
        if (requestedCustomerId != null && requestedCustomerId != principal.sub) {
            throw ForbiddenException("cannot reassign a payment to another customer")
        }
        return requestedCustomerId
    }

    private fun buildContext() = ActorContext(actorId = 0L, ActorType.USER)
}
