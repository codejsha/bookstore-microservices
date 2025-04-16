package com.codejsha.bookstore.order.domain.model

data class FilterCondition(
    val userId: String?,
    val filter: QueryFilter
) {
    data class QueryFilter(
        val sort: String,
        val order: String,
        val limit: Int,
        val offset: Int
    )
}
