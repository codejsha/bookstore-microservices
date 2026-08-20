package com.codejsha.bookstore.admin.domain.model.option

import java.time.LocalDate

data class SettlementQueryOption(
    val settlementDate: LocalDate? = null,
    val status: String? = null,
)
