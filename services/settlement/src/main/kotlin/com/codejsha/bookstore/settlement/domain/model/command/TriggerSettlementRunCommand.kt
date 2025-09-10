package com.codejsha.bookstore.settlement.domain.model.command

import java.time.LocalDate

data class TriggerSettlementRunCommand(
    val targetDate: LocalDate,
    val rerun: Boolean,
)
