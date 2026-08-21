package com.codejsha.bookstore.order.infrastructure.adapter.mysql

import com.codejsha.bookstore.order.application.port.repo.OrderExtendedRepo
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.order.domain.model.FilterCondition

import org.springframework.data.r2dbc.core.R2dbcEntityTemplate
import org.springframework.data.r2dbc.core.select
import org.springframework.data.relational.core.query.Criteria
import org.springframework.data.relational.core.query.Query
import org.springframework.stereotype.Repository
import reactor.core.publisher.Flux

@Repository
class OrderExtendedRepository(
    private val r2dbcEntityTemplate: R2dbcEntityTemplate
) : OrderExtendedRepo {
    override fun findAllPaged(cond: FilterCondition): Flux<OrderEntity> {
        val pageable = cond.filter.toPageable()
        val criteria = Criteria.empty()
        if (!cond.userId.isNullOrBlank()) {
            criteria.and("user_id").`is`(cond.userId)
        }

        val entityFlux =
            r2dbcEntityTemplate
                .select<OrderEntity>()
                .matching(Query.query(criteria).with(pageable))
                .all()
        return entityFlux
    }
}
