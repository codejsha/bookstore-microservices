package com.codejsha.bookstore.order.application.port.repo

import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.UUID

data class OrderResult(
    val id: Long,
    val uid: UUID,
    val userUid: UUID,
    val orderNumber: String,
    val status: String,
    val currency: String,
    val itemsAmount: BigDecimal,
    val discountAmount: BigDecimal,
    val shippingAmount: BigDecimal,
    val taxAmount: BigDecimal,
    val totalAmount: BigDecimal,
    val idempotencyKey: String,
    val paymentUid: UUID?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class OrderItemResult(
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

data class OrderAdjustmentResult(
    val id: Long,
    val uid: UUID,
    val orderId: Long,
    val type: String,
    val label: String?,
    val amount: BigDecimal,
    val meta: Map<String, Any>?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class OrderShippingResult(
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

data class CartResult(
    val id: Long,
    val uid: UUID,
    val userUid: UUID,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class CartItemResult(
    val id: Long,
    val uid: UUID,
    val cartId: Long,
    val productId: Long,
    val productName: String?,
    val quantity: Int,
    val currency: String,
    val price: BigDecimal,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)
