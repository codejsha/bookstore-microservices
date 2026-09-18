package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.application.SettlementRunLockLostException
import com.codejsha.bookstore.settlement.application.port.SettlementRunLock
import org.springframework.stereotype.Component
import java.time.LocalDate

@Component
class SettlementRunLockGuard(
    private val runLock: SettlementRunLock,
) {

    fun verifyOwnership(targetDate: LocalDate, ownerToken: String?) {
        if (ownerToken.isNullOrBlank()) {
            throw SettlementRunLockLostException(
                "Settlement execution for $targetDate carries no runUid and cannot prove lock ownership",
            )
        }
        if (!runLock.isHeldBy(targetDate, ownerToken)) {
            throw SettlementRunLockLostException(
                "Settlement run $ownerToken no longer holds the lock for $targetDate; " +
                    "refusing to mutate settlement data",
            )
        }
    }
}
