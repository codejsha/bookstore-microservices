package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.PaymentUseCase
import com.codejsha.bookstore.admin.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.admin.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.admin.infrastructure.support.auth.Principal
import com.codejsha.bookstore.admin.infrastructure.support.auth.assertAdmin
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminPaymentApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminPaymentFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminPaymentResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRefundFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRefundResponse
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.springframework.data.domain.Pageable
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
class AdminPaymentController(
    private val paymentUseCase: PaymentUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : AdminPaymentApi {

    override fun adminPaymentsListPayments(
        customerId: String?,
        status: String?,
        connector: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminPaymentFindAllResponse> = runBlocking {
        val principal = requireAdmin()
        val option = PaymentQueryOption(customerId = customerId, status = status, connector = connector)
        val context = buildContext(principal)
        val result = paymentUseCase.findAllPayments(option, pageable ?: Pageable.unpaged(), context)
        ResponseEntity.ok(
            AdminPaymentFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminPaymentResponse(it) },
            )
        )
    }

    override fun adminPaymentsReadPayment(uid: String): ResponseEntity<AdminPaymentResponse> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        ResponseEntity.ok(toAdminPaymentResponse(paymentUseCase.findPayment(uid, context)))
    }

    override fun adminPaymentsListRefunds(
        paymentId: String?,
        status: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminRefundFindAllResponse> = runBlocking {
        val principal = requireAdmin()
        val option = RefundQueryOption(paymentId = paymentId, status = status)
        val context = buildContext(principal)
        val result = paymentUseCase.findAllRefunds(option, pageable ?: Pageable.unpaged(), context)
        ResponseEntity.ok(
            AdminRefundFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminRefundResponse(it) },
            )
        )
    }

    override fun adminPaymentsReadRefund(uid: String): ResponseEntity<AdminRefundResponse> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        ResponseEntity.ok(toAdminRefundResponse(paymentUseCase.findRefund(uid, context)))
    }

    private fun requireAdmin(): Principal = principalResolver.require().also { it.assertAdmin() }

    private fun buildContext(principal: Principal) =
        ActorContext(actorId = 0L, ActorType.USER)
}
