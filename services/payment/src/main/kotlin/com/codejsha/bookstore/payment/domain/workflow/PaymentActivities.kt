package com.codejsha.bookstore.payment.domain.workflow

import io.temporal.activity.ActivityInterface
import io.temporal.activity.ActivityMethod
import java.math.BigDecimal

data class ProcessPaymentRequest(
    val orderUid: String,
    val userUid: String,
    val amount: BigDecimal,
    val currency: String,
)

data class ProcessPaymentResult(
    val paymentUid: String,
    val status: String,
    val gatewayPaymentId: String? = null,
)

@ActivityInterface
interface PaymentActivities {
    @ActivityMethod
    fun processPayment(request: ProcessPaymentRequest): ProcessPaymentResult

    @ActivityMethod
    fun refundPayment(paymentUid: String): Boolean

    @ActivityMethod
    fun refundPaymentByOrder(orderUid: String): Boolean
}
