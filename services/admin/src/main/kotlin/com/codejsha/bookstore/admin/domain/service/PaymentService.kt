package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.PaymentClient
import com.codejsha.bookstore.admin.application.usecase.PaymentUseCase
import com.codejsha.bookstore.admin.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.admin.domain.model.option.RefundQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service

@Service
class PaymentService(
    private val paymentClient: PaymentClient,
) : PaymentUseCase {

    override suspend fun findAllPayments(option: PaymentQueryOption, pageable: Pageable, context: ActorContext) =
        paymentClient.findAllPayments(option, pageable)

    override suspend fun findPayment(uid: String, context: ActorContext) = paymentClient.findPayment(uid)

    override suspend fun findAllRefunds(option: RefundQueryOption, pageable: Pageable, context: ActorContext) =
        paymentClient.findAllRefunds(option, pageable)

    override suspend fun findRefund(uid: String, context: ActorContext) = paymentClient.findRefund(uid)
}
