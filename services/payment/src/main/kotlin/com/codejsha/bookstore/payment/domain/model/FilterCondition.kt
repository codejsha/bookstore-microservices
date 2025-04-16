package com.codejsha.bookstore.payment.domain.model

import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Pageable
import org.springframework.data.domain.Sort

data class FilterCondition(
    val filter: QueryFilter
) {
    data class QueryFilter(
        val sort: String,
        val order: String,
        val limit: Int,
        val offset: Int
    ) {
        fun toPageable(): Pageable {
            return PageRequest.of(offset, limit, Sort.by(Sort.Direction.valueOf(order), sort))
        }
    }
}
