package com.codejsha.bookstore.order.domain.aggregate.entity

import org.springframework.data.annotation.CreatedDate
import org.springframework.data.annotation.Id
import org.springframework.data.annotation.LastModifiedDate
import org.springframework.data.relational.core.mapping.Column
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime

@Table(name = "book_order_item")
data class OrderItemEntity(
    @Id var id: Long?,
    @Column var orderId: Long?,
    @Column var bookId: Long,
    @Column var quantity: Int,
    @Column @CreatedDate var createdAt: LocalDateTime?,
    @Column @LastModifiedDate var updatedAt: LocalDateTime?
) {

    fun toMap(): Map<String, Any?> {
        val map = mapOf(
            "bookId" to bookId,
            "quantity" to quantity,
        )
        return map
    }
}
