package com.codejsha.bookstore.payment.application.port.saga

import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType
import java.math.BigDecimal

data class MakePaymentPayload(
    val requestId: String,
    val orderId: Long,
    val userId: String,
    val paymentType: PaymentType,
    val cardNumber: String,
    val amount: BigDecimal
)
