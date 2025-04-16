package com.codejsha.bookstore.order.domain.aggregate.entity

import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus
import org.springframework.data.annotation.CreatedDate
import org.springframework.data.annotation.Id
import org.springframework.data.annotation.LastModifiedDate
import org.springframework.data.relational.core.mapping.Column
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime

@Table(name = "book_order")
data class OrderEntity(
    @Id var id: Long?,
    @Column var userId: String,
    @Column var totalPrice: Double,
    @Column var status: OrderStatus,
    @Column @CreatedDate var createdAt: LocalDateTime?,
    @Column @LastModifiedDate var updatedAt: LocalDateTime?
) {
    fun toMap(): Map<String, Any?> {
        val map = mapOf(
            "id" to id,
            "userId" to userId,
            "totalPrice" to totalPrice,
            "status" to status,
        )
        return map
    }

    fun changeStatus(status: OrderStatus): OrderEntity {
        this.status = status
        return this
    }
}
