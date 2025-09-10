package com.codejsha.bookstore.order.domain.aggregate

import com.codejsha.bookstore.order.application.port.repo.OrderAdjustmentResult
import com.codejsha.bookstore.order.application.port.repo.OrderItemResult
import com.codejsha.bookstore.order.application.port.repo.OrderResult
import com.codejsha.bookstore.order.application.port.repo.OrderShippingResult
import com.codejsha.bookstore.order.domain.constant.OrderAdjustmentType
import com.codejsha.bookstore.order.domain.constant.OrderStatus
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.UUID

data class OrderAggregate(
    val id: Long,
    val uid: UUID,
    val userUid: UUID,
    val orderNumber: String,
    val status: OrderStatus,
    val currency: String,
    val itemsAmount: BigDecimal,
    val discountAmount: BigDecimal,
    val shippingAmount: BigDecimal,
    val taxAmount: BigDecimal,
    val totalAmount: BigDecimal,
    val idempotencyKey: String,
    val paymentUid: UUID?,
    val items: List<OrderItemEntity>,
    val adjustments: List<OrderAdjustmentEntity>,
    val shipping: OrderShippingEntity?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class OrderItemEntity(
    val id: Long,
    val uid: UUID,
    val orderId: Long,
    val productId: Long,
    val sku: String?,
    val productName: String?,
    val options: String?,
    val quantity: Int,
    val currency: String,
    val price: BigDecimal,
    val taxRate: BigDecimal,
    val subtotal: BigDecimal,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class OrderAdjustmentEntity(
    val id: Long,
    val uid: UUID,
    val orderId: Long,
    val type: OrderAdjustmentType,
    val label: String?,
    val amount: BigDecimal,
    val meta: Map<String, Any>?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class OrderShippingEntity(
    val id: Long,
    val uid: UUID,
    val orderId: Long,
    val recipientName: String,
    val recipientPhone: String,
    val addressLine1: String,
    val addressLine2: String?,
    val city: String,
    val state: String,
    val postalCode: String,
    val country: String,
    val shippingMethod: String,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

fun OrderResult.toAggregate() = OrderAggregate(
    id = id,
    uid = uid,
    userUid = userUid,
    orderNumber = orderNumber,
    status = OrderStatus.fromValue(status),
    currency = currency,
    itemsAmount = itemsAmount,
    discountAmount = discountAmount,
    shippingAmount = shippingAmount,
    taxAmount = taxAmount,
    totalAmount = totalAmount,
    idempotencyKey = idempotencyKey,
    paymentUid = paymentUid,
    items = emptyList(),
    adjustments = emptyList(),
    shipping = null,
    createdAt = createdAt,
    updatedAt = updatedAt,
)

fun OrderItemResult.toEntity() = OrderItemEntity(
    id = id,
    uid = uid,
    orderId = orderId,
    productId = productId,
    sku = sku,
    productName = productName,
    options = options,
    quantity = quantity,
    currency = currency,
    price = price,
    taxRate = taxRate,
    subtotal = subtotal,
    createdAt = createdAt,
    updatedAt = updatedAt,
)

fun OrderAdjustmentResult.toEntity() = OrderAdjustmentEntity(
    id = id,
    uid = uid,
    orderId = orderId,
    type = OrderAdjustmentType.fromValue(type),
    label = label,
    amount = amount,
    meta = meta,
    createdAt = createdAt,
    updatedAt = updatedAt,
)

fun OrderShippingResult.toEntity() = OrderShippingEntity(
    id = id,
    uid = uid,
    orderId = orderId,
    recipientName = recipientName,
    recipientPhone = recipientPhone,
    addressLine1 = addressLine1,
    addressLine2 = addressLine2,
    city = city,
    state = state,
    postalCode = postalCode,
    country = country,
    shippingMethod = shippingMethod,
    createdAt = createdAt,
    updatedAt = updatedAt,
)
