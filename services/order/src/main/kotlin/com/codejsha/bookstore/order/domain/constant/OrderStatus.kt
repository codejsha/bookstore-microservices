package com.codejsha.bookstore.order.domain.constant

enum class OrderStatus(val value: String) {
    PENDING("PENDING"),
    PAID("PAID"),
    SHIPPED("SHIPPED"),
    DELIVERED("DELIVERED"),
    CANCELLED("CANCELLED"),
    REFUNDED("REFUNDED");

    companion object {
        fun fromValue(value: String): OrderStatus = entries.first { it.value == value }

        fun fromValueOrNull(value: String): OrderStatus? = entries.firstOrNull { it.value == value }
    }
}
