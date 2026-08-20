package com.codejsha.bookstore.admin.domain.model.external

import java.time.OffsetDateTime

data class Author(
    val uid: String,
    val name: String,
)

data class Subject(
    val uid: String,
    val name: String,
)

data class Work(
    val uid: String,
    val title: String,
    val description: String?,
    val firstPublishDate: String?,
    val olKey: String?,
    val authors: List<Author>,
    val subjects: List<Subject>,
    val createdAt: OffsetDateTime,
    val updatedAt: OffsetDateTime?,
)
