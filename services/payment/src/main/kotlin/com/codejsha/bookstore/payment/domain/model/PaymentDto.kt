package com.codejsha.bookstore.payment.domain.model

import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType
import java.math.BigDecimal
import java.time.LocalDateTime

data class PaymentDto(
    val id: Long?,
    val orderId: Long?,
    val userId: String?,
    val paymentType: PaymentType?,
    val cardNumber: String?,
    val amount: BigDecimal?,
    val paymentDate: LocalDateTime? = null
)
