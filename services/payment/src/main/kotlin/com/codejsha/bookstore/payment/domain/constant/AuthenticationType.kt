package com.codejsha.bookstore.payment.domain.constant

enum class AuthenticationType(
    val value: String
) {
    THREE_DS("three_ds"),
    NO_THREE_DS("no_three_ds");

    companion object {
        fun fromValue(value: String): AuthenticationType =
            entries.first { it.value == value }
    }
}
