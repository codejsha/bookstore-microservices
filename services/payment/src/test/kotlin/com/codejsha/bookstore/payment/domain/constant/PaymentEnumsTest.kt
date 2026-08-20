package com.codejsha.bookstore.payment.domain.constant

import org.junit.jupiter.api.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class PaymentEnumsTest {

    @Test
    fun `PaymentStatus round-trip and lookup failure`() {
        PaymentStatus.entries.forEach { status ->
            assertEquals(status, PaymentStatus.fromValue(status.value))
        }
        assertFailsWith<NoSuchElementException> { PaymentStatus.fromValue("does_not_exist") }
    }

    @Test
    fun `RefundStatus round-trip and lookup failure`() {
        RefundStatus.entries.forEach { assertEquals(it, RefundStatus.fromValue(it.value)) }
        assertFailsWith<NoSuchElementException> { RefundStatus.fromValue("nope") }
    }

    @Test
    fun `RefundType round-trip`() {
        RefundType.entries.forEach { assertEquals(it, RefundType.fromValue(it.value)) }
    }

    @Test
    fun `CaptureMethod round-trip`() {
        CaptureMethod.entries.forEach { assertEquals(it, CaptureMethod.fromValue(it.value)) }
    }

    @Test
    fun `MandateStatus and MandateType round-trip`() {
        MandateStatus.entries.forEach { assertEquals(it, MandateStatus.fromValue(it.value)) }
        MandateType.entries.forEach { assertEquals(it, MandateType.fromValue(it.value)) }
    }

    @Test
    fun `AuthenticationType round-trip`() {
        AuthenticationType.entries.forEach { assertEquals(it, AuthenticationType.fromValue(it.value)) }
    }

    @Test
    fun `PaymentMethodType round-trip`() {
        PaymentMethodType.entries.forEach { assertEquals(it, PaymentMethodType.fromValue(it.value)) }
    }

    @Test
    fun `FutureUsage round-trip`() {
        FutureUsage.entries.forEach { assertEquals(it, FutureUsage.fromValue(it.value)) }
    }
}
