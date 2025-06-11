package com.codejsha.bookstore.order.domain.aggregate.entity

import com.google.common.base.Objects
import org.springframework.data.annotation.CreatedDate
import org.springframework.data.annotation.Id
import org.springframework.data.annotation.LastModifiedDate
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime

@Table(name = "book_order_item")
data class OrderItemEntity(
    @Id var id: Long?,
    var orderId: Long,
    var bookId: Long,
    var quantity: Int,
    @CreatedDate var createdAt: LocalDateTime? = null,
    @LastModifiedDate var updatedAt: LocalDateTime? = null,
) {

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is OrderItemEntity) return false

        if (id != other.id) return false
        if (orderId != other.orderId) return false
        if (bookId != other.bookId) return false
        if (quantity != other.quantity) return false

        return true
    }

    override fun hashCode(): Int {
        return Objects.hashCode(id, orderId, bookId, quantity)
    }
}
