package com.codejsha.bookstore.order.domain.constant

enum class OrderAdjustmentType(val value: String) {
    COUPON("COUPON"),
    POINT("POINT"),
    MANUAL("MANUAL"),
    SHIPPING("SHIPPING"),
    TAX("TAX");

    companion object {
        fun fromValue(value: String): OrderAdjustmentType = entries.first { it.value == value }
    }
}
