package com.codejsha.bookstore.payment.domain.aggregate

import com.codejsha.bookstore.payment.domain.model.Address
import com.codejsha.bookstore.payment.domain.model.ConnectorInfo
import com.codejsha.bookstore.payment.domain.model.Money
import com.codejsha.bookstore.payment.domain.constant.AuthenticationType
import com.codejsha.bookstore.payment.domain.constant.CaptureMethod
import com.codejsha.bookstore.payment.domain.constant.FutureUsage
import com.codejsha.bookstore.payment.domain.constant.PaymentMethodType
import com.codejsha.bookstore.payment.domain.constant.PaymentStatus
import java.time.LocalDateTime
import java.util.UUID

data class PaymentAggregate(
    val id: Long,
    val uid: UUID,
    val paymentId: String,
    val merchantId: String?,
    val profileId: String?,
    val customerId: String?,
    val paymentMethodId: String?,
    val mandateId: String?,
    val connector: ConnectorInfo?,
    val amount: Money,
    val amountCapturable: Long?,
    val amountCaptured: Long?,
    val surchargeAmount: Long?,
    val taxAmount: Long?,
    val status: PaymentStatus,
    val captureMethod: CaptureMethod,
    val authenticationType: AuthenticationType?,
    val paymentMethod: PaymentMethodType?,
    val paymentMethodType: String?,
    val clientSecret: String?,
    val setupFutureUsage: FutureUsage?,
    val offSession: Boolean,
    val description: String?,
    val returnUrl: String?,
    val statementDescriptor: String?,
    val billingAddress: Address?,
    val shippingAddress: Address?,
    val metadata: Map<String, Any>?,
    val errorCode: String?,
    val errorMessage: String?,
    val confirmedAt: LocalDateTime?,
    val capturedAt: LocalDateTime?,
    val cancelledAt: LocalDateTime?,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime?,
    val attempts: List<PaymentAttemptEntity>,
)

data class PaymentAttemptEntity(
    val id: Long,
    val uid: UUID,
    val attemptId: String,
    val paymentId: String,
    val connector: ConnectorInfo?,
    val amount: Money,
    val status: PaymentStatus,
    val authenticationType: AuthenticationType?,
    val paymentMethod: PaymentMethodType?,
    val paymentMethodType: String?,
    val errorCode: String?,
    val errorMessage: String?,
    val createdAt: LocalDateTime,
)
