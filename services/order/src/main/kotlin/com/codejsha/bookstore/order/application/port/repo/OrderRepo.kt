package com.codejsha.bookstore.order.application.port.repo

import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.order.domain.model.FilterCondition
import org.springframework.data.repository.reactive.ReactiveCrudRepository
import reactor.core.publisher.Flux

interface OrderRepo : ReactiveCrudRepository<OrderEntity, Long>, OrderExtendedRepo

interface OrderExtendedRepo {
    fun findAllPaged(cond: FilterCondition): Flux<OrderEntity>
}
