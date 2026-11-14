package com.codejsha.bookstore.settlement.infrastructure.adapter.hyperswitch

import com.codejsha.bookstore.settlement.config.properties.HyperswitchReconcileProperties
import com.codejsha.bookstore.settlement.config.properties.ReconcileProperties
import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import org.junit.jupiter.api.Test
import java.time.LocalDate
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

class NoOpPgReconciliationAdapterTest {

    private val date: LocalDate = LocalDate.of(2026, 7, 1)

    @Test
    fun `rejectEnabledFlag_reconciliationDisabled_passes`() {
        adapter(enabled = false).rejectEnabledFlag()
    }

    @Test
    fun `rejectEnabledFlag_reconciliationEnabled_failsStartup`() {
        val error = assertFailsWith<IllegalStateException> { adapter(enabled = true).rejectEnabledFlag() }

        assertTrue(error.message!!.contains("settlement.reconcile.hyperswitch.enabled=true"))
    }

    @Test
    fun `fetchDailyTotals_anyWindow_throwsUnsupported`() {
        assertFailsWith<UnsupportedOperationException> {
            adapter(enabled = false).fetchDailyTotals(date, date.atStartOfDay(), date.plusDays(1).atStartOfDay())
        }
    }

    private fun adapter(enabled: Boolean) = NoOpPgReconciliationAdapter(
        SettlementBatchProperties(
            reconcile = ReconcileProperties(hyperswitch = HyperswitchReconcileProperties(enabled = enabled)),
        ),
    )
}
