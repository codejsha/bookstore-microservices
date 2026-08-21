package com.codejsha.bookstore.order.application.port.saga

import com.codejsha.bookstore.order.domain.model.OrderDto
import com.codejsha.bookstore.order.domain.model.OrderItemDto
import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType

import java.math.BigDecimal
import java.util.*

class PlaceOrderPayload(
    val requestId: String,
    val userId: String,
    val totalPrice: BigDecimal,
    val status: OrderStatus,
    val orderItems: List<OrderItemDto>,
    val paymentType: PaymentType,
    val cardNumber: String
) {
    companion object {
        fun create(dto: OrderDto): PlaceOrderPayload =
            PlaceOrderPayload(
                requestId = UUID.randomUUID().toString(),
                userId = requireNotNull(dto.userId),
                totalPrice = requireNotNull(dto.totalPrice),
                status = requireNotNull(dto.status),
                orderItems = requireNotNull(dto.orderItems),
                paymentType = requireNotNull(dto.paymentType),
                cardNumber = requireNotNull(dto.cardNumber)
            )
    }

    fun toWorkflowPayload(): (Map<String, Any>) =
        mapOf(
            "requestId" to requestId,
            "userId" to userId,
            "totalPrice" to totalPrice,
            "status" to status,
            "orderItems" to orderItems,
            "paymentType" to paymentType,
            "cardNumber" to cardNumber
        )
}
