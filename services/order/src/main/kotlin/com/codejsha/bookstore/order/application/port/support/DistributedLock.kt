package com.codejsha.bookstore.order.application.port.support

import java.time.Duration

interface DistributedLock {

    fun tryLock(key: String, ttl: Duration? = null): String?

    fun unlock(key: String, token: String)

    fun <T> withLock(key: String, ttl: Duration? = null, action: () -> T): T
}
