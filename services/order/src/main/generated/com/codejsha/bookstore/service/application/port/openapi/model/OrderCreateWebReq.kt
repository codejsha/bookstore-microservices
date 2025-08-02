package com.codejsha.bookstore.service.application.port.openapi.model

import com.codejsha.bookstore.service.application.port.openapi.model.OrderItem
import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType

import com.fasterxml.jackson.annotation.JsonProperty
import jakarta.validation.Valid

/**
 *
 * @param userId
 * @param totalPrice
 * @param status
 * @param orderItems
 * @param paymentType
 * @param cardNumber
 */
data class OrderCreateWebReq(
    @get:JsonProperty("user_id") val userId: kotlin.String? = null,
    @get:JsonProperty("total_price") val totalPrice: java.math.BigDecimal? = null,
    @field:Valid
    @get:JsonProperty("status") val status: OrderStatus? = null,
    @field:Valid
    @get:JsonProperty("order_items") val orderItems: kotlin.collections.List<OrderItem>? = null,
    @field:Valid
    @get:JsonProperty("PaymentType") val paymentType: PaymentType? = null,
    @get:JsonProperty("card_number") val cardNumber: kotlin.String? = null
)
