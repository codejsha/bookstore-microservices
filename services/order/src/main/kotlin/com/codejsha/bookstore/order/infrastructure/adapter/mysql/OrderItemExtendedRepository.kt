package com.codejsha.bookstore.order.infrastructure.adapter.mysql

import com.codejsha.bookstore.order.application.port.repo.OrderItemExtendedRepo
import com.codejsha.bookstore.order.domain.aggregate.entity.OrderItemEntity
import com.codejsha.bookstore.order.domain.model.FilterCondition
import org.springframework.r2dbc.core.DatabaseClient
import org.springframework.stereotype.Repository
import reactor.core.publisher.Flux
import java.time.LocalDateTime

@Repository
class OrderItemExtendedRepository(
    private val databaseClient: DatabaseClient
) : OrderItemExtendedRepo {

    override fun findAllPaged(cond: FilterCondition): Flux<OrderItemEntity> {
        val sql = buildQuery(cond)

        val entityFlux = databaseClient.sql(sql)
            .map { row, _ ->
                OrderItemEntity(
                    id = row["id"] as Long,
                    orderId = row["order_id"] as Long,
                    bookId = row["book_id"] as Long,
                    quantity = row["quantity"] as Int,
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

        val sql = """
            SELECT * FROM book_order_item
            ORDER BY $sort $order
            LIMIT ${cond.filter.limit} OFFSET ${cond.filter.offset}
        """.trimIndent()
        return sql
    }

    override fun findAllByOrderId(orderId: Long): Flux<OrderItemEntity> {
        val sql = """
            SELECT * FROM book_order_item
            WHERE order_id = :orderId
        """.trimIndent()

        val entityFlux = databaseClient.sql(sql)
            .bind("orderId", orderId)
            .map { row, _ ->
                OrderItemEntity(
                    id = row["id"] as Long,
                    orderId = row["order_id"] as Long,
                    bookId = row["book_id"] as Long,
                    quantity = row["quantity"] as Int,
                    createdAt = row["created_at"] as LocalDateTime?,
                    updatedAt = row["updated_at"] as LocalDateTime?
                )
            }
            .all()
        return entityFlux
    }
}
