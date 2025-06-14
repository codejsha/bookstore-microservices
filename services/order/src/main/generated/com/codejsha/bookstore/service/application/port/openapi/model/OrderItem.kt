package com.codejsha.bookstore.service.application.port.openapi.model

import java.util.Objects
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
 * @param bookId 
 * @param quantity 
 */
data class OrderItem(

    @get:JsonProperty("book_id") val bookId: kotlin.Long? = null,

    @get:JsonProperty("quantity") val quantity: kotlin.Int? = null
    ) {

}

