package com.codejsha.bookstore.payment.domain.model

import org.junit.jupiter.api.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotEquals

class MoneyTest {

    @Test
    fun `moneyEquals_whenAmountAndCurrencyMatch_returnsTrueWithSameHashCode`() {
        val a = Money(amount = 1_000L, currency = "KRW")
        val b = Money(amount = 1_000L, currency = "KRW")
        assertEquals(a, b)
        assertEquals(a.hashCode(), b.hashCode())
    }

    @Test
    fun `moneyEquals_whenAmountDiffers_returnsFalse`() {
        assertNotEquals(Money(1_000L, "KRW"), Money(2_000L, "KRW"))
    }

    @Test
    fun `moneyEquals_whenCurrencyDiffers_returnsFalse`() {
        assertNotEquals(Money(1_000L, "KRW"), Money(1_000L, "USD"))
    }

    @Test
    fun `moneyCopy_whenOneFieldGiven_changesOnlyThatField`() {
        val original = Money(1_000L, "KRW")
        val withDifferentCurrency = original.copy(currency = "USD")
        assertEquals(1_000L, withDifferentCurrency.amount)
        assertEquals("USD", withDifferentCurrency.currency)
    }
}

class ConnectorInfoTest {

    @Test
    fun `connectorInfo_whenTransactionIdGiven_holdsBothFields`() {
        val info = ConnectorInfo(connector = "stripe", transactionId = "tx_001")
        assertEquals("stripe", info.connector)
        assertEquals("tx_001", info.transactionId)
    }

    @Test
    fun `connectorInfo_whenTransactionIdNull_holdsNull`() {
        val info = ConnectorInfo(connector = "stripe", transactionId = null)
        assertEquals(null, info.transactionId)
    }

    @Test
    fun `connectorInfoEquals_whenFieldsMatch_returnsTrue`() {
        assertEquals(
            ConnectorInfo("stripe", "tx_1"),
            ConnectorInfo("stripe", "tx_1"),
        )
    }
}
