package com.codejsha.bookstore.settlement.domain.model.option

import java.time.LocalDate

data class SettlementQueryOption(
    val settlementDate: LocalDate? = null,
    val status: String? = null,
)
