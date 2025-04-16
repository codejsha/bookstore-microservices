package com.codejsha.bookstore.service.application.port.openapi.model

import java.util.Objects
import com.fasterxml.jackson.annotation.JsonValue
import com.fasterxml.jackson.annotation.JsonCreator
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
* Values: UNKNOWN,PENDING,PAID,SHIPPING,COMPLETED,CANCELLED
*/
enum class OrderStatus(@get:JsonValue val value: kotlin.String) {

    UNKNOWN("unknown"),
    PENDING("pending"),
    PAID("paid"),
    SHIPPING("shipping"),
    COMPLETED("completed"),
    CANCELLED("cancelled");

    companion object {
        @JvmStatic
        @JsonCreator
        fun forValue(value: kotlin.String): OrderStatus {
                return values().first{it -> it.value == value}
        }
    }
}

