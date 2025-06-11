package com.codejsha.bookstore.order.domain.model

import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType
import java.math.BigDecimal

data class OrderDto(
    val id: Long?,
    val userId: String?,
    val totalPrice: BigDecimal?,
    val status: OrderStatus?,
    val orderItems: List<OrderItemDto>?,
    val paymentType: PaymentType?,
    val cardNumber: String?
)
