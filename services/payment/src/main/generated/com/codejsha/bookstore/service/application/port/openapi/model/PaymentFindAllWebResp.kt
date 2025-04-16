package com.codejsha.bookstore.service.application.port.openapi.model

import java.util.Objects
import com.codejsha.bookstore.service.application.port.openapi.model.PaymentFindWebResp
import com.fasterxml.jackson.annotation.JsonProperty
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
 * @param total 
 * @param items 
 */
data class PaymentFindAllWebResp(

    @get:Min(0L)
    @get:JsonProperty("total", required = true) val total: kotlin.Long,

    @field:Valid
    @get:Size(min=0)
    @get:JsonProperty("items", required = true) val items: kotlin.collections.List<PaymentFindWebResp>
    ) {

}

