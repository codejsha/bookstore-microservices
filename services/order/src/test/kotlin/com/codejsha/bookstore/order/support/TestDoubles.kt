package com.codejsha.bookstore.order.support

import com.codejsha.bookstore.order.application.port.support.DistributedLock
import com.codejsha.bookstore.order.application.port.support.IdempotencyService
import com.codejsha.bookstore.order.application.port.support.TransactionRunner
import java.time.Duration

class FakeDistributedLock : DistributedLock {
    var lastKey: String? = null

    var lastTtl: Duration? = null

    var invocationCount: Int = 0

    var lastUnlockedToken: String? = null

    override fun <T> withLock(key: String, ttl: Duration?, action: () -> T): T {
        lastKey = key
        lastTtl = ttl
        invocationCount++
        return action()
    }

    override fun tryLock(key: String, ttl: Duration?): String? {
        lastKey = key
        lastTtl = ttl
        invocationCount++
        return TEST_TOKEN
    }

    override fun unlock(key: String, token: String) {
        lastUnlockedToken = token
    }

    companion object {
        const val TEST_TOKEN: String = "test-token"
    }
}

class FakeIdempotencyService : IdempotencyService {
    private val store: MutableMap<String, String> = mutableMapOf()

    override fun isProcessed(idempotencyKey: String): Boolean = store.containsKey(idempotencyKey)

    override fun getResult(idempotencyKey: String): String? = store[idempotencyKey]

    override fun tryMarkProcessed(idempotencyKey: String, resultUid: String): Boolean =
        store.putIfAbsent(idempotencyKey, resultUid) == null

    override fun markProcessed(idempotencyKey: String, resultUid: String): Boolean =
        tryMarkProcessed(idempotencyKey, resultUid)
}

class FakeTransactionRunner : TransactionRunner {
    override suspend fun <T> tx(block: () -> T): T = block()
}
