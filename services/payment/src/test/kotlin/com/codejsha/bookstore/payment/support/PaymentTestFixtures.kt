package com.codejsha.bookstore.payment.support

import com.codejsha.bookstore.payment.application.port.repo.CustomerResult
import com.codejsha.bookstore.payment.application.port.repo.MandateResult
import com.codejsha.bookstore.payment.application.port.repo.PaymentAttemptResult
import com.codejsha.bookstore.payment.application.port.repo.PaymentMethodResult
import com.codejsha.bookstore.payment.application.port.repo.PaymentResult
import com.codejsha.bookstore.payment.application.port.repo.RefundResult
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAggregate
import com.codejsha.bookstore.payment.domain.aggregate.RefundAggregate
import com.codejsha.bookstore.payment.domain.constant.CaptureMethod
import com.codejsha.bookstore.payment.domain.constant.PaymentStatus
import com.codejsha.bookstore.payment.domain.constant.RefundStatus
import com.codejsha.bookstore.payment.domain.constant.RefundType
import com.codejsha.bookstore.payment.domain.model.Money
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import java.time.LocalDateTime
import java.util.UUID

object PaymentTestFixtures {

    val FIXED_TIME: LocalDateTime = LocalDateTime.of(2026, 1, 1, 0, 0, 0)

    val DEFAULT_CONTEXT: ActorContext = ActorContext(
        actorId = 1L,
        actorType = ActorType.USER,
    )

    val PAYMENT_UID: UUID = UUID.fromString("11111111-1111-1111-1111-111111111111")
    val ATTEMPT_UID: UUID = UUID.fromString("22222222-2222-2222-2222-222222222222")
    val REFUND_UID: UUID = UUID.fromString("33333333-3333-3333-3333-333333333333")
    val CUSTOMER_UID: UUID = UUID.fromString("44444444-4444-4444-4444-444444444444")
    val PAYMENT_METHOD_UID: UUID = UUID.fromString("55555555-5555-5555-5555-555555555555")
    val MANDATE_UID: UUID = UUID.fromString("66666666-6666-6666-6666-666666666666")

    fun paymentResult(
        id: Long = 1L,
        uid: UUID = PAYMENT_UID,
        paymentId: String = "pay_001",
        amount: Long = 10_000L,
        currency: String = "KRW",
        status: String = "succeeded",
        captureMethod: String = "automatic",
        connector: String? = "stripe",
        authenticationType: String? = "no_three_ds",
        paymentMethod: String? = "card",
        setupFutureUsage: String? = null,
    ): PaymentResult = PaymentResult(
        id = id,
        uid = uid,
        paymentId = paymentId,
        merchantId = "merchant_1",
        profileId = "profile_1",
        customerId = "cus_1",
        paymentMethodId = "pm_1",
        mandateId = null,
        connector = connector,
        connectorTransactionId = "tx_1",
        amount = amount,
        amountCapturable = amount,
        amountCaptured = amount,
        surchargeAmount = null,
        taxAmount = null,
        currency = currency,
        status = status,
        captureMethod = captureMethod,
        authenticationType = authenticationType,
        paymentMethod = paymentMethod,
        paymentMethodType = "credit",
        clientSecret = "secret_xyz",
        setupFutureUsage = setupFutureUsage,
        offSession = false,
        description = "test payment",
        returnUrl = null,
        statementDescriptor = null,
        billingAddress = null,
        shippingAddress = null,
        metadata = mapOf("k" to "v"),
        errorCode = null,
        errorMessage = null,
        confirmedAt = FIXED_TIME,
        capturedAt = FIXED_TIME,
        cancelledAt = null,
        createdAt = FIXED_TIME,
        updatedAt = FIXED_TIME,
    )

    fun paymentAttemptResult(
        id: Long = 1L,
        uid: UUID = ATTEMPT_UID,
        attemptId: String = "att_001",
        paymentId: String = "pay_001",
        amount: Long = 10_000L,
        currency: String = "KRW",
        status: String = "succeeded",
    ): PaymentAttemptResult = PaymentAttemptResult(
        id = id,
        uid = uid,
        attemptId = attemptId,
        paymentId = paymentId,
        connector = "stripe",
        connectorTransactionId = "tx_1",
        amount = amount,
        currency = currency,
        status = status,
        authenticationType = "no_three_ds",
        paymentMethod = "card",
        paymentMethodType = "credit",
        errorCode = null,
        errorMessage = null,
        createdAt = FIXED_TIME,
    )

    fun refundResult(
        id: Long = 1L,
        uid: UUID = REFUND_UID,
        refundId: String = "ref_001",
        paymentId: String = "pay_001",
        amount: Long = 5_000L,
        currency: String = "KRW",
        status: String = "succeeded",
        refundType: String = "instant",
    ): RefundResult = RefundResult(
        id = id,
        uid = uid,
        refundId = refundId,
        paymentId = paymentId,
        connector = "stripe",
        connectorRefundId = "rf_xx",
        amount = amount,
        currency = currency,
        status = status,
        reason = "requested_by_customer",
        refundType = refundType,
        errorCode = null,
        errorMessage = null,
        metadata = null,
        createdAt = FIXED_TIME,
        updatedAt = FIXED_TIME,
    )

    fun refundAggregate(
        id: Long = 1L,
        uid: UUID = REFUND_UID,
        refundId: String = "ref_001",
        paymentId: String = "pay_001",
        amount: Long = 5_000L,
        currency: String = "KRW",
        status: RefundStatus = RefundStatus.SUCCEEDED,
        reason: String? = "requested_by_customer",
        refundType: RefundType = RefundType.INSTANT,
        errorCode: String? = null,
        errorMessage: String? = null,
    ): RefundAggregate = RefundAggregate(
        id = id,
        uid = uid,
        refundId = refundId,
        paymentId = paymentId,
        connector = "stripe",
        connectorRefundId = "rf_xx",
        amount = Money(amount = amount, currency = currency),
        status = status,
        reason = reason,
        refundType = refundType,
        errorCode = errorCode,
        errorMessage = errorMessage,
        metadata = null,
        createdAt = FIXED_TIME,
        updatedAt = FIXED_TIME,
    )

    fun paymentAggregate(
        uid: UUID = PAYMENT_UID,
        paymentId: String = "pay_001",
        amount: Long = 10_000L,
        currency: String = "KRW",
        status: PaymentStatus = PaymentStatus.SUCCEEDED,
        customerId: String = "cus_1",
    ): PaymentAggregate = PaymentAggregate(
        id = 1L,
        uid = uid,
        paymentId = paymentId,
        merchantId = "merchant_1",
        profileId = "profile_1",
        customerId = customerId,
        paymentMethodId = "pm_1",
        mandateId = null,
        connector = null,
        amount = Money(amount = amount, currency = currency),
        amountCapturable = amount,
        amountCaptured = amount,
        surchargeAmount = null,
        taxAmount = null,
        status = status,
        captureMethod = CaptureMethod.AUTOMATIC,
        authenticationType = null,
        paymentMethod = null,
        paymentMethodType = "credit",
        clientSecret = "secret_xyz",
        setupFutureUsage = null,
        offSession = false,
        description = "test payment",
        returnUrl = null,
        statementDescriptor = null,
        billingAddress = null,
        shippingAddress = null,
        metadata = null,
        errorCode = null,
        errorMessage = null,
        confirmedAt = FIXED_TIME,
        capturedAt = FIXED_TIME,
        cancelledAt = null,
        createdAt = FIXED_TIME,
        updatedAt = FIXED_TIME,
        attempts = emptyList(),
    )

    fun customerResult(
        id: Long = 1L,
        uid: UUID = CUSTOMER_UID,
        customerId: String = "cus_001",
        email: String? = "alice@example.com",
    ): CustomerResult = CustomerResult(
        id = id,
        uid = uid,
        customerId = customerId,
        name = "Alice",
        email = email,
        phone = "1234",
        phoneCountryCode = "+82",
        description = "test",
        metadata = null,
        defaultBillingAddress = null,
        defaultShippingAddress = null,
        createdAt = FIXED_TIME,
        updatedAt = FIXED_TIME,
    )

    fun paymentMethodResult(
        id: Long = 1L,
        uid: UUID = PAYMENT_METHOD_UID,
        paymentMethodId: String = "pm_001",
        customerId: String = "cus_001",
        paymentMethod: String = "card",
        isDefault: Boolean = true,
    ): PaymentMethodResult = PaymentMethodResult(
        id = id,
        uid = uid,
        paymentMethodId = paymentMethodId,
        customerId = customerId,
        paymentMethod = paymentMethod,
        paymentMethodType = "credit",
        paymentMethodIssuer = "Visa",
        cardNetwork = "visa",
        cardLast4 = "4242",
        cardExpMonth = 12,
        cardExpYear = 2030,
        cardHolderName = "Alice",
        isDefault = isDefault,
        metadata = null,
        createdAt = FIXED_TIME,
        updatedAt = FIXED_TIME,
    )

    fun mandateResult(
        id: Long = 1L,
        uid: UUID = MANDATE_UID,
        mandateId: String = "mand_001",
        customerId: String = "cus_001",
        mandateType: String = "multi_use",
        mandateStatus: String = "active",
        setupFutureUsage: String? = "off_session",
    ): MandateResult = MandateResult(
        id = id,
        uid = uid,
        mandateId = mandateId,
        customerId = customerId,
        paymentMethodId = "pm_001",
        mandateType = mandateType,
        mandateStatus = mandateStatus,
        mandateAmount = 100_000L,
        mandateCurrency = "KRW",
        startDate = FIXED_TIME,
        endDate = null,
        setupFutureUsage = setupFutureUsage,
        customerAcceptanceType = "online",
        customerAcceptedAt = FIXED_TIME,
        metadata = null,
        createdAt = FIXED_TIME,
        updatedAt = FIXED_TIME,
    )
}
