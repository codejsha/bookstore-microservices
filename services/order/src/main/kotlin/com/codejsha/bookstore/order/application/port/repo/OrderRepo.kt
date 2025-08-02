package com.codejsha.bookstore.order.application.port.repo

import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.order.domain.model.FilterCondition

import org.springframework.data.r2dbc.repository.R2dbcRepository
import reactor.core.publisher.Flux

interface OrderRepo : R2dbcRepository<OrderEntity, Long>

interface OrderExtendedRepo {
    fun findAllPaged(cond: FilterCondition): Flux<OrderEntity>
}
