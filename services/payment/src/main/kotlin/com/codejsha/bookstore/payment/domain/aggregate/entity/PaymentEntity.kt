package com.codejsha.bookstore.payment.domain.aggregate.entity

import com.codejsha.bookstore.service.application.port.openapi.model.PaymentFindWebResp
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType
import com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoResp
import com.google.protobuf.Timestamp
import jakarta.persistence.*
import org.springframework.data.annotation.CreatedDate
import org.springframework.data.annotation.LastModifiedDate
import java.time.LocalDateTime
import java.time.ZoneId

@Entity
@Table(name = "payment")
data class PaymentEntity(
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    var id: Long,
    @Column var orderId: Long,
    @Column var userId: String,
    @Column var paymentType: PaymentType,
    @Column var cardNumber: String,
    @Column var amount: Double,
    @Column var paymentDate: LocalDateTime,
    @Column @CreatedDate var createdAt: LocalDateTime?,
    @Column @LastModifiedDate var updatedAt: LocalDateTime?
) {

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
            .setId(id)
            .setOrderId(orderId)
            .setUserId(userId)
            .setPaymentType(paymentType.value)
            .setCardNumber(cardNumber)
            .setAmount(amount)
            .setPaymentDate(timestamp)
            .build()
        return response
    }
}
