package com.codejsha.bookstore.order.domain.model

import com.codejsha.bookstore.order.application.port.saga.PlaceOrderPayload
import com.codejsha.bookstore.service.application.port.openapi.model.OrderCreateWebReq
import com.codejsha.bookstore.service.application.port.openapi.model.OrderFindAllWebParam
import com.codejsha.bookstore.service.application.port.openapi.model.OrderUpdateWebReq
import com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq

sealed class OrderMapper {
    companion object {
        fun toOrderDto(req: OrderCreateWebReq): OrderDto = req.toOrderDto()
        fun toOrderDto(id: Long, req: OrderUpdateWebReq): OrderDto = req.toOrderDto(id)
        fun toOrderDto(payload: PlaceOrderPayload): OrderDto = payload.toOrderDto()

        fun toFilterCondition(param: OrderFindAllWebParam): FilterCondition = param.toFilterCondition()
        fun toFilterCondition(req: OrderFindAllProtoReq): FilterCondition = req.toFilterCondition()
    }
}

private fun OrderCreateWebReq.toOrderDto(): OrderDto {
    return OrderDto(
        id = null,
        userId = this.userId,
        totalPrice = this.totalPrice,
        status = this.status,
        orderItems = this.orderItems?.map {
            OrderItemDto(
                bookId = requireNotNull(it.bookId),
                quantity = requireNotNull(it.quantity)
            )
        },
        paymentType = this.paymentType,
        cardNumber = this.cardNumber
    )
}

private fun OrderUpdateWebReq.toOrderDto(id: Long): OrderDto {
    return OrderDto(
        id = id,
        userId = this.userId,
        totalPrice = this.totalPrice,
        status = this.status,
        orderItems = this.orderItems?.map {
            OrderItemDto(
                bookId = requireNotNull(it.bookId),
                quantity = requireNotNull(it.quantity)
            )
        },
        paymentType = this.paymentType,
        cardNumber = this.cardNumber
    )
}

private fun PlaceOrderPayload.toOrderDto(): OrderDto {
    return OrderDto(
        id = null,
        userId = this.userId,
        totalPrice = this.totalPrice,
        status = this.status,
        orderItems = this.orderItems.map {
            OrderItemDto(
                bookId = requireNotNull(it.bookId),
                quantity = requireNotNull(it.quantity)
            )
        },
        paymentType = this.paymentType,
        cardNumber = this.cardNumber
    )
}

private fun OrderFindAllWebParam.toFilterCondition(): FilterCondition {
    return FilterCondition(
        userId = this.userId,
        filter = FilterCondition.QueryFilter(
            sort = this.sort,
            order = this.order,
            limit = this.limit,
            offset = this.offset
        )
    )
}

private fun OrderFindAllProtoReq.toFilterCondition(): FilterCondition {
    return FilterCondition(
        userId = this.userId,
        filter = FilterCondition.QueryFilter(
            sort = this.filter.sort,
            order = this.filter.order,
            limit = this.filter.limit,
            offset = this.filter.offset
        )
    )
}
