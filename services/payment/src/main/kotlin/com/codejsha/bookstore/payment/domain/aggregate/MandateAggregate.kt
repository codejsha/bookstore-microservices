package com.codejsha.bookstore.payment.domain.aggregate

import com.codejsha.bookstore.payment.domain.constant.FutureUsage
import com.codejsha.bookstore.payment.domain.constant.MandateStatus
import com.codejsha.bookstore.payment.domain.constant.MandateType
import java.time.LocalDateTime
import java.util.UUID

data class MandateAggregate(
    val id: Long,
    val uid: UUID,
    val mandateId: String,
    val customerId: String,
    val paymentMethodId: String?,
    val mandateType: MandateType,
    val mandateStatus: MandateStatus,
    val mandateAmount: Long?,
    val mandateCurrency: String?,
    val startDate: LocalDateTime?,
    val endDate: LocalDateTime?,
    val setupFutureUsage: FutureUsage?,
    val customerAcceptanceType: String?,
    val customerAcceptedAt: LocalDateTime?,
    val metadata: Map<String, Any>?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)
