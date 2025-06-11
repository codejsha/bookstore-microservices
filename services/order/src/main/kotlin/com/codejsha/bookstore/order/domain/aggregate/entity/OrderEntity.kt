package com.codejsha.bookstore.order.domain.aggregate.entity

import com.codejsha.bookstore.order.application.usecase.OrderCommand
import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus
import com.google.common.base.Objects
import org.springframework.data.annotation.CreatedDate
import org.springframework.data.annotation.Id
import org.springframework.data.annotation.LastModifiedDate
import org.springframework.data.relational.core.mapping.Table
import java.math.BigDecimal
import java.time.LocalDateTime

@Table(name = "book_order")
data class OrderEntity(
    @Id var id: Long? = null,
    var userId: String,
    var totalPrice: BigDecimal,
    var status: OrderStatus,
    @CreatedDate var createdAt: LocalDateTime? = null,
    @LastModifiedDate var updatedAt: LocalDateTime? = null,
) {

    fun update(command: OrderCommand.UpdateOrderCommand) =
        apply {
            userId = command.userId ?: userId
            totalPrice = command.totalPrice ?: totalPrice
            status = command.status ?: status
        }

    fun update(command: OrderCommand.ChangeOrderStatusCommand): OrderEntity {
        this.status = status
        return this
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is OrderEntity) return false

        if (id != other.id) return false
        if (userId != other.userId) return false
        if (totalPrice != other.totalPrice) return false
        if (status != other.status) return false

        return true
    }

    override fun hashCode(): Int {
        return Objects.hashCode(id, userId, totalPrice, status)
    }
}
