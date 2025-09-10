package com.codejsha.bookstore.payment.domain.constant

enum class MandateStatus(
    val value: String
) {
    ACTIVE("active"),
    REVOKED("revoked"),
    INACTIVE("inactive"),
    PENDING("pending");

    companion object {
        fun fromValue(value: String): MandateStatus =
            entries.first { it.value == value }
    }
}
