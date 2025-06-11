package com.codejsha.bookstore.order.application.usecase

import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.model.FilterCondition
import com.codejsha.bookstore.order.domain.model.OrderDto
import reactor.core.publisher.Mono

interface OrderUseCase {
    fun startPlaceOrderWorkflow(dto: OrderDto): Mono<Unit>
    fun findAllOrders(cond: FilterCondition): Mono<List<OrderAggregate>>
    fun findOrder(id: Long): Mono<OrderAggregate>
    fun createOrder(dto: OrderDto): Mono<OrderAggregate>
    fun updateOrder(dto: OrderDto): Mono<OrderAggregate>
    fun deleteOrder(id: Long): Mono<Void>
}
