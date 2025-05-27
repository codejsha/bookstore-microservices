package com.codejsha.bookstore.payment.domain.aggregate.entity

import com.codejsha.bookstore.payment.application.usecase.PaymentCommand
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentFindWebResp
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType
import com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoResp
import com.google.common.base.Objects
import com.google.protobuf.Timestamp
import jakarta.persistence.*
import org.hibernate.Hibernate
import org.springframework.data.annotation.CreatedDate
import org.springframework.data.annotation.LastModifiedDate
import org.springframework.data.jpa.domain.support.AuditingEntityListener
import java.math.BigDecimal
import java.time.LocalDateTime
import java.time.ZoneId

@Entity
@Table(name = "payment")
@EntityListeners(AuditingEntityListener::class)
data class PaymentEntity(
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    var id: Long? = null,
    var orderId: Long,
    var userId: String,
    @Enumerated(EnumType.STRING) var paymentType: PaymentType,
    var cardNumber: String,
    var amount: BigDecimal,
    var paymentDate: LocalDateTime,
    @CreatedDate var createdAt: LocalDateTime? = null,
    @LastModifiedDate var updatedAt: LocalDateTime? = null,
    @Version var version: Long? = null
) {

    fun update(command: PaymentCommand.UpdatePayment) =
        apply {
            orderId = command.orderId ?: orderId
            userId = command.userId ?: userId
            paymentType = command.paymentType ?: paymentType
            cardNumber = command.cardNumber ?: cardNumber
            amount = command.amount ?: amount
        }

    fun toMap(): Map<String, Any?> {
        val map = mapOf(
            "id" to id,
            "orderId" to orderId,
            "userId" to userId,
            "paymentType" to paymentType,
            "cardNumber" to cardNumber,
            "amount" to amount,
            "paymentDate" to paymentDate,
            "createdAt" to createdAt,
            "updatedAt" to updatedAt
        )
        return map
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other == null || Hibernate.getClass(this) != Hibernate.getClass(other)) return false
        other as PaymentEntity
        return id == other.id &&
                orderId == other.orderId &&
                userId == other.userId
    }

    override fun hashCode(): Int {
        return Objects.hashCode(id, orderId, userId)
    }

    fun toPaymentFindWebResp(): PaymentFindWebResp {
        val zoneOffset = ZoneId.systemDefault().rules.getOffset(paymentDate)
        val odt = paymentDate.atOffset(zoneOffset)
        val response = PaymentFindWebResp(
            id = id,
            orderId = orderId,
            userId = userId,
            paymentType = paymentType,
            cardNumber = cardNumber,
            amount = amount,
            paymentDate = odt
        )
        return response
    }

    fun toPaymentFindProtoResp(): PaymentFindProtoResp {
        val zoneOffset = ZoneId.systemDefault().rules.getOffset(paymentDate)
        val instant = paymentDate.toInstant(zoneOffset)
        val timestamp = Timestamp.newBuilder().setSeconds(instant.epochSecond).setNanos(instant.nano).build()
        val response = PaymentFindProtoResp.newBuilder()
            .setId(requireNotNull(id))
            .setOrderId(orderId)
            .setUserId(userId)
            .setPaymentType(paymentType.value)
            .setCardNumber(cardNumber)
            .setAmount(amount.toPlainString())
            .setPaymentDate(timestamp)
            .build()
        return response
    }
}
