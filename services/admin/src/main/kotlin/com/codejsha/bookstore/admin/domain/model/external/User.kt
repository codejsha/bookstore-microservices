package com.codejsha.bookstore.admin.domain.model.external

import java.time.OffsetDateTime

data class User(
    val uid: String,
    val email: String,
    val firstName: String,
    val lastName: String,
    val phone: String?,
    val status: String,
    val roles: List<String>,
    val lastLoginAt: OffsetDateTime?,
    val createdAt: OffsetDateTime,
    val updatedAt: OffsetDateTime?,
)
