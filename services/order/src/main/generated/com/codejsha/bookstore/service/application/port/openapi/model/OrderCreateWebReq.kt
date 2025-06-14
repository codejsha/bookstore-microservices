package com.codejsha.bookstore.service.application.port.openapi.model

import java.util.Objects
import com.codejsha.bookstore.service.application.port.openapi.model.OrderItem
import com.codejsha.bookstore.service.application.port.openapi.model.OrderStatus
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType
import com.fasterxml.jackson.annotation.JsonCreator
import com.fasterxml.jackson.annotation.JsonProperty
import com.fasterxml.jackson.annotation.JsonValue
import jakarta.validation.constraints.DecimalMax
import jakarta.validation.constraints.DecimalMin
import jakarta.validation.constraints.Email
import jakarta.validation.constraints.Max
import jakarta.validation.constraints.Min
import jakarta.validation.constraints.NotNull
import jakarta.validation.constraints.Pattern
import jakarta.validation.constraints.Size
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
    ) {

}

