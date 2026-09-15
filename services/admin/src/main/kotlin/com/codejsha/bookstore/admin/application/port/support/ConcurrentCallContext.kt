package com.codejsha.bookstore.admin.application.port.support

import kotlin.coroutines.CoroutineContext

interface ConcurrentCallContext {

    fun capture(): CoroutineContext
}
