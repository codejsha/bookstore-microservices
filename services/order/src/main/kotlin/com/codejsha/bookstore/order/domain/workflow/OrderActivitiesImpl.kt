package com.codejsha.bookstore.order.domain.workflow

import com.codejsha.bookstore.order.application.port.repo.OrderItemRepo
import com.codejsha.bookstore.order.application.port.repo.OrderRepo
import com.codejsha.bookstore.order.application.port.repo.OrderShippingRepo
import com.codejsha.bookstore.order.application.port.support.TransactionRunner
import com.codejsha.bookstore.order.domain.constant.OrderStatus
import com.codejsha.bookstore.order.domain.model.command.OrderCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderShippingCreateCommand
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import io.temporal.failure.ApplicationFailure
import kotlinx.coroutines.runBlocking
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Component
import java.math.BigDecimal
import java.util.UUID

@Component
class OrderActivitiesImpl(
    private val orderRepo: OrderRepo,
    private val orderItemRepo: OrderItemRepo,
    private val orderShippingRepo: OrderShippingRepo,
    private val txRunner: TransactionRunner,
) : OrderActivities {

    override fun createOrder(request: CreateOrderRequest): CreateOrderResult = runBlocking {
        val context = ActorContext(actorId = 0L, ActorType.USER)
        txRunner.tx {
            orderRepo.findByIdempotencyKey(request.idempotencyKey, context)?.let { existing ->
                return@tx CreateOrderResult(
                    orderUid = existing.uid.toString(),
                    status = existing.status,
                    paymentUid = existing.paymentUid?.toString(),
                )
            }

            val command = OrderCreateCommand(
                userUid = UUID.fromString(request.userUid),
                currency = request.currency,
                itemsAmount = request.itemsAmount,
                discountAmount = BigDecimal.ZERO,
                shippingAmount = BigDecimal.ZERO,
                taxAmount = BigDecimal.ZERO,
                totalAmount = request.totalAmount,
                idempotencyKey = request.idempotencyKey,
            )
            val order = orderRepo.create(command, context)

            for (item in request.items) {
                orderItemRepo.create(
                    order.uid,
                    OrderItemCreateCommand(
                        productId = item.productId,
                        sku = null,
                        productName = item.productName,
                        options = null,
                        quantity = item.quantity,
                        currency = item.currency,
                        price = item.price,
                        taxRate = BigDecimal.ZERO,
                    ),
                    context,
                )
            }

            request.shipping?.let { shipping ->
                orderShippingRepo.create(
                    order.uid,
                    OrderShippingCreateCommand(
                        recipientName = shipping.recipientName,
                        recipientPhone = shipping.recipientPhone,
                        addressLine1 = shipping.addressLine1,
                        addressLine2 = shipping.addressLine2,
                        city = shipping.city,
                        state = shipping.state,
                        postalCode = shipping.postalCode,
                        country = shipping.country,
                        shippingMethod = shipping.shippingMethod,
                    ),
                    context,
                )
            }

            CreateOrderResult(orderUid = order.uid.toString(), status = order.status, paymentUid = null)
        }
    }

    override fun confirmOrder(orderUid: String) {
        transitionStatus(orderUid, from = setOf(OrderStatus.PENDING), to = OrderStatus.PAID)
    }

    override fun cancelOrder(orderUid: String) {
        transitionStatus(orderUid, from = setOf(OrderStatus.PENDING), to = OrderStatus.CANCELLED)
    }

    override fun recordPayment(orderUid: String, paymentUid: String) {
        val uid = UUID.fromString(orderUid)
        val context = ActorContext(actorId = 0L, ActorType.SYSTEM)
        orderRepo.setPaymentUid(uid, UUID.fromString(paymentUid), context)
    }

    override fun loadCancellationState(orderUid: String): OrderCancellationState {
        val uid = UUID.fromString(orderUid)
        val context = ActorContext(actorId = 0L, ActorType.SYSTEM)
        val order = orderRepo.findOne(uid, context)
        return OrderCancellationState(
            status = order.status,
            paymentUid = order.paymentUid?.toString(),
        )
    }

    override fun loadOrderItems(orderUid: String): List<StockReservationItem> {
        val uid = UUID.fromString(orderUid)
        val context = ActorContext(actorId = 0L, ActorType.SYSTEM)
        return orderItemRepo.findAllByOrder(uid, Pageable.unpaged(), context)
            .content
            .map { StockReservationItem(productId = it.productId, quantity = it.quantity) }
    }

    override fun markOrderRefunded(orderUid: String) {
        transitionStatus(
            orderUid,
            from = setOf(OrderStatus.PENDING, OrderStatus.PAID, OrderStatus.CANCELLED),
            to = OrderStatus.REFUNDED,
        )
    }

    override fun markOrderShipped(orderUid: String) {
        transitionStatus(orderUid, from = setOf(OrderStatus.PAID), to = OrderStatus.SHIPPED)
    }

    override fun restoreOrderPaid(orderUid: String) {
        transitionStatus(orderUid, from = setOf(OrderStatus.SHIPPED, OrderStatus.REFUNDED), to = OrderStatus.PAID)
    }

    private fun transitionStatus(orderUid: String, from: Set<OrderStatus>, to: OrderStatus) {
        val uid = UUID.fromString(orderUid)
        val context = ActorContext(actorId = 0L, ActorType.SYSTEM)
        val transitioned = orderRepo.transitionStatus(uid, from.map { it.value }.toSet(), to.value, context)
        if (!transitioned) {
            val current = orderRepo.findOne(uid, context).status
            if (current == to.value) return
            throw ApplicationFailure.newNonRetryableFailure(
                "order $orderUid cannot transition to ${to.value} from status $current",
                "OrderStateConflict",
            )
        }
    }
}
