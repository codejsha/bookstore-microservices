package com.codejsha.bookstore.order.application.port.repo

import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import com.codejsha.bookstore.order.domain.model.FilterCondition

import org.springframework.data.r2dbc.repository.R2dbcRepository
import reactor.core.publisher.Flux

interface OrderItemRepo : R2dbcRepository<OrderItemEntity, Long>

interface OrderItemExtendedRepo {
    fun findAllPaged(cond: FilterCondition): Flux<OrderItemEntity>

    fun findAllByOrderId(orderId: Long): Flux<OrderItemEntity>
}
