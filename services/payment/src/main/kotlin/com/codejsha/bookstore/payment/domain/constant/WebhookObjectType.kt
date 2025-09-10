package com.codejsha.bookstore.payment.domain.constant

enum class WebhookObjectType(
    val value: String
) {
    PAYMENT("payment"),
    REFUND("refund"),
    MANDATE("mandate"),
    DISPUTE("dispute"),
    OTHER("other");
}
