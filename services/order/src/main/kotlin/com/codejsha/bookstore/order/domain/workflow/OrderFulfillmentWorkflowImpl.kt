package com.codejsha.bookstore.order.domain.workflow

import com.codejsha.bookstore.order.domain.constant.OrderStatus
import io.temporal.activity.ActivityOptions
import io.temporal.common.RetryOptions
import io.temporal.failure.ActivityFailure
import io.temporal.failure.ApplicationFailure
import io.temporal.workflow.Workflow
import java.time.Duration

class OrderFulfillmentWorkflowImpl : OrderFulfillmentWorkflow {

    private val orderActivities = Workflow.newActivityStub(
        OrderActivities::class.java,
        ActivityOptions.newBuilder()
            .setStartToCloseTimeout(Duration.ofSeconds(30))
            .setScheduleToCloseTimeout(Duration.ofMinutes(5))
            .setRetryOptions(RetryOptions.newBuilder().setMaximumAttempts(3).build())
            .build(),
    )

    private val deliveryActivities = Workflow.newActivityStub(
        DeliveryActivities::class.java,
        ActivityOptions.newBuilder()
            .setStartToCloseTimeout(Duration.ofSeconds(60))
            .setScheduleToCloseTimeout(Duration.ofMinutes(10))
            .setTaskQueue("delivery-task-queue")
            .setRetryOptions(RetryOptions.newBuilder().setMaximumAttempts(3).build())
            .build(),
    )

    private val notificationActivities = Workflow.newActivityStub(
        NotificationActivities::class.java,
        ActivityOptions.newBuilder()
            .setStartToCloseTimeout(Duration.ofSeconds(30))
            .setScheduleToCloseTimeout(Duration.ofMinutes(5))
            .setTaskQueue("notification-task-queue")
            .setRetryOptions(RetryOptions.newBuilder().setMaximumAttempts(5).build())
            .build(),
    )

    override fun fulfillOrder(request: FulfillOrderWorkflowRequest): FulfillOrderWorkflowResult {
        val destination = request.destination

        val state = orderActivities.loadCancellationState(request.orderUid)
        val status = OrderStatus.fromValueOrNull(state.status)
            ?: throw ApplicationFailure.newNonRetryableFailure(
                "order ${request.orderUid} has unrecognized status ${state.status}",
                "UnknownOrderStatus",
            )
        if (status != OrderStatus.PAID) {
            throw ApplicationFailure.newNonRetryableFailure(
                "order ${request.orderUid} is not fulfillable in status ${state.status}; must be PAID",
                "OrderNotFulfillable",
            )
        }

        orderActivities.markOrderShipped(request.orderUid)

        val shipment = try {
            deliveryActivities.createShipment(
                CreateShipmentRequest(
                    orderUid = request.orderUid,
                    originAddress = "warehouse",
                    destinationAddress = destination.addressLine,
                    destinationCity = destination.city,
                    destinationState = destination.state,
                    destinationCountryCode = destination.countryCode,
                    destinationPostalCode = destination.postalCode,
                ),
            )
        } catch (e: ActivityFailure) {
            orderActivities.restoreOrderPaid(request.orderUid)
            throw e
        }

        try {
            notificationActivities.sendShipmentUpdate(
                ShipmentNotificationRequest(
                    userUid = request.userUid,
                    orderUid = request.orderUid,
                    notificationType = "SHIPMENT_DISPATCHED",
                    trackingNumber = shipment.trackingNumber,
                ),
            )
        } catch (e: ActivityFailure) {
            Workflow.getLogger(OrderFulfillmentWorkflowImpl::class.java)
                .warn("shipment notification failed for order ${request.orderUid}; continuing fulfillment", e)
        }

        return FulfillOrderWorkflowResult(
            orderUid = request.orderUid,
            shipmentUid = shipment.shipmentUid,
            status = OrderStatus.SHIPPED.value,
        )
    }
}
