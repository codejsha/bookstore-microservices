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
) {
    init {
        requireNonBlank("recipient_name", recipientName)
        requireMaxLength("recipient_name", recipientName, 255)
        requireNonBlank("recipient_phone", recipientPhone)
        requireMaxLength("recipient_phone", recipientPhone, 20)
        requireNonBlank("address_line1", addressLine1)
        requireMaxLength("address_line1", addressLine1, 255)
        requireNonBlankIfPresent("address_line2", addressLine2)
        requireMaxLengthIfPresent("address_line2", addressLine2, 255)
        requireNonBlank("city", city)
        requireMaxLength("city", city, 100)
        requireNonBlank("state", state)
        requireMaxLength("state", state, 100)
        requireNonBlank("postal_code", postalCode)
        requireMaxLength("postal_code", postalCode, 20)
        requireCountry("country", country)
        requireNonBlank("shipping_method", shippingMethod)
        requireMaxLength("shipping_method", shippingMethod, 100)
    }
}

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
) {
    init {
        requireNonBlankIfPresent("recipient_name", recipientName)
        requireMaxLengthIfPresent("recipient_name", recipientName, 255)
        requireNonBlankIfPresent("recipient_phone", recipientPhone)
        requireMaxLengthIfPresent("recipient_phone", recipientPhone, 20)
        requireNonBlankIfPresent("address_line1", addressLine1)
        requireMaxLengthIfPresent("address_line1", addressLine1, 255)
        requireNonBlankIfPresent("address_line2", addressLine2)
        requireMaxLengthIfPresent("address_line2", addressLine2, 255)
        requireNonBlankIfPresent("city", city)
        requireMaxLengthIfPresent("city", city, 100)
        requireNonBlankIfPresent("state", state)
        requireMaxLengthIfPresent("state", state, 100)
        requireNonBlankIfPresent("postal_code", postalCode)
        requireMaxLengthIfPresent("postal_code", postalCode, 20)
        requireCountryIfPresent("country", country)
        requireNonBlankIfPresent("shipping_method", shippingMethod)
        requireMaxLengthIfPresent("shipping_method", shippingMethod, 100)
    }
}
