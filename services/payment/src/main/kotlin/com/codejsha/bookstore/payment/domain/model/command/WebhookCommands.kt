package com.codejsha.bookstore.payment.domain.model.command

data class WebhookEventCreateCommand(
    val eventId: String,
    val eventType: String,
    val objectType: String,
    val objectId: String,
    val payload: String,
    val signature: String?,
)
