package com.codejsha.bookstore.payment.domain.constant

enum class FutureUsage(
    val value: String
) {
    ON_SESSION("on_session"),
    OFF_SESSION("off_session");

    companion object {
        fun fromValue(value: String): FutureUsage =
            entries.first { it.value == value }
    }
}
