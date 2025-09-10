package com.codejsha.bookstore.payment.domain.constant

enum class RefundStatus(
    val value: String
) {
    PENDING("pending"),
    SUCCEEDED("succeeded"),
    FAILED("failed"),
    REVIEW("review");

    companion object {
        fun fromValue(value: String): RefundStatus =
            entries.first { it.value == value }

        fun fromValueOrNull(value: String): RefundStatus? =
            entries.firstOrNull { it.value == value }
    }
}
