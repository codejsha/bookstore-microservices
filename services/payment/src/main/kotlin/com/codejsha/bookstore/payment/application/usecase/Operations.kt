package com.codejsha.bookstore.payment.application.usecase

import com.codejsha.bookstore.payment.domain.aggregate.entity.PaymentEntity
import com.codejsha.bookstore.payment.domain.model.FilterCondition
import com.codejsha.bookstore.payment.domain.model.PaymentDto
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType
import java.math.BigDecimal
import java.time.LocalDateTime

sealed class PaymentCommand {

    data class CreatePayment(
        val orderId: Long,
        val userId: String,
        val paymentType: PaymentType,
        val cardNumber: String,
        val amount: BigDecimal
    ) : PaymentCommand() {

        constructor(dto: PaymentDto) : this(
            orderId = dto.orderId!!,
            userId = dto.userId!!,
            paymentType = dto.paymentType!!,
            cardNumber = dto.cardNumber!!,
            amount = dto.amount!!,
        )

        fun newEntity(): PaymentEntity {
            return PaymentEntity(
                id = null,
                orderId = this.orderId,
                userId = this.userId,
                paymentType = this.paymentType,
                cardNumber = this.cardNumber,
                amount = this.amount,
                paymentDate = LocalDateTime.now()
            )
        }
    }

    data class UpdatePayment(
        val id: Long,
        val orderId: Long?,
        val userId: String?,
        val paymentType: PaymentType?,
        val cardNumber: String?,
        val amount: BigDecimal?
    ) : PaymentCommand() {

        constructor(dto: PaymentDto) : this(
            id = dto.id!!,
            orderId = dto.orderId,
            userId = dto.userId,
            paymentType = dto.paymentType,
            cardNumber = dto.cardNumber,
            amount = dto.amount
        )
    }

    data class DeletePayment(val id: Long) : PaymentCommand()
}

sealed class PaymentQuery {
    data class FindAllPayments(val cond: FilterCondition) : PaymentQuery()
    data class FindPayment(val id: Long) : PaymentQuery()
}
