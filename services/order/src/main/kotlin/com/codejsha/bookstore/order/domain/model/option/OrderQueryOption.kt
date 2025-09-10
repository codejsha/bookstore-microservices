package com.codejsha.bookstore.order.domain.model.option

import java.util.UUID

data class OrderQueryOption(
    val userUid: UUID? = null,
    val status: String? = null,
)
