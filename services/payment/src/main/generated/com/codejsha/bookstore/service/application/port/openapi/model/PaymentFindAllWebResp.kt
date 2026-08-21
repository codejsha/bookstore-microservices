package com.codejsha.bookstore.service.application.port.openapi.model

import com.codejsha.bookstore.service.application.port.openapi.model.PaymentFindWebResp

import com.fasterxml.jackson.annotation.JsonProperty
import jakarta.validation.Valid
import jakarta.validation.constraints.Min
import jakarta.validation.constraints.Size

/**
 *
 * @param total
 * @param items
 */
data class PaymentFindAllWebResp(
    @get:Min(0L)
    @get:JsonProperty("total", required = true) val total: kotlin.Long,
    @field:Valid
    @get:Size(min = 0)
    @get:JsonProperty("items", required = true) val items: kotlin.collections.List<PaymentFindWebResp>
)
