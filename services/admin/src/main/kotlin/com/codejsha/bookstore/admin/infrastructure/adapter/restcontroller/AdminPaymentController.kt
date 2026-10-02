package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.PaymentUseCase
import com.codejsha.bookstore.admin.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.admin.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminPaymentApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminPaymentFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminPaymentResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRefundFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRefundResponse
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
    ): ResponseEntity<AdminPaymentFindAllResponse> {
        principalResolver.requireStaff()
        val option = PaymentQueryOption(customerId = customerId, status = status, connector = connector)
        val context = buildContext()
        val result = paymentUseCase.findAllPayments(option, pageable ?: Pageable.unpaged(), context)
        return ResponseEntity.ok(
            AdminPaymentFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminPaymentResponse(it) },
            )
        )
    }

    override fun adminPaymentsReadPayment(uid: String): ResponseEntity<AdminPaymentResponse> {
        principalResolver.requireStaff()
        val context = buildContext()
        return ResponseEntity.ok(toAdminPaymentResponse(paymentUseCase.findPayment(uid, context)))
    }

    override fun adminPaymentsListRefunds(
        paymentId: String?,
        status: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminRefundFindAllResponse> {
        principalResolver.requireStaff()
        val option = RefundQueryOption(paymentId = paymentId, status = status)
        val context = buildContext()
        val result = paymentUseCase.findAllRefunds(option, pageable ?: Pageable.unpaged(), context)
        return ResponseEntity.ok(
            AdminRefundFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminRefundResponse(it) },
            )
        )
    }

    override fun adminPaymentsReadRefund(uid: String): ResponseEntity<AdminRefundResponse> {
        principalResolver.requireStaff()
        val context = buildContext()
        return ResponseEntity.ok(toAdminRefundResponse(paymentUseCase.findRefund(uid, context)))
    }
}
