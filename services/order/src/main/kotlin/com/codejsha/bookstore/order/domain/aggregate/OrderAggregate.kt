package com.codejsha.bookstore.order.domain.aggregate

import com.codejsha.bookstore.order.application.usecase.OrderCommand
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import com.codejsha.bookstore.service.application.port.openapi.model.OrderFindWebResp
import com.codejsha.bookstore.service.application.port.openapi.model.OrderItem
import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus
import com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp
import java.math.BigDecimal
import com.codejsha.bookstore.service.application.port.pb.orderpb.OrderItem as OrderItemProto
import com.codejsha.bookstore.service.application.port.pb.orderpb.OrderStatus as OrderStatusProto

data class OrderAggregate(
    val id: Long,
    val userId: String,
    val totalPrice: BigDecimal,
    var status: OrderStatus,
    val orderItems: List<OrderItem>
) {

    companion object {
        fun create(id: Long, order: OrderEntity, orderItems: List<OrderItemEntity>): OrderAggregate {
            return OrderAggregate(
                id = id,
                userId = order.userId,
                totalPrice = order.totalPrice,
                status = order.status,
                orderItems = orderItems.map {
                    OrderItem(
                        bookId = requireNotNull(it.bookId),
                        quantity = requireNotNull(it.quantity)
                    )
                }
            )
        }
    }

    fun update(command: OrderCommand.ChangeOrderStatusCommand) =
        apply {
            status = command.status
        }

    fun toOrderFindWebResp(): OrderFindWebResp {
        val response = OrderFindWebResp(
            id = id,
            userId = userId,
            totalPrice = totalPrice,
            status = status,
            orderItems = orderItems.map {
                OrderItem(
                    bookId = it.bookId,
                    quantity = it.quantity,
                )
            }
        )
        return response
    }

    fun toOrderFindProtoResp(): OrderFindProtoResp {
        val response = OrderFindProtoResp.newBuilder()
            .setId(id)
            .setUserId(userId)
            .setTotalPrice(totalPrice.toPlainString())
            .setStatus(OrderStatusProto.valueOf(status.name))
            .addAllOrderItems(orderItems.map {
                OrderItemProto.newBuilder().apply {
                    setBookId(requireNotNull(bookId))
                    setQuantity(requireNotNull(quantity))
                }.build()
            })
            .build()
        return response
    }
}
