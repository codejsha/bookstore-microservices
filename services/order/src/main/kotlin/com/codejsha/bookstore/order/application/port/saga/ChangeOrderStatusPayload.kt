package com.codejsha.bookstore.order.application.port.saga

class ChangeOrderStatusPayload(
    val requestId: String,
    val orderId: Long,
)
