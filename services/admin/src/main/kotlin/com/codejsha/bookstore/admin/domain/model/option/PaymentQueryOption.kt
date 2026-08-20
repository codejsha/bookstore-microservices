package com.codejsha.bookstore.admin.domain.model.option

data class PaymentQueryOption(
    val customerId: String? = null,
    val status: String? = null,
    val connector: String? = null,
)
