package com.codejsha.bookstore.admin.domain.model.external

import java.time.OffsetDateTime

data class RiskEntry(
    val userUid: String,
    val level: String,
    val reason: String,
    val flaggedBy: String?,
    val flaggedAt: OffsetDateTime,
    val expiresAt: OffsetDateTime,
)
