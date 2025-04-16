package com.codejsha.bookstore.order.application.usecase

import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import com.codejsha.bookstore.order.domain.model.FilterCondition

sealed class OrderCommand {
    data class CreateOrder(val ent: OrderEntity) : OrderCommand()
    data class UpdateOrder(val ent: OrderEntity) : OrderCommand()
    data class DeleteOrder(val id: Long) : OrderCommand()
}

sealed class OrderQuery {
    data class FindAllOrders(val cond: FilterCondition) : OrderQuery()
    data class FindOrder(val id: Long) : OrderQuery()
}

sealed class OrderItemCommand {
    data class CreateOrderItem(val ents: List<OrderItemEntity>) : OrderItemCommand()
    data class UpdateOrderItem(val ents: List<OrderItemEntity>) : OrderItemCommand()
    data class DeleteOrderItem(val id: Long) : OrderItemCommand()
}
