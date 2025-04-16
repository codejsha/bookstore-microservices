package com.codejsha.bookstore.order.domain.aggregate

import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import com.codejsha.bookstore.service.application.port.openapi.model.OrderFindWebResp
import com.codejsha.bookstore.service.application.port.openapi.model.OrderItem
import com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp
import com.codejsha.bookstore.service.application.port.pb.orderpb.OrderStatus
import com.codejsha.bookstore.service.application.port.pb.orderpb.OrderItem as OrderItemProto

data class OrderAggregate(
    // order id
    val id: Long,
    val order: OrderEntity,
    val orderItems: List<OrderItemEntity>
) {

    fun toMap(): Map<String, Any?> {
        val map = mapOf(
            "id" to id,
            "userId" to order.userId,
            "orderItems" to orderItems.map { it.toMap() },
            "totalPrice" to order.totalPrice,
            "status" to order.status
        )
        return map
    }

    fun toOrderFindWebResp(): OrderFindWebResp {
        val response = OrderFindWebResp(
            id = id,
            userId = order.userId,
            orderItems = orderItems.map {
                OrderItem(
                    bookId = it.bookId,
                    quantity = it.quantity,
                )
            },
            totalPrice = order.totalPrice,
            status = order.status,
        )
        return response
    }

    fun toOrderFindProtoResp(): OrderFindProtoResp {
        val response = OrderFindProtoResp.newBuilder()
            .setId(id)
            .setUserId(order.userId)
            .addAllOrderItems(orderItems.map {
                OrderItemProto.newBuilder()
                    .setBookId(it.bookId)
                    .setQuantity(it.quantity)
                    .build()
            })
            .setTotalPrice(order.totalPrice)
            .setStatus(OrderStatus.valueOf(order.status.name))
            .build()
        return response
    }
}
