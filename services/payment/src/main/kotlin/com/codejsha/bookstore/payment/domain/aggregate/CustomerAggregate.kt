package com.codejsha.bookstore.payment.domain.aggregate

import com.codejsha.bookstore.payment.domain.model.Address
import com.codejsha.bookstore.payment.domain.constant.PaymentMethodType
import java.time.LocalDateTime
import java.util.UUID

data class CustomerAggregate(
    val id: Long,
    val uid: UUID,
    val customerId: String,
    val name: String?,
    val email: String?,
    val phone: String?,
    val phoneCountryCode: String?,
    val description: String?,
    val metadata: Map<String, Any>?,
    val defaultBillingAddress: Address?,
    val defaultShippingAddress: Address?,
    val paymentMethods: List<PaymentMethodEntity>,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class PaymentMethodEntity(
    val id: Long,
    val uid: UUID,
    val paymentMethodId: String,
    val customerId: String,
    val paymentMethod: PaymentMethodType,
    val paymentMethodType: String?,
    val paymentMethodIssuer: String?,
    val cardNetwork: String?,
    val cardLast4: String?,
    val cardExpMonth: Int?,
    val cardExpYear: Int?,
    val cardHolderName: String?,
    val isDefault: Boolean,
    val metadata: Map<String, Any>?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)
