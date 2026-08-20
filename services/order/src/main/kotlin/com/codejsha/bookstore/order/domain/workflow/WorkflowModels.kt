package com.codejsha.bookstore.order.domain.workflow

import java.math.BigDecimal

data class PlaceOrderWorkflowRequest(
    val userUid: String,
    val currency: String,
    val items: List<OrderItemWorkflowRequest>,
    val shipping: ShippingWorkflowRequest?,
    val idempotencyKey: String,
)

data class OrderItemWorkflowRequest(
    val productId: Long,
    val productName: String?,
    val quantity: Int,
    val currency: String,
    val price: BigDecimal,
)

data class ShippingWorkflowRequest(
    val recipientName: String,
    val recipientPhone: String,
    val addressLine1: String,
    val addressLine2: String?,
    val city: String,
    val state: String,
    val postalCode: String,
    val country: String,
    val shippingMethod: String,
)

data class PlaceOrderWorkflowResult(
    val orderUid: String,
    val paymentUid: String?,
    val status: String,
)

data class CreateOrderRequest(
    val userUid: String,
    val currency: String,
    val itemsAmount: BigDecimal,
    val totalAmount: BigDecimal,
    val idempotencyKey: String,
    val items: List<OrderItemWorkflowRequest>,
    val shipping: ShippingWorkflowRequest?,
)

data class CreateOrderResult(
    val orderUid: String,
    val status: String,
    val paymentUid: String?,
)

data class ReserveStockRequest(
    val orderUid: String,
    val items: List<StockReservationItem>,
)

data class StockReservationItem(
    val productId: Long,
    val quantity: Int,
)

data class ProcessPaymentRequest(
    val orderUid: String,
    val userUid: String,
    val amount: BigDecimal,
    val currency: String,
)

data class ProcessPaymentResult(
    val paymentUid: String,
    val status: String,
    val gatewayPaymentId: String? = null,
)

// ─── Cancellation / Refund saga ──────────────────────────────────────────────

data class CancelOrderWorkflowRequest(
    val orderUid: String,
    val reason: String?,
)

data class CancelOrderWorkflowResult(
    val orderUid: String,
    val status: String,
)

data class OrderCancellationState(
    val status: String,
    val paymentUid: String?,
)

// ─── Fulfillment / Delivery saga ─────────────────────────────────────────────

data class FulfillOrderWorkflowRequest(
    val orderUid: String,
    val userUid: String,
    val destination: ShipmentDestination,
)

data class ShipmentDestination(
    val addressLine: String,
    val city: String,
    val state: String,
    val countryCode: String,
    val postalCode: String,
)

data class FulfillOrderWorkflowResult(
    val orderUid: String,
    val shipmentUid: String,
    val status: String,
)

data class CreateShipmentRequest(
    val orderUid: String,
    val originAddress: String,
    val destinationAddress: String,
    val destinationCity: String,
    val destinationState: String,
    val destinationCountryCode: String,
    val destinationPostalCode: String,
)

data class CreateShipmentResult(
    val shipmentUid: String,
    val trackingNumber: String?,
    val status: String,
)

data class ShipmentNotificationRequest(
    val userUid: String,
    val orderUid: String,
    val notificationType: String,
    val trackingNumber: String?,
)
