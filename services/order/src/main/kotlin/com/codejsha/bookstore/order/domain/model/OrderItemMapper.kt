package com.codejsha.bookstore.order.domain.model

import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import com.codejsha.bookstore.service.application.port.openapi.model.OrderCreateWebReq
import com.codejsha.bookstore.service.application.port.openapi.model.OrderUpdateWebReq

sealed class OrderItemMapper {
    companion object {
        fun toOrderItemEntity(req: OrderCreateWebReq): List<OrderItemEntity> = req.toOrderItemEntity()
        fun toOrderItemEntity(id: Long, req: OrderUpdateWebReq): List<OrderItemEntity> = req.toOrderItemEntity(id)
    }
}

private fun OrderCreateWebReq.toOrderItemEntity(): List<OrderItemEntity> {
    return this.orderItems.orEmpty().map { orderItem ->
        OrderItemEntity(
            id = null,
            orderId = null,
            bookId = requireNotNull(orderItem.bookId) { "bookId cannot be null" },
            quantity = requireNotNull(orderItem.quantity) { "orderItem cannot be null" },
            createdAt = null,
            updatedAt = null
        )
    }
}

private fun OrderUpdateWebReq.toOrderItemEntity(id: Long): List<OrderItemEntity> {
    return this.orderItems.orEmpty().map { orderItem ->
        OrderItemEntity(
            id = null,
            orderId = id,
            bookId = requireNotNull(orderItem.bookId) { "bookId cannot be null" },
            quantity = requireNotNull(orderItem.quantity) { "orderItem cannot be null" },
            createdAt = null,
            updatedAt = null
        )
    }
}
