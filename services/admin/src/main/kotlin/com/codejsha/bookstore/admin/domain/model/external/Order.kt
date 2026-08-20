package com.codejsha.bookstore.admin.domain.model.external

import java.time.OffsetDateTime

data class Order(
    val uid: String,
    val orderNumber: String,
    val userUid: String,
    val status: String,
    val currency: String,
    val itemsAmount: Double,
    val discountAmount: Double,
    val shippingAmount: Double,
    val taxAmount: Double,
    val totalAmount: Double,
    val lines: List<OrderLine>,
    val shipping: OrderShipping?,
    val createdAt: OffsetDateTime,
    val updatedAt: OffsetDateTime?,
)

data class OrderLine(
    val uid: String,
    val productId: Long,
    val productName: String?,
    val sku: String?,
    val quantity: Int,
    val price: Double,
    val subtotal: Double,
    val taxRate: Double,
    val currency: String,
    val options: String?,
)

data class OrderShipping(
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
