package com.codejsha.bookstore.order.application.usecase

import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import com.codejsha.bookstore.order.domain.model.FilterCondition
import com.codejsha.bookstore.order.domain.model.OrderDto
import com.codejsha.bookstore.order.domain.model.OrderItemDto
import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus
import java.math.BigDecimal

sealed class OrderCommand {
    data class CreateOrderCommand(
        val userId: String,
        val totalPrice: BigDecimal,
        val status: OrderStatus,
        val orderItems: List<OrderItemDto>
    ) : OrderCommand() {

        constructor(dto: OrderDto) : this(
            userId = requireNotNull(dto.userId),
            totalPrice = requireNotNull(dto.totalPrice),
            status = requireNotNull(dto.status),
            orderItems = requireNotNull(dto.orderItems)
        )

        fun newEntity(): OrderEntity {
            return OrderEntity(
                id = null,
                userId = this.userId,
                totalPrice = this.totalPrice,
                status = this.status
            )
        }
    }

    data class UpdateOrderCommand(
        val id: Long,
        val userId: String?,
        val totalPrice: BigDecimal?,
        val status: OrderStatus?
    ) : OrderCommand() {

        constructor(dto: OrderDto) : this(
            id = requireNotNull(dto.id),
            userId = dto.userId,
            totalPrice = dto.totalPrice,
            status = dto.status
        )
    }

    data class ChangeOrderStatusCommand(
        val id: Long,
        val status: OrderStatus
    ) : OrderCommand() {

        constructor(dto: OrderDto) : this(
            id = requireNotNull(dto.id),
            status = requireNotNull(dto.status)
        )
    }

    data class DeleteOrderCommand(val id: Long) : OrderCommand()
}

sealed class OrderRead {
    data class FindAllOrdersRead(val cond: FilterCondition) : OrderRead()
    data class FindOrderRead(val id: Long) : OrderRead()
}

sealed class OrderItemCommand {

    data class CreateOrderItemsCommand(
        val orderId: Long,
        val orderItems: List<OrderItemDto>
    ) : OrderItemCommand() {

        fun newEntities(): List<OrderItemEntity> {
            return orderItems.map { item ->
                OrderItemEntity(
                    id = null,
                    orderId = orderId,
                    bookId = item.bookId,
                    quantity = item.quantity
                )
            }
        }
    }

    data class UpdateOrderItemsCommand(
        val orderId: Long,
        val orderItems: List<OrderItemDto>
    ) : OrderItemCommand()

    data class DeleteOrderItemsCommand(val id: Long) : OrderItemCommand()
}
