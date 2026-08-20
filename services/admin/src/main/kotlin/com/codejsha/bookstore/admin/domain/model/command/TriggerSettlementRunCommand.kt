package com.codejsha.bookstore.admin.domain.model.command

import java.time.LocalDate

data class TriggerSettlementRunCommand(
    val targetDate: LocalDate,
    val rerun: Boolean,
)
