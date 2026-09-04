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
) {
    init {
        requireNonBlank("customer_id", customerId)
        requireMaxLength("customer_id", customerId, 64)
        requireNonBlankIfPresent("name", name)
        requireMaxLengthIfPresent("name", name, 255)
        requireEmailIfPresent("email", email)
        requireMaxLengthIfPresent("email", email, 255)
        requireNonBlankIfPresent("phone", phone)
        requireMaxLengthIfPresent("phone", phone, 32)
        requireNonBlankIfPresent("phone_country_code", phoneCountryCode)
        requireMaxLengthIfPresent("phone_country_code", phoneCountryCode, 8)
        requireMaxLengthIfPresent("description", description, 500)
    }
}

data class CustomerUpdateCommand(
    val name: String?,
    val email: String?,
    val phone: String?,
    val phoneCountryCode: String?,
    val description: String?,
    val metadata: Map<String, Any>?,
    val defaultBillingAddress: Map<String, Any>?,
    val defaultShippingAddress: Map<String, Any>?,
) {
    init {
        requireNonBlankIfPresent("name", name)
        requireMaxLengthIfPresent("name", name, 255)
        requireEmailIfPresent("email", email)
        requireMaxLengthIfPresent("email", email, 255)
        requireNonBlankIfPresent("phone", phone)
        requireMaxLengthIfPresent("phone", phone, 32)
        requireNonBlankIfPresent("phone_country_code", phoneCountryCode)
        requireMaxLengthIfPresent("phone_country_code", phoneCountryCode, 8)
        requireMaxLengthIfPresent("description", description, 500)
    }
}
