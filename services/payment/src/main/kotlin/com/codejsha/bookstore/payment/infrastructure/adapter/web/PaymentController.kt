package com.codejsha.bookstore.payment.infrastructure.adapter.web

import com.codejsha.bookstore.payment.application.usecase.PaymentUseCase
import com.codejsha.bookstore.payment.domain.model.FilterCondition
import com.codejsha.bookstore.service.application.port.openapi.api.PaymentApi
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentFindAllWebResp
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentFindWebResp
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
class PaymentController(
    private val paymentUseCase: PaymentUseCase
) : PaymentApi {

    override fun apiV1PaymentsGet(
        sort: String,
        order: String,
        limit: Int,
        offset: Int
    ): ResponseEntity<PaymentFindAllWebResp> {
        val cond = FilterCondition(
            filter = FilterCondition.QueryFilter(
                sort = sort,
                order = order,
                limit = limit,
                offset = offset
            )
        )
        val payments = paymentUseCase.findAllPayments(cond)
        val response = PaymentFindAllWebResp(
            total = payments.size.toLong(),
            items = payments.map { it.toPaymentFindWebResp() }
        )
        return ResponseEntity.ok(response)
    }

    override fun apiV1PaymentsIdGet(id: Long): ResponseEntity<PaymentFindWebResp> {
        val payment = paymentUseCase.findPayment(id)
        val response = payment.toPaymentFindWebResp()
        return ResponseEntity.ok(response)
    }
}
