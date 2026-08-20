package com.codejsha.bookstore.order.application.port.support

interface IdempotencyService {

    fun isProcessed(idempotencyKey: String): Boolean

    fun getResult(idempotencyKey: String): String?

    fun tryMarkProcessed(idempotencyKey: String, resultUid: String): Boolean

    fun markProcessed(idempotencyKey: String, resultUid: String): Boolean
}
