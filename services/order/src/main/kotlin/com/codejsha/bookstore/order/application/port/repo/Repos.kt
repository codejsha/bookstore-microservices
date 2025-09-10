package com.codejsha.bookstore.order.application.port.repo

import com.codejsha.bookstore.order.domain.model.command.CartAddItemCommand
import com.codejsha.bookstore.order.domain.model.command.OrderAdjustmentCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemUpdateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderShippingCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderShippingUpdateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderUpdateCommand
import com.codejsha.bookstore.order.domain.model.option.OrderQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import java.util.UUID

interface OrderRepo {
    fun findAll(option: OrderQueryOption, pageable: Pageable, context: ActorContext): Page<OrderResult>
    fun findOne(uid: UUID, context: ActorContext): OrderResult
    fun findByIdempotencyKey(idempotencyKey: String, context: ActorContext): OrderResult?
    fun create(command: OrderCreateCommand, context: ActorContext): OrderResult
    fun update(uid: UUID, command: OrderUpdateCommand, context: ActorContext): OrderResult
    fun transitionStatus(uid: UUID, from: Set<String>, to: String, context: ActorContext): Boolean
    fun setPaymentUid(uid: UUID, paymentUid: UUID, context: ActorContext)
    fun delete(uid: UUID, context: ActorContext)
}

interface OrderItemRepo {
    fun findAllByOrder(orderUid: UUID, pageable: Pageable, context: ActorContext): Page<OrderItemResult>
    fun findOne(orderUid: UUID, uid: UUID, context: ActorContext): OrderItemResult
    fun create(orderUid: UUID, command: OrderItemCreateCommand, context: ActorContext): OrderItemResult
    fun update(orderUid: UUID, uid: UUID, command: OrderItemUpdateCommand, context: ActorContext): OrderItemResult
    fun delete(orderUid: UUID, uid: UUID, context: ActorContext)
}

interface OrderAdjustmentRepo {
    fun findAllByOrder(orderUid: UUID, pageable: Pageable, context: ActorContext): Page<OrderAdjustmentResult>
    fun findOne(orderUid: UUID, uid: UUID, context: ActorContext): OrderAdjustmentResult
    fun create(orderUid: UUID, command: OrderAdjustmentCreateCommand, context: ActorContext): OrderAdjustmentResult
    fun delete(orderUid: UUID, uid: UUID, context: ActorContext)
}

interface OrderShippingRepo {
    fun findByOrder(orderUid: UUID, context: ActorContext): OrderShippingResult?
    fun create(orderUid: UUID, command: OrderShippingCreateCommand, context: ActorContext): OrderShippingResult
    fun update(orderUid: UUID, command: OrderShippingUpdateCommand, context: ActorContext): OrderShippingResult
    fun delete(orderUid: UUID, context: ActorContext)
}

interface CartRepo {
    fun findByUser(userUid: UUID, context: ActorContext): CartResult?
    fun createForUser(userUid: UUID, context: ActorContext): CartResult
    fun delete(userUid: UUID, context: ActorContext)
}

interface CartItemRepo {
    fun findAllByCart(cartId: Long, context: ActorContext): List<CartItemResult>
    fun findOne(cartId: Long, uid: UUID, context: ActorContext): CartItemResult
    fun findByProduct(cartId: Long, productId: Long, context: ActorContext): CartItemResult?
    fun create(cartId: Long, command: CartAddItemCommand, context: ActorContext): CartItemResult
    fun updateQuantity(cartId: Long, uid: UUID, quantity: Int, context: ActorContext): CartItemResult
    fun delete(cartId: Long, uid: UUID, context: ActorContext)
    fun deleteAllByCart(cartId: Long, context: ActorContext)
}
