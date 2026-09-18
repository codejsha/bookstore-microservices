package com.codejsha.bookstore.settlement.application.port

import java.time.LocalDate

interface SettlementRunLock {

    fun tryAcquire(targetDate: LocalDate, ownerToken: String): Boolean

    fun isHeldBy(targetDate: LocalDate, ownerToken: String): Boolean

    fun release(targetDate: LocalDate, ownerToken: String): Boolean
}
