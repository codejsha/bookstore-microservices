package com.codejsha.bookstore.order.domain.workflow

import io.temporal.activity.ActivityInterface
import io.temporal.activity.ActivityMethod

@ActivityInterface
interface PaymentActivities {
    @ActivityMethod
    fun processPayment(request: ProcessPaymentRequest): ProcessPaymentResult

    @ActivityMethod
    fun refundPayment(paymentUid: String): Boolean

    @ActivityMethod
    fun refundPaymentByOrder(orderUid: String): Boolean
}
