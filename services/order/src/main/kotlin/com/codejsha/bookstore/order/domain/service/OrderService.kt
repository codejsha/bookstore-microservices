package com.codejsha.bookstore.order.domain.service

import com.codejsha.bookstore.order.application.port.repo.*
import com.codejsha.bookstore.order.application.port.support.TransactionRunner
import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.aggregate.*
import com.codejsha.bookstore.order.domain.constant.OrderStatus
import com.codejsha.bookstore.order.domain.model.OrderStateConflictException
import com.codejsha.bookstore.order.domain.model.command.*
import com.codejsha.bookstore.order.domain.model.option.OrderQueryOption
import com.codejsha.platform.shared.data.ActorContext
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service
import java.util.UUID

@Service
class OrderService(
    private val orderRepo: OrderRepo,
    private val orderItemRepo: OrderItemRepo,
    private val orderAdjustmentRepo: OrderAdjustmentRepo,
    private val orderShippingRepo: OrderShippingRepo,
    private val txRunner: TransactionRunner,
) : OrderUseCase {

    // ─── Query ──────────────────────────────────────────────────────────────

    @WithSpan
    override fun findAllOrders(
        option: OrderQueryOption, pageable: Pageable, context: ActorContext
    ): Page<OrderAggregate> = txRunner.tx {
        orderRepo.findAll(option, pageable, context).map { it.toAggregate() }
    }

    @WithSpan
    override fun findOrder(uid: UUID, context: ActorContext): OrderAggregate = txRunner.tx {
        val order = orderRepo.findOne(uid, context).toAggregate()
        val items = orderItemRepo.findAllByOrder(uid, Pageable.unpaged(), context)
            .content.map { it.toEntity() }
        val adjustments = orderAdjustmentRepo.findAllByOrder(uid, Pageable.unpaged(), context)
            .content.map { it.toEntity() }
        val shipping = orderShippingRepo.findByOrder(uid, context)?.toEntity()

        order.copy(items = items, adjustments = adjustments, shipping = shipping)
    }

    // ─── Order lifecycle ────────────────────────────────────────────────────

    @WithSpan
    override fun placeOrder(
        command: OrderCreateCommand,
        items: List<OrderItemCreateCommand>,
        shipping: OrderShippingCreateCommand?,
        context: ActorContext,
    ): OrderAggregate = txRunner.tx {
        val computedTotal = command.itemsAmount - command.discountAmount + command.shippingAmount + command.taxAmount
        if (command.totalAmount.compareTo(computedTotal) != 0) {
            throw InvalidCommandException(
                "total_amount must equal items - discount + shipping + tax, was ${command.totalAmount}, expected $computedTotal",
            )
        }

        orderRepo.findByIdempotencyKey(command.idempotencyKey, context)?.let { existing ->
            val existingItems = orderItemRepo.findAllByOrder(existing.uid, Pageable.unpaged(), context)
                .content.map { it.toEntity() }
            val existingShipping = orderShippingRepo.findByOrder(existing.uid, context)?.toEntity()
            return@tx existing.toAggregate().copy(items = existingItems, shipping = existingShipping)
        }

        val orderResult = orderRepo.create(command, context)
        val orderUid = orderResult.uid

        val itemEntities = items.map { itemCmd ->
            orderItemRepo.create(orderUid, itemCmd, context).toEntity()
        }

        val shippingEntity = shipping?.let {
            orderShippingRepo.create(orderUid, it, context).toEntity()
        }

        orderResult.toAggregate().copy(
            items = itemEntities,
            shipping = shippingEntity,
        )
    }

    @WithSpan
    override fun cancelOrder(uid: UUID, context: ActorContext): OrderAggregate = txRunner.tx {
        val order = orderRepo.findOne(uid, context)
        requirePending(order.status)

        val updated = orderRepo.update(
            uid,
            OrderUpdateCommand(
                userUid = null,
                status = OrderStatus.CANCELLED.value,
                currency = null,
                itemsAmount = null,
                discountAmount = null,
                shippingAmount = null,
                taxAmount = null,
                totalAmount = null,
            ),
            context,
        )

        updated.toAggregate()
    }

    // ─── Item management ────────────────────────────────────────────────────

    @WithSpan
    override fun addItem(
        orderUid: UUID, command: OrderItemCreateCommand, context: ActorContext
    ): OrderItemEntity = txRunner.tx {
        requirePending(orderRepo.findOne(orderUid, context).status)
        orderItemRepo.create(orderUid, command, context).toEntity()
    }

    @WithSpan
    override fun updateItem(
        orderUid: UUID, itemUid: UUID, command: OrderItemUpdateCommand, context: ActorContext
    ): OrderItemEntity = txRunner.tx {
        requirePending(orderRepo.findOne(orderUid, context).status)
        orderItemRepo.update(orderUid, itemUid, command, context).toEntity()
    }

    @WithSpan
    override fun removeItem(orderUid: UUID, itemUid: UUID, context: ActorContext) {
        txRunner.tx {
            requirePending(orderRepo.findOne(orderUid, context).status)
            orderItemRepo.delete(orderUid, itemUid, context)
        }
    }

    // ─── Shipping ───────────────────────────────────────────────────────────

    @WithSpan
    override fun setShipping(
        orderUid: UUID, command: OrderShippingCreateCommand, context: ActorContext
    ): OrderShippingEntity = txRunner.tx {
        requirePending(orderRepo.findOne(orderUid, context).status)
        val existing = orderShippingRepo.findByOrder(orderUid, context)
        if (existing != null) {
            orderShippingRepo.update(
                orderUid,
                OrderShippingUpdateCommand(
                    recipientName = command.recipientName,
                    recipientPhone = command.recipientPhone,
                    addressLine1 = command.addressLine1,
                    addressLine2 = command.addressLine2,
                    city = command.city,
                    state = command.state,
                    postalCode = command.postalCode,
                    country = command.country,
                    shippingMethod = command.shippingMethod,
                ),
                context,
            ).toEntity()
        } else {
            orderShippingRepo.create(orderUid, command, context).toEntity()
        }
    }

    @WithSpan
    override fun updateShipping(
        orderUid: UUID, command: OrderShippingUpdateCommand, context: ActorContext
    ): OrderShippingEntity = txRunner.tx {
        requirePending(orderRepo.findOne(orderUid, context).status)
        orderShippingRepo.update(orderUid, command, context).toEntity()
    }

    // ─── Adjustments ────────────────────────────────────────────────────────

    @WithSpan
    override fun applyAdjustment(
        orderUid: UUID, command: OrderAdjustmentCreateCommand, context: ActorContext
    ): OrderAdjustmentEntity = txRunner.tx {
        requirePending(orderRepo.findOne(orderUid, context).status)
        orderAdjustmentRepo.create(orderUid, command, context).toEntity()
    }

    @WithSpan
    override fun removeAdjustment(orderUid: UUID, adjustmentUid: UUID, context: ActorContext) {
        txRunner.tx {
            requirePending(orderRepo.findOne(orderUid, context).status)
            orderAdjustmentRepo.delete(orderUid, adjustmentUid, context)
        }
    }

    // ─── Business rules ─────────────────────────────────────────────────────

    private fun requirePending(status: String) {
        if (OrderStatus.fromValue(status) != OrderStatus.PENDING) {
            throw OrderStateConflictException(
                "Operation allowed only when order status is PENDING, current: $status",
            )
        }
    }
}
