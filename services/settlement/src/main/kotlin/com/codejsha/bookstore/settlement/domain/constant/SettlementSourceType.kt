package com.codejsha.bookstore.settlement.domain.constant

enum class SettlementSourceType(val value: String) {
    PAYMENT("PAYMENT"),
    REFUND("REFUND");

    companion object {
        fun fromValue(value: String): SettlementSourceType = entries.first { it.value == value }
    }
}
