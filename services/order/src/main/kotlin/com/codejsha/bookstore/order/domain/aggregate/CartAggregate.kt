package com.codejsha.bookstore.order.domain.aggregate

import com.codejsha.bookstore.order.application.port.repo.CartItemResult
import com.codejsha.bookstore.order.application.port.repo.CartResult
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.UUID

data class CartAggregate(
    val id: Long,
    val uid: UUID,
    val userUid: UUID,
    val items: List<CartItemEntity>,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class CartItemEntity(
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

fun CartResult.toAggregate() = CartAggregate(
    id = id,
    uid = uid,
    userUid = userUid,
    items = emptyList(),
    createdAt = createdAt,
    updatedAt = updatedAt,
)

fun CartItemResult.toEntity() = CartItemEntity(
    id = id,
    uid = uid,
    cartId = cartId,
    productId = productId,
    productName = productName,
    quantity = quantity,
    currency = currency,
    price = price,
    createdAt = createdAt,
    updatedAt = updatedAt,
)
