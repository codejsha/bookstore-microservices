package com.codejsha.bookstore.order.support

import com.codejsha.bookstore.order.domain.workflow.CreateOrderRequest
import com.codejsha.bookstore.order.domain.workflow.CreateOrderResult
import com.codejsha.bookstore.order.domain.workflow.CreateShipmentRequest
import com.codejsha.bookstore.order.domain.workflow.CreateShipmentResult
import com.codejsha.bookstore.order.domain.workflow.DeliveryActivities
import com.codejsha.bookstore.order.domain.workflow.InventoryActivities
import com.codejsha.bookstore.order.domain.workflow.NotificationActivities
import com.codejsha.bookstore.order.domain.workflow.OrderActivities
import com.codejsha.bookstore.order.domain.workflow.OrderCancellationState
import com.codejsha.bookstore.order.domain.workflow.PaymentActivities
import com.codejsha.bookstore.order.domain.workflow.ProcessPaymentRequest
import com.codejsha.bookstore.order.domain.workflow.ProcessPaymentResult
import com.codejsha.bookstore.order.domain.workflow.ReserveStockRequest
import com.codejsha.bookstore.order.domain.workflow.ShipmentNotificationRequest
import com.codejsha.bookstore.order.domain.workflow.StockReservationItem
import java.util.Collections

class FakeOrderActivities : OrderActivities {
    var createOrderReturn: CreateOrderResult =
        CreateOrderResult(orderUid = "00000000-0000-0000-0000-000000000000", status = "PENDING", paymentUid = null)
    @Volatile var cancellationState: OrderCancellationState =
        OrderCancellationState(status = "PENDING", paymentUid = null)
    @Volatile var orderItems: List<StockReservationItem> = emptyList()

    val createOrderRequests: MutableList<CreateOrderRequest> = synced()
    val confirmedOrderUids: MutableList<String> = synced()
    val cancelledOrderUids: MutableList<String> = synced()
    val recordedPayments: MutableList<Pair<String, String>> = synced()
    val refundedOrderUids: MutableList<String> = synced()
    val shippedOrderUids: MutableList<String> = synced()
    val restoredPaidOrderUids: MutableList<String> = synced()

    override fun createOrder(request: CreateOrderRequest): CreateOrderResult {
        createOrderRequests.add(request)
        return createOrderReturn
    }

    override fun confirmOrder(orderUid: String) {
        confirmedOrderUids.add(orderUid)
    }

    override fun cancelOrder(orderUid: String) {
        cancelledOrderUids.add(orderUid)
    }

    override fun recordPayment(orderUid: String, paymentUid: String) {
        recordedPayments.add(orderUid to paymentUid)
    }

    override fun loadCancellationState(orderUid: String): OrderCancellationState = cancellationState

    override fun loadOrderItems(orderUid: String): List<StockReservationItem> = orderItems

    override fun markOrderRefunded(orderUid: String) {
        refundedOrderUids.add(orderUid)
    }

    override fun markOrderShipped(orderUid: String) {
        shippedOrderUids.add(orderUid)
    }

    override fun restoreOrderPaid(orderUid: String) {
        restoredPaidOrderUids.add(orderUid)
    }
}

class FakeInventoryActivities : InventoryActivities {
    @Volatile var reserveStockFailure: RuntimeException? = null

    val reserveStockRequests: MutableList<ReserveStockRequest> = synced()
    val releaseStockRequests: MutableList<ReserveStockRequest> = synced()

    override fun reserveStock(request: ReserveStockRequest) {
        reserveStockRequests.add(request)
        reserveStockFailure?.let { throw it }
    }

    override fun releaseStock(request: ReserveStockRequest) {
        releaseStockRequests.add(request)
    }
}

class FakePaymentActivities : PaymentActivities {
    @Volatile var processPaymentReturn: ProcessPaymentResult =
        ProcessPaymentResult(paymentUid = "00000000-0000-0000-0000-000000000000", status = "succeeded")
    @Volatile var processPaymentFailure: RuntimeException? = null
    @Volatile var refundPaymentReturn: Boolean = true
    @Volatile var refundPaymentByOrderReturn: Boolean = false

    val processPaymentRequests: MutableList<ProcessPaymentRequest> = synced()
    val refundedPaymentUids: MutableList<String> = synced()
    val refundedByOrderUids: MutableList<String> = synced()

    override fun processPayment(request: ProcessPaymentRequest): ProcessPaymentResult {
        processPaymentRequests.add(request)
        processPaymentFailure?.let { throw it }
        return processPaymentReturn
    }

    override fun refundPayment(paymentUid: String): Boolean {
        refundedPaymentUids.add(paymentUid)
        return refundPaymentReturn
    }

    override fun refundPaymentByOrder(orderUid: String): Boolean {
        refundedByOrderUids.add(orderUid)
        return refundPaymentByOrderReturn
    }
}

class FakeDeliveryActivities : DeliveryActivities {
    @Volatile var createShipmentReturn: CreateShipmentResult =
        CreateShipmentResult(shipmentUid = "00000000-0000-0000-0000-000000000000", trackingNumber = null, status = "CREATED")

    val createShipmentRequests: MutableList<CreateShipmentRequest> = synced()

    override fun createShipment(request: CreateShipmentRequest): CreateShipmentResult {
        createShipmentRequests.add(request)
        return createShipmentReturn
    }
}

class FakeNotificationActivities : NotificationActivities {
    val sentNotifications: MutableList<ShipmentNotificationRequest> = synced()

    override fun sendShipmentUpdate(request: ShipmentNotificationRequest) {
        sentNotifications.add(request)
    }
}

private fun <T> synced(): MutableList<T> = Collections.synchronizedList(mutableListOf())
