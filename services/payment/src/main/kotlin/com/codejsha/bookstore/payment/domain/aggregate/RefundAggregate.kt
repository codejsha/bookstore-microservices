package com.codejsha.bookstore.payment.domain.aggregate

import com.codejsha.bookstore.payment.domain.model.Money
import com.codejsha.bookstore.payment.domain.constant.RefundStatus
import com.codejsha.bookstore.payment.domain.constant.RefundType
import java.time.LocalDateTime
import java.util.UUID

data class RefundAggregate(
    val id: Long,
    val uid: UUID,
    val refundId: String,
    val paymentId: String,
    val connector: String?,
    val connectorRefundId: String?,
    val amount: Money,
    val status: RefundStatus,
    val reason: String?,
    val refundType: RefundType,
    val errorCode: String?,
    val errorMessage: String?,
    val metadata: Map<String, Any>?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)
