package com.codejsha.bookstore.payment.application.port.repo

import java.time.LocalDateTime
import java.util.UUID

// ─── Customer Aggregate ─────────────────────────────────────────────────────

data class CustomerResult(
    val id: Long,
    val uid: UUID,
    val customerId: String,
    val name: String?,
    val email: String?,
    val phone: String?,
    val phoneCountryCode: String?,
    val description: String?,
    val metadata: Map<String, Any>?,
    val defaultBillingAddress: Map<String, Any>?,
    val defaultShippingAddress: Map<String, Any>?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class PaymentMethodResult(
    val id: Long,
    val uid: UUID,
    val paymentMethodId: String,
    val customerId: String,
    val paymentMethod: String,
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

// ─── Payment Aggregate ──────────────────────────────────────────────────────

data class PaymentResult(
    val id: Long,
    val uid: UUID,
    val paymentId: String,
    val merchantId: String?,
    val profileId: String?,
    val customerId: String?,
    val paymentMethodId: String?,
    val mandateId: String?,
    val connector: String?,
    val connectorTransactionId: String?,
    val amount: Long,
    val amountCapturable: Long?,
    val amountCaptured: Long?,
    val surchargeAmount: Long?,
    val taxAmount: Long?,
    val currency: String,
    val status: String,
    val captureMethod: String,
    val authenticationType: String?,
    val paymentMethod: String?,
    val paymentMethodType: String?,
    val clientSecret: String?,
    val setupFutureUsage: String?,
    val offSession: Boolean,
    val description: String?,
    val returnUrl: String?,
    val statementDescriptor: String?,
    val billingAddress: Map<String, Any>?,
    val shippingAddress: Map<String, Any>?,
    val metadata: Map<String, Any>?,
    val errorCode: String?,
    val errorMessage: String?,
    val confirmedAt: LocalDateTime?,
    val capturedAt: LocalDateTime?,
    val cancelledAt: LocalDateTime?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

data class PaymentAttemptResult(
    val id: Long,
    val uid: UUID,
    val attemptId: String,
    val paymentId: String,
    val connector: String?,
    val connectorTransactionId: String?,
    val amount: Long,
    val currency: String,
    val status: String,
    val authenticationType: String?,
    val paymentMethod: String?,
    val paymentMethodType: String?,
    val errorCode: String?,
    val errorMessage: String?,
    val createdAt: LocalDateTime,
)

// ─── Refund Aggregate ───────────────────────────────────────────────────────

data class RefundResult(
    val id: Long,
    val uid: UUID,
    val refundId: String,
    val paymentId: String,
    val connector: String?,
    val connectorRefundId: String?,
    val amount: Long,
    val currency: String,
    val status: String,
    val reason: String?,
    val refundType: String,
    val errorCode: String?,
    val errorMessage: String?,
    val metadata: Map<String, Any>?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

// ─── Mandate Aggregate ──────────────────────────────────────────────────────

data class MandateResult(
    val id: Long,
    val uid: UUID,
    val mandateId: String,
    val customerId: String,
    val paymentMethodId: String?,
    val mandateType: String,
    val mandateStatus: String,
    val mandateAmount: Long?,
    val mandateCurrency: String?,
    val startDate: LocalDateTime?,
    val endDate: LocalDateTime?,
    val setupFutureUsage: String?,
    val customerAcceptanceType: String?,
    val customerAcceptedAt: LocalDateTime?,
    val metadata: Map<String, Any>?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)

// ─── Webhook Event ──────────────────────────────────────────────────────────

data class WebhookEventResult(
    val id: Long,
    val uid: UUID,
    val eventId: String,
    val eventType: String,
    val objectType: String,
    val objectId: String,
    val processed: Boolean,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
)
