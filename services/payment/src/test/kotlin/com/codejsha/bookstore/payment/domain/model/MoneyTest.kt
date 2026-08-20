package com.codejsha.bookstore.payment.domain.model

import org.junit.jupiter.api.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotEquals

class MoneyTest {

    @Test
    fun `Money equality is value-based`() {
        val a = Money(amount = 1_000L, currency = "KRW")
        val b = Money(amount = 1_000L, currency = "KRW")
        assertEquals(a, b)
        assertEquals(a.hashCode(), b.hashCode())
    }

    @Test
    fun `Money differs when amount differs`() {
        assertNotEquals(Money(1_000L, "KRW"), Money(2_000L, "KRW"))
    }

    @Test
    fun `Money differs when currency differs`() {
        assertNotEquals(Money(1_000L, "KRW"), Money(1_000L, "USD"))
    }

    @Test
    fun `Money copy mutates only requested fields`() {
        val original = Money(1_000L, "KRW")
        val withDifferentCurrency = original.copy(currency = "USD")
        assertEquals(1_000L, withDifferentCurrency.amount)
        assertEquals("USD", withDifferentCurrency.currency)
    }
}

class ConnectorInfoTest {

    @Test
    fun `ConnectorInfo holds connector and optional transaction id`() {
        val info = ConnectorInfo(connector = "stripe", transactionId = "tx_001")
        assertEquals("stripe", info.connector)
        assertEquals("tx_001", info.transactionId)
    }

    @Test
    fun `ConnectorInfo allows null transaction id`() {
        val info = ConnectorInfo(connector = "stripe", transactionId = null)
        assertEquals(null, info.transactionId)
    }

    @Test
    fun `ConnectorInfo equality is value-based`() {
        assertEquals(
            ConnectorInfo("stripe", "tx_1"),
            ConnectorInfo("stripe", "tx_1"),
        )
    }
}
