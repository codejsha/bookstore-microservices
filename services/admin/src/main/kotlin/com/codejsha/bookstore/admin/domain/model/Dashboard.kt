package com.codejsha.bookstore.admin.domain.model

data class Metric(
    val count: Long,
    val available: Boolean,
) {
    companion object {
        fun of(count: Long) = Metric(count = count, available = true)

        fun unavailable() = Metric(count = 0, available = false)
    }
}

data class Dashboard(
    val works: Metric,
    val orders: Metric,
    val warehouses: Metric,
)
