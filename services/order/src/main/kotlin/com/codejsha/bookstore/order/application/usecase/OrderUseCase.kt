package com.codejsha.bookstore.order.application.usecase

import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.aggregate.OrderAdjustmentEntity
import com.codejsha.bookstore.order.domain.aggregate.OrderItemEntity
import com.codejsha.bookstore.order.domain.aggregate.OrderShippingEntity
import com.codejsha.bookstore.order.domain.model.command.*
import com.codejsha.bookstore.order.domain.model.option.OrderQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import java.util.UUID

interface OrderUseCase {
    // ─── Query ──────────────────────────────────────────────────────────────
    fun findAllOrders(option: OrderQueryOption, pageable: Pageable, context: ActorContext): Page<OrderAggregate>
    fun findOrder(uid: UUID, context: ActorContext): OrderAggregate

    // ─── Order lifecycle ────────────────────────────────────────────────────
    fun placeOrder(command: OrderCreateCommand, items: List<OrderItemCreateCommand>, shipping: OrderShippingCreateCommand?, context: ActorContext): OrderAggregate
    fun cancelOrder(uid: UUID, context: ActorContext): OrderAggregate

    // ─── Item management (PENDING only) ─────────────────────────────────────
    fun addItem(orderUid: UUID, command: OrderItemCreateCommand, context: ActorContext): OrderItemEntity
    fun updateItem(orderUid: UUID, itemUid: UUID, command: OrderItemUpdateCommand, context: ActorContext): OrderItemEntity
    fun removeItem(orderUid: UUID, itemUid: UUID, context: ActorContext)

    // ─── Shipping ───────────────────────────────────────────────────────────
    fun setShipping(orderUid: UUID, command: OrderShippingCreateCommand, context: ActorContext): OrderShippingEntity
    fun updateShipping(orderUid: UUID, command: OrderShippingUpdateCommand, context: ActorContext): OrderShippingEntity

    // ─── Adjustments (coupon, point, tax, etc.) ─────────────────────────────
    fun applyAdjustment(orderUid: UUID, command: OrderAdjustmentCreateCommand, context: ActorContext): OrderAdjustmentEntity
    fun removeAdjustment(orderUid: UUID, adjustmentUid: UUID, context: ActorContext)
}
