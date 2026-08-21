package com.codejsha.bookstore.service.application.port.openapi.model

import com.fasterxml.jackson.annotation.JsonCreator
import com.fasterxml.jackson.annotation.JsonValue

/**
 *
 * Values: CARD
 */
enum class PaymentType(
    @get:JsonValue val value: kotlin.String
) {
    CARD("card");

    companion object {
        @JvmStatic
        @JsonCreator
        fun forValue(value: kotlin.String): PaymentType = values().first { it -> it.value == value }
    }
}
