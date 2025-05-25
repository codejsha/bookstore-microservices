package com.codejsha.bookstore.payment.domain.model

import com.google.common.base.CaseFormat
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
            val sortingField = CaseFormat.LOWER_UNDERSCORE
                .to(CaseFormat.LOWER_CAMEL, sort)
            val direction = Sort.Direction.valueOf(order.uppercase())
            val sortOption = Sort.by(direction, sortingField)
            return PageRequest.of(offset, limit, sortOption)
        }
    }
}
