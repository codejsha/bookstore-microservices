package com.codejsha.bookstore.admin.domain.model.command

data class FlagRiskCommand(
    val level: String,
    val reason: String,
    val ttlSeconds: Long?,
)
