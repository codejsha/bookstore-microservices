package com.codejsha.bookstore.order.domain.model

import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity

data class OrderItemDto(
    val id: Long?,
    val orderId: Long?,
    val bookId: Long,
    val quantity: Int,
) {

    fun toEntity() = OrderItemEntity(
        id = id!!,
        orderId = orderId!!,
        bookId = bookId,
        quantity = quantity,
        createdAt = null,
        updatedAt = null
    )
}
