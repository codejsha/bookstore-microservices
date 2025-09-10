package com.codejsha.bookstore.payment.domain.constant

enum class CaptureMethod(
    val value: String
) {
    AUTOMATIC("automatic"),
    MANUAL("manual"),
    MANUAL_MULTIPLE("manual_multiple"),
    SCHEDULED("scheduled");

    companion object {
        fun fromValue(value: String): CaptureMethod =
            entries.first { it.value == value }
    }
}
