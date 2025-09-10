package com.codejsha.bookstore.order.domain.workflow

import com.codejsha.bookstore.order.domain.constant.OrderStatus
import io.temporal.activity.ActivityOptions
import io.temporal.common.RetryOptions
import io.temporal.failure.ActivityFailure
import io.temporal.failure.ApplicationFailure
import io.temporal.workflow.Workflow
import java.time.Duration

class OrderPlacementWorkflowImpl : OrderPlacementWorkflow {

    private val log = Workflow.getLogger(OrderPlacementWorkflowImpl::class.java)

    private val orderActivities = Workflow.newActivityStub(
        OrderActivities::class.java,
        ActivityOptions.newBuilder()
            .setStartToCloseTimeout(Duration.ofSeconds(30))
            .setScheduleToCloseTimeout(Duration.ofMinutes(5))
            .setRetryOptions(RetryOptions.newBuilder().setMaximumAttempts(3).build())
            .build(),
    )

    private val inventoryActivities = Workflow.newActivityStub(
        InventoryActivities::class.java,
        ActivityOptions.newBuilder()
            .setStartToCloseTimeout(Duration.ofSeconds(30))
            .setScheduleToCloseTimeout(Duration.ofMinutes(5))
            .setTaskQueue("inventory-task-queue")
            .setRetryOptions(RetryOptions.newBuilder().setMaximumAttempts(3).build())
            .build(),
    )

    private val paymentActivities = Workflow.newActivityStub(
        PaymentActivities::class.java,
        ActivityOptions.newBuilder()
            .setStartToCloseTimeout(Duration.ofSeconds(60))
            .setScheduleToCloseTimeout(Duration.ofMinutes(10))
            .setTaskQueue("payment-task-queue")
            .setRetryOptions(
                RetryOptions.newBuilder()
                    .setInitialInterval(Duration.ofSeconds(2))
                    .setMaximumInterval(Duration.ofSeconds(30))
                    .setMaximumAttempts(PAYMENT_MAX_ATTEMPTS)
                    .build(),
            )
            .build(),
    )

    private val compensationInventoryActivities = Workflow.newActivityStub(
        InventoryActivities::class.java,
        ActivityOptions.newBuilder()
            .setStartToCloseTimeout(Duration.ofSeconds(30))
            .setScheduleToCloseTimeout(Duration.ofMinutes(15))
            .setTaskQueue("inventory-task-queue")
            .setRetryOptions(RetryOptions.newBuilder().setMaximumAttempts(COMPENSATION_MAX_ATTEMPTS).build())
            .build(),
    )

    private val compensationPaymentActivities = Workflow.newActivityStub(
        PaymentActivities::class.java,
        ActivityOptions.newBuilder()
            .setStartToCloseTimeout(Duration.ofSeconds(60))
            .setScheduleToCloseTimeout(Duration.ofMinutes(15))
            .setTaskQueue("payment-task-queue")
            .setRetryOptions(RetryOptions.newBuilder().setMaximumAttempts(COMPENSATION_MAX_ATTEMPTS).build())
            .build(),
    )

    override fun placeOrder(request: PlaceOrderWorkflowRequest): PlaceOrderWorkflowResult {
        val totalAmount = request.items.sumOf { it.price * it.quantity.toBigDecimal() }

        val idempotencyKey = Workflow.getInfo().workflowId

        val created = orderActivities.createOrder(
            CreateOrderRequest(
                userUid = request.userUid,
                currency = request.currency,
                itemsAmount = totalAmount,
                totalAmount = totalAmount,
                idempotencyKey = idempotencyKey,
                items = request.items,
                shipping = request.shipping,
            ),
        )
        val orderUid = created.orderUid
        when (OrderStatus.fromValueOrNull(created.status)) {
            OrderStatus.PENDING -> Unit
            OrderStatus.PAID -> return PlaceOrderWorkflowResult(
                orderUid = orderUid,
                paymentUid = created.paymentUid,
                status = OrderStatus.PAID.value,
            )
            else -> throw ApplicationFailure.newNonRetryableFailure(
                "order $orderUid for idempotency key $idempotencyKey is already finalized with status ${created.status}",
                "OrderAlreadyFinalized",
            )
        }

        try {
            inventoryActivities.reserveStock(
                ReserveStockRequest(
                    orderUid = orderUid,
                    items = request.items.map { StockReservationItem(it.productId, it.quantity) },
                ),
            )
        } catch (e: ActivityFailure) {
            compensate(orderUid, request, paymentUid = null, paymentAttempted = false)
            throw e
        }

        val paymentResult: ProcessPaymentResult
        try {
            paymentResult = paymentActivities.processPayment(
                ProcessPaymentRequest(
                    orderUid = orderUid,
                    userUid = request.userUid,
                    amount = totalAmount,
                    currency = request.currency,
                ),
            )
        } catch (e: ActivityFailure) {
            compensate(orderUid, request, paymentUid = null, paymentAttempted = true)
            throw e
        }

        if (!isPaymentSuccessful(paymentResult.status)) {
            compensate(orderUid, request, paymentUid = paymentResult.paymentUid, paymentAttempted = true)
            throw ApplicationFailure.newNonRetryableFailure(
                "payment did not succeed for order $orderUid: status=${paymentResult.status}",
                "PaymentNotSuccessful",
            )
        }

        try {
            orderActivities.recordPayment(orderUid, paymentResult.paymentUid)
            orderActivities.confirmOrder(orderUid)
        } catch (e: ActivityFailure) {
            compensate(orderUid, request, paymentUid = paymentResult.paymentUid, paymentAttempted = true)
            throw e
        }

        return PlaceOrderWorkflowResult(
            orderUid = orderUid,
            paymentUid = paymentResult.paymentUid,
            status = "PAID",
        )
    }

    private fun compensate(
        orderUid: String,
        request: PlaceOrderWorkflowRequest,
        paymentUid: String?,
        paymentAttempted: Boolean,
    ) {
        var refunded = false
        if (paymentAttempted) {
            try {
                refunded = if (paymentUid != null) {
                    compensationPaymentActivities.refundPayment(paymentUid)
                } else {
                    compensationPaymentActivities.refundPaymentByOrder(orderUid)
                }
            } catch (e: ActivityFailure) {
                log.error(
                    "compensation refund failed for order $orderUid (paymentUid=$paymentUid); " +
                        "manual payment reconciliation required",
                    e,
                )
            }
        }
        try {
            compensationInventoryActivities.releaseStock(
                ReserveStockRequest(
                    orderUid = orderUid,
                    items = request.items.map { StockReservationItem(it.productId, it.quantity) },
                ),
            )
        } catch (e: ActivityFailure) {
            log.error(
                "compensation stock release failed for order $orderUid; reservations may be leaked",
                e,
            )
        }
        try {
            if (refunded) {
                orderActivities.markOrderRefunded(orderUid)
            } else {
                orderActivities.cancelOrder(orderUid)
            }
        } catch (e: ActivityFailure) {
            log.error("compensation status transition failed for order $orderUid; order left in inconsistent state", e)
        }
    }

    private fun isPaymentSuccessful(status: String): Boolean =
        status.equals("succeeded", ignoreCase = true)

    companion object {
        private const val COMPENSATION_MAX_ATTEMPTS = 10
        private const val PAYMENT_MAX_ATTEMPTS = 10
    }
}
