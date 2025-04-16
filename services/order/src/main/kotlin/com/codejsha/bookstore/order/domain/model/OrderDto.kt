package com.codejsha.bookstore.order.domain.model

import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus

data class OrderDto(
    val id: Long?,
    val userId: String?,
    val orderItems: List<OrderItemDto>?,
    val totalPrice: Double?,
    val status: OrderStatus?
) {

    fun toEntity() = OrderEntity(
        id = id!!,
        userId = userId!!,
        totalPrice = totalPrice!!,
        status = status!!,
        createdAt = null,
        updatedAt = null
    )
}
