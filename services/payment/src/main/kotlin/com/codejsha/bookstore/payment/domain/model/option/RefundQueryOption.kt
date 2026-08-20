package com.codejsha.bookstore.payment.domain.model.option

data class RefundQueryOption(
    val paymentId: String? = null,
    val status: String? = null,
)
