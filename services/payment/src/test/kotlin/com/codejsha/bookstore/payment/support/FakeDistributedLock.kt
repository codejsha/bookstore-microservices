package com.codejsha.bookstore.payment.support

import com.codejsha.bookstore.payment.application.port.support.DistributedLock
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
