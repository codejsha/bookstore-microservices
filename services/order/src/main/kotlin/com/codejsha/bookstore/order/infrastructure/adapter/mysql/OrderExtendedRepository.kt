package com.codejsha.bookstore.order.infrastructure.adapter.mysql

import com.codejsha.bookstore.order.application.port.repo.OrderExtendedRepo
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderEntity
import com.codejsha.bookstore.order.domain.model.FilterCondition
import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus
import org.springframework.r2dbc.core.DatabaseClient
import org.springframework.stereotype.Repository
import reactor.core.publisher.Flux
import java.math.BigDecimal
import java.time.LocalDateTime

@Repository
class OrderExtendedRepository(
    private val databaseClient: DatabaseClient
) : OrderExtendedRepo {

    override fun findAllPaged(cond: FilterCondition): Flux<OrderEntity> {
        val sql = buildQuery(cond)

        val entityFlux = databaseClient.sql(sql)
            .map { row, _ ->
                OrderEntity(
                    id = row["id"] as Long,
                    userId = row["user_id"] as String,
                    totalPrice = row["total_price"] as BigDecimal,
                    status = row["status"] as OrderStatus,
                    createdAt = row["created_at"] as LocalDateTime?,
                    updatedAt = row["updated_at"] as LocalDateTime?
                )
            }
            .all()
        return entityFlux
    }

    private fun buildQuery(cond: FilterCondition): String {
        val sort = cond.filter.sort.lowercase()
        val order = if (listOf("ASC", "DESC").contains(cond.filter.order.uppercase()))
            cond.filter.order.uppercase() else "DESC"
        val where = if (!cond.userId.isNullOrBlank()) "WHERE user_id = ${cond.userId}" else ""

        val sql = """
            SELECT * FROM book_order
            $where
            ORDER BY $sort $order
            LIMIT ${cond.filter.limit} OFFSET ${cond.filter.offset}
        """.trimIndent()
        return sql
    }
}
