package com.codejsha.bookstore.service.application.port.openapi.model

import com.fasterxml.jackson.annotation.JsonProperty

/**
 *
 * @param bookId
 * @param quantity
 */
data class OrderItem(
    @get:JsonProperty("book_id") val bookId: kotlin.Long? = null,
    @get:JsonProperty("quantity") val quantity: kotlin.Int? = null
)
