package com.codejsha.bookstore.payment.domain.model.command

data class WebhookEventCreateCommand(
    val eventId: String,
    val eventType: String,
    val objectType: String,
    val objectId: String,
    val payload: String,
    val signature: String?,
) {
    init {
        requireNonBlank("event_id", eventId)
        requireMaxLength("event_id", eventId, 64)
        requireNonBlank("event_type", eventType)
        requireMaxLength("event_type", eventType, 64)
        requireNonBlank("object_type", objectType)
        requireMaxLength("object_type", objectType, 32)
        requireNonBlank("object_id", objectId)
        requireMaxLength("object_id", objectId, 64)
        requireNonBlank("payload", payload)
    }
}
