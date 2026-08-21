package com.codejsha.bookstore.order.domain.model

import com.codejsha.bookstore.service.application.port.openapi.model.OrderCreateWebReq
import com.codejsha.bookstore.service.application.port.openapi.model.OrderUpdateWebReq

sealed class OrderItemMapper {
    companion object {
        fun toOrderItemDto(req: OrderCreateWebReq): List<OrderItemDto> = req.toOrderItemDto()

        fun toOrderItemDto(
            id: Long,
            req: OrderUpdateWebReq
        ): List<OrderItemDto> = req.toOrderItemDto(id)
    }
}

private fun OrderCreateWebReq.toOrderItemDto(): List<OrderItemDto> =
    this.orderItems.orEmpty().map { orderItem ->
        OrderItemDto(
            bookId = requireNotNull(orderItem.bookId),
            quantity = requireNotNull(orderItem.quantity)
        )
    }

private fun OrderUpdateWebReq.toOrderItemDto(id: Long): List<OrderItemDto> =
    this.orderItems.orEmpty().map { orderItem ->
        OrderItemDto(
            bookId = requireNotNull(orderItem.bookId),
            quantity = requireNotNull(orderItem.quantity)
        )
    }
