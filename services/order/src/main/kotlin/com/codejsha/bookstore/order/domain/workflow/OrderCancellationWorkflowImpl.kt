package com.codejsha.bookstore.order.domain.workflow

import com.codejsha.bookstore.order.domain.constant.OrderStatus
import io.temporal.activity.ActivityOptions
import io.temporal.common.RetryOptions
import io.temporal.failure.ActivityFailure
import io.temporal.failure.ApplicationFailure
import io.temporal.workflow.Workflow
import java.time.Duration

class OrderCancellationWorkflowImpl : OrderCancellationWorkflow {

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
            .setScheduleToCloseTimeout(Duration.ofMinutes(10))
            .setTaskQueue("inventory-task-queue")
            .setRetryOptions(RetryOptions.newBuilder().setMaximumAttempts(5).build())
            .build(),
    )

    private val paymentActivities = Workflow.newActivityStub(
        PaymentActivities::class.java,
        ActivityOptions.newBuilder()
            .setStartToCloseTimeout(Duration.ofSeconds(60))
            .setScheduleToCloseTimeout(Duration.ofMinutes(10))
            .setTaskQueue("payment-task-queue")
            .setRetryOptions(RetryOptions.newBuilder().setMaximumAttempts(5).build())
            .build(),
    )

    override fun cancelOrder(request: CancelOrderWorkflowRequest): CancelOrderWorkflowResult {
        val state = orderActivities.loadCancellationState(request.orderUid)
        val status = OrderStatus.fromValueOrNull(state.status)
            ?: throw ApplicationFailure.newNonRetryableFailure(
                "order ${request.orderUid} has unrecognized status ${state.status}",
                "UnknownOrderStatus",
            )
        if (status != OrderStatus.PENDING && status != OrderStatus.PAID) {
            throw ApplicationFailure.newNonRetryableFailure(
                "order ${request.orderUid} is not cancellable in status ${state.status}",
                "OrderNotCancellable",
            )
        }

        val paymentUid = state.paymentUid
        return if (paymentUid != null) {
            orderActivities.markOrderRefunded(request.orderUid)
            val refunded = try {
                paymentActivities.refundPayment(paymentUid)
            } catch (e: ActivityFailure) {
                orderActivities.restoreOrderPaid(request.orderUid)
                throw e
            }
            if (!refunded) {
                orderActivities.restoreOrderPaid(request.orderUid)
                throw ApplicationFailure.newNonRetryableFailure(
                    "payment $paymentUid for order ${request.orderUid} is not refundable; " +
                        "stock reservations left intact; manual payment reconciliation required",
                    "PaymentNotRefundable",
                )
            }
            releaseReservedStock(request.orderUid)
            CancelOrderWorkflowResult(orderUid = request.orderUid, status = OrderStatus.REFUNDED.value)
        } else {
            orderActivities.cancelOrder(request.orderUid)
            releaseReservedStock(request.orderUid)
            CancelOrderWorkflowResult(orderUid = request.orderUid, status = OrderStatus.CANCELLED.value)
        }
    }

    private fun releaseReservedStock(orderUid: String) {
        val items = orderActivities.loadOrderItems(orderUid)
        if (items.isNotEmpty()) {
            inventoryActivities.releaseStock(ReserveStockRequest(orderUid = orderUid, items = items))
        }
    }
}
