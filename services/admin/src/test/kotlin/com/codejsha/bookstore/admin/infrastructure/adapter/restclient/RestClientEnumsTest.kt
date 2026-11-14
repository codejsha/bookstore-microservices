package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException
import com.codejsha.bookstore.generated.application.port.restclient.payment.model.PaymentStatus
import org.junit.jupiter.api.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class RestClientEnumsTest {

    @Test
    fun `parseEnumValue_knownWireValue_resolvesConstant`() {
        val status = parseEnumValue("payment status", "requires_capture", PaymentStatus.Companion::fromValue)

        assertEquals(PaymentStatus.REQUIRES_CAPTURE, status)
    }

    @Test
    fun `parseEnumValue_constantNameInsteadOfWireValue_failsAsUnsupportedValue`() {
        val failure = assertFailsWith<UnsupportedValueException> {
            parseEnumValue("payment status", "REQUIRES_CAPTURE", PaymentStatus.Companion::fromValue)
        }

        assertEquals("unknown payment status: REQUIRES_CAPTURE", failure.message)
    }

    @Test
    fun `parseEnumValue_parserThrowingIllegalArgument_failsAsUnsupportedValue`() {
        val failure = assertFailsWith<UnsupportedValueException> {
            parseEnumValue("payment status", "bogus") { PaymentStatus.valueOf(it) }
        }

        assertEquals("unknown payment status: bogus", failure.message)
    }
}
