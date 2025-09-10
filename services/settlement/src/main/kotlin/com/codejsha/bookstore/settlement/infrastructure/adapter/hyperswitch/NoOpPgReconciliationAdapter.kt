package com.codejsha.bookstore.settlement.infrastructure.adapter.hyperswitch

import com.codejsha.bookstore.settlement.application.port.PgReconciliationPort
import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import jakarta.annotation.PostConstruct
import org.springframework.stereotype.Component
import java.time.LocalDate
import java.time.LocalDateTime

@Component
class NoOpPgReconciliationAdapter(
    private val properties: SettlementBatchProperties,
) : PgReconciliationPort {

    @PostConstruct
    fun rejectEnabledFlag() {
        check(!properties.reconcile.hyperswitch.enabled) {
            "settlement.reconcile.hyperswitch.enabled=true but Hyperswitch reconciliation " +
                "is not wired yet (Phase 5); disable the flag or provide a real PgReconciliationPort"
        }
    }

    override fun fetchDailyTotals(
        targetDate: LocalDate,
        startUtc: LocalDateTime,
        endUtc: LocalDateTime,
    ): List<PgReconciliationPort.CurrencyTotal> =
        throw UnsupportedOperationException(
            "Hyperswitch reconciliation is not wired yet (Phase 5). " +
                "Keep settlement.reconcile.hyperswitch.enabled=false.",
        )
}
