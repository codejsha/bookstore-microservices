package com.codejsha.bookstore.payment.domain.model

import com.codejsha.bookstore.payment.domain.aggregate.entity.PaymentEntity
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType
import java.time.LocalDateTime

data class PaymentDto(
    val id: Long?,
    val orderId: Long?,
    val userId: String?,
    val paymentType: PaymentType?,
    val cardNumber: String?,
    val amount: Double?,
    val paymentDate: LocalDateTime?
) {

    fun toEntity() = PaymentEntity(
        id = this.id!!,
        orderId = this.orderId!!,
        userId = this.userId!!,
        paymentType = this.paymentType!!,
        cardNumber = this.cardNumber!!,
        amount = this.amount!!,
        paymentDate = this.paymentDate!!,
        createdAt = null,
        updatedAt = null
    )
}
