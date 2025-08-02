package com.codejsha.bookstore.service.application.port.openapi.model

import com.fasterxml.jackson.annotation.JsonCreator
import com.fasterxml.jackson.annotation.JsonValue

/**
 *
 * Values: UNKNOWN,PENDING,PAID,SHIPPING,COMPLETED,CANCELLED
 */
enum class OrderStatus(
    @get:JsonValue val value: kotlin.String
) {
    UNKNOWN("unknown"),
    PENDING("pending"),
    PAID("paid"),
    SHIPPING("shipping"),
    COMPLETED("completed"),
    CANCELLED("cancelled");

    companion object {
        @JvmStatic
        @JsonCreator
        fun forValue(value: kotlin.String): OrderStatus = values().first { it -> it.value == value }
    }
}
