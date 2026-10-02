package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.application.port.restclient.PaymentClient
import com.codejsha.bookstore.admin.domain.model.external.Payment
import com.codejsha.bookstore.admin.domain.model.external.Refund
import com.codejsha.bookstore.admin.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.admin.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.generated.application.port.restclient.payment.api.PaymentApi
import com.codejsha.bookstore.generated.application.port.restclient.payment.api.RefundApi
import com.codejsha.bookstore.generated.application.port.restclient.payment.model.PaymentFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.payment.model.PaymentStatus
import com.codejsha.bookstore.generated.application.port.restclient.payment.model.RefundFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.payment.model.RefundStatus
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Component

@Component
class PaymentRestClient(
    private val paymentApi: PaymentApi,
    private val refundApi: RefundApi,
) : PaymentClient {

    override fun findAllPayments(option: PaymentQueryOption, pageable: Pageable): Page<Payment> {
        val response = paymentApi.paymentsGetAll(
            option.customerId,
            option.status?.toPaymentStatus(),
            option.connector,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toPayment() }, pageable, response.total)
    }

    override fun findPayment(uid: String): Payment = paymentApi.paymentsRead(uid).toPayment()

    override fun findAllRefunds(option: RefundQueryOption, pageable: Pageable): Page<Refund> {
        val response = refundApi.refundsGetAll(
            option.paymentId,
            option.status?.toRefundStatus(),
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toRefund() }, pageable, response.total)
    }

    override fun findRefund(uid: String): Refund = refundApi.refundsRead(uid).toRefund()

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun String.toPaymentStatus(): PaymentStatus =
        parseEnumValue("payment status", this, PaymentStatus.Companion::fromValue)

    private fun String.toRefundStatus(): RefundStatus =
        parseEnumValue("refund status", this, RefundStatus.Companion::fromValue)

    private fun PaymentFindResponse.toPayment() = Payment(
        uid = uid,
        paymentId = paymentId,
        customerId = customerId,
        status = status.value,
        amount = amount,
        amountCaptured = amountCaptured,
        amountCapturable = amountCapturable,
        currency = currency,
        paymentMethod = paymentMethod?.value,
        connector = connector,
        errorCode = errorCode,
        errorMessage = errorMessage,
        confirmedAt = confirmedAt,
        capturedAt = capturedAt,
        cancelledAt = cancelledAt,
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun RefundFindResponse.toRefund() = Refund(
        uid = uid,
        refundId = refundId,
        paymentId = paymentId,
        status = status.value,
        refundType = refundType.value,
        amount = amount,
        currency = currency,
        reason = reason,
        connector = connector,
        errorCode = errorCode,
        errorMessage = errorMessage,
        createdAt = createdAt,
        updatedAt = updatedAt,
    )
}
