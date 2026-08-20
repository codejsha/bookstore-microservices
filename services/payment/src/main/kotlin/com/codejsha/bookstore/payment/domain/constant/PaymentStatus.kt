package com.codejsha.bookstore.payment.domain.constant

enum class PaymentStatus(
    val value: String
) {
    REQUIRES_PAYMENT_METHOD("requires_payment_method"),
    REQUIRES_CONFIRMATION("requires_confirmation"),
    REQUIRES_CUSTOMER_ACTION("requires_customer_action"),
    REQUIRES_CAPTURE("requires_capture"),
    PROCESSING("processing"),
    SUCCEEDED("succeeded"),
    FAILED("failed"),
    CANCELLED("cancelled"),
    PARTIALLY_CAPTURED("partially_captured"),
    EXPIRED("expired");

    companion object {
        fun fromValue(value: String): PaymentStatus =
            entries.first { it.value == value }

        fun fromValueOrNull(value: String): PaymentStatus? =
            entries.firstOrNull { it.value == value }
    }
}
