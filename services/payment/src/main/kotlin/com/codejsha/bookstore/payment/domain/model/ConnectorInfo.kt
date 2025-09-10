package com.codejsha.bookstore.payment.domain.model

data class ConnectorInfo(
    val connector: String,
    val transactionId: String?,
)
