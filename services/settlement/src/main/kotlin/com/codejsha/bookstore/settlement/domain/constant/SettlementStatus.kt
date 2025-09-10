package com.codejsha.bookstore.settlement.domain.constant

enum class SettlementStatus(val value: String) {
    OPEN("OPEN"),
    CONFIRMED("CONFIRMED"),
    DISCREPANCY("DISCREPANCY");

    companion object {
        fun fromValue(value: String): SettlementStatus = entries.first { it.value == value }
    }
}
