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
* Values: CARD
*/
enum class PaymentType(@get:JsonValue val value: kotlin.String) {

    CARD("card");

    companion object {
        @JvmStatic
        @JsonCreator
        fun forValue(value: kotlin.String): PaymentType {
                return values().first{it -> it.value == value}
        }
    }
}

