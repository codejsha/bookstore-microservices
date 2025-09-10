package com.codejsha.bookstore.admin.domain.model.external

import java.time.OffsetDateTime

data class Payment(
    val uid: String,
    val paymentId: String,
    val customerId: String?,
    val status: String,
    val amount: Long,
    val amountCaptured: Long?,
    val amountCapturable: Long?,
    val currency: String,
    val paymentMethod: String?,
    val connector: String?,
    val errorCode: String?,
    val errorMessage: String?,
    val confirmedAt: OffsetDateTime?,
    val capturedAt: OffsetDateTime?,
    val cancelledAt: OffsetDateTime?,
    val createdAt: OffsetDateTime,
    val updatedAt: OffsetDateTime?,
)


data class Refund(
    val uid: String,
    val refundId: String,
    val paymentId: String,
    val status: String,
    val refundType: String,
    val amount: Long,
    val currency: String,
    val reason: String?,
    val connector: String?,
    val errorCode: String?,
    val errorMessage: String?,
    val createdAt: OffsetDateTime,
    val updatedAt: OffsetDateTime?,
)
