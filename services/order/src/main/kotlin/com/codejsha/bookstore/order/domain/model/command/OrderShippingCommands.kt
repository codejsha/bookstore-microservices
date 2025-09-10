package com.codejsha.bookstore.order.domain.model.command

data class OrderShippingCreateCommand(
    val recipientName: String,
    val recipientPhone: String,
    val addressLine1: String,
    val addressLine2: String?,
    val city: String,
    val state: String,
    val postalCode: String,
    val country: String,
    val shippingMethod: String,
)

data class OrderShippingUpdateCommand(
    val recipientName: String?,
    val recipientPhone: String?,
    val addressLine1: String?,
    val addressLine2: String?,
    val city: String?,
    val state: String?,
    val postalCode: String?,
    val country: String?,
    val shippingMethod: String?,
)
