package com.codejsha.bookstore.payment.application.usecase

import com.codejsha.bookstore.payment.domain.aggregate.entity.PaymentEntity
import com.codejsha.bookstore.payment.domain.model.FilterCondition
import com.codejsha.bookstore.payment.domain.model.PaymentDto
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType
import java.math.BigDecimal
import java.time.LocalDateTime

sealed class PaymentCommand {

    data class CreatePaymentCommand(
        val orderId: Long,
        val userId: String,
        val paymentType: PaymentType,
        val cardNumber: String,
        val amount: BigDecimal
    ) : PaymentCommand() {

        constructor(dto: PaymentDto) : this(
            orderId = requireNotNull(dto.orderId),
            userId = requireNotNull(dto.userId),
            paymentType = requireNotNull(dto.paymentType),
            cardNumber = requireNotNull(dto.cardNumber),
            amount = requireNotNull(dto.amount),
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

    data class UpdatePaymentCommand(
        val id: Long,
        val orderId: Long?,
        val userId: String?,
        val paymentType: PaymentType?,
        val cardNumber: String?,
        val amount: BigDecimal?
    ) : PaymentCommand() {

        constructor(dto: PaymentDto) : this(
            id = requireNotNull(dto.id),
            orderId = dto.orderId,
            userId = dto.userId,
            paymentType = dto.paymentType,
            cardNumber = dto.cardNumber,
            amount = dto.amount
        )
    }

    data class DeletePaymentCommand(val id: Long) : PaymentCommand()
}

sealed class PaymentRead {
    data class FindAllPaymentsRead(val cond: FilterCondition) : PaymentRead()
    data class FindPaymentRead(val id: Long) : PaymentRead()
}
