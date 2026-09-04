package com.codejsha.bookstore.payment.domain.constant

import org.junit.jupiter.api.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class PaymentEnumsTest {

    // Every known value is round-tripped through fromValue first.
    @Test
    fun `paymentStatusFromValue_whenValueUnknown_throwsNoSuchElementException`() {
        PaymentStatus.entries.forEach { status ->
            assertEquals(status, PaymentStatus.fromValue(status.value))
        }
        assertFailsWith<NoSuchElementException> { PaymentStatus.fromValue("does_not_exist") }
    }

    // Every known value is round-tripped through fromValue first.
    @Test
    fun `refundStatusFromValue_whenValueUnknown_throwsNoSuchElementException`() {
        RefundStatus.entries.forEach { assertEquals(it, RefundStatus.fromValue(it.value)) }
        assertFailsWith<NoSuchElementException> { RefundStatus.fromValue("nope") }
    }

    @Test
    fun `refundTypeFromValue_whenValueKnown_returnsMatchingEntry`() {
        RefundType.entries.forEach { assertEquals(it, RefundType.fromValue(it.value)) }
    }

    @Test
    fun `captureMethodFromValue_whenValueKnown_returnsMatchingEntry`() {
        CaptureMethod.entries.forEach { assertEquals(it, CaptureMethod.fromValue(it.value)) }
    }

    @Test
    fun `mandateStatusAndTypeFromValue_whenValueKnown_returnMatchingEntry`() {
        MandateStatus.entries.forEach { assertEquals(it, MandateStatus.fromValue(it.value)) }
        MandateType.entries.forEach { assertEquals(it, MandateType.fromValue(it.value)) }
    }

    @Test
    fun `authenticationTypeFromValue_whenValueKnown_returnsMatchingEntry`() {
        AuthenticationType.entries.forEach { assertEquals(it, AuthenticationType.fromValue(it.value)) }
    }

    @Test
    fun `paymentMethodTypeFromValue_whenValueKnown_returnsMatchingEntry`() {
        PaymentMethodType.entries.forEach { assertEquals(it, PaymentMethodType.fromValue(it.value)) }
    }

    @Test
    fun `futureUsageFromValue_whenValueKnown_returnsMatchingEntry`() {
        FutureUsage.entries.forEach { assertEquals(it, FutureUsage.fromValue(it.value)) }
    }
}
