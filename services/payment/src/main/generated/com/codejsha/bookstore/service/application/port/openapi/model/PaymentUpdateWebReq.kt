package com.codejsha.bookstore.service.application.port.openapi.model

import com.codejsha.bookstore.service.application.port.openapi.model.PaymentType

import com.fasterxml.jackson.annotation.JsonProperty
import jakarta.validation.Valid

/**
 *
 * @param orderId
 * @param userId
 * @param paymentType
 * @param cardNumber
 * @param amount
 * @param paymentDate
 */
data class PaymentUpdateWebReq(
    @get:JsonProperty("order_id") val orderId: kotlin.Long? = null,
    @get:JsonProperty("user_id") val userId: kotlin.String? = null,
    @field:Valid
    @get:JsonProperty("PaymentType") val paymentType: PaymentType? = null,
    @get:JsonProperty("card_number") val cardNumber: kotlin.String? = null,
    @get:JsonProperty("amount") val amount: java.math.BigDecimal? = null,
    @get:JsonProperty("payment_date") val paymentDate: java.time.OffsetDateTime? = null
)
