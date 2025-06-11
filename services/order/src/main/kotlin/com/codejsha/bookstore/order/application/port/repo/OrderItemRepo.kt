package com.codejsha.bookstore.order.application.port.repo

import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import com.codejsha.bookstore.order.domain.model.FilterCondition
import org.springframework.data.repository.reactive.ReactiveCrudRepository
import reactor.core.publisher.Flux

interface OrderItemRepo : ReactiveCrudRepository<OrderItemEntity, Long>, OrderItemExtendedRepo

interface OrderItemExtendedRepo {
    fun findAllPaged(cond: FilterCondition): Flux<OrderItemEntity>
    fun findAllByOrderId(orderId: Long): Flux<OrderItemEntity>
}
