package com.codejsha.bookstore.payment.domain.constant

enum class MandateType(
    val value: String
) {
    SINGLE_USE("single_use"),
    MULTI_USE("multi_use");

    companion object {
        fun fromValue(value: String): MandateType =
            entries.first { it.value == value }
    }
}
