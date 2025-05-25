package com.codejsha.bookstore.service.application.port.openapi.model

import java.util.Objects
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
 * @param id 
 * @param orderId 
 * @param userId 
 * @param paymentType 
 * @param cardNumber 
 * @param amount 
 * @param paymentDate 
 */
data class PaymentFindWebResp(

    @get:JsonProperty("id") val id: kotlin.Long? = null,

    @get:JsonProperty("order_id") val orderId: kotlin.Long? = null,

    @get:JsonProperty("user_id") val userId: kotlin.String? = null,

    @field:Valid
    @get:JsonProperty("PaymentType") val paymentType: PaymentType? = null,

    @get:JsonProperty("card_number") val cardNumber: kotlin.String? = null,

    @get:JsonProperty("amount") val amount: java.math.BigDecimal? = null,

    @get:JsonProperty("payment_date") val paymentDate: java.time.OffsetDateTime? = null
    ) {

}

