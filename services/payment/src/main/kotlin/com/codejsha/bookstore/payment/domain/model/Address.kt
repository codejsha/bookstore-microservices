package com.codejsha.bookstore.payment.domain.model

data class Address(
    val city: String?,
    val country: String?,
    val line1: String?,
    val line2: String?,
    val state: String?,
    val zip: String?,
    val firstName: String?,
    val lastName: String?,
)
