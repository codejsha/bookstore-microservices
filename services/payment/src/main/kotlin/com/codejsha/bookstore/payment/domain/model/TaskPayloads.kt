package com.codejsha.bookstore.payment.domain.model

import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType
import java.math.BigDecimal

data class MakePaymentPayload(
    val orderId: Long,
    val userId: String,
    val paymentType: PaymentType,
    val cardNumber: String,
    val amount: BigDecimal
) {

    fun toPaymentDto(): PaymentDto {
        return PaymentDto(
            id = null,
            orderId = orderId,
            userId = userId,
            paymentType = paymentType,
            cardNumber = cardNumber,
            amount = amount
        )
    }
}
