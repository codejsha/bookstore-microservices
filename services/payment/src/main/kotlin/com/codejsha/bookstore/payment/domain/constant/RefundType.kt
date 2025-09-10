package com.codejsha.bookstore.payment.domain.constant

enum class RefundType(
    val value: String
) {
    INSTANT("instant"),
    SCHEDULED("scheduled");

    companion object {
        fun fromValue(value: String): RefundType =
            entries.first { it.value == value }
    }
}
