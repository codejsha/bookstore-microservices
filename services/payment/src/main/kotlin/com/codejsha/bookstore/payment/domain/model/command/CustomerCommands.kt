package com.codejsha.bookstore.payment.domain.model.command

data class CustomerCreateCommand(
    val customerId: String,
    val name: String?,
    val email: String?,
    val phone: String?,
    val phoneCountryCode: String?,
    val description: String?,
    val metadata: Map<String, Any>?,
    val defaultBillingAddress: Map<String, Any>?,
    val defaultShippingAddress: Map<String, Any>?,
)

data class CustomerUpdateCommand(
    val name: String?,
    val email: String?,
    val phone: String?,
    val phoneCountryCode: String?,
    val description: String?,
    val metadata: Map<String, Any>?,
    val defaultBillingAddress: Map<String, Any>?,
    val defaultShippingAddress: Map<String, Any>?,
)
