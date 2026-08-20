package com.codejsha.bookstore.admin.domain.model.command

data class WorkCreateCommand(
    val title: String,
    val description: String?,
    val firstPublishDate: String?,
    val olKey: String?,
    val authorUids: List<String>,
    val subjectNames: List<String>?,
)
